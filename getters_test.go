package env

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetString(t *testing.T) {
	assert.Equal(t, "default_fallback", GetString(`TEST_VAL`, `default_fallback`), `Unexpected value for TEST_VAL`)

	_ = os.Setenv("TEST_VAL", "test_value")
	assert.Equal(t, "test_value", GetString(`TEST_VAL`, `different`), `Unexpected value for TEST_VAL`)
	assert.Equal(t, "test_value", MustString(`TEST_VAL`), `Unexpected value for TEST_VAL`)

	assert.True(t, Has(`TEST_VAL`))

	_ = os.Unsetenv(`TEST_VAL`)
	assert.False(t, Has(`TEST_VAL`))
	assert.Panics(t, func() { MustString(`TEST_VAL`) }, `MustString should panic when env var is not set`)
}

func TestGetInt(t *testing.T) {
	assert.Equal(t, 123, GetInt(`TEST_VAL`, 123), `Unexpected value for TEST_VAL`)

	_ = os.Setenv("TEST_VAL", `123`)
	assert.Equal(t, 123, GetInt(`TEST_VAL`, 321), `Unexpected value for TEST_VAL`)
	assert.Equal(t, 123, MustInt(`TEST_VAL`), `Unexpected value for TEST_VAL`)

	_ = os.Setenv("TEST_VAL", `123A`)
	assert.Equal(t, 321, GetInt(`TEST_VAL`, 321), `Unexpected fallback value for TEST_VAL`)
	_ = os.Setenv("TEST_VAL", `123Aa`)
	assert.Panics(t, func() { MustInt(`TEST_VAL`) }, `MustInt should panic when value is invalid`)

	_ = os.Unsetenv(`TEST_VAL`)
	assert.Panics(t, func() { MustInt(`TEST_VAL`) }, `MustInt should panic when env var is not set`)
}

// TestGetUInt - uses GetInt as its backbone so cases shouldn't fail
func TestGetUInt(t *testing.T) {
	var expected uint64 = 123

	assert.Equal(t, expected, GetUInt(`TEST_VAL`, 123), `Unexpected value for TEST_VAL`)

	_ = os.Setenv("TEST_VAL", `123`)
	assert.Equal(t, expected, GetUInt(`TEST_VAL`, 321), `Unexpected value for TEST_VAL`)
	assert.Equal(t, expected, MustUInt(`TEST_VAL`), `Unexpected value for TEST_VAL`)

	_ = os.Setenv("TEST_VAL", `123A`)
	expected = 321
	assert.Equal(t, expected, GetUInt(`TEST_VAL`, 321), `Unexpected fallback value for TEST_VAL`)

	_ = os.Unsetenv(`TEST_VAL`)
}

func TestGetFloat(t *testing.T) {
	assert.Equal(t, 3.14159, GetFloat(`TEST_VAL`, 3.14159), `Unexpected value for TEST_VAL`)

	_ = os.Setenv("TEST_VAL", `3.14159`)
	assert.Equal(t, 3.14159, GetFloat(`TEST_VAL`, 3.21), `Unexpected value for TEST_VAL`)
	assert.Equal(t, 3.14159, MustFloat(`TEST_VAL`), `Unexpected value for TEST_VAL`)

	_ = os.Setenv("TEST_VAL", `3.14159A`)
	assert.Equal(t, 1.337, GetFloat(`TEST_VAL`, 1.337), `Unexpected fallback value for TEST_VAL`)
	_ = os.Setenv("TEST_VAL", `3.14159Aa`)
	assert.Panics(t, func() { MustFloat(`TEST_VAL`) }, `MustFloat should panic when value is invalid`)

	_ = os.Unsetenv(`TEST_VAL`)
	assert.Panics(t, func() { MustFloat(`TEST_VAL`) }, `MustFloat should panic when env var is not set`)
}

func TestGetBool(t *testing.T) {
	assert.Equal(t, false, GetBool(`TEST_VAL`, false), `Unexpected value for TEST_VAL`)

	_ = os.Setenv("TEST_VAL", `true`)
	assert.Equal(t, true, GetBool(`TEST_VAL`, false), `Unexpected value for TEST_VAL`)
	assert.Equal(t, true, MustBool(`TEST_VAL`), `Unexpected value for TEST_VAL`)

	_ = os.Setenv("TEST_VAL", `yes`)
	assert.Equal(t, false, GetBool(`TEST_VAL`, false), `Unexpected fallback value for TEST_VAL`)
	assert.Panics(t, func() { MustBool(`TEST_VAL`) }, `MustBool should panic when env var is not set`)

	_ = os.Unsetenv(`TEST_VAL`)
	assert.Panics(t, func() { MustBool(`TEST_VAL`) }, `MustBool should panic when env var is not set`)
}

