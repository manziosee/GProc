package main

import (
	"fmt"
	"time"
)

func main() {
	for i := 0; i < 100; i++ {
		fmt.Printf("Test app running - iteration %d\n", i)
		time.Sleep(2 * time.Second)
	}
}