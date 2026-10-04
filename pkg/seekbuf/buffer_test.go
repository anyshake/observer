package seekbuf_test

import (
	"bytes"
	"errors"
	"io"
	"testing"

	"github.com/anyshake/observer/pkg/seekbuf"
)

var _ io.ReadWriteSeeker = (*seekbuf.Buffer)(nil)

func TestReadAdvancesPositionAndReturnsEOF(t *testing.T) {
	t.Parallel()
	var buffer seekbuf.Buffer
	if _, err := buffer.WriteString("abcdef"); err != nil {
		t.Fatal(err)
	}
	if _, err := buffer.Seek(0, io.SeekStart); err != nil {
		t.Fatal(err)
	}
	p := make([]byte, 4)
	if n, err := buffer.Read(p); n != 4 || err != nil || string(p) != "abcd" {
		t.Fatalf("first Read() = %d, %v, %q", n, err, p)
	}
	if n, err := buffer.Read(p); n != 2 || err != nil || string(p[:n]) != "ef" {
		t.Fatalf("second Read() = %d, %v, %q", n, err, p)
	}
	if n, err := buffer.Read(p); n != 0 || !errors.Is(err, io.EOF) {
		t.Fatalf("Read() at end = %d, %v, want EOF", n, err)
	}
	if _, err := buffer.Seek(10, io.SeekStart); err != nil {
		t.Fatal(err)
	}
	if n, err := buffer.Read(p); n != 0 || !errors.Is(err, io.EOF) {
		t.Fatalf("Read() beyond end = %d, %v, want EOF", n, err)
	}
	if n, err := buffer.Read(nil); n != 0 || err != nil {
		t.Fatalf("zero-length Read() = %d, %v", n, err)
	}
}

func TestSeekOriginsAndInvalidOffsets(t *testing.T) {
	t.Parallel()
	var buffer seekbuf.Buffer
	if _, err := buffer.WriteString("abcdef"); err != nil {
		t.Fatal(err)
	}
	tests := []struct {
		name    string
		offset  int64
		whence  int
		want    int64
		wantErr bool
	}{
		{"start", 2, io.SeekStart, 2, false},
		{"current", 1, io.SeekCurrent, 3, false},
		{"end", -1, io.SeekEnd, 5, false},
		{"negative start", -1, io.SeekStart, 5, true},
		{"negative current", -6, io.SeekCurrent, 5, true},
		{"negative end", -7, io.SeekEnd, 5, true},
		{"overflow", int64(^uint(0) >> 1), io.SeekCurrent, 5, true},
		{"invalid origin", 0, 99, 5, true},
		{"beyond end", 4, io.SeekEnd, 10, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			pos, err := buffer.Seek(tt.offset, tt.whence)
			if (err != nil) != tt.wantErr || (!tt.wantErr && pos != tt.want) {
				t.Fatalf("Seek() = %d, %v, want %d, error=%v", pos, err, tt.want, tt.wantErr)
			}
			if pos, err := buffer.Seek(0, io.SeekCurrent); err != nil || pos != tt.want {
				t.Fatalf("current position = %d, %v, want %d", pos, err, tt.want)
			}
		})
	}
}

func TestWritePreservesDataAcrossGrowthAndOverwrite(t *testing.T) {
	t.Parallel()
	var buffer seekbuf.Buffer
	if _, err := buffer.WriteString("head"); err != nil {
		t.Fatal(err)
	}
	payload := bytes.Repeat([]byte{0x5a}, 4096)
	if n, err := buffer.Write(payload); n != len(payload) || err != nil {
		t.Fatalf("Write() = %d, %v", n, err)
	}
	if _, err := buffer.Seek(0, io.SeekStart); err != nil {
		t.Fatal(err)
	}
	if n, err := buffer.WriteString("WAVE"); n != 4 || err != nil {
		t.Fatalf("WriteString() = %d, %v", n, err)
	}
	want := append([]byte("WAVE"), payload...)
	if !bytes.Equal(buffer.Bytes(), want) || buffer.Len() != len(want) {
		t.Fatal("growing or overwriting the buffer corrupted data")
	}
}

func TestSparseWriteAfterResetClearsOldData(t *testing.T) {
	t.Parallel()
	var buffer seekbuf.Buffer
	if _, err := buffer.WriteString("old-secret-data"); err != nil {
		t.Fatal(err)
	}
	buffer.Reset()
	if buffer.Len() != 0 || buffer.String() != "" {
		t.Fatal("Reset() retained visible data")
	}
	if _, err := buffer.Seek(4, io.SeekStart); err != nil {
		t.Fatal(err)
	}
	if n, err := buffer.Write(nil); n != 0 || err != nil || buffer.Len() != 0 {
		t.Fatalf("empty Write() extended buffer: %d, %v, len=%d", n, err, buffer.Len())
	}
	if _, err := buffer.WriteString("x"); err != nil {
		t.Fatal(err)
	}
	if got := buffer.Bytes(); !bytes.Equal(got, []byte{0, 0, 0, 0, 'x'}) {
		t.Fatalf("sparse Write() = %v, want zero-filled gap", got)
	}
}

func TestWriteRejectsPositionOverflow(t *testing.T) {
	t.Parallel()
	var buffer seekbuf.Buffer
	if _, err := buffer.Seek(int64(^uint(0)>>1), io.SeekStart); err != nil {
		t.Fatal(err)
	}
	if n, err := buffer.Write([]byte{1}); n != 0 || err == nil || buffer.Len() != 0 {
		t.Fatalf("overflowing Write() = %d, %v, len=%d", n, err, buffer.Len())
	}
}
