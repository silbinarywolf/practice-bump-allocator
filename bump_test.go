package bump

import (
	"testing"
	"unsafe"
)

type Player struct {
	X, Y          int
	Width, Height int
}

// BenchmarkNativeAlloc is here so we can compare native Go
// allocations to our bump allocator
func BenchmarkNativeAlloc1000(b *testing.B) {
	b.ReportAllocs()

	var p1 *Player

	for n := 0; n < b.N; n++ {
		for i := 0; i < 100; i++ {
			p1 = &Player{}
			p1.X = 1 + n
			p1.Y = 2 + n
			p1.Width = 30 + n
			p1.Height = 40 + n
		}
	}
}

// BenchmarkAlloc just tests out how things go with a naive allocator
func BenchmarkAlloc1000(b *testing.B) {
	b.ReportAllocs()

	var p1 *Player

	alloc := New(make([]byte, 32678))
	for n := 0; n < b.N; n++ {
		for i := 0; i < 100; i++ {
			p1 = alloc.New(&Player{}).(*Player)
			p1.X = 1 + n
			p1.Y = 2 + n
			p1.Width = 30 + n
			p1.Height = 40 + n
		}
		alloc.Reset()
	}
}

func TestOutOfMemory(t *testing.T) {
	defer func() {
		hasOutOfMemoryError := false
		if r := recover(); r != nil {
			hasOutOfMemoryError = r == ErrOutOfMemory
		}
		if !hasOutOfMemoryError {
			t.Fatalf(`expected "out of memory" panic as we shouldn't have enough memory to allocate`)
		}
	}()
	alloc := New(make([]byte, 4))
	alloc.New(&Player{})
}

func TestReset(t *testing.T) {
	alloc := New(make([]byte, 100))

	p1 := alloc.New(&Player{}).(*Player)
	p1.X = 1
	p1.Y = 2
	p1.Width = 30
	p1.Height = 40

	alloc.Reset()

	p2 := alloc.New(&Player{}).(*Player)
	p2.X = 300
	p2.Y = 302

	// Test if first allocated object (p1) has now been stomped over with (p2)
	// values as we'd expect after a call to alloc.Reset()
	if unsafe.Pointer(p1) != unsafe.Pointer(p2) {
		t.Fatalf("expected p1 and p2 to be pointing to same memory location, not %v and %v", unsafe.Pointer(p1), unsafe.Pointer(p2))
	}
	if expected := 300; p1.X != expected {
		t.Errorf("expected X to be %v not %v", expected, p1.X)
	}
	if expected := 302; p1.Y != expected {
		t.Errorf("expected Y to be %v not %v", expected, p1.Width)
	}
	// We test out these values because we expect the allocator to zero
	// these out for us, like Go would.
	if expected := 0; p1.Width != expected {
		t.Errorf("expected Width to be %v not %v", expected, p1.Width)
	}
	if expected := 0; p1.Height != expected {
		t.Errorf("expected Height to be %v not %v", expected, p1.Height)
	}
}
