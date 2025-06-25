package main

import "fmt"

func main() {
	ch := make(chan string)

	// Sender
	go func() {
		ch <- "hello"
	}()

	// Receiver
	msg := <-ch
	fmt.Println(msg) // Output: hello
}
