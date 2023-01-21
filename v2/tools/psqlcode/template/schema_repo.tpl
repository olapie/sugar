{{ define `schema_repo` }}package repo

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/lib/pq"
)

type {{.Entity.Name}} struct {
{{range .Entity.Fields}}   {{toStructName .Name}} {{.Type}} `json:"{{toSnake .Name}}"` 
{{end}}}

type {{.Name}} struct {
    db *sql.DB
    schema string
}

func New{{.Name}}(db *sql.DB, schema string) *{{.Name}} {
    r := new({{.Name}})
    r.db = db
    r.schema = schema
    return r
}

func (r *{{.Name}}) Insert(ctx context.Context, v *{{.Entity.Name}}) error {
    query := fmt.Sprintf(`INSERT INTO %s.{{.Table}}({{.Columns}}) VALUES({{.Placeholders}})`, r.schema)
	_, err := r.db.ExecContext(ctx, query, {{.Args}})
	return err 
}

func (r *{{.Name}}) Update(ctx context.Context, v *{{.Entity.Name}}) error {
    query := fmt.Sprintf(`UPDATE %s.{{.Table}} SET {{.UpdateColumns}} WHERE {{.KeyConditions}}`, r.schema)
	_, err := r.db.ExecContext(ctx, query, {{.Args}})
	return err
}

func (r *{{.Name}}) Save(ctx context.Context, v *{{.Entity.Name}}) error {
    query := fmt.Sprintf(`INSERT INTO %s.{{.Table}}({{.Columns}}) VALUES({{.Placeholders}}) ON CONFLICT({{.Keys}}) DO UPDATE SET {{.UpdateColumns}}`, r.schema)
	_, err := r.db.ExecContext(ctx, query, {{.Args}})
	return err
}

func (r *{{.Name}}) Delete(ctx context.Context, {{.KeyParams}}) error {
    query := fmt.Sprintf(`DELETE FROM %s.{{.Table}} WHERE {{.KeyConditions}}`, r.schema)
	_, err := r.db.ExecContext(ctx, query, {{.KeyArgs}})
	return err
}

func (r *{{.Name}}) Get(ctx context.Context, {{.KeyParams}}) (v *{{.Entity.Name}}, err error) {
    query := fmt.Sprintf(`SELECT {{.Columns}} FROM %s.{{.Table}} WHERE {{.KeyConditions}}`, r.schema)
	row := r.db.QueryRowContext(ctx, query, {{.KeyArgs}})
	v = new({{.Entity.Name}})
    err = row.Scan({{.ScanHolders}})
    if err != nil {
        return nil, err
    }
    return v, ni
}

{{if eq .NumKeys 1}}
func (r *{{.Name}}) BatchGet(ctx context.Context, {{.BatchKeyParams}}) (list []*{{.Entity.Name}}, err error) {
    query := fmt.Sprintf(`SELECT {{.Columns}} FROM %s.{{.Table}} WHERE {{.BatchKeyConditions}}`, r.schema)
	rows, err :=  r.db.QueryContext(ctx, query, pq.Array({{.BatchKeyArgs}}))
	if err != nil {
        return nil, err
    }
    defer rows.Close()
    for rows.Next() {
        v := new({{.Entity.Name}})
        err = rows.Scan({{.ScanHolders}})
        if err != nil {
            return nil, err
        }
        list = append(list, v)
    }
	return list, nil
}
{{end}}

func (r *{{.Name}}) List(ctx context.Context) (list []*{{.Entity.Name}}, err error) {
    query := fmt.Sprintf(`SELECT {{.Columns}} FROM %s.{{.Table}}`, r.schema)
	rows, err := r.db.QueryContext(ctx, query)
    if err != nil {
        return nil, err
    }
    defer rows.Close()
    for rows.Next() {
        v := new({{.Entity.Name}})
        err = rows.Scan({{.ScanHolders}})
        if err != nil {
            return nil, err
        }
        list = append(list, v)
    }
	return list, nil
}

{{end}}