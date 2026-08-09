package series

import (
	"testing"

	"github.com/apache/arrow-go/v18/arrow/memory"
)

func TestArrowFloat64sBasics(t *testing.T) {
	s := NewArrowFloat64s([]float64{1, 2, 3}, memory.DefaultAllocator)
	defer s.Release()
	if s.Len() != 3 {
		t.Fatalf("Len: got %d, want 3", s.Len())
	}
	if s.Get(0) != 1 || s.Get(2) != 3 {
		t.Fatalf("Get: got %v,%v want 1,3", s.Get(0), s.Get(2))
	}
}

func TestArrowFloat64sAdd(t *testing.T) {
	a := NewArrowFloat64s([]float64{1, 2, 3}, memory.DefaultAllocator)
	defer a.Release()
	b := NewArrowFloat64s([]float64{10, 20, 30}, memory.DefaultAllocator)
	defer b.Release()
	c := a.Add(b)
	defer c.Release()
	if c.Len() != 3 || c.Get(0) != 11 || c.Get(2) != 33 {
		t.Fatalf("Add: got [%v ... %v], want 11..33", c.Get(0), c.Get(2))
	}
}

func TestArrowFloat64sGreaterThanFilter(t *testing.T) {
	a := NewArrowFloat64s([]float64{1, 5, 2, 8, 3}, memory.DefaultAllocator)
	defer a.Release()
	mask := a.GreaterThan(3)
	defer mask.Release()
	out := a.Filter(mask)
	defer out.Release()
	if out.Len() != 2 || out.Get(0) != 5 || out.Get(1) != 8 {
		t.Fatalf("GreaterThan+Filter: got len=%d [%v %v], want len=2 [5 8]", out.Len(), out.Get(0), out.Get(1))
	}
}
