package web

import (
	"bufio"
	"bytes"
	"errors"
	"fmt"
	"io"
	"net"
	"strings"
	"time"
)

func NewRequest(conn net.Conn) (*Request, error) {

	request := &Request{
		Base: Base{
			IP:       conn.RemoteAddr().String(),
			Headers:  make(Headers),
			Cookies:  make(Cookies),
			Params:   make(Params),
			UrlParms: make(UrlParms),
			conn:     conn,
		},
		Reader: conn.(io.Reader),
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
					r.Params[parm[0]] = []string{parm[1]}
				}
			} else {
				if parm := strings.Split(r.FullUrl, "="); len(parm) > 0 {
					r.Params[parm[0]] = []string{parm[1]}
				}
			}
		} else {
			r.Url = string(firstLine[1])
		}
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

			switch string(split[0]) {
			case "cookie":
				var cookieValue string
				splitCookie := strings.Split(value, "=")
				if len(splitCookie) > 1 {
					cookieValue = splitCookie[1]
				}
				r.Cookies[string(split[0])] = Cookie{Name: splitCookie[0], Value: cookieValue}
			default:
				r.Headers[HeaderType(string(split[0]))] = Header{value}
			}
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
