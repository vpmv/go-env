package env

import (
	"os"
	"reflect"
	"strconv"
	"strings"
)

type BasicType interface {
	~string | ~int | ~float64 | ~uint64 | ~bool | ~[]string | ~[]int
}

// Set will convert all basic types to string and set the environment variable
func Set[T BasicType](key string, value T) {
	var (
		s string
		v = reflect.ValueOf(value)
	)

	switch v.Kind() {
	case reflect.String:
		s = v.String()
	case reflect.Int:
		s = strconv.Itoa(int(v.Int()))
	case reflect.Float64:
		s = strconv.FormatFloat(v.Float(), 'f', -1, 64) // no trailing zeros
	case reflect.Uint64:
		s = strconv.FormatUint(v.Uint(), 10)
	case reflect.Bool:
		s = strconv.FormatBool(v.Bool())
	case reflect.Slice:
		s = strings.Join(v.Slice(0, v.Len()).Interface().([]string), delimiter)
	}

	_ = os.Setenv(key, s)
}
