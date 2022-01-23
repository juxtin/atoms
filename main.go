package main

import (
	"fmt"

	"github.com/juxtin/atoms"
)

func main() {
	a := atoms.New[int](1)
	a.Swap(func(i int) int { return i + 1 })
	fmt.Println("Value of a is:", a.Deref())
}
