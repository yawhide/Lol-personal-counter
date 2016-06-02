package main

import (
	"log"

	"github.com/juju/ratelimit"
)

func main() {
	bucket := ratelimit.NewBucketWithRate(5, 5)
	counter := 0
	go func() {
		for {
			bucket.Wait(1)
			log.Println("[1]done ", counter)
			counter++
		}
	}()
	go func() {
		for {
			bucket.Wait(1)
			log.Println("[2]done ", counter)
			counter++
		}
	}()

	// bucket.Wait(1)
	// log.Println("done 2")
	//
	// bucket.Wait(1)
	// log.Println("done 3")
	select {}
}
