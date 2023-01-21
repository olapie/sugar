package main

import (
	"bytes"
	"go/format"
	"os"

	"code.olapie.com/log"
	"code.olapie.com/sugar/v2/xname"
)

func Generate(filename string) {
	os.Mkdir("generate", 0755)
	var model struct {
		Entities []Entity
	}
	for _, e := range ParseYAML(filename) {
		generateSQLForEntity(e)
		model.Entities = append(model.Entities, getEntity(e))
	}

	var b bytes.Buffer
	err := globalTemplate.ExecuteTemplate(&b, "model", model)
	if err != nil {
		log.Fatalln(err)
	}

	data, err := format.Source(b.Bytes())
	if err != nil {
		log.Fatalln(err)
	}

	err = os.WriteFile("generate/model.go", data, 0644)
	if err != nil {
		log.Fatalln(err)
	}
}

func getEntity(r *RepoModel) Entity {
	var e Entity
	e.Name = xname.ToClassName(r.Name)
	for _, c := range r.Columns {
		field := &Field{
			Name: c.Key.(string),
			Type: c.Value.(string),
		}
		e.Fields = append(e.Fields, field)
	}
	return e
}

func generateSQLForEntity(r *RepoModel) {
	type Model struct {
		Name               string
		Table              string
		Entity             Entity
		Columns            string
		KeyParams          string
		KeyConditions      string
		KeyArgs            string
		UpdateColumns      string
		ScanHolders        string
		Args               string
		Placeholders       string
		Keys               string
		BatchKeyParams     string
		BatchKeyConditions string
		BatchKeyArgs       string
		NumKeys            int
	}

	m := &Model{
		Name:               xname.ToClassName(r.Name) + "Repo",
		Table:              r.Table,
		Columns:            r.GetColumns(),
		KeyParams:          r.KeyParams(),
		KeyConditions:      r.KeyConditions(),
		KeyArgs:            r.KeyArgs(),
		UpdateColumns:      r.UpdateColumns(),
		ScanHolders:        r.ScanHolders(),
		Args:               r.Args(),
		Placeholders:       r.Placeholders(),
		Keys:               r.GetKeys(),
		BatchKeyParams:     r.BatchKeyParams(),
		BatchKeyConditions: r.BatchKeyConditions(),
		BatchKeyArgs:       r.BatchKeyArgs(),
		NumKeys:            len(r.PrimaryKey),
	}

	tplName := "repo"
	testTplName := "repotest"
	m.Entity = getEntity(r)
	var b bytes.Buffer
	err := globalTemplate.ExecuteTemplate(&b, tplName, m)
	if err != nil {
		log.Fatalln(err)
	}

	data, err := format.Source(b.Bytes())
	if err != nil {
		log.Fatalln(err)
	}

	err = os.WriteFile("generate/"+r.Name+".go", data, 0644)
	if err != nil {
		log.Fatalln(err)
	}
	log.Infof("Generate repo %s", r.Name)

	b.Reset()
	err = globalTemplate.ExecuteTemplate(&b, testTplName, m)
	if err != nil {
		log.Fatalln(err)
	}

	data, err = format.Source(b.Bytes())
	if err != nil {
		log.Fatalln(err)
	}

	err = os.WriteFile("generate/"+r.Name+"_test.go", data, 0644)
	if err != nil {
		log.Fatalln(err)
	}
	log.Infof("Generate repo test %s", r.Name)
}
