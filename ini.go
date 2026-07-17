package env

import (
	"errors"
	"regexp"
	"strings"

	"gopkg.in/ini.v1"
)

// loadIniData overloads files and returns *ini.File
func loadIniData(baseDir string, files []string) (*ini.File, error) {
	var fi = make([]interface{}, len(files))
	for i, file := range files {
		fi[i] = safeDir(baseDir) + file
	}

	data, err := ini.LooseLoad(fi[0], fi[1:]...) // LooseLoad ignores unknown files
	if err != nil {
		return nil, err
	}
	return data, nil
}

// mapIniToEnv sets environment variables from ini file
func mapIniToEnv(data *ini.File) {
	keyName := func(k string) string {
		re := regexp.MustCompile(`\W+`)
		return strings.ToUpper(re.ReplaceAllString(k, "_"))
	}

	var k string
	for _, section := range data.Sections() {
		for _, key := range section.Keys() {
			if section.Name() == ini.DefaultSection {
				k = keyName(key.Name())
			} else {
				k = keyName(section.Name() + `_` + key.Name())
			}
			Set(k, key.Value())
		}
	}
}

// LoadIni loads environment variables from ini files
//
// The order of the files is important; subsequent files will overload previously set variables.
// The default order is: env.ini, env.local.ini, env.<env>.ini, env.<env>.local.ini
func LoadIni(baseDir string, files ...string) {
	data, err := LoadIniFile(baseDir, files...)
	if err != nil {
		panic(`Error processing INI file: ` + err.Error())
	}
	mapIniToEnv(data)
}

// LoadIniFile loads variables from environment ini files
// returning *ini.File
//
// The order of the files is important; subsequent files will overload previously set variables.
// The default order is: env.ini, env.local.ini, env.<env>.ini, env.<env>.local.ini
func LoadIniFile(baseDir string, files ...string) (*ini.File, error) {
	SetEnv(false)
	env := GetEnv().String()

	files = append([]string{
		`env.ini`,
		`env.local.ini`,
		`env.` + env + `.ini`,
		`env.` + env + `.local.ini`,
	}, files...)

	return loadIniData(baseDir, files)
}

// MapIni maps environment ini files to user-defined interface
//
// The order of the files is important; subsequent files will overload previously set variables.
// The default order is: env.ini, env.local.ini, env.<env>.ini, env.<env>.local.ini
func MapIni(dest any, baseDir string, files ...string) error {
	var (
		data *ini.File
		err  error
	)

	data, err = LoadIniFile(baseDir, files...)
	if err != nil {
		return err
	}

	if err = data.MapTo(dest); err != nil {
		return errors.New(`Error mapping INI data to struct: ` + err.Error())
	}
	return nil
}
