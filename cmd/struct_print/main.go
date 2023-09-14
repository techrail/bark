package main

import (
	`fmt`
)

type Xyz struct {
	Field1 string
	Field2 string
	F3     []int
	F4     []string
}

func main() {
	x := Xyz{
		Field1: "abc",
		Field2: "def",
		F3:     []int{1, 2, 3, 4, 5},
		F4:     []string{"a", "b", "D"},
	}

	fmt.Printf("%v\n", x)
}
