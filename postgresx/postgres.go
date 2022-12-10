package postgresx

import (
	"database/sql"
	"fmt"
	"os/user"

	"code.olapie.com/sugar/must"
)

func GetConnectionString(name, host string, port int, user, password string, sslEnabled bool) string {
	if host == "" {
		host = "localhost"
	}

	if port == 0 {
		port = 5432
	}

	url := fmt.Sprintf("%s:%d/%s", host, port, name)
	if user == "" {
		url = "postgres://" + url
	} else {
		if password == "" {
			url = "postgres://" + user + "@" + url
		} else {
			url = "postgres://" + user + ":" + password + "@" + url
		}
	}
	if !sslEnabled {
		url = url + "?sslmode=disable"
	}
	return url
}

func Open(connectionString string) (*sql.DB, error) {
	db, err := sql.Open("postgres", connectionString)
	if err != nil {
		return nil, fmt.Errorf("open: %w", err)
	}

	err = db.Ping()
	if err != nil {
		return nil, fmt.Errorf("ping: %w", err)
	}
	return db, nil
}

func MustOpen(connectionString string) *sql.DB {
	return must.Get(Open(connectionString))
}

func GetLocalConnectionString(unixSocket bool) string {
	u, err := user.Current()
	if err != nil {
		fmt.Println(err)
		return ""
	}
	if unixSocket {
		return fmt.Sprintf("postgres:///%s?host=/var/run/postgresql/", u.Username)
	}
	return GetConnectionString(u.Username, "localhost", 5432, u.Username, "", false)
}

func OpenLocal() (*sql.DB, error) {
	if db, err := Open(GetLocalConnectionString(false)); err == nil {
		fmt.Println("Connected via unix socket")
		return db, nil
	}
	db, err := Open(GetLocalConnectionString(true))
	if err == nil {
		fmt.Println("Connected via tcp socket")
	}
	return db, err
}

func MustOpenLocal() *sql.DB {
	return must.Get(OpenLocal())
}
