# tannin
[![Go Reference](https://pkg.go.dev/badge/github.com/floatdrop/tannin.svg)](https://pkg.go.dev/github.com/floatdrop/tannin)
[![build](https://github.com/floatdrop/tannin/actions/workflows/ci.yml/badge.svg)](https://github.com/floatdrop/tannin/actions/workflows/ci.yml)
![Coverage](https://img.shields.io/badge/Coverage-67.5%25-yellow)
[![Go Report Card](https://goreportcard.com/badge/github.com/floatdrop/tannin)](https://goreportcard.com/report/github.com/floatdrop/tannin)


Tannin is implementation of [W-TinyLRU](https://arxiv.org/pdf/1512.00727.pdf) cache, that was implemented in [Caffeine](https://github.com/ben-manes/caffeine) as eviction policy. Later it was adopted to [Ristretto](https://github.com/dgraph-io/ristretto) cache library.

## Example

```go
import (
	"fmt"

	tannin "github.com/floatdrop/tannin"
)

func main() {
	cache := tannin.New[string, int](256, 100)

	cache.Set("Hello", 5)

	if e := cache.Get("Hello"); e != nil {
		fmt.Println(*e)
		// Output: 5
	}
}
```

## Benchmarks

```
floatdrop/tannin
	BenchmarkTannin_Rand-8   	 3957745	       295.8 ns/op	      84 B/op	       7 allocs/op
	BenchmarkTannin_Freq-8   	 4598648	       255.7 ns/op	      78 B/op	       7 allocs/op
```