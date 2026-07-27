package dotenv

import (
	"errors"
	"os"
	"path/filepath"

	"github.com/joho/godotenv"
	"github.com/vpmv/go-env"
)

// Load loads environment variables from .env files.
//
// The order of the files is important; subsequent files will overload previously set variables.
// The default order is: .env, .env.local, .env.<env>, .env.<env>.local
func Load(baseDir string, files ...string) {
	env.SetEnv(false)
	e := env.GetEnv().String()

	files = append([]string{
		`.env`,
		`.env.local`,
		`.env.` + e,
		`.env.` + e + `.local`,
	}, files...)
	for _, file := range files {
		if err := godotenv.Overload(filepath.Join(baseDir, file)); err != nil && !errors.Is(err, os.ErrNotExist) {
			panic(`Error loading environment file(s):` + err.Error())
		}
	}
}
