package env

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestEnv(t *testing.T) {
	assert.Equal(t, ``, GetEnv().String(), `ENV should not be set; is `, GetEnv().String())

	SetEnv(true)
	assert.Equal(t, Development, GetEnv(), `ENV should be development`)

	SetEnv(true, Testing)
	assert.Equal(t, Testing, GetEnv(), `ENV should be testing`)

	SetEnv(true, `test`)
	assert.Equal(t, Testing, GetEnv(), `ENV should be testing`)

	if !IsTesting() {
		t.Error(`ENV should be testing`)
		t.Fail()
	}

	SetEnv(true, `devel`)
	assert.Equal(t, Development, GetEnv(), `ENV should be development`)

	if !IsDevelopment() {
		t.Error(`ENV should be development`)
		t.Fail()
	}

	if IsProduction() {
		t.Error(`ENV should not be production`)
		t.Fail()
	}

	SetEnv(true, `prod`)
	assert.Equal(t, Production, GetEnv(), `ENV should be production`)

	SetEnv(true, Staging)
	assert.Equal(t, Staging, GetEnv(), `ENV should be staging`)
	if !IsStaging() {
		t.Error(`ENV should be staging`)
		t.Fail()
	}

	SetEnv(false)
	if IsDevelopment() {
		t.Error(`ENV should not be overridden; must remain staging'`)
		t.Fail()
	}

	_ = os.Unsetenv(`ENV`)
}
