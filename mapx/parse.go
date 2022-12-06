package mapx

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
