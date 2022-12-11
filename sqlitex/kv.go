package sqlitex

import (
	"database/sql"
	"fmt"
	"sync"

	"code.olapie.com/sugar/conv"
	"code.olapie.com/sugar/sqlx"
	"code.olapie.com/sugar/timing"
	"google.golang.org/protobuf/proto"
)

type KVTable struct {
	ID       any
	clock    timing.Clock
	db       *sql.DB
	mu       sync.RWMutex
	filename string
}

func NewKVTable(db *sql.DB, clock timing.Clock) *KVTable {
	r := &KVTable{
		clock:    clock,
		db:       db,
	}

	if r.clock == nil {
		r.clock = timing.LocalClock{}
	}

	_, err := db.Exec(`
CREATE TABLE IF NOT EXISTS kv(
k VARCHAR(255) PRIMARY KEY, 
v BLOB NOT NULL,
updated_at BIGINT NOT NULL
)`)
	if err != nil {
		panic(err)
	}
	return r
}

func (s *KVTable) Filename() string {
	return s.filename
}

func (s *KVTable) SaveInt64(key string, val int64) error{
	s.mu.Lock()
	_, err := s.db.Exec("REPLACE INTO kv(k,v,updated_at) VALUES(?1,?2,?3)",
		key, fmt.Sprint(val), s.clock.Now())
	s.mu.Unlock()
	return err
}

func (s *KVTable) GetInt64(key string) (int64, error) {
	var v string
	s.mu.RLock()
	err := s.db.QueryRow("SELECT v FROM kv WHERE k=?", key).Scan(&v)
	s.mu.RUnlock()
	if err != nil {
		return 0, err
	}

	n, err := conv.ToInt64(v)
	if err != nil {
		return 0, err
	}
	return n, nil
}

func (s *KVTable) SaveData(key string, data []byte) error {
	s.mu.Lock()
	_, err := s.db.Exec("REPLACE INTO kv(k,v,updated_at) VALUES(?1,?2,?3)", key, data, s.clock.Now())
	s.mu.Unlock()
	return err
}

func (s *KVTable) GetData(key string) ([]byte, error) {
	var v []byte
	s.mu.RLock()
	err := s.db.QueryRow("SELECT v FROM kv WHERE k=?", key).Scan(&v)
	s.mu.RUnlock()
	if err != nil {
		return nil, fmt.Errorf("query: %w", err)
	}
	return v, nil
}

func (s *KVTable) SaveString(key string, str string)error {
	return s.SaveData(key, []byte(str))
}

func (s *KVTable) GetString(key string) (string, error) {
	data, err := s.GetData(key)
	if err != nil {
		return "", err
	}
	return string(data), nil
}

func (s *KVTable) SavePB(key string, msg proto.Message) error{
	data, err := proto.Marshal(msg)
	if err != nil {
		return err
	}
	s.mu.Lock()
	_, err = s.db.Exec("REPLACE INTO kv(k,v,updated_at) VALUES(?1,?2,?3)", key, data, s.clock.Now())
	s.mu.Unlock()
	return err
}

func (s *KVTable) GetPB(key string, msg proto.Message) error {
	var v []byte
	s.mu.RLock()
	err := s.db.QueryRow("SELECT v FROM kv WHERE k=?", key).Scan(&v)
	s.mu.RUnlock()
	if err != nil {
		return err
	}
	return proto.Unmarshal(v, msg)
}

func (s *KVTable) SaveJSON(key string, obj any) error{
	v := sqlx.JSON(obj)
	s.mu.Lock()
	var err error
	if v == nil {
		_, err = s.db.Exec("DELETE FROM kv WHERE k=?1", key)
	} else {
		_, err = s.db.Exec("REPLACE INTO kv(k,v,updated_at) VALUES(?1,?2,?3)", key, v, s.clock.Now())
	}
	s.mu.Unlock()
	return err
}

func (s *KVTable) GetJSON(key string, ptrToObj any) error {
	s.mu.RLock()
	err := s.db.QueryRow("SELECT v FROM kv WHERE k=?", key).Scan(sqlx.JSON(ptrToObj))
	s.mu.RUnlock()
	return err
}

func (s *KVTable) Exists(key string) (bool, error) {
	s.mu.RLock()
	var exists bool
	err := s.db.QueryRow("SELECT EXISTS(SELECT * FROM kv WHERE k=?)", key).Scan(&exists)
	s.mu.RUnlock()
	return exists, err
}

func (s *KVTable) Close() error {
	s.mu.Lock()
	err := s.db.Close()
	s.mu.Unlock()
	return err
}
