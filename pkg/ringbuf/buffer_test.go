package ringbuf_test

import (
	"slices"
	"testing"

	"github.com/anyshake/observer/pkg/ringbuf"
)

func TestBufferRetainsNewestValuesInOrder(t *testing.T) {
	t.Parallel()

	buffer := ringbuf.New[int](3)
	steps := []struct {
		push []int
		want []int
	}{
		{want: []int{}},
		{push: []int{1, 2}, want: []int{1, 2}},
		{push: []int{3}, want: []int{1, 2, 3}},
		{push: []int{4}, want: []int{2, 3, 4}},
		{push: []int{5, 6, 7, 8}, want: []int{6, 7, 8}},
		{want: []int{6, 7, 8}},
	}
	for i, step := range steps {
		buffer.Push(step.push...)
		if got := buffer.Values(); !slices.Equal(got, step.want) {
			t.Fatalf("step %d: Values() = %v, want %v", i, got, step.want)
		}
		if got := buffer.Len(); got != len(step.want) {
			t.Fatalf("step %d: Len() = %d, want %d", i, got, len(step.want))
		}
	}
	values := buffer.Values()
	values[0] = 99
	if got := buffer.Values(); !slices.Equal(got, []int{6, 7, 8}) {
		t.Fatalf("Values() exposed internal storage: %v", got)
	}
	buffer.Reset()
	if buffer.Len() != 0 || len(buffer.Values()) != 0 {
		t.Fatal("Reset() retained values")
	}
	buffer.Push(9)
	if got := buffer.Values(); !slices.Equal(got, []int{9}) {
		t.Fatalf("Values() after reuse = %v, want [9]", got)
	}
}

func TestSingleSlotBuffer(t *testing.T) {
	t.Parallel()

	buffer := ringbuf.New[int](1)
	buffer.Push(1, 2, 3)
	if got := buffer.Values(); !slices.Equal(got, []int{3}) || buffer.Len() != 1 {
		t.Fatalf("single-slot buffer = %v, length = %d", got, buffer.Len())
	}
}
