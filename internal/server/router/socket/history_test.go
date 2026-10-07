package socket

import (
	"testing"
	"time"

	"github.com/anyshake/observer/internal/hardware/explorer"
)

func TestHistoryRingKeepsNewestInOrder(t *testing.T) {
	t.Parallel()

	s := &socket{historyBuffer: make([]buffer, HISTORY_BUFFER_SIZE)}
	const extra = 5
	total := HISTORY_BUFFER_SIZE + extra
	for i := 0; i < total; i++ {
		s.storeHistory(explorer.Event{
			Timestamp:  time.UnixMilli(int64(i + 1)),
			SampleRate: 10,
			ChannelData: []explorer.ChannelData{{
				ChannelCode: "EHZ",
				Data:        []int32{int32(i + 1)},
			}},
		})
	}

	if s.historyLen != HISTORY_BUFFER_SIZE {
		t.Fatalf("history length = %d, want %d", s.historyLen, HISTORY_BUFFER_SIZE)
	}

	oldest := int64(extra + 1)
	for i := 0; i < s.historyLen; i++ {
		got := s.historyAt(i)
		want := oldest + int64(i)
		if got.Timestamp != want || got.SampleRate != 10 || len(got.ChannelData) != 1 || got.ChannelData[0].Data[0] != int32(want) {
			t.Fatalf("history[%d] = timestamp %d data %v, want %d", i, got.Timestamp, got.ChannelData, want)
		}
	}
}

func TestHistoryRingPreservesPartialOrder(t *testing.T) {
	t.Parallel()

	s := &socket{historyBuffer: make([]buffer, HISTORY_BUFFER_SIZE)}
	for i := 0; i < 3; i++ {
		s.storeHistory(explorer.Event{
			Timestamp:  time.UnixMilli(int64(100 + i)),
			SampleRate: 20,
		})
	}

	if s.historyLen != 3 {
		t.Fatalf("history length = %d, want 3", s.historyLen)
	}
	for i := 0; i < 3; i++ {
		if got := s.historyAt(i).Timestamp; got != int64(100+i) {
			t.Fatalf("history[%d] timestamp = %d, want %d", i, got, 100+i)
		}
	}
}
