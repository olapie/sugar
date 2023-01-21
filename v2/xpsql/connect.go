package xpsql

import (
	"database/sql"
	"fmt"
	"net/url"
	"os/user"

	"code.olapie.com/sugar/v2/must"
)

type OpenOptions struct {
	UnixSocket bool
	Host       string
	Port       int
	User       string
	Password   string
	Database   string
	Schema     string
	SSL        bool
}

func NewOpenOptions() *OpenOptions {
	return &OpenOptions{
		Host: "localhost",
		Port: 5432,
	}
}

func (c *OpenOptions) String() string {
	if c.UnixSocket {
		u, err := user.Current()
		if err != nil {
			fmt.Println(err)
			return ""
		}
		if c.Schema == "" {
			return fmt.Sprintf("postgres:///%s?host=/var/run/postgresql/", u.Username)
		} else {
			return fmt.Sprintf("postgres:///%s?host=/var/run/postgresql/&search_path=%s", u.Username, c.Schema)
		}
	}
	host := c.Host
	port := c.Port
	if host == "" {
		host = "localhost"
	}

	if port == 0 {
		port = 5432
	}

	connStr := fmt.Sprintf("%s:%d", host, port)
	if c.Database != "" {
		connStr += "/" + c.Database
	}
	if c.User == "" {
		connStr = "postgres://" + connStr
	} else {
		if c.Password == "" {
			connStr = "postgres://" + c.User + "@" + connStr
		} else {
			connStr = "postgres://" + c.User + ":" + c.Password + "@" + connStr
		}
	}
	query := url.Values{}
	if !c.SSL {
		query.Add("sslmode", "disable")
	}
	if c.Schema != "" {
		query.Add("search_path", c.Schema)
	}
	if len(query) == 0 {
		return connStr
	}
	return connStr + "?" + query.Encode()
}

func Open(options *OpenOptions) (*sql.DB, error) {
	if options == nil {
		options = NewOpenOptions()
	}
	db, err := sql.Open("postgres", options.String())
	if err != nil {
		return nil, fmt.Errorf("open: %w", err)
	}

	err = db.Ping()
	if err != nil {
		return nil, fmt.Errorf("ping: %w", err)
	}
	return db, nil
}

func MustOpen(options *OpenOptions) *sql.DB {
	return must.Get(Open(options))
}

func OpenLocal() (*sql.DB, error) {
	if db, err := Open(&OpenOptions{UnixSocket: true}); err == nil {
		fmt.Println("Connected via unix socket")
		return db, nil
	}
	db, err := Open(NewOpenOptions())
	if err == nil {
		fmt.Println("Connected via tcp socket")
	}
	return db, err
}

func MustOpenLocal() *sql.DB {
	return must.Get(OpenLocal())
}
