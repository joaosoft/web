package web

import (
	"reflect"
	"strings"
)

func (r *Base) SetHeader(name string, header []string) {
	r.Headers[strings.Title(name)] = header
}

func (r *Base) GetHeader(name string) string {
	if header, ok := r.Headers[strings.Title(name)]; ok {
		return header[0]
	}
	return ""
}

func (r *Base) SetContentType(contentType ContentType) {
	r.ContentType = contentType
}

func (r *Base) GetContentType() *ContentType {
	if value, ok := r.Headers[string(HeaderContentType)]; ok {
		contentType := ContentType(value[0])
		return &contentType
	}
	return nil
}

func (r *Base) SetCharset(charset Charset) {
	r.Charset = charset
}

func (r *Base) GetCharset() Charset {
	return r.Charset
}

func (r *Base) SetCookie(name string, cookie Cookie) {
	r.Cookies[name] = cookie
}

func (r *Base) GetCookie(name string) *Cookie {
	if cookie, ok := r.Cookies[name]; ok {
		return &cookie
	}
	return nil
}

func (r *Base) SetParam(name string, param []string) {
	r.Params[name] = param
}

func (r *Base) GetParam(name string) string {
	if param, ok := r.Params[name]; ok {
		return param[0]
	}
	return ""
}

func (r *Base) GetParams(name string) []string {
	if param, ok := r.Params[name]; ok {
		return param
	}
	return nil
}

func (r *Base) SetUrlParam(name string, urlParam []string) {
	r.UrlParams[name] = urlParam
}

func (r *Base) GetUrlParam(name string) string {
	if urlParam, ok := r.UrlParams[name]; ok {
		return urlParam[0]
	}
	return ""
}

func (r *Base) GetUrlParams(name string) []string {
	if urlParam, ok := r.UrlParams[name]; ok {
		return urlParam
	}
	return nil
}

func (r *Base) BindParams(obj interface{}) error {
	if len(r.Params) == 0 {
		return nil
	}

	data := make(map[string]string)
	for name, values := range r.Params {
		data[name] = values[0]
	}

	return readData(reflect.ValueOf(obj), data)
}

func (r *Base) BindUrlParams(obj interface{}) error {
	if len(r.UrlParams) == 0 {
		return nil
	}

	data := make(map[string]string)
	for name, values := range r.UrlParams {
		data[name] = values[0]
	}

	return readData(reflect.ValueOf(obj), data)
}

func readData(obj reflect.Value, data map[string]string) error {
	types := reflect.TypeOf(obj)

	if !obj.CanInterface() {
		return nil
	}

	if obj.Kind() == reflect.Ptr && !obj.IsNil() {
		obj = obj.Elem()

		if obj.IsValid() {
			types = obj.Type()
		} else {
			return nil
		}
	}

	switch obj.Kind() {
	case reflect.Struct:
		for i := 0; i < types.NumField(); i++ {
			nextValue := obj.Field(i)
			nextType := types.Field(i)

			if nextValue.Kind() == reflect.Ptr {
				if !nextValue.IsNil() {
					nextValue = nextValue.Elem()
				} else {
					isSlice := nextValue.Kind() == reflect.Slice
					isMap := nextValue.Kind() == reflect.Map
					isMapOfSlices := isMap && nextValue.Type().Elem().Kind() == reflect.Slice

					if isMapOfSlices {
						nextValue = reflectAlloc(nextValue.Type().Elem().Elem())
					} else if isSlice || isMap {
						nextValue = reflectAlloc(nextValue.Type().Elem())
					} else {
						nextValue = reflect.Value{}
					}
				}
			}

			if !nextValue.CanInterface() {
				continue
			}

			var tagName string
			jsonName, exists := nextType.Tag.Lookup("json")
			if exists {
				tagName = strings.SplitN(jsonName, ",", 2)[0]
			}

			if value, ok := data[tagName]; ok {
				if err := setValue(nextValue.Kind(), nextValue, value); err != nil {
					return err
				}
			}

			if err := readData(nextValue, data); err != nil {
				return err
			}
		}

	case reflect.Array, reflect.Slice:
		for i := 0; i < obj.Len(); i++ {
			nextValue := obj.Index(i)

			if !nextValue.CanInterface() {
				continue
			}

			if err := readData(nextValue, data); err != nil {
				return err
			}
		}
	case reflect.Map:

	default:
		// do nothing ...
	}
	return nil
}