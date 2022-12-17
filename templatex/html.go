package templatex

import "html/template"

var HTMLFuncMap = template.FuncMap{
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

type HTMLMeta struct {
	Name    string
	Content string
}

type HTMLHeader struct {
	Title          string
	Meta           []*HTMLMeta
	CSSLinks       []string
	JSLinks        []string
	BodyAttributes string
}

type HTMLFooter struct {
	JSLinks []string
}
