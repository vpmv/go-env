package env

import "strings"

func safeDir(dir string) string {
	if dir != `` {
		dir = strings.TrimRight(dir, `/`) + `/`
	}
	return dir
}
