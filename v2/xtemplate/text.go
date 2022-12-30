package xtemplate

import "text/template"

var TextFuncMap = template.FuncMap{
	"plus":       Plus,
	"minus":      Minus,
	"multiple":   Multiple,
	"divide":     Divide,
	"join":       Join,
	"lower":      ToLower,
	"upper":      ToUpper,
	"concat":     Concat,
	"capitalize": Capitalize,
}
