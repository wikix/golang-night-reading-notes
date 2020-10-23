package main

import (
  "sync/atomic"
  "log"
)

type WillPass struct {
  uncounted int64
}

type WillAlsoPass struct {
  init      int64
  uncounted int64
}

type WillPanic struct {
  init      bool
  uncounted int64
}

func main() {
  willPass := &WillPass{}
  willAlsoPass := &WillAlsoPass{}
  willPanic := &WillPanic{}
  var n int64 = 2

  atomic.AddInt64(&willPass.uncounted, n)
  log.Printf("willPass count is %d", willPass.uncounted)

  atomic.AddInt64(&willAlsoPass.uncounted, n)
  log.Printf("willAlsoPass count is %d", willAlsoPass.uncounted)

  // Kaboom
  atomic.AddInt64(&willPanic.uncounted, n)
  log.Printf("willPanic count is %d", willPanic.uncounted)
}
