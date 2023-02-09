package maps_test

import (
	"testing"

	"code.olapie.com/sugar/v2/testutil"
)

func TestToEnvMap(t *testing.T) {
	m := map[string]any{
		"k1": "v1",
		"k2": map[string]any{
			"k21": "v21",
			"k22": 22,
			"k23": map[string]any{
				"k231": "v231",
				"k232": 232,
			},
		},
	}

	m1 := maputil.ToEnvironMap(m)
	m2 := map[string]any{"k1": "v1", "k2.k21": "v21", "k2.k22": 22, "k2.k23.k231": "v231", "k2.k23.k232": 232}
	testutil.Equal(t, m2, m1)
}

func TestOSEnvsToEnvMap(t *testing.T) {
	envs := []string{
		"debug=true",
		"test=1",
		"db_password=123",
		"db_user=user",
		"DB_PASS=456",
		"DB_URL=localhost:4436",
	}

	m := maputil.FromEnvirons(envs)
	expected := map[string]string{
		"debug":       "true",
		"test":        "1",
		"db.password": "123",
		"db.user":     "user",
		"db.pass":     "456",
		"db.url":      "localhost:4436",
	}

	testutil.Equal(t, expected, m)
}

func TestOSArgsToEnvMap(t *testing.T) {
	args := []string{
		"enabled",
		"test=1",
		"--db_password=123",
		"-db_user=user",
		"--DB_PASS",
		"456",
		"-DB_URL",
		"localhost:4436",
		"-flag",
	}

	m := maputil.ArgsToEnvironMap(args)
	expected := map[string]string{
		"enabled":     "",
		"test=1":      "",
		"db.password": "123",
		"db.user":     "user",
		"db.pass":     "456",
		"db.url":      "localhost:4436",
		"flag":        "",
	}

	testutil.Equal(t, expected, m)
}
