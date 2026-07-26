package env

import (
	"fmt"
	"os"
	"strconv"
)

// Has checks if the environment variable exists
func Has(key string) bool {
	_, ok := os.LookupEnv(key)
	return ok
}

// GetString - get string with default fallback
func GetString(key, defaultValue string) string {
	v := os.Getenv(key)
	if v == `` {
		return defaultValue
	}
	return v
}

// GetInt - get int with default fallback
func GetInt(key string, defaultValue int) int {
	v := os.Getenv(key)
	if v == `` {
		return defaultValue
	}

	i, err := strconv.Atoi(v)
	if err != nil {
		return defaultValue
	}
	return i
}

// GetUInt - get uint with default fallback
func GetUInt(key string, defaultValue uint64) uint64 {
	return uint64(GetInt(key, int(defaultValue)))
}

// GetFloat - get float64 with default fallback
func GetFloat(key string, defaultValue float64) float64 {
	v := os.Getenv(key)
	if v == `` {
		return defaultValue
	}

	i, err := strconv.ParseFloat(v, 64)
	if err != nil {
		return defaultValue
	}
	return i
}

// GetBool - get boolean with default fallback
func GetBool(key string, defaultValue bool) bool {
	v := os.Getenv(key)
	if v == `` {
		return defaultValue
	}

	b, err := strconv.ParseBool(v)
	if err != nil {
		return defaultValue
	}
	return b
}

// GetStringSlice - get string slice with default fallback
func GetStringSlice(key string, defaultValue []string) []string {
	v := os.Getenv(key)
	if v == `` {
		return defaultValue
	}

	s := stringer.Split(v)
	return s
}

// GetIntSlice - get int slice with default fallback
//
// may produce zero-values when unable to convert
func GetIntSlice(key string, defaultValue []int) []int {
	v := os.Getenv(key)
	if v == `` {
		return defaultValue
	}

	s := stringer.Split(v)
	o := make([]int, len(s))
	for i := range s {
		iv, _ := strconv.Atoi(s[i])
		o[i] = iv
	}
	return o
}

// MustString - get string or panic
func MustString(key string) string {
	v := os.Getenv(key)
	if v == `` {
		panic(fmt.Sprintf(`Missing required environment variable: %s`, key))
	}
	return v
}

// MustInt - get int or panic
func MustInt(key string) int {
	v := MustString(key)
	i, err := strconv.Atoi(v)
	if err != nil {
		panic(fmt.Sprintf(`Invalid integer value for environment variable: %s`, key))
	}
	return i
}

// MustUInt - get uint64 or panic
func MustUInt(key string) uint64 {
	return uint64(MustInt(key))
}

// MustFloat - get float64 or panic
func MustFloat(key string) float64 {
	v := MustString(key)
	i, err := strconv.ParseFloat(v, 64)
	if err != nil {
		panic(fmt.Sprintf(`Invalid float value for environment variable: %s`, key))
	}
	return i
}

// MustBool - get boolean or panic
func MustBool(key string) bool {
	v := MustString(key)
	b, err := strconv.ParseBool(v)
	if err != nil {
		panic(fmt.Sprintf(`Invalid boolean value for environment variable: %s`, key))
	}
	return b
}

// MustStringSlice - get string slice or panic
func MustStringSlice(key string) []string {
	v := MustString(key)
	return stringer.Split(v)
}

// MustIntSlice - get integer slice or panic
func MustIntSlice(key string) []int {
	v := MustStringSlice(key)
	o := make([]int, len(v))
	for i := range v {
		iv, err := strconv.Atoi(v[i])
		if err != nil {
			panic(fmt.Sprintf(`Invalid integer value for environment variable: %s`, key))
		}
		o[i] = iv
	}
	return o
}
