package env

import (
	"os"

	"github.com/vpmv/go-env/internal/format"
)

type BasicType interface {
	~string |
		~int |
		~int8 |
		~int16 |
		~int32 |
		~int64 |
		~float32 |
		~float64 |
		~uint |
		~uint8 |
		~uint16 |
		~uint32 |
		~uint64 |
		~bool |
		~[]any |
		~[]string |
		~[]float32 |
		~[]float64 |
		~[]int |
		~[]uint |
		~[]uint8 |
		~[]uint16 |
		~[]uint32 |
		~[]uint64
}

// Set will convert all basic types to string and set the environment variable
func Set[T BasicType](key string, value T) {
	s, _ := format.Stringer.ToString(value)
	_ = os.Setenv(key, s)
}
