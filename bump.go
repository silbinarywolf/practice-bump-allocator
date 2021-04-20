package bump

import (
	"unsafe"

	"github.com/silbinarywolf/practice-bump-allocator/internal/reflectlite"
)

type outOfMemoryError struct {
}

func (err *outOfMemoryError) Error() string {
	return "out of memory"
}

// ErrOutOfMemory is used if a bump allocator runs out of space
var ErrOutOfMemory error = &outOfMemoryError{}

// memoryZeroer needs to be the size of the largest allocation you can make
// its used to clear memory with copy()
//
// not an ideal approach to solve the problem
var memoryZeroer = make([]byte, 4294967295)

type Allocator struct {
	pos int32
	buf []byte
}

func New(buf []byte) *Allocator {
	alloc := &Allocator{}
	alloc.buf = buf
	return alloc
}

func (alloc *Allocator) Reset() {
	alloc.pos = 0
}

// Pos returns the current position of the allocator
func (alloc *Allocator) Pos() int {
	return int(alloc.pos)
}

// Len returns how many bytes of data can be allocated into this allocator
func (alloc *Allocator) Len() int {
	return len(alloc.buf)
}

func (alloc *Allocator) New(v interface{}) interface{} {
	// note(jae): 2021-04-20
	// noescape() ensures the given "v" stays on the stack and doesn't escape
	// to the heap where it will be allocated
	eface := *(*reflectlite.InterfaceHeader)(noescape(unsafe.Pointer(&v)))
	size := int32(eface.Type.Size)
	{
		kind := eface.Type.Kind & reflectlite.KindMask // reflect.TypeOf().Kind()
		if kind == reflectlite.Ptr {
			tt := (*reflectlite.PtrType)(unsafe.Pointer(eface.Type)) // reflect.TypeOf().Elem()
			size = int32(tt.Elem.Size)
		}
	}
	endSlice := alloc.pos + size
	if int(endSlice) >= len(alloc.buf) {
		panic(ErrOutOfMemory)
	}
	// clear memory, we use this instead of a for-loop for speed.
	// In benchmarks we go from ~1700ns in "BenchmarkAlloc1000" to ~700ns
	copy(alloc.buf[alloc.pos:endSlice], memoryZeroer)
	r := reflectlite.InterfaceHeader{
		Type: eface.Type,
		Data: unsafe.Pointer(&alloc.buf[alloc.pos]),
	}
	alloc.pos += size
	castR := *(*interface{})(unsafe.Pointer(&r))
	return castR
}

// noescape hides a pointer from escape analysis. noescape is the identity
// function but escape analysis doesn't think the output depends on the input.
// noescape is inlined and currently compiles down to zero instructions.
// USE CAREFULLY!
//
// - noescape is copy/pasted from Go's runtime/stubs.go:noescape()
//
//go:nosplit
func noescape(p unsafe.Pointer) unsafe.Pointer {
	x := uintptr(p)
	return unsafe.Pointer(x ^ 0)
}
