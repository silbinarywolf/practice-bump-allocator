// inlinetest exists so that inlinetest_go.m
package inlinetest

import bump "github.com/silbinarywolf/practice-bump-allocator"

func inlineable() {
	alloc := bump.New(make([]byte, 256))
	_ = alloc.New(&structTest{}).(*structTest)
}

type structTest struct {
	X int32
}
