package main

type T6 struct {
	M *int
}

func main() {
	var x T6
	var i int
	ref6(&i, &x)
}

func ref6(y *int, z *T6) {
	z.M = y
}
