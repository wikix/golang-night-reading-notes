package main

type T5 struct {
	M *int
}

func main() {
	var i int
	ref5(&i)
}

func ref5(y *int) (z T5) {
	z.M = y
	return z
}
