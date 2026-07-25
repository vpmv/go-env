package env

import (
	"os"
	"path/filepath"
	"strings"

	"github.com/goccy/go-yaml"
)

// MapYAML maps YAML environment files to interface
//
// The order of files is important; subsequent files will overload previously set variables.
// The default order is: env.yaml, env.local.yaml, env.<env>.yaml, env.<env>.local.yaml
//
// All files are expected to be relative to baseDir
func MapEnvYAML(dest any, baseDir string, files ...string) error {
	return MapEnvYAMLWithReferences(dest, baseDir, files)
}

// MapYAML maps YAML environment files to interface.
// Reference files are used to parse YAML anchors within the source files.
// Reference files should be in order of precedence,
// i.e. anchors must be defined before they're referenced
//
// The order of files is important; subsequent files will overload previously set variables.
// The default order is: env.yaml, env.local.yaml, env.<env>.yaml, env.<env>.local.yaml
//
// All files are expected to be relative to baseDir
func MapEnvYAMLWithReferences(dest any, baseDir string, referenceFiles []string, files ...string) error {
	SetEnv(false)
	env := GetEnv().String()

	files = append([]string{
		`env.yaml`,
		`env.local.yaml`,
		`env.` + env + `.yaml`,
		`env.` + env + `.local.yaml`,
	}, files...)

	return MapYAML(dest, baseDir, files, referenceFiles...)
}

// MapYAML maps YAML file(s) contents to interface.
// Reference files are used to parse YAML anchors within the source files.
// Reference files should be in order of precedence,
// i.e. anchors must be defined before they're referenced
//
// All files are expected to be relative to baseDir
func MapYAML(dest any, baseDir string, files []string, referenceFiles ...string) error {
	for i, file := range referenceFiles {
		if !strings.HasPrefix(file, baseDir) {
			referenceFiles[i] = filepath.Join(baseDir, file)
		}
	}

	for _, file := range files {
		f, err := fileSystem.Open(filepath.Join(baseDir, file))
		defer fileSystem.CloseFile(f)

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
