package reflectlite

import (
	"unsafe"
)

const (
	// KindMask was taken from the "reflect" package
	KindMask = (1 << 5) - 1
)

const (
	Ptr uint8 = 22
	//Struct uint8 = 25
)

// tflag is used by an rtype to signal what extra type information is
// available in the memory directly following the rtype value.
//
// tflag values must be kept in sync with copies in:
//	cmd/compile/internal/gc/reflect.go
//	cmd/link/internal/ld/decodesym.go
//	runtime/type.go
type tflag uint8

type nameOff int32 // offset to a name
type typeOff int32 // offset to an *rtype

// Rtype is the common implementation of most values.
// It is embedded in other struct types.
//
// Rtype must be kept in sync with ../runtime/type.go:/^type._type.
type Rtype struct {
	Size uintptr
	_    uintptr // number of bytes in the type that can contain pointers
	_    uint32  // hash of type; avoids computation in hash tables
	_    tflag   // extra type information flags
	_    uint8   // alignment of variable with this type
	_    uint8   // alignment of struct field with this type
	Kind uint8   // enumeration for C
	// function for comparing objects of this type
	// (ptr to object A, ptr to object B) -> ==?
	_ func(unsafe.Pointer, unsafe.Pointer) bool
	_ *byte   // garbage collection data
	_ nameOff // string form
	_ typeOff // type for pointer to this type, may be zero
}

// InterfaceHeader is the header for an interface{} value
type InterfaceHeader struct {
	Type *Rtype
	Data unsafe.Pointer
}

// PtrType represents a pointer type.
type PtrType struct {
	Rtype
	Elem *Rtype // pointer element (pointed at) type
}