func TestGetStringSlice(t *testing.T) {
	expected := []string{`foo`, `bar`, `baz`}
	assert.Equal(t, expected, GetStringSlice(`TEST_VAL`, expected), `Unexpected value for TEST_VAL`)

	_ = os.Setenv("TEST_VAL", `foo;bar;baz`)
	assert.Equal(t, expected, GetStringSlice(`TEST_VAL`, []string{`hello`}), `Unexpected value for TEST_VAL`)
	assert.Equal(t, expected, MustStringSlice(`TEST_VAL`), `Unexpected value for TEST_VAL`)

	_ = os.Setenv("TEST_VAL", `hello`)
	expected = []string{`hello`}
	assert.Equal(t, expected, GetStringSlice(`TEST_VAL`, []string{`world`}), `Unexpected fallback value for TEST_VAL`)

	_ = os.Setenv("TEST_VAL", ``)
	expected = []string{`hello`}
	assert.Equal(t, expected, GetStringSlice(`TEST_VAL`, expected), `Unexpected fallback value for TEST_VAL`)

	_ = os.Unsetenv(`TEST_VAL`)
	assert.Panics(t, func() { MustStringSlice(`TEST_VAL`) }, `MustStringSlice should panic when env var is not set`)
}

func TestGetIntSlice(t *testing.T) {
	expected := []int{101, 202, 303}
	assert.Equal(t, expected, GetIntSlice(`TEST_VAL`, expected), `Unexpected value for TEST_VAL`)

	_ = os.Setenv("TEST_VAL", `101;202;303`)
	assert.Equal(t, expected, GetIntSlice(`TEST_VAL`, []int{404}), `Unexpected value for TEST_VAL`)
	assert.Equal(t, expected, MustIntSlice(`TEST_VAL`), `Unexpected value for TEST_VAL`)

	// test erroneous values
	_ = os.Setenv("TEST_VAL", `101;aaa;303`)
	expected = []int{101, 0, 303}
	assert.Equal(t, expected, GetIntSlice(`TEST_VAL`, []int{404}), `Unexpected value for TEST_VAL`)
	assert.Panics(t, func() { MustIntSlice(`TEST_VAL`) }, `MustIntSlice should panic when value is invalid`)

	_ = os.Setenv("TEST_VAL", `101`)
	expected = []int{101}
	assert.Equal(t, expected, GetIntSlice(`TEST_VAL`, []int{202}), `Unexpected value for TEST_VAL`)

	_ = os.Setenv("TEST_VAL", ``)
	expected = []int{101}
	assert.Equal(t, expected, GetIntSlice(`TEST_VAL`, expected), `Unexpected fallback value for TEST_VAL`)

	_ = os.Unsetenv(`TEST_VAL`)
	assert.Panics(t, func() { MustIntSlice(`TEST_VAL`) }, `MustIntSlice should panic when env var is not set`)
}

func TestGetFloatSlice(t *testing.T) {
	expected := []float64{101.011, 202.022, 303.033}
	assert.Equal(t, expected, GetFloatSlice(`TEST_VAL`, expected), `Unexpected value for TEST_VAL`)

	_ = os.Setenv("TEST_VAL", `101.011;202.022;303.033`)
	assert.Equal(t, expected, GetFloatSlice(`TEST_VAL`, []float64{404}), `Unexpected value for TEST_VAL`)
	assert.Equal(t, expected, MustFloatSlice(`TEST_VAL`), `Unexpected value for TEST_VAL`)

	// test erroneous values
	_ = os.Setenv("TEST_VAL", `101.011;aaa;303.033`)
	expected = []float64{101.011, 0, 303.033}
	assert.Equal(t, expected, GetFloatSlice(`TEST_VAL`, []float64{404}), `Unexpected value for TEST_VAL`)
	assert.Panics(t, func() { MustFloatSlice(`TEST_VAL`) }, `MustFloatSlice should panic when value is invalid`)

	_ = os.Setenv("TEST_VAL", `101.011`)
	expected = []float64{101.011}
	assert.Equal(t, expected, GetFloatSlice(`TEST_VAL`, []float64{202.022}), `Unexpected value for TEST_VAL`)

	_ = os.Setenv("TEST_VAL", ``)
	expected = []float64{101.011}
	assert.Equal(t, expected, GetFloatSlice(`TEST_VAL`, expected), `Unexpected fallback value for TEST_VAL`)

	_ = os.Unsetenv(`TEST_VAL`)
	assert.Panics(t, func() { MustFloatSlice(`TEST_VAL`) }, `MustFloatSlice should panic when env var is not set`)
}

func TestDelimiter(t *testing.T) {
	expected := []string{`foo`, `bar`, `baz`}

	SetDelimiter(`,`)
	Set(`TEST_VAL`, expected)
	assert.Equal(t, `foo,bar,baz`, GetString(`TEST_VAL`, `fallback`), `Unexpected value for TEST_VAL`)
	assert.Equal(t, expected, MustStringSlice(`TEST_VAL`), `Unexpected value for TEST_VAL`)

	_ = os.Unsetenv(`TEST_VAL`)
	SetDelimiter(`;`)
}
