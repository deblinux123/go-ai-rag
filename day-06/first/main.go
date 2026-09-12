package main

import (
	"fmt"
	"time"
)

func main() {
	fmt.Println("Main started")

	go func() {
		fmt.Println("Background task started.")

		time.Sleep(3 * time.Second)

		fmt.Println("Background task finished.")
	}()

	fmt.Println("Main finished.")

	time.Sleep(4 * time.Second)
}
