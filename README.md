# atoms (for Go)

An implementation of [Clojure's atoms](https://clojure.org/reference/atoms) for Go, using generics.

Arguably much better than the [`sync/atomic`](https://golang.org/pkg/sync/atomic/) package, but with the following limitations:

* Work _very much_ in progress. Tests are not yet written. Use at your own risk!
* This package's `Atom`s currently use locks, which somewhat defeats the purpose. I'll fix this in the future.

## Usage

Note that you must use go 1.18 or later to use this package. Run `script/bootstrap` to install it alongside your current Go version.

### Working example
```go
package main

import (
	"fmt"
	"time"
	"github.com/juxtin/atoms"
)

func inc(i int)int {
	return i + 1
}

func main() {
	a := atoms.New[int](0)
	a.AddWatch("printer", func(_ string, _ *atoms.Atom[int], old, new int) {
		fmt.Println("Updated", old, "to", new)
	})
	a.SetValidator(func(i int)error {
		if i > 68 {
			return fmt.Errorf("%d must never exceed 68!", i)
		}
		return nil
	})
	for i := 0; i < 100; i++ {
		go a.Swap(inc)
	}
	// wait for the goroutines to finish
	time.Sleep(10 * time.Millisecond)
	fmt.Println("Value of a is:", a.Deref()) // prints "Value of a is: 68"
}
```
