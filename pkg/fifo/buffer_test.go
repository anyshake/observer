package fifo_test

import (
	"bytes"
	"sync"
	"testing"

	"github.com/anyshake/observer/pkg/fifo"
)

func TestBufferReadAcrossWraparound(t *testing.T) {
	t.Parallel()

	buffer := fifo.New[byte](6)
	writeBytes(t, buffer, []byte{1, 2, 3, 4})
	if got, err := buffer.Read(3); err != nil || !bytes.Equal(got, []byte{1, 2, 3}) {
		t.Fatalf("Read(3) = %v, %v", got, err)
	}
	writeBytes(t, buffer, []byte{5, 6, 7})
	if got := buffer.Len(); got != 4 {
		t.Fatalf("Len() = %d, want 4", got)
	}
	if got, err := buffer.Read(5); err == nil || got != nil {
		t.Fatalf("oversized Read() = %v, %v, want nil and error", got, err)
	}
	if got, err := buffer.Read(4); err != nil || !bytes.Equal(got, []byte{4, 5, 6, 7}) {
		t.Fatalf("Read() after insufficient data = %v, %v", got, err)
	}
	if buffer.Len() != 0 {
		t.Fatal("buffer is not empty after reading all data")
	}
}

func TestBufferOverflowAndReset(t *testing.T) {
	t.Parallel()

	// One slot distinguishes a full FIFO from an empty FIFO.
	buffer := fifo.New[byte](5)
	writeBytes(t, buffer, []byte{1, 2, 3, 4, 5, 6, 7})
	if got := buffer.Len(); got != 4 {
		t.Fatalf("Len() after overflow = %d, want 4", got)
	}
	if got, err := buffer.Read(4); err != nil || !bytes.Equal(got, []byte{4, 5, 6, 7}) {
		t.Fatalf("Read() after overflow = %v, %v", got, err)
	}
	writeBytes(t, buffer, []byte{8, 9})
	buffer.Reset()
	if buffer.Len() != 0 {
		t.Fatal("Reset() retained buffered data")
	}
	if _, err := buffer.Read(1); err == nil {
		t.Fatal("Read() succeeded after Reset()")
	}
	writeBytes(t, buffer, []byte{10, 11})
	if got, err := buffer.Read(2); err != nil || !bytes.Equal(got, []byte{10, 11}) {
		t.Fatalf("Read() after reuse = %v, %v", got, err)
	}
}

func TestBufferPeekResynchronizesAcrossWraparound(t *testing.T) {
	t.Parallel()

	buffer := fifo.New[byte](8)
	writeBytes(t, buffer, []byte{1, 2, 3, 4, 5, 6})
	if _, err := buffer.Read(6); err != nil {
		t.Fatal(err)
	}
	// Noise occupies slot 6; the two-byte header straddles the buffer boundary.
	writeBytes(t, buffer, []byte{0x00, 0xaa, 0x55, 0x01})
	if got, err := buffer.Peek([]byte{0xaa, 0x55}, 4); err == nil || got != nil {
		t.Fatalf("Peek() with incomplete packet = %v, %v", got, err)
	}
	if got := buffer.Len(); got != 3 {
		t.Fatalf("Len() after discarding noise = %d, want 3", got)
	}
	writeBytes(t, buffer, []byte{0x02, 0x99})
	if got, err := buffer.Peek([]byte{0xaa, 0x55}, 4); err != nil || !bytes.Equal(got, []byte{0xaa, 0x55, 0x01, 0x02}) {
		t.Fatalf("Peek() after completing packet = %v, %v", got, err)
	}
	if got, err := buffer.Peek(nil, 1); err != nil || !bytes.Equal(got, []byte{0x99}) {
		t.Fatalf("Peek() remaining data = %v, %v", got, err)
	}
	if buffer.Len() != 0 {
		t.Fatal("Peek() did not consume the packet")
	}
}

func writeBytes(t *testing.T, buffer *fifo.Buffer[byte], data []byte) {
	t.Helper()
	if n, err := buffer.Write(data...); err != nil || n != len(data) {
		t.Fatalf("Write() = %d, %v, want %d, nil", n, err, len(data))
	}
}

func TestInvalidReadSizesDoNotConsumeData(t *testing.T) {
	t.Parallel()
	buffer := fifo.New[byte](8)
	writeBytes(t, buffer, []byte{1, 2, 3})
	if _, err := buffer.Read(-1); err == nil {
		t.Fatal("negative Read() size accepted")
	}
	if _, err := buffer.Peek(nil, -1); err == nil {
		t.Fatal("negative Peek() size accepted")
	}
	if _, err := buffer.Peek([]byte{1, 2}, 1); err == nil {
		t.Fatal("header larger than packet accepted")
	}
	if got, err := buffer.Read(3); err != nil || !bytes.Equal(got, []byte{1, 2, 3}) {
		t.Fatalf("invalid reads consumed data: %v, %v", got, err)
	}
}

func TestConcurrentConsumersDeliverEachValueOnce(t *testing.T) {
	t.Parallel()
	const count = 2048
	buffer := fifo.New[int](count + 1)
	values := make([]int, count)
	for i := range values {
		values[i] = i
	}
	if _, err := buffer.Write(values...); err != nil {
		t.Fatal(err)
	}
	received := make(chan int, count)
	var wg sync.WaitGroup
	for i := 0; i < 8; i++ {
		wg.Add(1)
		go func(usePeek bool) {
			defer wg.Done()
			for {
				var data []int
				var err error
				if usePeek {
					data, err = buffer.Peek(nil, 1)
				} else {
					data, err = buffer.Read(1)
				}
				if err != nil {
					return
				}
				received <- data[0]
				_ = buffer.Len()
			}
		}(i%2 == 0)
	}
	go func() { wg.Wait(); close(received) }()
	seen := make(map[int]int)
	for value := range received {
		seen[value]++
	}
	for _, value := range values {
		if seen[value] != 1 {
			t.Fatalf("value %d delivered %d times, want once", value, seen[value])
		}
	}
	if buffer.Len() != 0 {
		t.Fatal("consumers left unread values")
	}
}
