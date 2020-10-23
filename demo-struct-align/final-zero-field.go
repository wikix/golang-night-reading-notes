package main

import (
	"fmt"
	"unsafe"
)

type T1 struct {
	a struct{}
	x int64
}

type T2 struct {
	x int64
	a struct{}
}

func main() {
	a1 := T1{}
	a2 := T2{}
	fmt.Printf("zero size of struct with T1 size: %d; T2(as final-zero-field) size: %d", unsafe.Sizeof(a1), unsafe.Sizeof(a2))
}
