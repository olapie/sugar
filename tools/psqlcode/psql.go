package main

import (
	"embed"
	"text/template"

	"code.olapie.com/sugar/v2/naming"
)

//go:embed template
var tplFS embed.FS
var globalTemplate = template.New("")

func init() {
	globalTemplate = globalTemplate.Funcs(template.FuncMap{
		"toStructName": naming.ToClassName,
		"toCamel":      naming.ToCamel,
		"toSnake":      naming.ToSnake,
		"toEntityName": func(s string) string {
			return naming.ToClassName(s) + "Entity"
		},
		"toBuilderName": func(s string) string {
			return naming.ToCamel(s) + "EntityBuilder"
		},
		"toModifierName": func(s string) string {
			return naming.ToCamel(s) + "EntityModifier"
		},
		"first": func(s string) string {
			return s[:1]
		},
	})
	globalTemplate = template.Must(globalTemplate.ParseFS(tplFS, "template/*.tpl"))
}
