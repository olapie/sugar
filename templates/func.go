package templates

import "strings"

func Plus(a, b int) int {
	return a + b
}

func Minus(a, b int) int {
	return a - b
}

func Multiple(a, b int) int {
	return a * b
}

func Divide(a, b int) int {
	return a / b
}

func Join(strs []string, sep string) string {
	return strings.Join(strs, sep)
}

func Concat(strs ...string) string {
	return strings.Join(strs, "")
}

func ToLower(s string) string {
	return strings.ToLower(s)
}

func ToUpper(s string) string {
	return strings.ToUpper(s)
}

func Capitalize(s string, n int) string {
	left := s[:n]
	right := s[n:]
	return strings.ToUpper(left) + right
}
