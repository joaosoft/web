package client

import (
	"bufio"
	"bytes"
	"encoding/json"
	"encoding/xml"
	"errors"
	"fmt"
	"io"
	"net"
	"strconv"
	"strings"
	"time"
	"web/common"
)

func (c *Client) NewResponse(method common.Method, conn net.Conn) (*Response, error) {

	response := &Response{
		Base: Base{
			Method:    method,
			Headers:   make(common.Headers),
			Cookies:   make(common.Cookies),
			Params:    make(common.Params),
			UrlParams: make(common.UrlParams),
			Charset:   common.CharsetUTF8,
			conn:      conn,
		},
		Attachments: make(map[string]Attachment),
		Reader:      conn.(io.Reader),
	}

	return response, response.read()
}

func (r *Response) Bind(i interface{}) error {
	contentType := r.GetContentType()

	if len(r.Body) == 0 || contentType == nil {
		return nil
	}

	switch *contentType {
	case common.ContentTypeApplicationJSON:
		if err := json.Unmarshal(r.Body, i); err != nil {
			return err
		}
	case common.ContentTypeApplicationXML:
		if err := xml.Unmarshal(r.Body, i); err != nil {
			return err
		}
	default:
		tmp := string(r.Body)
		i = &tmp
	}
	return nil
}

func (r *Response) read() error {
	reader := bufio.NewReader(r.Reader)

	// header
	if err := r.readHeader(reader); err != nil {
		return err
	}

	// headers
	if err := r.readHeaders(reader); err != nil {
		return err
	}

	// body
	if _, ok := common.MethodHasBody[r.Method]; ok {

		// boundary
		if r.Boundary != "" {
			r.handleBoundary(reader)
		} else {

			// body
			if err := r.readBody(reader); err != nil {
				return err
			}
		}
	}

	return nil
}

func (r *Response) readHeader(reader *bufio.Reader) error {

	// read one line (ended with \n or \r\n)
	r.conn.SetReadDeadline(time.Now().Add(time.Second * 1))
	line, _, err := reader.ReadLine()
	if err != nil {
		return fmt.Errorf("invalid http send: %s", err)
	}

	if firstLine := bytes.SplitN(line, []byte(` `), 3); len(firstLine) < 3 {
		return errors.New("invalid http send")
	} else {
		status, err := strconv.Atoi(string(firstLine[1]))
		if err != nil {
			return err
		}

		r.Protocol = common.Protocol(firstLine[0])
		r.Status = common.Status(status)
		r.StatusText = string(firstLine[2])
	}

	return nil
}

func (r *Response) readHeaders(reader *bufio.Reader) error {
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
				r.Cookies[strings.Title(string(split[0]))] = common.Cookie{Name: string(splitCookie[0]), Value: cookieValue}
			case "Content-Type":
				if args := bytes.Split(split[1], []byte(`;`)); len(args) > 0 {
					split[1] = bytes.TrimSpace(args[0])
					r.ContentType = common.ContentType(split[1])
					for _, arg := range args {
						parm := bytes.Split(arg, []byte(`=`))
						switch string(bytes.TrimSpace(parm[0])) {
						case "boundary":
							r.Boundary = string(bytes.Replace(parm[1], []byte(`"`), []byte(``), -1))
						case "charset":
							r.Charset = common.Charset(bytes.Replace(parm[1], []byte(`"`), []byte(``), -1))
						}
					}
				}
				fallthrough
			default:
				r.Headers[common.HeaderType(strings.Title(string(split[0])))] = []string{string(split[1])}
			}
		}
	}

	return nil
}

func (r *Response) handleBoundary(reader *bufio.Reader) error {
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
				attachment.ContentType = common.ContentType(content[1])

			case "Content-Disposition":
				contentDisposition := bytes.Split(content[1], []byte(`;`))
				attachment.ContentDisposition = common.ContentDisposition(string(contentDisposition[0]))
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

func (r *Response) readBody(reader *bufio.Reader) error {
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
