package main

import (
	"embed"
	"text/template"

	"code.olapie.com/sugar/v2/xname"
)

//go:embed template
var tplFS embed.FS
var globalTemplate = template.New("")

func init() {
	globalTemplate = globalTemplate.Funcs(template.FuncMap{
		"toStructName": xname.ToClassName,
		"toCamel":      xname.ToCamel,
		"toSnake":      xname.ToSnake,
		"toEntityName": func(s string) string {
			return xname.ToClassName(s) + "Entity"
		},
		"toBuilderName": func(s string) string {
			return xname.ToCamel(s) + "EntityBuilder"
		},
		"toModifierName": func(s string) string {
			return xname.ToCamel(s) + "EntityModifier"
		},
		"first": func(s string) string {
			return s[:1]
		},
	})
	globalTemplate = template.Must(globalTemplate.ParseFS(tplFS, "template/*.tpl"))
}
