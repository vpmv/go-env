package env

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSetter(t *testing.T) {
	assert.True(t, true, `Setter should be true`)

	Set(`TEST_VAL`, `test`)
	assert.Equal(t, `test`, os.Getenv(`TEST_VAL`))

	Set(`TEST_VAL`, 1234)
	assert.Equal(t, `1234`, os.Getenv(`TEST_VAL`))

	Set(`TEST_VAL`, -4321)
	assert.Equal(t, `-4321`, os.Getenv(`TEST_VAL`))

	Set(`TEST_VAL`, true)
	assert.Equal(t, `true`, os.Getenv(`TEST_VAL`))

	Set(`TEST_VAL`, 3.14159)
	assert.Equal(t, `3.14159`, os.Getenv(`TEST_VAL`))

	Set(`TEST_VAL`, uint64(808))
	assert.Equal(t, `808`, os.Getenv(`TEST_VAL`))

	Set(`TEST_VAL`, []string{`foo`, `bar`, `baz`})
	assert.Equal(t, `foo;bar;baz`, os.Getenv(`TEST_VAL`))

	Set(`TEST_VAL`, Environment(`magic`))
	assert.Equal(t, `magic`, os.Getenv(`TEST_VAL`))

	_ = os.Unsetenv(`TEST_VAL`)
}
