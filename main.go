package main

import (
	"fmt"
	"io"
	"os"
	"strings"
)

func main() {
	messages, err := os.Open("messages.txt")
	if err != nil {
		fmt.Printf("error reading file: %v", err)
		return
	}
	buffer := make([]byte, 8)

	currentLine := ""

	for {
		read, err := messages.Read(buffer)
		if err != nil {
			if err == io.EOF && currentLine == "" {
				fmt.Printf("read: %s\n", currentLine)
			}
			return
		}

		text := string(buffer[:read])

		parts := strings.Split(text, "\n")

		for i := 0; i < len(parts)-1; i++ {
			currentLine += parts[i]
			fmt.Printf("read: %s\n", currentLine)
			currentLine = ""
		}

		if len(text) > 0 && text[len(text)-1] == '\n' {
			currentLine += parts[len(parts)-1]
			fmt.Printf("read: %s\n", currentLine)
			currentLine = ""
		} else {
			currentLine += parts[len(parts)-1]

		}
	}
}
