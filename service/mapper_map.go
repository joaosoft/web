package gomapper

import (
	"fmt"
	"reflect"
	"strings"
)

// GoMapper ...
type GoMapper struct{}

// NewMapper ...
func NewMapper(options ...GoMapperOption) *GoMapper {
	gomapper := &GoMapper{}

	gomapper.Reconfigure(options...)

	return gomapper
}

// ToMap ...
func (mapper *GoMapper) ToMap(value interface{}) (map[string]interface{}, error) {
	mapping := make(map[string]interface{})

	if err := execute(value, "", mapping); err != nil {
		return nil, err
	}

	return mapping, nil
}

func execute(obj interface{}, path string, mapping map[string]interface{}) error {
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

			if nextValue.CanInterface() {
				execute(nextValue.Interface(), newPath, mapping)
			} else {
				mapping[newPath] = nextValue
				log.Debugf("%s=%+v", newPath, nextValue)
			}
		}

	case reflect.Array, reflect.Slice:
		path = addPoint(path)
		for i := 0; i < value.Len(); i++ {
			nextValue := value.Index(i)
			newPath := fmt.Sprintf("%s[%d]", path, i)

			if nextValue.CanInterface() {
				execute(nextValue.Interface(), newPath, mapping)
			} else {
				mapping[newPath] = nextValue
				log.Debugf("%s=%+v", newPath, nextValue)
			}
		}

	case reflect.Map:
		path = addPoint(path)
		for _, key := range value.MapKeys() {
			nextValue := value.MapIndex(key)
			newPath := fmt.Sprintf("%s{%+v}", path, key)

			if nextValue.CanInterface() {
				execute(nextValue.Interface(), newPath, mapping)
			} else {
				mapping[newPath] = nextValue
				log.Debugf("%s=%+v", newPath, value)
			}
		}

	default:
		if value.CanInterface() {
			mapping[path] = value.Interface()
			log.Debugf(fmt.Sprintf("%s=%+v", path, value.Interface()))
		} else {
			mapping[path] = value
			log.Debugf(fmt.Sprintf("%s=%+v", path, value.Interface()))
		}

	}
	return nil
}

func addPoint(path string) string {
	if path != "" {
		path += "."
	}
	return path
}
