package env

import (
	"errors"
	"os"

	"github.com/joho/godotenv"
)

// LoadDotEnv loads environment variables from .env files.
//
// The order of the files is important; subsequent files will overload previously set variables.
// The default order is: .env, .env.local, .env.<env>, .env.<env>.local
func LoadDotEnv(baseDir string, files ...string) {
	SetEnvironment(false)
	env := GetEnv().String()

	files = append([]string{
		`.env`,
		`.env.local`,
		`.env.` + env,
		`.env.` + env + `.local`,
	}, files...)
	for _, file := range files {
		if err := godotenv.Overload(safeDir(baseDir) + file); err != nil && !errors.Is(err, os.ErrNotExist) {
			panic(`Error loading environment file(s):` + err.Error())
		}
	}
}
