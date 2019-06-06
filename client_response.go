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
	"reflect"
	"strconv"
	"strings"
	"time"
)

func (c *Client) NewResponse(method Method, address *Address, conn net.Conn) (*Response, error) {

	response := &Response{
		Base: Base{
			Client:    c,
			Method:    method,
			Address:   address,
			Headers:   make(Headers),
			Cookies:   make(Cookies),
			Params:    make(Params),
			UrlParams: make(UrlParams),
			Charset:   CharsetUTF8,
			conn:      conn,
		},
		FormData:    make(map[string]*FormData),
		Attachments: make(map[string]*Attachment),
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

func (r *Response) BindFormData(obj interface{}) error {
	if len(r.FormData) == 0 {
		return nil
	}

	data := make(map[string]string)
	for _, item := range r.FormData {
		if item.IsAttachment {
			continue
		}

		data[item.Name] = string(item.Body)
	}

	return readData(reflect.ValueOf(obj), data)
}

func setValue(kind reflect.Kind, obj reflect.Value, newValue string) error {

	if !obj.CanAddr() {
		return nil
	}

	switch kind {
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		v, _ := strconv.Atoi(newValue)
		obj.SetInt(int64(v))
	case reflect.Float32, reflect.Float64:
		v, _ := strconv.ParseFloat(newValue, 64)
		obj.SetFloat(v)
	case reflect.String:
		obj.SetString(newValue)
	case reflect.Bool:
		v, _ := strconv.ParseBool(newValue)
		obj.SetBool(v)
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
	if _, ok := MethodHasBody[r.Method]; ok {

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
	r.conn.SetReadDeadline(time.Now().Add(time.Second * 5))
	line, _, err := reader.ReadLine()
	if err != nil {
		return fmt.Errorf("invalid http send: %s", err)
	}

	if firstLine := bytes.SplitN(line, []byte(` `), 3); len(firstLine) < 3 {
		return errors.New("invalid http send")
	} else {
		status, err := strconv.Atoi(string(firstLine[1]))
		if err != nil {
			return fmt.Errorf("invalid http response [%s]", string(line))
		}

		r.Protocol = Protocol(firstLine[0])
		r.Status = Status(status)
		r.StatusText = string(firstLine[2])
	}

	return nil
}

func (r *Response) readHeaders(reader *bufio.Reader) error {
	for {
		r.conn.SetReadDeadline(time.Now().Add(time.Millisecond * 5))
		line, _, err := reader.ReadLine()
		if err != nil || len(line) == 0 {
			break
		}

		if split := bytes.SplitN(line, []byte(`: `), 2); len(split) > 0 {
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
					r.ContentType = ContentType(split[1])
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
				r.Headers[strings.Title(string(split[0]))] = []string{string(split[1])}
			}
		}
	}

	return nil
}

func (r *Response) handleBoundary(reader *bufio.Reader) error {

	// read next line
	r.conn.SetReadDeadline(time.Now().Add(time.Millisecond * 5))
	line, _, err := reader.ReadLine()
	if err != nil {
		return err
	}

	for {
		data := &Data{}
		formDataBody := bytes.NewBuffer(nil)

		for {
			content := bytes.SplitN(line, []byte(`: `), 2)
			switch string(bytes.Title(bytes.TrimSpace(content[0]))) {
			case "Content-Type":
				bytes.Split(content[1], []byte(`;`))
				data.ContentType = ContentType(content[1])

			case "Content-Disposition":
				contentDisposition := bytes.Split(content[1], []byte(`;`))
				data.ContentDisposition = ContentDisposition(string(contentDisposition[0]))
				for i := 1; i < len(contentDisposition); i++ {
					parms := bytes.Split(contentDisposition[i], []byte(`=`))
					switch string(bytes.TrimSpace(parms[0])) {
					case "name":
						data.Name = string(bytes.Replace(parms[1], []byte(`"`), []byte(""), 2))
					case "filename":
						data.FileName = string(bytes.Replace(parms[1], []byte(`"`), []byte(""), 2))
						data.IsAttachment = true
					}
				}
			}

			// read next line
			r.conn.SetReadDeadline(time.Now().Add(time.Millisecond * 5))
			line, _, err = reader.ReadLine()
			if err != nil || len(line) == 0 {
				break
			}
		}

		for {
			r.conn.SetReadDeadline(time.Now().Add(time.Millisecond * 5))
			line, err = reader.ReadSlice('\n')

			if !bytes.HasPrefix(line, []byte(fmt.Sprintf("--%s--", r.Boundary))) { // here we dont have a new line

				if err != nil {
					data.Body = formDataBody.Bytes()
					return err
				}

				if !data.IsAttachment {
					if line[len(line)-1] == '\n' {
						drop := 1
						if len(line) > 1 && line[len(line)-2] == '\r' {
							drop = 2
						}
						line = line[:len(line)-drop]
					}
				}
			}

			// is another boundary ?
			if bytes.HasPrefix(line, []byte(fmt.Sprintf("--%s", r.Boundary))) ||
				bytes.HasPrefix(line, []byte(fmt.Sprintf("--%s--", r.Boundary))) {
				// save data
				data.Body = formDataBody.Bytes()
				key := data.Name
				if key == "" {
					key = data.FileName
				}

				if data.IsAttachment {
					r.Attachments[key] = &Attachment{
						Data: data,
					}
				} else {
					r.FormData[key] = &FormData{
						Data: data,
					}
				}

				// next data
				data = &Data{}
				formDataBody = bytes.NewBuffer(nil)

				break
			} else {
				formDataBody.Write(line)
			}
		}

		if bytes.HasPrefix(line, []byte(fmt.Sprintf("--%s--", r.Boundary))) {
			return nil
		}
	}

	return nil
}

func (r *Response) readBody(reader *bufio.Reader) error {
	var buf bytes.Buffer
	var encoding = EncodingNone

	if enc, ok := r.Headers[HeaderTransferEncoding]; ok {
		encoding = Encoding(enc[0])
	}

	switch encoding {
	case EncodingChunked:
		var size uint64

		for {
			r.conn.SetReadDeadline(time.Now().Add(time.Millisecond * 5))
			line, _, err := reader.ReadLine()
			if err != nil {
				break
			}

			size, _ = parseHexUint(line)
			if size == 0 {
				break
			}

			chunk := make([]byte, size)
			_, err = reader.Read(chunk)
			if err != nil {
				break
			}

			buf.Write(chunk)
		}
	default:
		for {
			r.conn.SetReadDeadline(time.Now().Add(time.Millisecond * 5))
			line, _, err := reader.ReadLine()
			if err != nil {
				break
			}

			buf.Write(line)
		}
	}

	r.Body = buf.Bytes()

	return nil
}
