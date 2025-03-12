package request

import (
	"fmt"
	"io"
	"strings"
)

type Request struct {
	RequestLine RequestLine
}

type RequestLine struct {
	HttpVersion   string
	RequestTarget string
	Method        string
}

func RequestFromReader(reader io.Reader) (*Request, error) {
	b, err := io.ReadAll(reader)
	if err != nil {
		fmt.Println("error reading data")
		return &Request{}, err
	}

	req := strings.Split(string(b), "\r\n")

	parts := strings.Split(req[0], " ")
	if len(parts) != 3 {
		return nil, fmt.Errorf("too many or too few elements")
	}
	if parts[0] != strings.ToUpper(parts[0]) {
		return nil, fmt.Errorf("method not upper case")
	}
	if parts[2] != "HTTP/1.1" {
		return nil, fmt.Errorf("not correct http version")
	}

	version := strings.Split(parts[2], "HTTP/")

	reqLine := RequestLine{
		HttpVersion:   version[1],
		RequestTarget: parts[1],
		Method:        parts[0],
	}

	request := &Request{RequestLine: reqLine}

	fmt.Printf("%v\n", version[1])
	return request, nil
}
