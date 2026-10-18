package seedlink

import (
	"context"
	"testing"
	"time"

	"github.com/anyshake/observer/internal/dao/model"
	"github.com/anyshake/observer/internal/hardware/explorer"
	"github.com/anyshake/observer/internal/testsupport"
	"github.com/anyshake/observer/pkg/logger"
	"github.com/anyshake/observer/pkg/ringbuf"
	"github.com/bclswl0827/slgo/handlers"
	"gorm.io/gorm"
)

func TestQueryHistoryIncludesUncommittedTail(t *testing.T) {
	logger.Init()
	database, handler := testsupport.OpenDAO(t)
	started := time.Date(2026, time.October, 18, 1, 0, 0, 0, time.UTC)
	testsupport.EnsureSeisTable(t, database, started)
	recentBuffer := ringbuf.New[explorer.Event](seedLinkHistoryBufferSize)
	for second := -1; second <= 10; second++ {
		event := explorer.Event{
			Timestamp: started.Add(time.Duration(second) * time.Second), SampleRate: 1,
			ChannelData: []explorer.ChannelData{
				{ChannelCode: "EHZ", Data: []int32{int32(second)}},
				{ChannelCode: "EHN", Data: []int32{int32(-second)}},
				{ChannelCode: "EHE", Data: []int32{999}},
			},
		}
		if second <= 1 {
			var record model.SeisRecord
			if err := record.Encode(event.Timestamp, event.SampleRate, event.ChannelData); err != nil {
				t.Fatal(err)
			}
			if err := handler.SeisRecordsCreate(record); err != nil {
				t.Fatal(err)
			}
		}
		if second >= 1 {
			event.ChannelData[0].Data = []int32{int32(second + 100)}
			recentBuffer.Push(event)
		}
	}
	source := &provider{ctx: context.Background(), actionHandler: handler, recentBuffer: recentBuffer}
	packets, err := source.QueryHistory(started, started.Add(10*time.Second), []handlers.SeedLinkChannel{{ChannelName: "EHZ"}, {ChannelName: "EHN"}})
	if err != nil {
		t.Fatal(err)
	}
	if len(packets) != 20 {
		t.Fatalf("history contains %d packets, want 20 without gaps or duplicates", len(packets))
	}
	for second := 0; second < 10; second++ {
		for index, channel := range []string{"EHN", "EHZ"} {
			packet := packets[second*2+index]
			want := int32(-second)
			if channel == "EHZ" {
				want = int32(second)
				if second >= 1 {
					want += 100
				}
			}
			if packet.Timestamp != started.Add(time.Duration(second)*time.Second).UnixMilli() || packet.Channel != channel || len(packet.DataArr) != 1 || packet.DataArr[0] != want {
				t.Fatalf("second %d channel %s: %#v, want sample %d", second, channel, packet, want)
			}
		}
	}
}

func TestQueryHistoryIncludesRecordCrossingStart(t *testing.T) {
	logger.Init()
	database, handler := testsupport.OpenDAO(t)
	started := time.Date(2026, time.October, 18, 0, 0, 0, 0, time.UTC)
	recordTime := started.Add(-250 * time.Millisecond)
	testsupport.EnsureSeisTable(t, database, recordTime)
	testsupport.EnsureSeisTable(t, database, started)
	var record model.SeisRecord
	if err := record.Encode(recordTime, 4, []explorer.ChannelData{{ChannelCode: "EHZ", Data: []int32{1, 2, 3, 4}}}); err != nil {
		t.Fatal(err)
	}
	if err := handler.SeisRecordsCreate(record); err != nil {
		t.Fatal(err)
	}
	source := &provider{ctx: context.Background(), actionHandler: handler}
	// The lookback must also work at midnight and at the DAO's maximum window.
	for _, duration := range []time.Duration{500 * time.Millisecond, handler.SeisRecordsGetQueryWindow()} {
		packets, err := source.QueryHistory(started, started.Add(duration), []handlers.SeedLinkChannel{{ChannelName: "EHZ"}})
		if err != nil {
			t.Fatal(err)
		}
		if len(packets) != 1 || packets[0].Timestamp != recordTime.UnixMilli() || len(packets[0].DataArr) != 4 {
			t.Fatalf("duration %s: history omitted the record crossing start: %#v", duration, packets)
		}
	}
}

func TestQueryHistoryRetainsSnapshotDuringQuery(t *testing.T) {
	logger.Init()
	database, handler := testsupport.OpenDAO(t)
	started := time.Date(2026, time.October, 18, 1, 0, 0, 0, time.UTC)
	testsupport.EnsureSeisTable(t, database, started)
	recentBuffer := ringbuf.New[explorer.Event](2)
	push := func(second int) {
		recentBuffer.Push(explorer.Event{
			Timestamp: started.Add(time.Duration(second) * time.Second), SampleRate: 1,
			ChannelData: []explorer.ChannelData{{ChannelCode: "EHZ", Data: []int32{int32(second)}}},
		})
	}
	push(0)
	push(1)
	// Advance the live ring after QueryHistory takes its snapshot and before
	// the empty database query finishes, without using sleeps.
	if err := database.Database.Callback().Row().Before("gorm:row").Register("seedlink:advance_ring", func(*gorm.DB) {
		push(10)
		push(11)
	}); err != nil {
		t.Fatal(err)
	}
	source := &provider{ctx: context.Background(), actionHandler: handler, recentBuffer: recentBuffer}
	packets, err := source.QueryHistory(started, started.Add(2*time.Second), []handlers.SeedLinkChannel{{ChannelName: "EHZ"}})
	if err != nil {
		t.Fatal(err)
	}
	if len(packets) != 2 || packets[0].Timestamp != started.UnixMilli() || packets[1].Timestamp != started.Add(time.Second).UnixMilli() {
		t.Fatalf("uncommitted history was evicted during the query: %#v", packets)
	}
}
