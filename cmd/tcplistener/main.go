package main

import (
	"errors"
	"fmt"
	"io"
	"net"
	"strings"
)

func main() {
	listener, err := net.Listen("tcp", ":42069")
	if err != nil {
		fmt.Printf("error starting tcp server: %v", err)
		return
	}
	defer listener.Close()

	for {
		conn, err := listener.Accept()
		if err != nil {
			fmt.Printf("error starting connection: %v", err)
			return
		}
		fmt.Println("Connection started......")
		lines := getLinesChannel(conn)
		for line := range lines {
			fmt.Printf("%s\n", line)
		}
		conn.Close()
		fmt.Println("Connection closed......")
	}

}

func getLinesChannel(f io.ReadCloser) <-chan string {
	lines := make(chan string)
	go func() {
		defer f.Close()
		defer close(lines)
		currentLineContents := ""
		for {
			buffer := make([]byte, 8, 8)
			read, err := f.Read(buffer)
			if err != nil {
				if currentLineContents != "" {
					lines <- currentLineContents
				}
				if errors.Is(err, io.EOF) {
					break
				}
				fmt.Printf("error: %s\n", err.Error())
				return
			}
			str := string(buffer[:read])
			parts := strings.Split(str, "\n")
			for i := 0; i < len(parts)-1; i++ {
				lines <- fmt.Sprintf("%s%s", currentLineContents, parts[i])
				currentLineContents = ""
			}
			currentLineContents += parts[len(parts)-1]
		}
	}()
	return lines
}
