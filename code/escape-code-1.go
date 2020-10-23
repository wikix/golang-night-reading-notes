package main

type T1 struct{}

func main() {
	var x T1
	_ = identity1(x)
}

func identity1(x T1) T1 {
	return x
}
