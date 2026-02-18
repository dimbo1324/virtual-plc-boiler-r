package main

import (
	"fmt"
	"time"
)

func main() {
	fmt.Println("PLC Service Started")

	go func() {
		ticker := time.NewTicker(100 * time.Millisecond)
		for range ticker.C {
			fmt.Println("Tick...")
		}
	}()

	select {}
}
