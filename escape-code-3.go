package main

type T3 struct{}

func main() {
	var x T3
	_ = *ref3(x)
}

func ref3(z T3) *T3 {
	return &z
}
