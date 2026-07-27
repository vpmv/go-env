package dotenv

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/vpmv/go-env"
)

func TestLoadDotEnv(t *testing.T) {
	// should set env to "development" and only load ".env"
	Load(`../fixtures`)
	assert.Equal(t, `foo`, os.Getenv(`VALUE1`), `VALUE1 should be "foo"`)
	assert.Equal(t, `world`, os.Getenv(`SECTION2_VALUE`), `SECTION2_VALUE should be "world"`)

	unsetEnvVars()

	env.SetEnv(true, env.Testing)
	// should set env to "testing" and load ".env" & ".env.testing"
	Load(`../fixtures`)
	assert.Equal(t, `dlrow`, os.Getenv(`SECTION2_VALUE`), `SECTION2_VALUE should be "dlrow"`)

	unsetEnvVars()
}

func TestLoadDotEnvPanic(t *testing.T) {
	assert.Panics(t, func() {
		Load(`../fixtures`, `env.ini`)
	})

	unsetEnvVars()
}

func unsetEnvVars() {
	_ = os.Unsetenv(`ENV`)
	_ = os.Unsetenv(`VALUE1`)
	_ = os.Unsetenv(`VALUE2`)
	_ = os.Unsetenv(`SECTION1_VALUE`)
	_ = os.Unsetenv(`SECTION2_VALUE`)
}
