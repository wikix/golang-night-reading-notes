package main

type T4 struct {
	M *int
}

func main() {
	var i int
	ref4(i)
}

func ref4(y int) (z T4) {
	z.M = &y
	return z
}
