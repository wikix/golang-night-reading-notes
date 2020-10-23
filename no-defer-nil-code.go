package main

import (
	"os"
)

func main() {
	file := "a.txt"
	f, err := os.Open(file)
	if err != nil {
		panic(err)
	}
	f.Close()
	//f.process()
}
