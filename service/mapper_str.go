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

	if err := convertToString(value, "", spaces, "\n", &print); err != nil {
		return "", err
	}

	return print, nil
}

func convertToString(obj interface{}, path string, spaces int, delimiter string, print *string) error {
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
		for i := 0; i < types.NumField(); i++ {
			nextValue := value.Field(i)
			newPath := fmt.Sprintf("%s%s", path, types.Field(i).Name)
			*print += fmt.Sprintf("%s%s%s", delimiter, strings.Repeat(" ", spaces), types.Field(i).Name)
			convertToString(nextValue.Interface(), newPath, spaces+2, delimiter, print)
		}

	case reflect.Array, reflect.Slice:
		for i := 0; i < value.Len(); i++ {
			nextValue := value.Index(i)
			newPath := fmt.Sprintf("[%d]", i)
			*print += fmt.Sprintf("%s%s%s", delimiter, strings.Repeat(" ", spaces), newPath)
			convertToString(nextValue.Interface(), newPath, spaces+2, delimiter, print)
		}

	case reflect.Map:
		for _, key := range value.MapKeys() {
			var keyValue string
			nextValue := value.MapIndex(key)
			convertToString(key.Interface(), "", 0, " ", &keyValue)
			newPath := fmt.Sprintf("{%s}", keyValue)
			*print += fmt.Sprintf("%s%s%s", delimiter, strings.Repeat(" ", spaces), newPath)
			convertToString(nextValue.Interface(), newPath, spaces+2, delimiter, print)
		}

	default:
		var rtnValue interface{}
		if value.IsValid() {
			if value.CanInterface() {
				rtnValue = value.Interface()
			} else {
				rtnValue = value
			}

			if path != "" {
				*print += ": "
			}

			newPath := fmt.Sprintf("%s=%+v", path, rtnValue)
			*print += fmt.Sprintf("%+v", rtnValue)

			log.Debugf(fmt.Sprintf("%s", newPath))
		}
	}
	return nil
}
