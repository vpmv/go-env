package env

import (
	"os"
)

type BasicType interface {
	~string | ~int | ~float64 | ~uint64 | ~bool | ~[]string | ~[]int
}

// Set will convert all basic types to string and set the environment variable
func Set[T BasicType](key string, value T) {
	s, _ := stringer.ToString(value)
	_ = os.Setenv(key, s)
}
