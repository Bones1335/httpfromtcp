package headers

import (
	"fmt"
	"strings"
)

type Headers map[string]string

func NewHeaders() Headers {
	return make(map[string]string)
}

func (h Headers) Parse(data []byte) (n int, done bool, err error) {
	dataString := string(data)

	if strings.Contains(dataString, "\r\n") {
		idx := strings.Index(dataString, "\r\n")
		if idx == 0 {
			return idx, true, nil
		}
		if idx > 0 {
			header := dataString[:idx]
			slice := strings.SplitN(header, ":", 2)
			colonIdx := strings.Index(header, ":")
			if colonIdx > 0 && header[colonIdx-1] == ' ' {
				return 0, false, fmt.Errorf("invalid space before colon")
			}
			if len(slice) != 2 {
				return 0, false, fmt.Errorf("missing colon")
			}
			key := strings.TrimSpace(slice[0])
			value := strings.TrimSpace(slice[1])

			h[key] = value
		}

		return idx + len("\r\n"), false, nil
	}

	return 0, false, nil
}
