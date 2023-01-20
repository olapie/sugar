package main

import (
	"bytes"
	"code.olapie.com/log"
	"code.olapie.com/sugar/v2/xname"
	"go/format"
	"os"
)

func Generate(filename string) {
	os.Mkdir("_generate", 0755)
	for _, e := range ParseYAML(filename) {
		generateSQLForEntity(e)
	}
}

func generateSQLForEntity(r *RepoModel) {
	type Model struct {
		Name           string
		Table          string
		Entity         Entity
		Columns        string
		KeyParams      string
		KeyConditions  string
		KeyArgs        string
		UpdateColumns  string
		ScanHolders    string
		Args           string
		Placeholders   string
		Keys           string
		BatchKeyParams string
		NumKeys        int
	}

	m := &Model{
		Name:           xname.ToClassName(r.Name) + "Repo",
		Table:          r.Table,
		Columns:        r.GetColumns(),
		KeyParams:      r.KeyParams(),
		KeyConditions:  r.KeyConditions(),
		KeyArgs:        r.KeyArgs(),
		UpdateColumns:  r.UpdateColumns(),
		ScanHolders:    r.ScanHolders(),
		Args:           r.Args(),
		Placeholders:   r.Placeholders(),
		Keys:           r.GetKeys(),
		BatchKeyParams: r.BatchKeyParams(),
		NumKeys:        len(r.PrimaryKey),
	}
	m.Entity.Name = xname.ToClassName(r.Name)
	for _, c := range r.Columns {
		field := &Field{
			Name: c.Key.(string),
			Type: c.Value.(string),
		}
		m.Entity.Fields = append(m.Entity.Fields, field)
	}
	var b bytes.Buffer
	err := globalTemplate.ExecuteTemplate(&b, "repo", m)
	if err != nil {
		log.Fatalln(err)
	}

	data, err := format.Source(b.Bytes())
	if err != nil {
		log.Fatalln(err)
	}

	err = os.WriteFile("_generate/"+r.Name+".gen.go", data, 0644)
	if err != nil {
		log.Fatalln(err)
	}
	log.Infof("Generate repo %s", r.Name)

	b.Reset()
	err = globalTemplate.ExecuteTemplate(&b, "repotest", m)
	if err != nil {
		log.Fatalln(err)
	}

	data, err = format.Source(b.Bytes())
	if err != nil {
		log.Fatalln(err)
	}

	err = os.WriteFile("_generate/"+r.Name+"_test.gen.go", data, 0644)
	if err != nil {
		log.Fatalln(err)
	}
	log.Infof("Generate repo test %s", r.Name)
}
