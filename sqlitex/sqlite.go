package sqlitex

import (
	"database/sql"
	"fmt"
	"os"
	"path/filepath"

	"code.olapie.com/sugar/must"
)

func Open(fileName string) (*sql.DB, error) {
	err := os.MkdirAll(filepath.Dir(fileName), 0755)
	if err != nil {
		return nil, err
	}

	_, err = os.OpenFile(fileName, os.O_RDWR|os.O_CREATE, 0644)
	if err != nil {
		return nil, err
	}

	dataSource := fmt.Sprintf("file:%s?cache=shared", fileName)
	return sql.Open("sqlite3", dataSource)
}

func MustOpen(filename string) *sql.DB {
	return must.Get(Open(filename))
}
