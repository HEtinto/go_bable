package main

import (
	"fmt"
	"time"
)

func main() {
	// ...create abort channel...
	abort := make(chan int)
	fmt.Println("Commencing countdown.  Press return to abort.")
	tick := time.Tick(1 * time.Second)
	for countdown := 10; countdown > 0; countdown-- {
		fmt.Println(countdown)
		select {
		case <-tick:
			// Do nothing.
		case <-abort:
			fmt.Println("Launch aborted:", countdown)
			return
		}
		go func() {
			if countdown%2 == 0 {
				abort <- 1
			}
		}()
	}
}
