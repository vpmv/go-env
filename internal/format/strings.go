package format

import (
	"errors"
	"reflect"
	"strconv"
	"strings"
)

var Stringer = NewStringFormatter()

type StringFormatter struct {
	delimiter string
}

type FormatterOption func(s *StringFormatter)

func WithDelimiter(d string) FormatterOption {
	return func(s *StringFormatter) {
		s.delimiter = d
	}
}

func (s StringFormatter) ToString(value any) (str string, err error) {
	var (
		v = reflect.ValueOf(value)
	)

	switch v.Kind() {
	case reflect.String:
		str = v.String()
	case reflect.Int,
		reflect.Int8,
		reflect.Int16,
		reflect.Int32,
		reflect.Int64:
		str = strconv.FormatInt(v.Int(), 10)
	case reflect.Float32,
		reflect.Float64:
		str = strconv.FormatFloat(v.Float(), 'f', -1, 64) // no trailing zeros
	case reflect.Uint,
		reflect.Uint8,
		reflect.Uint16,
		reflect.Uint32,
		reflect.Uint64:
		str = strconv.FormatUint(v.Uint(), 10)
	case reflect.Bool:
		str = strconv.FormatBool(v.Bool())
	case reflect.Array,
		reflect.Slice:
		str = strings.Join(s.stringSlice(v), s.delimiter)
	default:
		return ``, errors.New("unsupported type: " + v.Type().String())
	}

	return str, nil
}

func (s StringFormatter) stringSlice(v reflect.Value) []string {
	var parts = make([]string, v.Len())

	for i := 0; i < v.Len(); i++ {
		parts[i], _ = s.ToString(v.Index(i).Interface())
	}

	return parts
}

func (s StringFormatter) Split(str string) []string {
	return strings.Split(str, s.delimiter)
}

func NewStringFormatter(options ...FormatterOption) *StringFormatter {
	s := &StringFormatter{`;`}
	for _, option := range options {
		option(s)
	}
	return s
}
