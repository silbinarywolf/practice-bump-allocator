package reflectlite

import (
	"reflect"
	"testing"
	"unsafe"
)

type oneElementStruct struct {
	X int32
}

type twoElementStruct struct {
	X int32
	Y int16
}

var structsToTest = []interface{}{
	&oneElementStruct{},
	//&twoElementStruct{},
}

// TestValidity exists to ensure this package remains syncronized with upstream Go data layout changes
func TestReflectliteHasParityWithReflect(t *testing.T) {
	// WARNING / NOTE:
	//
	// If these tests are suddenly failing on a new version of Go, that likely means Go's internal
	// data structures have been updated.
	//
	// To resolve these issues, assuming they're not large sweeping changes,
	// I'd start with investigating the "reflect" package and seeing if "type rtype struct" has changed
	// vs this package.
	for _, v := range structsToTest {
		myInterface := *(*InterfaceHeader)(unsafe.Pointer(&v)) // reflect.TypeOf(v), equivalent
		theirTypeOf := reflect.TypeOf(v)

		// Test Kind() behaviour matches
		if theirTypeOf.Kind() != reflect.Ptr {
			t.Fatalf(`expected Kind in reflect to return reflect.Ptr`)
		}
		myKind := reflect.Kind(myInterface.Type.Kind & KindMask)
		if myKind != theirTypeOf.Kind() {
			t.Fatalf(`expected Kind in reflect and reflectlite to have same value, not "%d" (reflectlite) and "%d" (reflect)`, myKind, theirTypeOf.Kind())
		}

		// Test unsafe.Sizeof behaviour matches
		{
			const sizeOfPointer uintptr = 8
			theirSize := theirTypeOf.Size()
			if theirSize != sizeOfPointer {
				t.Fatalf(`expected ptr Size in reflect to return %d, not %d`, sizeOfPointer, theirSize)
			}
			mySize := myInterface.Type.Size
			if mySize != theirSize {
				t.Fatalf(`expected ptr Size in reflect and reflectlite to have same value, not "%d" (reflectlite) and "%d" (reflect)`, mySize, theirSize)
			}
		}

		// Test getting underlying element for ptr
		{
			tt := (*PtrType)(unsafe.Pointer(myInterface.Type))
			theirSize := theirTypeOf.Elem().Size()
			mySize := tt.Elem.Size
			if mySize != theirSize {
				t.Fatalf(`expected element %T Size in reflect and reflectlite to have same value, not "%d" (reflectlite) and "%d" (reflect)`, v, mySize, theirSize)
			}
		}
	}
}
