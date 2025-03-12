package main

import (
	"bufio"
	"fmt"
	"net"
	"os"
)

func main() {
	udpAddr, err := net.ResolveUDPAddr("udp", "localhost:42069")
	if err != nil {
		fmt.Printf("error starting UDP connection: %v", err)
		return
	}

	conn, err := net.DialUDP("udp", nil, udpAddr)
	if err != nil {
		fmt.Printf("error connecting to udp: %v", err)
		return
	}
	defer conn.Close()

	buffer := bufio.NewReader(os.Stdin)
	for {
		fmt.Println(">")
		input, err := buffer.ReadString('\n')
		if err != nil {
			fmt.Printf("error reading input: %v", err)
			return
		}
		conn.Write([]byte(input))
	}
}
