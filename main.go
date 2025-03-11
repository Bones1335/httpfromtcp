package main

import (
	"fmt"
	"os"
)

func main() {
	messages, err := os.Open("messages.txt")
	if err != nil {
		fmt.Printf("error reading file: %v", err)
		return
	}
	buffer := make([]byte, 8)
	for {
		read, err := messages.Read(buffer)
		if err != nil {
			return
		}
		fmt.Printf("read: %s\n", buffer[:read])

	}
}
