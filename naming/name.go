package naming

import "strings"

var abbreviations = map[string]bool{
	"tcp":   true,
	"http":  true,
	"udp":   true,
	"id":    true,
	"ssl":   true,
	"tls":   true,
	"cpu":   true,
	"dob":   true,
	"ttl":   true,
	"sso":   true,
	"https": true,
	"ip":    true,
	"xss":   true,
	"os":    true,
	"sip":   true,
	"xml":   true,
	"json":  true,
	"html":  true,
	"xhtml": true,
	"xsl":   true,
	"xslt":  true,
	"yaml":  true,
	"toml":  true,
	"wlan":  true,
	"wifi":  true,
	"vm":    true,
	"jvm":   true,
	"ui":    true,
	"uri":   true,
	"url":   true,
	"sla":   true,
	"scp":   true,
	"smtp":  true,
	"soa":   true,
	"oa":    true,
	"svg":   true,
	"png":   true,
	"jpg":   true,
	"jpeg":  true,
	"pdf":   true,
	"io":    true,
}

// ToSnake converts s from CamelCase to Snake
func ToSnake(s string) string {
	snake := make([]rune, 0, len(s)+1)
	flag := false
	k := 'a' - 'A'
	for i, c := range s {
		if c >= 'A' && c <= 'Z' {
			if !flag {
				flag = true
				if i > 0 {
					snake = append(snake, '_')
				}
			}
			snake = append(snake, c+k)
		} else {
			flag = false
			snake = append(snake, c)
		}
	}
	return string(snake)
}

// ToCamel converts s from Snake to Camel
func ToCamel(s string) string {
	s = strings.ReplaceAll(s, ".", "_")
	s = strings.ReplaceAll(s, "-", "_")
	a := strings.Split(s, "_")
	for i := 1; i < len(a); i++ {
		if abbreviations[a[i]] {
			a[i] = strings.ToUpper(a[i])
		} else if a[i] != "" {
			a[i] = strings.ToUpper(a[i][0:1]) + a[i][1:]
		}
	}
	return strings.Join(a, "")
}

func ToClassName(s string) string {
	if s == "" {
		return ""
	}
	s = ToCamel(s)
	stop := len(s)
	for i, c := range s {
		if c >= 'A' && c <= 'Z' {
			stop = i
			break
		}
	}

	if abbreviations[s[:stop]] {
		s = strings.ToUpper(s[:stop]) + s[stop:]
	} else {
		s = strings.ToUpper(s[:1]) + s[1:]
	}
	return s
}

// Checker check if names can be convertiable
type Checker interface {
	Check(src, dst string) bool
}

// CheckerFunc defines func type which implements Checker
type CheckerFunc func(src string, dst string) bool

// Check checks if srcName can be converted to dstName
func (f CheckerFunc) Check(srcName, dstName string) bool {
	return f(srcName, dstName)
}

var DefaultChecker = CheckerFunc(Check)

// Check is the default Checker
func Check(a, b string) bool {
	if a == b {
		return true
	}

	la := strings.ToLower(a)
	lb := strings.ToLower(b)
	switch {
	case la == lb:
		return true
	case strings.ToLower(ToSnake(a)) == lb:
		return true
	case la == strings.ToLower(ToSnake(b)):
		return true
	case strings.ToLower(ToCamel(a)) == lb:
		return true
	case la == strings.ToLower(ToCamel(b)):
		return true
	default:
		return false
	}
}
