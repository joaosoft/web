package webserver

import (
	"bufio"
	"bytes"
	"errors"
	"fmt"
	"net"
	"time"
)

func NewRequest(conn net.Conn) (*Request, error) {
	request := &Request{
		Base: Base{
			Headers: make(Headers),
			Cookies: make(Cookies),
		},
	}

	return request, request.read(conn)
}

func (r *Request) read(conn net.Conn) error {
	reader := bufio.NewReader(conn)

	// read one line (ended with \n or \r\n)
	conn.SetReadDeadline(time.Now().Add(time.Second * 1))
	line, _, err := reader.ReadLine()

	if err != nil {
		return fmt.Errorf("invalid http request: %s", err)
	}

	if split := bytes.SplitN(line, []byte(` `), 3); len(split) < 3 {
		return errors.New("invalid http request")
	} else {
		r.Method = Method(split[0])
		r.Url = string(split[1])
		r.Protocol = string(split[2])
	}

	// headers
	for {
		conn.SetReadDeadline(time.Now().Add(time.Millisecond * 1))
		line, _, err = reader.ReadLine()
		if err != nil || len(line) == 0 {
			break
		}

		if split := bytes.SplitN(line, []byte(`: `), 2); len(split) > 0 {
			var value string

			if len(split) == 2 {
				value = string(split[1])
			}

			r.Headers[string(split[0])] = Header{value}
		}
	}

	// body
	var buf bytes.Buffer
	for {
		conn.SetReadDeadline(time.Now().Add(time.Millisecond * 1))
		line, _, err = reader.ReadLine()
		if err != nil {
			break
		}

		buf.Write(line)
	}
	r.Body = buf.Bytes()

	return nil
}
