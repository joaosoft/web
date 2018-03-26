package gomapper

import (
	"fmt"
	"reflect"
)

// Map ...
func (mapper *GoMapper) Map(value interface{}) (map[string]interface{}, error) {
	mapping := make(map[string]interface{})

	if _, err := convertToMap(value, "", mapping, true); err != nil {
		return nil, err
	}

	return mapping, nil
}

func convertToMap(obj interface{}, path string, mapping map[string]interface{}, add bool) (string, error) {
	types := reflect.TypeOf(obj)
	value := reflect.ValueOf(obj)

	if value.Kind() == reflect.Ptr {
		value = reflect.ValueOf(value).Elem()

		if value.IsValid() {
			types = value.Type()
		} else {
			return "", nil
		}
	}

	switch value.Kind() {
	case reflect.Struct:
		path = addPoint(path)
		var innerPath string
		var len = value.NumField()
		for i := 0; i < types.NumField(); i++ {
			len--
			nextValue := value.Field(i)
			newPath := fmt.Sprintf("%s%s", path, types.Field(i).Name)
			tmp, _ := convertToMap(nextValue.Interface(), newPath, mapping, add)
			if len > 0 {
				tmp += ","
			}
			innerPath += fmt.Sprintf("%s", tmp)
		}

		if !add {
			return innerPath, nil
		}

	case reflect.Array, reflect.Slice:
		path = addPoint(path)
		var innerPath string
		var len = value.Len()
		for i := 0; i < value.Len(); i++ {
			len--
			nextValue := value.Index(i)
			newPath := fmt.Sprintf("%s[%d]", path, i)
			tmp, _ := convertToMap(nextValue.Interface(), newPath, mapping, add)
			if len > 0 {
				tmp += ","
			}
			innerPath += fmt.Sprintf("%s", tmp)
		}

		if !add {
			return innerPath, nil
		}

	case reflect.Map:
		path = addPoint(path)
		var innerPath string
		var len = value.Len()
		for _, key := range value.MapKeys() {
			len--
			nextValue := value.MapIndex(key)
			newPath := fmt.Sprintf("%s{", path)
			keyValue, _ := convertToMap(key.Interface(), "", mapping, false)
			newPath += fmt.Sprintf("%s}", keyValue)
			tmp, _ := convertToMap(nextValue.Interface(), newPath, mapping, add)
			if len > 0 {
				tmp += ","
			}
			innerPath += fmt.Sprintf("%s", tmp)
		}

		if !add {
			return innerPath, nil
		}

	default:
		if value.IsValid() {
			var rtnValue interface{}
			if value.CanInterface() {
				rtnValue = value.Interface()
				log.Debugf(fmt.Sprintf("%s=%+v", path, value.Interface()))
			} else {
				rtnValue = value
				log.Debugf(fmt.Sprintf("%s=%+v", path, value))
			}

			if add {
				mapping[path] = rtnValue
			} else {
				if path != "" {
					path += "="
				}
				return fmt.Sprintf("%s%+v", path, rtnValue), nil
			}
		}
	}
	return "", nil
}

func addPoint(path string) string {
	if path != "" {
		path += "."
	}
	return path
}
