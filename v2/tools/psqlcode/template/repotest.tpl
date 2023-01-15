{{ define `repotest` }}package repo_test

func setup{{.Name}}(t *testing.T) *repo.{{.Name}} {
    db := setupDB(t)
    _, err := db.Exec(`TRUNCATE TABLE {{.Table}}`)
    if err != nil {
        t.Error(err)
    }
    return repo.New{{.Name}}(db)
}

func new{{.Entity.Name}}() *repo.{{.Entity.Name}} {
    return new(repo.{{.Entity.Name}})
}

func TestInsert(t *testing.T) {
    ctx := context.TODO()
    v := new{{.Entity.Name}}()
    r := setup{{.Name}}(t)
    err := r.Insert(ctx, v)
    if err != nil {
        t.Error(err)
    }
}

func TestUpdate(t *testing.T) {
    ctx := context.TODO()
    v := new{{.Entity.Name}}()
    r := setup{{.Name}}(t)
    err := r.Insert(ctx, v)
    if err != nil {
        t.Error(err)
    }
    // TODO:
    // update columns
    err = r.Update(ctx, v)
    if err != nil {
        t.Error(err)
    }
}

func TestSave(t *testing.T) {
    ctx := context.TODO()
    v := new{{.Entity.Name}}()
    r := setup{{.Name}}(t)
    err := r.Save(ctx, v)
    if err != nil {
        t.Error(err)
    }
    // TODO:
    // update columns
    err = r.Save(ctx, v)
    if err != nil {
        t.Error(err)
    }
}

func TestGet(t *testing.T) {
    v := new{{.Entity.Name}}()
    r := setup{{.Name}}(t)
    // TODO:
}

func TestBatchGet(t *testing.T) {
    // TODO:
}

func TestList(t *testing.T) {
    // TODO:
}

{{end}}