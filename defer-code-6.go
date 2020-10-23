package main

import (
	"fmt"
	"time"
)

func main() {
	fmt.Println("The start of main function.")
	defer fmt.Println("defer main")
	var user = ""

	go func() {
		defer func() {
			fmt.Println("defer for goroutine caller.")
			if err := recover(); err != nil {
				fmt.Println("recover success. err:", err)
			}
		}()

		func() {
			defer func() {
				fmt.Println("panic occurs here. then we have panic defer.")
			}()

			if user == "" {
				panic("should set user env.")
			}

			//此处不会执行（终止于 panic）
			fmt.Println("after panic will be here~~~")
		}()
	}()

	time.Sleep(1000)

	fmt.Println("The end of main function.")
}
