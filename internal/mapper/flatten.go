package mapper

import (
	"reflect"
	"regexp"
	"strconv"
	"strings"
)

// Flatten flattens (nested) maps into a flat key-value array.
// Keys are dot-separated; lists are indexed using [n].
func Flatten(v map[string]any) map[string]any {
	m := make(map[string]any)
	flatten(``, v, m)
	return m
}

// flatten recursively flattens maps and appends flatmap
//
// map are flattened => "path.to.key = value"
// []map are scanned and flattened => "map[0].key = value"
// [][]any are appended => "key[0] => []values
//
// NOTE: do not change order
func flatten(prefix string, m map[string]any, flatmap map[string]any) {
	re := regexp.MustCompile(`\W+`)

	for key, value := range m {
		keyName := re.ReplaceAllString(key, "_")
		if prefix != "" {
			keyName = flattenKey(prefix, keyName)
		}

		switch {
		case isMap(value):
			flatten(keyName, value.(map[string]any), flatmap)
			continue
		case isMapSlice(value):
			for i, val := range value.([]any) {
				flatten(keyName+flattenIterator(i), val.(map[string]any), flatmap)
			}
			continue
		case isSlice(value):
			kind := reflect.ValueOf(value).Index(0).Elem().Kind()
			if kind == reflect.Slice || kind == reflect.Array {
				for i, val := range value.([]any) {
					flatmap[keyName+flattenIterator(i)] = val
				}
				continue
			}
		}

		flatmap[keyName] = value
	}
}

func flattenKey(parts ...string) string {
	return strings.Join(parts, ".")
}

func flattenIterator(i int) string {
	return `[` + strconv.Itoa(i) + `]`
}
