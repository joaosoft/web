package types

import (
	"database/sql/driver"
	"encoding/json"
	"fmt"
	"strings"
	"time"

	errors "github.com/joaosoft/errors"
)

func (d *Date) Scan(src interface{}) error {
	if d == nil {
		return nil
	}

	switch value := src.(type) {
	case string:
		*d = Date(value)
		return nil
	case time.Time:
		*d = Date(value.Format("02-01-2006"))
		return nil
	case nil:
		d = nil
		return nil
	}

	return errors.New("0", "pq: cannot convert %T to %T", src, *d)
}

func (d Date) Value() (driver.Value, error) {
	if d == "" {
		return nil, nil
	}
	return string(d), nil
}

func (t *Time) Scan(src interface{}) error {
	if t == nil {
		return nil
	}
	switch value := src.(type) {
	case string:
		*t = Time(value)
		return nil
	case time.Time:
		*t = Time(value.Format("15:04:05"))
		return nil
	case nil:
		t = nil
		return nil
	}

	return errors.New("0", "pq: cannot convert %T to %T", src, *t)
}

func (t Time) Value() (driver.Value, error) {
	if t == "" {
		return nil, nil
	}

	return string(t), nil
}

func (d *Day) Scan(src interface{}) error {
	if d == nil {
		return nil
	}
	switch value := src.(type) {
	case string:
		*d = Day(value)
		return nil
	case nil:
		d = nil
		return nil
	}

	return errors.New("0", "pq: cannot convert %T to %T", src, *d)
}

func (d Day) Value() (driver.Value, error) {
	if d == "" {
		return nil, nil
	}

	return string(d), nil
}

func (a *ListDay) Scan(src interface{}) error {
	if a == nil {
		return nil
	}

	var err error
	switch value := src.(type) {
	case string:
		value = strings.NewReplacer("{", "]", "}", "]").Replace(value)
		err = json.Unmarshal([]byte(value), a)
	case []byte:
		value = []byte(strings.NewReplacer("{", "[\"", "}", "\"]", ",", "\",\"").Replace(string(value)))
		err = json.Unmarshal(value, a)
	case nil:
		a = nil
		return nil
	}

	if err != nil {
		return errors.New("0", "pq: cannot convert %T to Time", src)
	}

	return nil
}

// Value implements the driver.Valuer interface.
func (a ListDay) Value() (driver.Value, error) {
	if a == nil {
		return nil, nil
	}

	return strings.NewReplacer(" ", ",", "[", "{", "]", "}").Replace(fmt.Sprintf("%+v", a)), nil
}
