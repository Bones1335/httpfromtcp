package request

import (
	"bytes"
	"fmt"
	"io"
	"strings"
)

type state int

const (
	initialized state = iota
	done
)

const bufferSize = 8

type Request struct {
	RequestLine RequestLine
	State       state
}

type RequestLine struct {
	HttpVersion   string
	RequestTarget string
	Method        string
}

const crlf = "\r\n"

func RequestFromReader(reader io.Reader) (*Request, error) {
	buf := make([]byte, bufferSize, bufferSize)
	readToIndex := 0

	req := &Request{State: initialized}

	for req.State != done {
		if readToIndex == cap(buf) {
			newBuf := make([]byte, cap(buf)*2, cap(buf)*2)
			copy(newBuf, buf[:readToIndex])
			buf = newBuf
		}
		newIdx, err := reader.Read(buf[readToIndex:cap(buf)])
		if err == io.EOF {
			req.State = done
			break
		} else if err != nil {
			return nil, err
		}
		readToIndex += newIdx

		parsedBytes, err := req.parse(buf[:readToIndex])
		if err != nil {
			return nil, err
		}

		copy(buf, buf[parsedBytes:])
		readToIndex -= parsedBytes

	}
	return req, nil
}

func parseRequestLine(data []byte) (*RequestLine, int, error) {
	idx := bytes.Index(data, []byte(crlf))
	if idx == -1 {
		return nil, 0, nil
	}
	requestLineText := string(data[:idx])
	requestLine, err := requestLineFromString(requestLineText)
	if err != nil {
		return nil, 0, err
	}
	return requestLine, idx + 2, nil
}

func requestLineFromString(str string) (*RequestLine, error) {
	parts := strings.Split(str, " ")
	if len(parts) != 3 {
		return nil, fmt.Errorf("poorly formatted request-line: %s", str)
	}

	method := parts[0]
	for _, c := range method {
		if c < 'A' || c > 'Z' {
			return nil, fmt.Errorf("invalid method: %s", method)
		}
	}

	requestTarget := parts[1]

	versionParts := strings.Split(parts[2], "/")
	if len(versionParts) != 2 {
		return nil, fmt.Errorf("malformed start-line: %s", str)
	}

	httpPart := versionParts[0]
	if httpPart != "HTTP" {
		return nil, fmt.Errorf("unrecognized HTTP-version: %s", httpPart)
	}
	version := versionParts[1]
	if version != "1.1" {
		return nil, fmt.Errorf("unrecognized HTTP-version: %s", version)
	}

	return &RequestLine{
		Method:        method,
		RequestTarget: requestTarget,
		HttpVersion:   versionParts[1],
	}, nil
}

func (r *Request) parse(data []byte) (int, error) {
	if r.State == initialized {
		reqLine, consumed, err := parseRequestLine(data)
		if err != nil {
			return 0, err
		}
		if consumed == 0 {
			return consumed, nil
		}
		r.RequestLine = *reqLine
		r.State = done
		return consumed, nil
	}

	if r.State == done {
		return 0, fmt.Errorf("trying to read data in a done state")
	}

	return 0, fmt.Errorf("unknown state")
}
