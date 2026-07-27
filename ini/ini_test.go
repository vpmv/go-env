package ini

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/vpmv/go-env"
)

type testIniAppConfig struct {
	App struct {
		Value1 string `ini:"value1"`
		Value2 string `ini:"value2"`
	} `ini:"app"`
	AppSection struct {
		Value1 string `ini:"value1"`
		Value2 string `ini:"value2"`
	} `ini:"app.section"`
	Database struct {
		Host     string `ini:"host"`
		Port     int    `ini:"port"`
		User     string `ini:"user"`
		Password string `ini:"password"`
	} `ini:"database"`
}

func TestLoadIni(t *testing.T) {
	// should set env to "development" and only load "env.ini"
	Load(`../fixtures`)

	// no section
	assert.Equal(t, `earth`, os.Getenv(`GLOBAL_VALUE1`), `GLOBAL_VALUE1 should be "earth"`)

	// [app] section
	assert.Equal(t, `foo`, os.Getenv(`APP_VALUE1`), `APP_VALUE1 should be "foo"`)
	assert.Equal(t, `bar`, os.Getenv(`APP_VALUE2`), `APP_VALUE2 should be "bar"`)

	assert.Equal(t, `oof`, os.Getenv(`APP_SECTION_VALUE1`), `APP_SECTION_VALUE1 should be "oof"`)
	assert.Equal(t, `rab`, os.Getenv(`APP_SECTION_VALUE2`), `APP_SECTION_VALUE2 should be "rab"`)

	// [database] section
	assert.Equal(t, `3306`, os.Getenv(`DATABASE_PORT`), `DATABASE_PORT should be "3306"`)

	unsetEnvVars()

	env.SetEnv(true, env.Testing)
	// should set env to "testing" and load "env.ini" & "env.testing.ini"
	Load(`../fixtures`)
	assert.Equal(t, `foobar`, os.Getenv(`APP_SECTION_VALUE1`), `APP_SECTION_VALUE1 should be "foobar"`)
	assert.Equal(t, `user2`, os.Getenv(`DATABASE_USER`), `DATABASE_USER should be "user2"`)
	assert.Equal(t, `verysecret`, os.Getenv(`DATABASE_PASSWORD`), `DATABASE_PASSWORD should be "verysecret"`)

	unsetEnvVars()
}

func TestLoadIniPanic(t *testing.T) {
	assert.Panics(t, func() {
		Load(`../fixtures`, `illegal_character.ini`)
	})
	unsetEnvVars()
}

func TestMapIni(t *testing.T) {
	config := new(testIniAppConfig)
	if err := Map(config, `../fixtures`); err != nil {
		t.Fatal(err)
	}

	// [app] section
	assert.Equal(t, `foo`, config.App.Value1, `App.Value1 should be "foo"`)
	assert.Equal(t, `bar`, config.App.Value2, `App.Value2 should be "bar"`)
	// [app.section] section
	assert.Equal(t, `oof`, config.AppSection.Value1, `AppSection.Value1 should be "oof"`)
	assert.Equal(t, `rab`, config.AppSection.Value2, `AppSection.Value2 should be "rab"`)
	// [database] section
	assert.Equal(t, `localhost`, config.Database.Host, `Database.Host should be "localhost"`)
	assert.Equal(t, 3306, config.Database.Port, `Database.Port should be int:3306`)
	assert.Equal(t, `user1`, config.Database.User, `Database.User should be "user"`)
	assert.Equal(t, `secret`, config.Database.Password, `Database.Password should be "secret"`)

	env.SetEnv(true, env.Testing)

	config = new(testIniAppConfig)
	if err := Map(config, `../fixtures`); err != nil {
		t.Fatal(err)
	}

	assert.Equal(t, `foo`, config.App.Value1, `App.Value1 should be "foo"`)
	assert.Equal(t, `user2`, config.Database.User, `Database.User should be "user2"`)
	assert.Equal(t, `verysecret`, config.Database.Password, `Database.Password should be "verysecret"`)

	unsetEnvVars()
}

func TestMapIniError(t *testing.T) {
	config := new(testIniAppConfig)
	err := Map(config, `../fixtures`, `illegal_character.ini`)
	if err == nil {
		assert.Fail(t, `MapIni should return error`)
	}

	illegalConfig := []string{}
	err = Map(illegalConfig, `../fixtures`)
	if err == nil {
		assert.Fail(t, `MapIni should return error`)
	}

	unsetEnvVars()
}

func TestEmptyIni(t *testing.T) {
	Load(`non_existent`)
	assert.False(t, env.Has(`APP_VALUE1`), `environment should not exist`)
	unsetEnvVars()
}

func unsetEnvVars() {
	_ = os.Unsetenv(`ENV`)
	_ = os.Unsetenv(`VALUE1`)
	_ = os.Unsetenv(`VALUE2`)
	_ = os.Unsetenv(`SECTION1_VALUE`)
	_ = os.Unsetenv(`SECTION2_VALUE`)

	_ = os.Unsetenv(`GLOBAL_VALUE1`)
	_ = os.Unsetenv(`APP_VALUE1`)
	_ = os.Unsetenv(`APP_VALUE2`)
	_ = os.Unsetenv(`APP_SECTION_VALUE1`)
	_ = os.Unsetenv(`APP_SECTION_VALUE2`)
	_ = os.Unsetenv(`DATABASE_HOST`)
	_ = os.Unsetenv(`DATABASE_PORT`)
	_ = os.Unsetenv(`DATABASE_USER`)
	_ = os.Unsetenv(`DATABASE_PASSWORD`)
}
