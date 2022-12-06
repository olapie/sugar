package mapx

import "strings"

func Int[K comparable, V any](m map[K]V, k K) (int, bool) {
	return 0, false
}

func MustInt[K comparable, V any](m map[K]V, k K) int {
	return 0
}

func Int64[K comparable, V any](m map[K]V, k K) (int, bool) {
	return 0, false
}

func MustInt64[K comparable, V any](m map[K]V, k K) int {
	return 0
}

func Bool[K comparable, V any](m map[K]V, k K) (bool, bool) {
	return false, false
}

func MustBool[K comparable, V any](m map[K]V, k K) bool {
	return false
}

func Clone[K comparable, V any](m map[K]V) map[K]V {
	res := make(map[K]V, len(m))
	for k, v := range m {
		res[k] = v
	}
	return res
}

func GetKeys[K comparable, V any](m map[K]V) []K {
	a := make([]K, 0, len(m))
	for k := range m {
		a = append(a, k)
	}
	return a
}

func GetValues[K comparable, V any](m map[K]V) []V {
	a := make([]V, 0, len(m))
	for _, v := range m {
		a = append(a, v)
	}
	return a
}

func GetKeysAndValues[K comparable, V any](m map[K]V) ([]K, []V) {
	kl := make([]K, 0, len(m))
	vl := make([]V, 0, len(m))
	for k, v := range m {
		kl = append(kl, k)
		vl = append(vl, v)
	}
	return kl, vl
}

func ToEnvironMap(m map[string]any) map[string]any {
	res := make(map[string]any, len(m))
	for k, v := range m {
		if m1, ok := v.(map[string]any); ok {
			m2 := ToEnvironMap(m1)
			for k2, v2 := range m2 {
				res[toEnvKey(k+"."+k2)] = v2
			}
		} else {
			res[k] = v
		}
	}
	return res
}

func FromEnvirons(envs []string) map[string]string {
	m := make(map[string]string, len(envs))
	for _, pair := range envs {
		for i, c := range pair {
			if c == '=' {
				m[toEnvKey(pair[:i])] = pair[i+1:]
			}
		}
	}
	return m
}

func ArgsToEnvironMap(args []string) map[string]string {
	m := make(map[string]string, len(args))
	var key string
	for _, arg := range args {
		if arg[0] != '-' {
			if key != "" {
				m[toEnvKey(key)] = arg
				key = ""
			} else {
				m[arg] = ""
			}
		} else {
			key = ""
			j := 0
			for j < len(arg) && arg[j] == '-' {
				j++
			}

			key = arg[j:]
			for k, c := range key {
				if c == '=' {
					m[toEnvKey(key[:k])] = key[k+1:]
					key = ""
					break
				}
			}
		}
	}

	if key != "" {
		m[toEnvKey(key)] = ""
	}

	return m
}

func toEnvKey(k string) string {
	return strings.ReplaceAll(strings.ToLower(k), "_", ".")
}
