package postgresx

import (
	"bytes"
	"fmt"
	"strings"
)

type compositeScanState int

const (
	compositeScanInit compositeScanState = iota
	compositeScanField
	compositeScanQuoted
)

func ToCompositeString(fields ...any) string {
	var builder strings.Builder
	builder.WriteString("(")
	n := len(fields)
	for i, field := range fields {
		switch v := field.(type) {
		case string:
			builder.WriteString(escapeCompositeField(v))
		case []byte:
			builder.WriteString(escapeCompositeField(string(v)))
		default:
			builder.WriteString(escapeCompositeField(fmt.Sprint(v)))
		}
		if i < n-1 {
			builder.WriteRune(',')
		}
	}
	builder.WriteString(")")
	return builder.String()
}

func ParseCompositeFields(column string) ([]string, error) {
	if len(column) == 0 {
		return nil, fmt.Errorf("empty column")
	}

	fields := make([]string, 0, 2)
	state := compositeScanInit
	var field bytes.Buffer
	chars := []rune(column)
	n := len(chars)
	errPos := -1
Loop:
	for i := 0; i < n; i++ {
		c := chars[i]
		switch state {
		case compositeScanInit:
			if c != '(' {
				//errPos = i
				//break Loop
				continue
			}
			state = compositeScanField
		case compositeScanField:
			switch c {
			case '"':
				if field.Len() == 0 {
					state = compositeScanQuoted
				} else {
					if i == len(chars)-1 || chars[i+1] != '"' {
						errPos = i
						break Loop
					}
					field.WriteRune('"')
					i++
				}
			case ')':
				fields = append(fields, field.String())
				if i != len(chars)-1 {
					errPos = i
					break Loop
				}
				return fields, nil
			case ',':
				fields = append(fields, field.String())
				field.Reset()
			default:
				field.WriteRune(c)
			}
		case compositeScanQuoted:
			switch c {
			case '"':
				if i == len(chars)-1 {
					errPos = i
					break Loop
				}
				i++
				switch chars[i] {
				case '"':
					// In quoted string, "" represents "
					field.WriteRune('"')
				case ',':
					fields = append(fields, field.String())
					field.Reset()
					state = compositeScanField
				case ')':
					fields = append(fields, field.String())
					if i != len(chars)-1 {
						errPos = i
						break Loop
					}
					return fields, nil
				default:
					errPos = i
					break Loop
				}
			default:
				field.WriteRune(c)
			}
		}
	}
	return nil, fmt.Errorf("syntax error at %d", errPos)
}

func MakeSQLPlaceholder(prefix string, num int) string {
	var b strings.Builder
	for i := 1; i <= num; i++ {
		if i > 1 {
			b.WriteString(",")
		}
		b.WriteString(fmt.Sprintf("%s%d", prefix, i))
	}
	return b.String()
}

func MakePlaceholder(num int) string {
	return MakeSQLPlaceholder("$", num)
}

func escapeCompositeField(s string) string {
	s = strings.Replace(s, ",", "\\,", -1)
	s = strings.Replace(s, "(", "\\(", -1)
	s = strings.Replace(s, ")", "\\)", -1)
	s = strings.Replace(s, "\"", "\\\"", -1)
	return s
}
