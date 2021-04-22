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

type noPtrError struct {
}

func (err *noPtrError) Error() string {
	return "invalid type, expected pointer"
}

// ErrOutOfMemory is used if a bump allocator runs out of space
var ErrOutOfMemory error = &outOfMemoryError{}

// ErrNotPtr is used if the incorrect type is used
var ErrNotPtr error = &noPtrError{}

type Allocator struct {
	pos uintptr
	buf []byte
}

func New(buf []byte) *Allocator {
	alloc := &Allocator{}
	alloc.buf = buf
	alloc.Reset()
	return alloc
}

func (alloc *Allocator) Reset() {
	// note(jae): 2021-04-21
	// As of Go 1.16, doing "for i := 0; i < len(alloc.buf); i++" will cause Go to
	// seemingly not optimize this to fast memset/memclr operations and Reset() will
	// become *very very* slow
	for i := range alloc.buf {
		alloc.buf[i] = 0
	}
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

// New will allocate data of the underlying type to the buffer and return a pointer
// to it
//
// v := alloc.New(&MyStruct{}).(*MyStruct)
//
// Do not use this with pointers to slices or maps.
func (alloc *Allocator) New(v interface{}) interface{} {
	// note(jae): 2021-04-21
	// this method is very sensitive to the Go inliner, so we have tests in "internal/inlinetest"
	// to see if its inlined or not.

	// note(jae): 2021-04-20
	// noescape() ensures the given "v" stays on the stack and doesn't escape
	// to the heap where it will be allocated
	efaceType := (*(*reflectlite.InterfaceHeader)(noescape(unsafe.Pointer(&v)))).Type
	if efaceType.Kind&reflectlite.KindMask != reflectlite.Ptr {
		panic(ErrNotPtr)
	}
	size := (*reflectlite.PtrType)(unsafe.Pointer(efaceType)).Elem.Size // reflect.TypeOf().Elem()
	if int(alloc.pos+size) >= len(alloc.buf) {
		panic(ErrOutOfMemory)
	}
	r := reflectlite.InterfaceHeader{
		Type: efaceType,
		Data: unsafe.Pointer(&alloc.buf[alloc.pos]),
	}
	alloc.pos += size
	return *(*interface{})(unsafe.Pointer(&r))
}
