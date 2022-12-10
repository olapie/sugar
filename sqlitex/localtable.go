package sqlitex

import (
	"context"
	"database/sql"
)

type LocalTable[R any] struct {
	db       *sql.DB
	password string
}

func NewLocalTable[R any](db *sql.DB, password string) *LocalTable[R] {
	t := &LocalTable[R]{
		db:       db,
		password: password,
	}

	// table remote_record: localID, recordData, updateTime, synced
	// table local_record: localID, recordData, createTime, updateTime
	// table deleted_record: localID, deleteTime

	return t
}

func (t *LocalTable[R]) SaveRemote(ctx context.Context, localID string, r R, updateTime int64) error {
	// check delete_record, if it's deleted, then ignore
	// if updateTime < remote_record.updateTime, then ignore
	// save: localID, recordData, udpateTime(new), synced(true)
	return nil
}

func (t *LocalTable[R]) SaveLocal(ctx context.Context, localID string, r R) error {
	// replace local_record
	return nil
}

func (t *LocalTable[R]) Delete(ctx context.Context, localID string) error {
	// delete from remote_record
	// delete from local_record
	// save in delete_record
	return nil
}

func (t *LocalTable[R]) Update(ctx context.Context, localID string, r R) error {
	// update remote_record with synced as false
	return nil
}

func (t *LocalTable[R]) ListUpdates(ctx context.Context) ([]R, error) {
	return nil, nil
}

func (t *LocalTable[R]) ListDeletions(ctx context.Context) ([]R, error) {
	return nil, nil
}

func (t *LocalTable[R]) marshal(localID string, r R) ([]byte, error) {
	// marshaler
	// gob
	// protocol buffer
	// json
	return nil, nil
}

func (t *LocalTable[R]) unmarshal(localID string, r R) ([]byte, error) {
	return nil, nil
}
