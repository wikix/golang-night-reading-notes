package main

type T2 struct{}

func main() {
	var x T2
	y := &x
	_ = *identity1(y)
}

func identity1(z *T2) *T2 {
	return z
}
