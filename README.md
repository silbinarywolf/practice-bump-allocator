# Practice Bump Allocator

⚠️ *This library is not supported and should not be used as-is for any purpose*

This is my first attempt at building a bump / pool / arena allocator in Go 1.16. There were no real goals in mind at the beginning other than to see if I...

- Could do it
- Make it faster
- Make it as safe as possible while retaining maintainability

Unfortunately in my investigation of building this, I found that unless you accessed reflect data directly and unsafely, you would pay heavy allocation and speed costs and so I put together a small subset of the `reflect` package so I could do this, with tests that check my implementation against the `reflect` package to ensure parity. I've found if you don't do something like this, your bump allocator will end up allocating anyway via the `reflect` package and you gain nothing, but I didn't dive too deep, so I'd be happy to be wrong about this.



## Results from testing

```sh
$ go test -bench=.
goos: windows
goarch: amd64
pkg: github.com/silbinarywolf/practice-bump-allocator
cpu: AMD Ryzen 5 3600 6-Core Processor
BenchmarkNativeAlloc1000-12       455080              2658 ns/op            3200 B/op        100 allocs/op
BenchmarkAlloc1000-12            1387950               868.5 ns/op             0 B/op          0 allocs/op
PASS
ok      github.com/silbinarywolf/practice-bump-allocator        3.692s
```

For readability convenience, relevant code copy-pasted out of [bump_test.go](bump_test.go)
```go
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
```
