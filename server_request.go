package web

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

func (w *Server) NewRequest(conn net.Conn, server *Server) (*Request, error) {

	request := &Request{
		Base: Base{
			IP:          conn.RemoteAddr().String(),
			Headers:     make(Headers),
			Cookies:     make(Cookies),
			Params:      make(Params),
			UrlParams:   make(UrlParams),
			ContentType: ContentTypeApplicationJSON,
			Charset:     CharsetUTF8,
			conn:        conn,
			server:      server,
		},
		Attachments: make(map[string]Attachment),
		Reader:      conn.(io.Reader),
	}

	return request, request.read()
}

func (r *Request) Bind(i interface{}) error {
	contentType := r.GetContentType()

	if len(r.Body) == 0 || contentType == nil {
		return nil
	}

	switch *contentType {
	case ContentTypeApplicationJSON:
		if err := json.Unmarshal(r.Body, i); err != nil {
			return err
		}
	case ContentTypeApplicationXML:
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

	// header
	if err := r.readHeader(reader); err != nil {
		return err
	}

	// headers
	if err := r.readHeaders(reader); err != nil {
		return err
	}

	// boundary
	if r.Boundary != "" {
		r.handleBoundary(reader)
	} else {
		// body
		if err := r.readBody(reader); err != nil {
			return err
		}
	}

	return nil
}

func (r *Request) readHeader(reader *bufio.Reader) error {

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
		r.Protocol = Protocol(firstLine[2])

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

	return nil
}

func (r *Request) readHeaders(reader *bufio.Reader) error {
	for {
		r.conn.SetReadDeadline(time.Now().Add(time.Millisecond * 1))
		line, _, err := reader.ReadLine()
		if err != nil || len(line) == 0 {
			break
		}

		if split := bytes.SplitN(line, []byte(`:`), 2); len(split) > 0 {
			switch string(bytes.TrimSpace(bytes.Title(split[0]))) {
			case "Cookie":
				var cookieValue string
				splitCookie := bytes.Split(split[1], []byte(`=`))
				if len(splitCookie) > 1 {
					cookieValue = string(splitCookie[1])
				}
				r.Cookies[strings.Title(string(split[0]))] = Cookie{Name: string(splitCookie[0]), Value: cookieValue}
			case "Content-Type":
				if args := bytes.Split(split[1], []byte(`;`)); len(args) > 0 {
					split[1] = bytes.TrimSpace(args[0])
					for _, arg := range args {
						parm := bytes.Split(arg, []byte(`=`))
						switch string(bytes.TrimSpace(parm[0])) {
						case "boundary":
							r.Boundary = string(bytes.Replace(parm[1], []byte(`"`), []byte(``), -1))
						case "charset":
							r.Charset = Charset(bytes.Replace(parm[1], []byte(`"`), []byte(``), -1))
						}
					}
				}
				fallthrough
			default:
				r.Headers[HeaderType(strings.Title(string(split[0])))] = []string{string(bytes.TrimSpace(split[1]))}
			}
		}
	}

	return nil
}

func (r *Request) handleBoundary(reader *bufio.Reader) error {
	var attachment Attachment
	var attachmentBody bytes.Buffer

	// read next line
	r.conn.SetReadDeadline(time.Now().Add(time.Millisecond * 1))
	line, _, err := reader.ReadLine()
	if err != nil {
		return err
	}

	for {
		for {
			content := bytes.SplitN(line, []byte(`:`), 2)
			switch string(bytes.Title(bytes.TrimSpace(content[0]))) {
			case "Content-Type":
				bytes.Split(content[1], []byte(`;`))
				attachment.ContentType = ContentType(content[1])

			case "Content-Disposition":
				contentDisposition := bytes.Split(content[1], []byte(`;`))
				attachment.ContentDisposition = ContentDisposition(string(contentDisposition[0]))
				for i := 1; i < len(contentDisposition); i++ {
					parms := bytes.Split(contentDisposition[i], []byte(`=`))
					switch string(bytes.TrimSpace(parms[0])) {
					case "name":
						attachment.Name = string(bytes.Replace(parms[1], []byte(`"`), []byte(""), 2))
					case "filename":
						attachment.File = string(bytes.Replace(parms[1], []byte(`"`), []byte(""), 2))
					}
				}
			}

			// read next line
			r.conn.SetReadDeadline(time.Now().Add(time.Millisecond * 1))
			line, _, err = reader.ReadLine()
			if err != nil || len(line) == 0 {
				break
			}
		}

		for {
			r.conn.SetReadDeadline(time.Now().Add(time.Millisecond * 1))
			line, _, err = reader.ReadLine()
			if err != nil {
				return err
			}

			// is another boundary ?
			if bytes.Compare(line, []byte(fmt.Sprintf("--%s", r.Boundary))) == 0 ||
				bytes.Compare(line, []byte(fmt.Sprintf("--%s--", r.Boundary))) == 0 {
				// save attachment
				attachment.Body = attachmentBody.Bytes()
				key := attachment.Name
				if key == "" {
					key = attachment.File
				}
				r.Attachments[key] = attachment

				// next attachment
				attachment = Attachment{}
				attachmentBody.Reset()

				break
			} else {
				attachmentBody.Write(line)
			}
		}

		if bytes.Compare(line, []byte(fmt.Sprintf("--%s--", r.Boundary))) == 0 {
			return nil
		}
	}

	return nil
}

func (r *Request) readBody(reader *bufio.Reader) error {
	var buf bytes.Buffer
	for {
		r.conn.SetReadDeadline(time.Now().Add(time.Millisecond * 1))
		line, _, err := reader.ReadLine()
		if err != nil {
			break
		}

		buf.Write(line)
	}
	r.Body = buf.Bytes()

	return nil
}
