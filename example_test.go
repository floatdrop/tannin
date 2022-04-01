package tannin_test

import (
	"fmt"

	tannin "github.com/floatdrop/tannin"
)

func ExampleTannin() {
	cache := tannin.New[string, int](256, 100)

	cache.Set("Hello", 5)

	if e := cache.Get("Hello"); e != nil {
		fmt.Println(*e)
		// Output: 5
	}
}
