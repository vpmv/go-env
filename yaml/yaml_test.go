package yaml

import (
	"io/fs"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/vpmv/go-env"
	fsys "github.com/vpmv/go-env/internal/fs"
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

func TestMapYAML(t *testing.T) {
	s := new(testYamlAppConfig)
	err := Map(s, `../fixtures`)
	assert.NoError(t, err)

	assert.Equal(t, `earth`, s.GlobalValue)
	assert.Equal(t, `foo`, s.App.Value1)
	assert.Equal(t, true, s.App.Value2)
	assert.Equal(t, []string{`emmentaler`, `gouda`, `roquefort`}, s.App.Value3)
	assert.Equal(t, int64(123), s.App.Section.Value2)

	env.SetEnv(true, env.Testing)

	s = new(testYamlAppConfig)
	err = Map(s, `../fixtures`)
	assert.NoError(t, err)

	assert.Equal(t, `world`, s.GlobalValue)
	assert.Equal(t, `foo`, s.App.Value1)
	assert.Equal(t, false, s.App.Value2)
	assert.Equal(t, []string{`human`, `dog`, `cat`}, s.App.Value3)
	assert.Equal(t, int64(321), s.App.Section.Value2)

	_ = os.Unsetenv(`ENV`)
}

func TestMapYAMLWithReferences(t *testing.T) {
	s := new(testYamlAppConfig)
	err := MapWithReferences(s, `../fixtures`, []string{`defaults.yaml`})
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

	env.SetEnv(true, env.Staging)

	s = new(testYamlAppConfig)
	err = MapWithReferences(s, `../fixtures`, []string{`defaults.yaml`})
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
	err := mapYAML(s, `../fixtures`, []string{`illegal_character.yaml`})
	assert.Error(t, err)

	fsys.Set(new(testFS))
	err = mapYAML(s, `../fixtures`, []string{`void.yaml`})
	assert.Error(t, err)

	fsys.Reset()
}

func TestLoadYAML(t *testing.T) {
	err := Load(`../fixtures`)
	assert.NoError(t, err)

	assert.Equal(t, `earth`, os.Getenv(`GLOBAL_VALUE1`))
	assert.Equal(t, `foo`, os.Getenv(`APP_VALUE1`))
	assert.Equal(t, `emmentaler;gouda;roquefort`, os.Getenv(`APP_VALUE3`))
	assert.Equal(t, `oof`, os.Getenv(`APP_SECTION_VALUE1`))
	assert.Equal(t, `localhost`, os.Getenv(`DATABASE_HOST`))

	assert.Equal(t, `cheese;cake;gherkin`, os.Getenv(`LIST[0]`))
	assert.Equal(t, `bread;wine;prayer`, os.Getenv(`LIST[1]`))

	assert.Equal(t, []string{`bread`, `wine`, `prayer`}, env.MustStringSlice(`LIST[1]`))

	assert.Equal(t, `postgres`, os.Getenv(`LISTMAP[1]_HOST`))
	assert.Equal(t, `5432`, os.Getenv(`LISTMAP[1]_PORT`))

	// reset env
	_ = os.Unsetenv(`ENV`)
	_ = os.Unsetenv(`GLOBAL_VALUE1`)
	_ = os.Unsetenv(`APP_VALUE1`)
	_ = os.Unsetenv(`APP_VALUE2`)
	_ = os.Unsetenv(`APP_VALUE3`)
	_ = os.Unsetenv(`APP_SECTION_VALUE1`)
	_ = os.Unsetenv(`APP_SECTION_VALUE2`)
	_ = os.Unsetenv(`DATABASE_HOST`)
	_ = os.Unsetenv(`DATABASE_PORT`)
	_ = os.Unsetenv(`DATABASE_USER`)
	_ = os.Unsetenv(`DATABASE_PASSWORD`)
	_ = os.Unsetenv(`LIST[0]`)
	_ = os.Unsetenv(`LIST[1]`)
	_ = os.Unsetenv(`LISTMAP[0]_HOST`)
	_ = os.Unsetenv(`LISTMAP[0]_PORT`)
	_ = os.Unsetenv(`LISTMAP[1]_HOST`)
	_ = os.Unsetenv(`LISTMAP[1]_PORT`)
}

func TestLoadYAML_Errors(t *testing.T) {
	err := Load(`../fixtures`, `illegal_character.yaml`)
	assert.Error(t, err)
}
