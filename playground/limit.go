package main

import (
  "fmt"

  "golang.org/x/time/rate"
)

func main() {
  var limit rate.Limit = 10 // 10/s
  l := rate.NewLimiter(limit, 10)
  fmt.Printf("limiter: %+v\n", l)
}
