{{ define `repo` }}package repo
type {{.Entity.Name}} struct {
{{range .Entity.Fields}}   {{toStructName .Name}} {{.Type}} `json:"{{toSnake .Name}}"` 
{{end}}}

type {{.Name}} struct {
    db *sql.DB 
}

func New{{.Name}}(db *sql.DB) *{{.Name}} {
    r := new({{.Name}})
    r.db = db 
    return r
}

func (r *{{.Name}}) Insert(ctx context.Context, v *{{.Entity.Name}}) error {
	_, err := r.db.ExecContext(ctx, `INSERT INTO {{.Table}}({{.Columns}}) VALUES({{.Placeholders}})`, {{.Args}})
	return err 
}

func (r *{{.Name}}) Update(ctx context.Context, v *{{.Entity.Name}}) error {
	_, err := r.db.ExecContext(ctx, `UPDATE {{.Table}} SET {{.UpdateColumns}} WHERE {{.KeyConditions}}`, {{.Args}})
	return err
}

func (r *{{.Name}}) Save(ctx context.Context, v *{{.Entity.Name}}) error {
	_, err := r.db.ExecContext(ctx, `INSERT INTO {{.Table}}({{.Columns}}) VALUES({{.Placeholders}}) ON CONFLICT({{.Keys}}) DO UPDATE SET {{.UpdateColumns}}`, {{.Args}})
	return err
}

func (r *{{.Name}}) Delete(ctx context.Context, {{.KeyParams}}) error {
	_, err := r.db.ExecContext(ctx, `DELETE FROM {{.Table}} WHERE {{.KeyConditions}}`, {{.KeyArgs}})
	return err
}

func (r *{{.Name}}) Get(ctx context.Context, {{.KeyParams}}) (v *{{.Entity.Name}}, err error) {
	row := r.db.QueryRowContext(ctx, `SELECT {{.Columns}} FROM {{.Table}} WHERE {{.KeyConditions}}`, {{.KeyArgs}})
	err = row.Scan({{.ScanHolders}})
	return
}

{{if eq .NumKeys 1}}
func (r *{{.Name}}) BatchGet(ctx context.Context, {{.BatchKeyParams}}) (list []*{{.Entity.Name}}, err error) {
	rows, err :=  r.db.QueryContext(ctx, `SELECT {{.Columns}} FROM {{.Table}} WHERE {{.KeyConditions}}`, {{.KeyArgs}})
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
	rows, err := r.db.QueryContext(ctx, `SELECT {{.Columns}} FROM {{.Table}}`)
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