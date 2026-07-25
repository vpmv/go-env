package env

import (
	"io/fs"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

type testFS struct{}

func (t testFS) Open(name string) (fs.File, error) {
	return nil, fs.ErrPermission
}

func (t testFS) CloseFile(file fs.File) {
}

type testYamlConnection struct {
	Host string `yaml:"host"`
	Port int    `yaml:"port"`
	SSL  bool   `yaml:"ssl"`
}

type testYamlDatabaseConfig struct {
	Host     string `yaml:"host"`
	Port     int    `yaml:"port"`
	User     string `yaml:"user"`
	Password string `yaml:"password"`
}

type testYamlAppConfig struct {
	App struct {
		testYamlConnection
		Value1  string   `yaml:"value1"`
		Value2  bool     `yaml:"value2"`
		Value3  []string `yaml:"value3"`
		Section struct {
			Value1 string `yaml:"value1"`
			Value2 int64  `yaml:"value2"`
		} `yaml:"section"`
	} `yaml:"app"`
	Database    *testYamlDatabaseConfig `yaml:"database"`
	GlobalValue string                  `yaml:"global_value1"`
}

func TestMapEnvYAML(t *testing.T) {
	s := new(testYamlAppConfig)
	err := MapEnvYAML(s, `fixtures`)
	assert.NoError(t, err)

	assert.Equal(t, `earth`, s.GlobalValue)
	assert.Equal(t, `foo`, s.App.Value1)
	assert.Equal(t, true, s.App.Value2)
	assert.Equal(t, []string{`emmentaler`, `gouda`, `roquefort`}, s.App.Value3)
	assert.Equal(t, int64(123), s.App.Section.Value2)

	SetEnv(true, Testing)

	s = new(testYamlAppConfig)
	err = MapEnvYAML(s, `fixtures`)
	assert.NoError(t, err)

	assert.Equal(t, `world`, s.GlobalValue)
	assert.Equal(t, `foo`, s.App.Value1)
	assert.Equal(t, false, s.App.Value2)
	assert.Equal(t, []string{`human`, `dog`, `cat`}, s.App.Value3)
	assert.Equal(t, int64(321), s.App.Section.Value2)

	_ = os.Unsetenv(`ENV`)
}

func TestMapEnvYAMLWithReferences(t *testing.T) {
	s := new(testYamlAppConfig)
	err := MapEnvYAMLWithReferences(s, `fixtures`, []string{`defaults.yaml`})
	assert.NoError(t, err)

	assert.Equal(t, `earth`, s.GlobalValue)
	assert.Equal(t, `foo`, s.App.Value1)
	assert.Equal(t, true, s.App.Value2)
	assert.Equal(t, []string{`emmentaler`, `gouda`, `roquefort`}, s.App.Value3)
	assert.Equal(t, int64(123), s.App.Section.Value2)

	assert.Equal(t, `localhost`, s.Database.Host)
	assert.Equal(t, 3306, s.Database.Port)

	// test unreferenced Host config
	assert.Equal(t, ``, s.App.Host)
	assert.Equal(t, false, s.App.SSL)

	SetEnv(true, Staging)

	s = new(testYamlAppConfig)
	err = MapEnvYAMLWithReferences(s, `fixtures`, []string{`defaults.yaml`})
	assert.NoError(t, err)

	assert.Equal(t, `world`, s.GlobalValue)
	assert.Equal(t, `foo`, s.App.Value1)
	assert.Equal(t, true, s.App.Value2)
	assert.Equal(t, []string{`bear`, `lynx`, `wolf`}, s.App.Value3)
	assert.Equal(t, int64(123), s.App.Section.Value2)

	assert.Equal(t, `postgres`, s.Database.Host)
	assert.Equal(t, 5432, s.Database.Port)

	_ = os.Unsetenv(`ENV`)
}

func TestMapYAML_Errors(t *testing.T) {
	s := new(testYamlAppConfig)
	err := MapYAML(s, `fixtures`, []string{`illegal_character.yaml`})
	assert.Error(t, err)

	fileSystem = new(testFS)
	err = MapYAML(s, `fixtures`, []string{`void.yaml`})
	assert.Error(t, err)
}
