package explorer

import (
	"context"
	"encoding/binary"
	"errors"
	"strings"
	"testing"
	"testing/synctest"
	"time"

	"github.com/anyshake/observer/pkg/fifo"
	"github.com/anyshake/observer/pkg/message"
	"github.com/anyshake/observer/pkg/timesource"
)

func TestV1OpenDecodesWirePackets(t *testing.T) {
	transport := newScriptTransport()
	transport.delay = 100 * time.Millisecond
	protocol := &ExplorerProtoImplV1{
		ChannelCodes:    []string{"HNZ", "HNE", "HNN"},
		ExplorerOptions: ExplorerOptions{Protocol: "v1", Model: "models/E-C111G", Latitude: 25, Longitude: 121, Elevation: 10},
		NtpOptions:      NtpOptions{Pool: []string{startTestNTP(t)}, Retry: 0, ReadTimeout: 1},
		Logger:          testLogger(nil),
		TimeSource:      timesource.New(time.Now),
		Transport:       transport,
	}
	t.Cleanup(func() { _ = protocol.Close() })

	ctx, cancel, err := protocol.Open(context.Background())
	if err != nil {
		t.Fatalf("Open() error = %v", err)
	}
	defer cancel()
	if ctx.Err() != nil {
		t.Fatalf("Open() returned a canceled context: %v", ctx.Err())
	}

	archived := make(chan Event, 4)
	realtime := make(chan Event, 16)
	if err := protocol.Subscribe("archive", func(event Event) { archived <- event }); err != nil {
		t.Fatal(err)
	}
	if err := protocol.SubscribeRealtime("live", func(event Event) { realtime <- event }); err != nil {
		t.Fatal(err)
	}
	transport.packets <- make([]byte, 66)
	wire := buildV1WirePacket(7)
	for i := 0; i < 10; i++ {
		transport.packets <- wire
	}

	select {
	case event := <-realtime:
		if event.SampleRate != 50 || len(event.ChannelData) != 3 || len(event.ChannelData[0].Data) != 5 || event.ChannelData[0].Data[0] != 7 || event.ChannelData[0].ChannelCode != "HNZ" {
			t.Fatalf("realtime event = %#v", event)
		}
	case <-time.After(3 * time.Second):
		config := protocol.GetConfig()
		t.Fatalf("timed out waiting for realtime data, sample rate %d", config.GetSampleRate())
	}
	select {
	case event := <-archived:
		if len(event.ChannelData) != 3 || len(event.ChannelData[2].Data) != 50 || event.ChannelData[2].ChannelCode != "HNN" {
			t.Fatalf("archived event channels = %d samples %d", len(event.ChannelData), len(event.ChannelData[0].Data))
		}
	case <-time.After(3 * time.Second):
		config := protocol.GetConfig()
		status := protocol.GetStatus()
		t.Fatalf("timed out waiting for archived data, sample rate %d frames %d", config.GetSampleRate(), status.GetFrames())
	}

	corrupt := append([]byte(nil), wire...)
	corrupt[2] ^= 0xFF
	transport.packets <- corrupt
	if !waitUntil(time.Second, func() bool {
		status := protocol.GetStatus()
		return status.GetErrors() > 0
	}) {
		t.Fatal("corrupt packet did not increment errors")
	}
	config := protocol.GetConfig()
	if protocol.GetDeviceId() != "012F81AC" || config.GetProtocol() != "v1" || config.GetModel() != "E-C111G" {
		t.Fatalf("device = %s protocol=%s model=%s", protocol.GetDeviceId(), config.GetProtocol(), config.GetModel())
	}
	if err := protocol.Close(); err != nil {
		t.Fatal(err)
	}
	assertStreamWorkersStopped(t, protocol.streamDone, protocol.decoderDone, protocol.ntpDone)
}

func TestV1OpenFlushFailure(t *testing.T) {
	transport := &fakeTransport{flushErr: errors.New("flush failed")}
	protocol := &ExplorerProtoImplV1{
		Logger:     testLogger(nil),
		TimeSource: timesource.New(time.Now),
		Transport:  transport,
		NtpOptions: NtpOptions{Pool: []string{startTestNTP(t)}, Retry: 0, ReadTimeout: 1},
	}
	_, _, err := protocol.Open(context.Background())
	if err == nil || !strings.Contains(err.Error(), "flush failed") || transport.closed != 1 || transport.flushed != 1 {
		t.Fatalf("cleanup flushed=%d closed=%d err=%v", transport.flushed, transport.closed, err)
	}
}

