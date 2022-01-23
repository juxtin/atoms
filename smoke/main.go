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
	fmt.Println("Value of a is:", a.Deref())
}
