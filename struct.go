package main

import (
    "fmt"
    "sync/atomic"
)

func main() {
    var c0 int64
    fmt.Println("64位字本身：", atomic.AddInt64(&c0, 1))
    c1 := [5]int64{}
    fmt.Println("64位字数组、切片：", atomic.AddInt64(&c1[:][0], 1))
}
