package gomapper

import (
	"fmt"
	"reflect"
	"strings"
)

// String ...
func (mapper *GoMapper) String(value interface{}) (string, error) {
	var spaces int
	var print string

	if err := convertToString(value, "", spaces, &print); err != nil {
		return "", err
	}

	return print, nil
}

func convertToString(obj interface{}, path string, spaces int, print *string) error {
	types := reflect.TypeOf(obj)
	value := reflect.ValueOf(obj)

	if value.Kind() == reflect.Ptr {
		value = reflect.ValueOf(value).Elem()

		if value.IsValid() {
			types = value.Type()
		} else {
			return nil
		}
	}

	switch value.Kind() {
	case reflect.Struct:
		path = addPoint(path)

		for i := 0; i < types.NumField(); i++ {
			nextValue := value.Field(i)
			newPath := fmt.Sprintf("%s%s", path, strings.ToLower(types.Field(i).Name))
			*print += fmt.Sprintf("\n%s%s", strings.Repeat(" ", spaces), strings.ToLower(types.Field(i).Name))
			convertToString(nextValue.Interface(), newPath, spaces+2, print)
		}

	case reflect.Array, reflect.Slice:
		path = addPoint(path)

		for i := 0; i < value.Len(); i++ {
			nextValue := value.Index(i)
			newPath := fmt.Sprintf("%s[%d]", path, i)
			*print += fmt.Sprintf("\n%s[%d]", strings.Repeat(" ", spaces), i)
			convertToString(nextValue.Interface(), newPath, spaces+2, print)
		}

	case reflect.Map:
		path = addPoint(path)

		for _, key := range value.MapKeys() {
			nextValue := value.MapIndex(key)
			newPath := fmt.Sprintf("%s{%+v}", path, key)
			*print += fmt.Sprintf("\n%s{%+v}", strings.Repeat(" ", spaces), key)
			convertToString(nextValue.Interface(), newPath, spaces+2, print)
		}

	default:
		if value.CanInterface() {
			*print += fmt.Sprintf(":%+v", value.Interface())
			log.Debugf(fmt.Sprintf("%s=%+v", path, value.Interface()))
		} else {
			*print += fmt.Sprintf(":%+v", value)
			log.Debugf(fmt.Sprintf("%s=%+v", path, value))
		}

	}
	return nil
}
