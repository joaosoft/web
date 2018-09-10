package webserver

import (
	"bufio"
	"bytes"
	"encoding/json"
	"encoding/xml"
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
			IP:        conn.RemoteAddr().String(),
			Headers:   make(Headers),
			Cookies:   make(Cookies),
			Params:    make(Params),
			UrlParams: make(UrlParams),
			conn:      conn,
		},
		Reader: conn.(io.Reader),
	}

	return request, request.read()
}

func (r *Request) Bind(i interface{}) error {
	contentType := r.GetContentType()

	if len(r.Body) == 0 || contentType == nil {
		return nil
	}

	switch *contentType {
	case ContentApplicationJSON:
		if err := json.Unmarshal(r.Body, i); err != nil {
			return err
		}
	case ContentApplicationXML:
		if err := xml.Unmarshal(r.Body, i); err != nil {
			return err
		}
	default:
		tmp := string(r.Body)
		i = &tmp
	}
	return nil
}

func (r *Request) read() error {
	reader := bufio.NewReader(r.conn)

	// read one line (ended with \n or \r\n)
	r.conn.SetReadDeadline(time.Now().Add(time.Second * 1))
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
		if split := strings.SplitN(r.FullUrl, "?", 2); len(split) > 1 {
			r.Url = string(split[0])
			if parms := strings.Split(split[1], "&"); len(parms) > 0 {
				for _, parm := range parms {
					if p := strings.Split(parm, "="); len(p) > 1 {
						if split := strings.SplitN(p[1], ",", -1); len(split) > 0 {
							for _, s := range split {
								r.Params[p[0]] = append(r.Params[p[0]], s)
							}
						}
						r.Params[p[0]] = append(r.Params[p[0]], p[1])
					}
				}
			}
		} else {
			r.Url = string(firstLine[1])
		}
	}

	// headers
	for {
		r.conn.SetReadDeadline(time.Now().Add(time.Millisecond * 1))
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
				r.Headers[HeaderType(string(split[0]))] = []string{value}
			}
		}
	}

	// body
	var buf bytes.Buffer
	for {
		r.conn.SetReadDeadline(time.Now().Add(time.Millisecond * 1))
		line, _, err = reader.ReadLine()
		if err != nil {
			break
		}

		buf.Write(line)
	}
	r.Body = buf.Bytes()

	return nil
}