func TestV1DecodeStreamArchivesInOrder(t *testing.T) {
	synctest.Test(t, func(t *testing.T) {
		ctx, cancel := context.WithCancel(context.Background())
		protocol := &ExplorerProtoImplV1{
			Logger:             testLogger(nil),
			ChannelCodes:       []string{"HNZ", "HNE", "HNN"},
			fifoBuffer:         fifo.New[byte](74 * 32),
			decoderDone:        make(chan struct{}),
			messageBus:         message.NewBus[Event](EXPLORER_STREAM_TOPIC),
			messageBusRealtime: message.NewBus[Event](EXPLORER_REALTIME_STREAM_TOPIC),
		}
		protocol.deviceConfig.SetSampleRate(50)
		protocol.deviceConfig.SetPacketInterval(100 * time.Millisecond)
		defer protocol.messageBus.Close()
		defer protocol.messageBusRealtime.Close()
		go protocol.decodeStream(ctx, 66)
		defer func() {
			cancel()
			<-protocol.decoderDone
		}()

		archive := make(chan Event, 2)
		live := make(chan Event, 32)
		if err := protocol.Subscribe("archive", func(event Event) { archive <- event }); err != nil {
			t.Fatal(err)
		}
		if err := protocol.SubscribeRealtime("live", func(event Event) { live <- event }); err != nil {
			t.Fatal(err)
		}
		push := func(timestamp int64, sample int32, corrupt bool) {
			t.Helper()
			samples := []int32{sample, sample, sample, sample, sample}
			packet := buildV1DataPacket(2, 5, samples, samples, samples)
			binary.BigEndian.PutUint64(packet[2:10], uint64(timestamp))
			if corrupt {
				packet[10] ^= 0xFF
			}
			_, _ = protocol.fifoBuffer.Write(packet...)
			time.Sleep(10 * time.Millisecond)
			synctest.Wait()
		}

		push(900, -1, true)
		status := protocol.GetStatus()
		if status.GetErrors() != 1 || status.GetFrames() != 0 || len(live) != 0 || len(archive) != 0 {
			t.Fatal("corrupt packet entered the stream")
		}
		for i := 0; i < 20; i++ {
			push(1000+int64(i)*100, int32(i), false)
			select {
			case event := <-live:
				if event.Timestamp.UnixMilli() != 1000+int64(i)*100 || event.SampleRate != 50 || len(event.ChannelData) != 3 || len(event.ChannelData[0].Data) != 5 || event.ChannelData[0].Data[0] != int32(i) {
					t.Fatalf("realtime packet %d = %#v", i, event)
				}
			default:
				t.Fatalf("missing realtime packet %d", i)
			}
		}
		for second := 0; second < 2; second++ {
			select {
			case event := <-archive:
				if event.Timestamp.UnixMilli() != 1000+int64(second)*1000 || len(event.ChannelData) != 3 {
					t.Fatalf("archive %d = %#v", second, event)
				}
				for _, channel := range event.ChannelData {
					if len(channel.Data) != 50 {
						t.Fatalf("archive has %d samples, want 50", len(channel.Data))
					}
					for i, sample := range channel.Data {
						if sample != int32(second*10+i/5) {
							t.Fatalf("archive %d sample %d = %d", second, i, sample)
						}
					}
				}
			default:
				t.Fatalf("missing archive %d", second)
			}
		}
	})
}

func TestV1ReadErrorStopsStreams(t *testing.T) {
	transport := newScriptTransport()
	transport.readErr = errors.New("read failed")
	protocol := &ExplorerProtoImplV1{
		Logger:     testLogger(nil),
		TimeSource: timesource.New(time.Now),
		Transport:  transport,
		NtpOptions: NtpOptions{Pool: []string{startTestNTP(t)}, ReadTimeout: 1},
	}
	t.Cleanup(func() { _ = protocol.Close() })
	ctx, _, err := protocol.Open(context.Background())
	if err != nil {
		t.Fatal(err)
	}
	transport.packets <- buildV1WirePacket(1)
	select {
	case <-ctx.Done():
	case <-time.After(2 * time.Second):
		t.Fatal("read error did not cancel the stream")
	}
	if err := protocol.Close(); err != nil {
		t.Fatal(err)
	}
	assertStreamWorkersStopped(t, protocol.streamDone, protocol.decoderDone, protocol.ntpDone)
}
