package main

import (
	"fmt"
	"net"

	"github.com/Bones1335/httpfromtcp/internal/request"
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
		r, err := request.RequestFromReader(conn)
		if err != nil {
			fmt.Printf("error: %s", err)
		}

		fmt.Println("Request line:")
		fmt.Printf("- Method: %+v\n", r.RequestLine.Method)
		fmt.Printf("- Target: %+v\n", r.RequestLine.RequestTarget)
		fmt.Printf("- Version: %+v\n", r.RequestLine.HttpVersion)
		fmt.Println("Headers:")
		for k, v := range r.Headers {
			fmt.Printf("- %+v: %+v\n", k, v)
		}
		fmt.Println("Body:")
		fmt.Printf("%s\n", r.Body)

		conn.Close()
		fmt.Println("Connection closed......")
	}

}
