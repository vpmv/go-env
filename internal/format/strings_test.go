package format

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

type testType string

func TestStringFormatter_ToString(t *testing.T) {
	var (
		str string
		err error
	)

	format := NewStringFormatter()

	str, err = format.ToString("localhost")
	assert.NoError(t, err)
	assert.Equal(t, `localhost`, str)

	str, err = format.ToString(1234)
	assert.NoError(t, err)
	assert.Equal(t, `1234`, str)

	str, err = format.ToString(uint(1234))
	assert.NoError(t, err)
	assert.Equal(t, `1234`, str)

	str, err = format.ToString(-4321)
	assert.NoError(t, err)
	assert.Equal(t, `-4321`, str)

	str, err = format.ToString(true)
	assert.NoError(t, err)
	assert.Equal(t, `true`, str)

	str, err = format.ToString(3.14159)
	assert.NoError(t, err)
	assert.Equal(t, `3.14159`, str)

	str, err = format.ToString(uint64(808))
	assert.NoError(t, err)
	assert.Equal(t, `808`, str)

	str, err = format.ToString([]string{`foo`, `bar`, `baz`})
	assert.NoError(t, err)
	assert.Equal(t, `foo;bar;baz`, str)

	str, err = format.ToString([]interface{}{`foo`, `bar`, `baz`})
	assert.NoError(t, err)
	assert.Equal(t, `foo;bar;baz`, str)

	str, err = format.ToString([]int{1, 2, 3})
	assert.NoError(t, err)
	assert.Equal(t, `1;2;3`, str)

	str, err = format.ToString(testType(`magic`))
	assert.NoError(t, err)
	assert.Equal(t, `magic`, str)

	str, err = format.ToString(map[string]string{"foo": "bar", "baz": "qux"})
	assert.Error(t, err)
	assert.Equal(t, ``, str)

}

func TestStringFormatter_Split(t *testing.T) {
	format := NewStringFormatter()
	arr := format.Split(`foo;bar;baz`)
	assert.Equal(t, []string{`foo`, `bar`, `baz`}, arr)

	format = NewStringFormatter(WithDelimiter(`,`))
	arr = format.Split(`foo;bar;baz`)
	assert.Equal(t, []string{`foo;bar;baz`}, arr)

	arr = format.Split(`foo,bar,baz`)
	assert.Equal(t, []string{`foo`, `bar`, `baz`}, arr)
}
