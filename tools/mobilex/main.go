package main

import (
	"bytes"
	"embed"
	_ "embed"
	"fmt"
	"html/template"
	"os"
	"strings"

	"code.olapie.com/sugar/must"
	"code.olapie.com/sugar/templatex"
	"code.olapie.com/sugar/types"
)

//go:embed templates/*
var templatesDir embed.FS

const header = `
package mobilex

import (
	"code.olapie.com/sugar/mobilex/nomobile"
	"code.olapie.com/sugar/types"
)

`

func main() {
	output := new(bytes.Buffer)
	output.WriteString(header)
	basicTypes := []string{"int", "int16", "int32", "int64", "float64", "bool", "string"}
	tpl := template.Must(template.New("global").Funcs(templatex.TextFuncMap).Funcs(
		template.FuncMap{
			"ttn": typeToName,
		},
	).ParseFS(templatesDir, "templates/*"))

	for _, elem := range basicTypes {
		fmt.Println(elem)
		must.NoError(tpl.ExecuteTemplate(output, "list", types.M{"Elem": elem}))
		must.NoError(tpl.ExecuteTemplate(output, "set", types.M{"Elem": elem}))
		must.NoError(tpl.ExecuteTemplate(output, "pair", types.M{"Elem": elem}))
		must.NoError(tpl.ExecuteTemplate(output, "result", types.M{"Elem": elem}))
		must.NoError(tpl.ExecuteTemplate(output, "result", types.M{"Elem": "*" + typeToName(elem) + "List"}))
	}
	must.NoError(tpl.ExecuteTemplate(output, "result", types.M{"Elem": "[]byte"}))

	keys := []string{"int", "int16", "int32", "int64", "string"}
	values := []string{"int", "int16", "int32", "int64", "float64", "bool", "string"}

	for _, key := range keys {
		for _, val := range values {
			fmt.Println(key, val)
			must.NoError(tpl.ExecuteTemplate(output, "map", types.M{"Key": key, "Value": val}))
		}
	}

	s := output.String()
	for i, c := range s {
		if c != '\n' {
			s = s[i:]
			break
		}
	}

	for j := len(s) - 1; j >= 0; j-- {
		if s[j] != '\n' {
			s = s[:j+1]
			break
		}
	}

	for {
		s2 := strings.Replace(s, "\n\n\n", "\n", -1)
		if s2 == s {
			break
		}
		s = s2
	}

	must.NoError(os.WriteFile("mobilex.gen.go", []byte(s), 0644))

	fmt.Println("Done")
}

func typeToName(typ string) string {
	if typ[0] == '*' {
		typ = typ[1:]
	}
	if typ[:2] == "[]" {
		return templatex.Capitalize(typ[2:], 1) + "Array"
	}
	return templatex.Capitalize(typ, 1)
}
