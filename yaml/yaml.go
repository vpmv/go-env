package yaml

import (
	"os"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/goccy/go-yaml"
	"github.com/vpmv/go-env"
	"github.com/vpmv/go-env/internal/format"
	"github.com/vpmv/go-env/internal/fs"
	"github.com/vpmv/go-env/internal/mapper"
)

// Load loads environment variables from YAML files
//
// The order of the files is important; subsequent files will overload previously set variables.
// The default order is: env.yaml, env.local.yaml, env.<env>.yaml, env.<env>.local.yaml
//
// All files are expected to be relative to baseDir
func Load(baseDir string, files ...string) error {
	return LoadWithReferences(baseDir, []string{}, files...)
}

// LoadWithReferences loads environment variables from YAML files
// Reference files are used to parse YAML anchors within the source files.
// Reference files should be in order of precedence,
// i.e. anchors must be defined before they're referenced
//
// The order of the files is important; subsequent files will overload previously set variables.
// The default order is: env.yaml, env.local.yaml, env.<env>.yaml, env.<env>.local.yaml
//
// All files are expected to be relative to baseDir
func LoadWithReferences(baseDir string, referenceFiles []string, files ...string) error {
	keyName := func(k string) string {
		re := regexp.MustCompile(`\.+`)
		return strings.ToUpper(re.ReplaceAllString(k, "_"))
	}

	env.SetEnv(false)
	e := env.GetEnv().String()

	files = append([]string{
		`env.yaml`,
		`env.local.yaml`,
		`env.` + e + `.yaml`,
		`env.` + e + `.local.yaml`,
	}, files...)

	dest := map[string]interface{}{}
	err := mapYAML(&dest, baseDir, files, referenceFiles...)
	if err != nil {
		return err
	}

	flatmap := mapper.Flatten(dest)
	for k, v := range flatmap {
		// unless goccy/go-yaml changes type assertion,
		// mapper.Flatten filters unsupported types and format.ToString returns silent
		str, _ := format.Stringer.ToString(v)
		env.Set(keyName(k), str)
	}

	return nil
}

// Map maps YAML environment files to interface
//
// The order of files is important; subsequent files will overload previously set variables.
// The default order is: env.yaml, env.local.yaml, env.<env>.yaml, env.<env>.local.yaml
//
// All files are expected to be relative to baseDir
func Map(dest any, baseDir string, files ...string) error {
	return MapWithReferences(dest, baseDir, files)
}

// MapWithReferences maps YAML environment files to interface.
// Reference files are used to parse YAML anchors within the source files.
// Reference files should be in order of precedence,
// i.e. anchors must be defined before they're referenced
//
// The order of files is important; subsequent files will overload previously set variables.
// The default order is: env.yaml, env.local.yaml, env.<env>.yaml, env.<env>.local.yaml
//
// All files are expected to be relative to baseDir
func MapWithReferences(dest any, baseDir string, referenceFiles []string, files ...string) error {
	env.SetEnv(false)
	e := env.GetEnv().String()

	files = append([]string{
		`env.yaml`,
		`env.local.yaml`,
		`env.` + e + `.yaml`,
		`env.` + e + `.local.yaml`,
	}, files...)

	return mapYAML(dest, baseDir, files, referenceFiles...)
}

// mapYAML maps YAML file(s) contents to interface.
// Reference files are used to parse YAML anchors within the source files.
// Reference files should be in order of precedence,
// i.e. anchors must be defined before they're referenced
//
// All files are expected to be relative to baseDir.
func mapYAML(dest any, baseDir string, files []string, referenceFiles ...string) error {
	for i, file := range referenceFiles {
		if !strings.HasPrefix(file, baseDir) {
			referenceFiles[i] = filepath.Join(baseDir, file)
		}
	}

	for _, file := range files {
		f, err := fs.FS().Open(filepath.Join(baseDir, file))
		defer fs.FS().CloseFile(f)

		if err != nil && os.IsNotExist(err) {
			continue
		} else if err != nil {
			return err
		}

		decoder := yaml.NewDecoder(f, yaml.ReferenceFiles(referenceFiles...), yaml.AllowDuplicateMapKey())
		if err := decoder.Decode(dest); err != nil {
			return err
		}
	}
	return nil
}
