package mapper

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestFlatten(t *testing.T) {
	src := map[string]any{
		"server": map[string]any{
			"host": "localhost",
			"port": 8080,
		},
		"databases": []any{
			map[string]any{
				"host": "mysql",
				"port": 3306,
			},
			map[string]any{
				"host": "postgres",
				"port": 5432,
			},
		},
		"debug": true,
		"cats": []any{
			"abyssinian",
			"bombay",
			"polydactyl",
		},
		"cool-animals": []any{
			[]any{
				"husky",
				"poodle",
				`schnauzer`,
			},
			[]any{
				"bear",
				"lynx",
				"wolf",
			},
			[]any{
				"elephant",
				"hippo",
				"zebra",
			},
		},
	}

	out := Flatten(src)

	assert.Len(t, out, 11)
	assert.Equal(t, "localhost", out["server.host"])
	assert.Equal(t, 8080, out["server.port"])
	assert.Equal(t, `mysql`, out["databases[0].host"])
	assert.Equal(t, 3306, out["databases[0].port"])
	assert.Equal(t, `postgres`, out["databases[1].host"])
	assert.Equal(t, 5432, out["databases[1].port"])
	assert.Equal(t, true, out["debug"])
	assert.Len(t, out["cats"], 3)
	assert.Len(t, out["cool_animals[0]"], 3)
}
