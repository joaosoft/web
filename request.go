package webserver

import (
	"bufio"
	"bytes"
	"errors"
	"fmt"
	"net"
	"strings"
	"time"
)

func NewRequest(conn net.Conn) (*Request, error) {
	request := &Request{
		Base: Base{
			Headers:  make(Headers),
			Cookies:  make(Cookies),
			Parms:    make(Parms),
			UrlParms: make(UrlParms),
		},
	}

	return request, request.read(conn)
}

func (r *Request) read(conn net.Conn) error {
	reader := bufio.NewReader(conn)

	// read one line (ended with \n or \r\n)
	conn.SetReadDeadline(time.Now().Add(time.Millisecond * 1))
	line, _, err := reader.ReadLine()
	fmt.Println(string(line))
	if err != nil {
		fmt.Println(err)
		return fmt.Errorf("invalid http request: %s", err)
	}

	if firstLine := bytes.SplitN(line, []byte(` `), 3); len(firstLine) < 3 {
		return errors.New("invalid http request")
	} else {
		r.Method = Method(firstLine[0])
		r.FullUrl = string(firstLine[1])
		r.Protocol = string(firstLine[2])

		// load query parameters
		if split := strings.SplitN(r.FullUrl, "?", 2); len(split) > 0 {
			r.Url = string(split[0])
			if parms := strings.Split(r.FullUrl, "&"); len(parms) > 0 {
				if parm := strings.Split(r.FullUrl, "="); len(parm) > 1 {
					r.Parms[parm[0]] = []string{parm[1]}
				}
			} else {
				if parm := strings.Split(r.FullUrl, "="); len(parm) > 0 {
					r.Parms[parm[0]] = []string{parm[1]}
				}
			}
		} else {
			r.Url = string(firstLine[1])
		}
	}

	fmt.Printf("%+v", *r)

	// headers
	for {
		conn.SetReadDeadline(time.Now().Add(time.Millisecond * 1))
		line, _, err = reader.ReadLine()
		fmt.Println(string(line))
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