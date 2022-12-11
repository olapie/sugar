package sqlitex

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"

	"code.olapie.com/sugar/bytex"
	"code.olapie.com/sugar/cryptox"
	"code.olapie.com/sugar/errorx"
	"code.olapie.com/sugar/must"
	"code.olapie.com/sugar/timing"
)

type LocalTable[R any] struct {
	db       *sql.DB
	password string
	clock    timing.Clock
}

func NewLocalTable[R any](db *sql.DB, password string, clock timing.Clock) *LocalTable[R] {
	if clock == nil {
		clock = timing.LocalClock{}
	}

	t := &LocalTable[R]{
		db:       db,
		password: password,
		clock:    clock,
	}

	// table remote_record: localID, recordData, updateTime, synced
	// table local_record: localID, recordData, createTime, updateTime
	// table deleted_record: localID, deleteTime

	must.Get(db.Exec(`CREATE TABLE IF NOT EXISTS remote_record(
    local_id VARCHAR PRIMARY KEY,
    data BLOB,
    update_time INTEGER,
    synced BOOL DEFAULT FALSE
)`))
	must.Get(db.Exec(`CREATE TABLE IF NOT EXISTS local_record(
    local_id VARCHAR PRIMARY KEY,
    data BLOB,
    create_time INTEGER,
    update_time INTEGER
)`))
	must.Get(db.Exec(`CREATE TABLE IF NOT EXISTS deleted_record(
    local_id VARCHAR PRIMARY KEY,
    data BLOB,
    delete_time INTEGER
)`))

	return t
}

func (t *LocalTable[R]) SaveRemote(ctx context.Context, localID string, record R, updateTime int64) error {
	// check delete_record, if it's deleted, then ignore
	// if updateTime < remote_record.updateTime, then ignore
	// save: localID, recordData, udpateTime(new), synced(true)
	var exists bool
	err := t.db.QueryRowContext(ctx, `SELECT EXISTS(SELECT * FROM deleted_record WHERE local_id=?)`, localID).Scan(&exists)
	if err != nil {
		return fmt.Errorf("query deleted_record: %w", err)
	}

	if exists {
		fmt.Println("Skipped locally deleted record", localID)
		return nil
	}

	err = t.db.QueryRowContext(ctx, `SELECT EXISTS(SELECT * FROM remote_record WHERE local_id=? AND update_time>?)`,
		localID, updateTime).Scan(&exists)
	if err != nil {
		return fmt.Errorf("query remote_record: %w", err)
	}

	if exists {
		fmt.Println("Don't overwrite newly updated local record", localID)
		return nil
	}

	data, err := t.encode(localID, record)
	if err != nil {
		return fmt.Errorf("encode: %s, %w", localID, err)
	}

	_, err = t.db.ExecContext(ctx, `REPLACE INTO remote_record(local_id, data, update_time, synced) VALUES(?,?,?,1)`, localID, data, updateTime)
	if err != nil {
		return fmt.Errorf("replace into remote_record: %s,%w", localID, err)
	}

	_, err = t.db.ExecContext(ctx, `DELETE FROM local_record WHERE local_id=? AND update_time<=?`, localID, updateTime)
	if err != nil {
		return fmt.Errorf("delete from local_record: %s, %w", localID, err)
	}

	return nil
}

func (t *LocalTable[R]) SaveLocal(ctx context.Context, localID string, record R) error {
	// replace local_record
	data, err := t.encode(localID, record)
	if err != nil {
		return fmt.Errorf("encode: %s, %w", localID, err)
	}

	_, err = t.db.ExecContext(ctx, `REPLACE INTO local_record(local_id, data, update_time) VALUES(?,?,?)`,
		localID, data, t.clock.Now().Unix())
	if err != nil {
		return fmt.Errorf("replace into remote_record: %s,%w", localID, err)
	}

	return nil
}

func (t *LocalTable[R]) Delete(ctx context.Context, localID string) error {
	// delete from local_record
	// delete from remote_record
	// save in delete_record
	_, err := t.db.ExecContext(ctx, `DELETE FROM local_record WHERE local_id=?`, localID)
	if err != nil {
		return fmt.Errorf("delete from local_record: %s, %w", localID, err)
	}

	var remoteData []byte
	err = t.db.QueryRowContext(ctx, `SELECT data FROM remote_record WHERE local_id=?`, localID).Scan(&remoteData)
	switch {
	case errors.Is(err, sql.ErrNoRows):
		// don't need to keep deleted record as it doesn't exist remotely
		break
	case err != nil:
		return fmt.Errorf("query remote_record: %s, %w", localID, err)
	default:
		_, err := t.db.ExecContext(ctx, `REPLACE INTO deleted_record(local_id, data, delete_time) VALUES (?,?,?)`,
			localID, remoteData, t.clock.Now().Unix())
		if err != nil {
			return fmt.Errorf("replace into deleted_record: %s, %w", localID, err)
		}

		_, err = t.db.ExecContext(ctx, `DELETE FROM remote_record WHERE local_id=?`, localID)
		if err != nil {
			return fmt.Errorf("delete from remote_record: %s, %w", localID, err)
		}
	}
	return nil
}

func (t *LocalTable[R]) Update(ctx context.Context, localID string, record R) error {
	data, err := t.encode(localID, record)
	if err != nil {
		return fmt.Errorf("encode: %s, %w", localID, err)
	}
	_, err = t.db.ExecContext(ctx, `REPLACE INTO remote_record(local_id, data, update_time, synced) VALUES(?,?,?,1)`,
		localID, data, t.clock.Now().Unix())
	if err != nil {
		return fmt.Errorf("replace into remote_record: %s,%w", localID, err)
	}
	return nil
}

func (t *LocalTable[R]) ListRemotes(ctx context.Context) ([]R, error) {
	rows, err := t.db.QueryContext(ctx, `SELECT local_id, data FROM remote_record`)
	if err != nil {
		return nil, fmt.Errorf("query remote_record: %w", err)
	}
	defer rows.Close()
	return t.scan(rows, "remote_record")
}

func (t *LocalTable[R]) ListUpdates(ctx context.Context) ([]R, error) {
	rows, err := t.db.QueryContext(ctx, `SELECT local_id, data FROM remote_record WHERE synced=0`)
	if err != nil {
		return nil, fmt.Errorf("query remote_record: %w", err)
	}
	defer rows.Close()
	return t.scan(rows, "remote_record")
}

func (t *LocalTable[R]) ListLocals(ctx context.Context) ([]R, error) {
	rows, err := t.db.QueryContext(ctx, `SELECT local_id, data FROM local_record`)
	if err != nil {
		return nil, fmt.Errorf("query deleted_record: %w", err)
	}
	defer rows.Close()
	return t.scan(rows, "local_record")
}

func (t *LocalTable[R]) ListDeletions(ctx context.Context) ([]R, error) {
	rows, err := t.db.QueryContext(ctx, `SELECT local_id, data FROM deleted_record`)
	if err != nil {
		return nil, fmt.Errorf("query deleted_record: %w", err)
	}
	defer rows.Close()
	return t.scan(rows, "remote_record")
}

func (t *LocalTable[R]) Get(ctx context.Context, localID string) (record R, err error) {
	var data []byte
	err = t.db.QueryRowContext(ctx, `SELECT data FROM remote_record WHERE local_id=?`, localID).Scan(&data)
	switch {
	case errors.Is(err, sql.ErrNoRows):
		break
	case err != nil:
		return record, fmt.Errorf("query remote_record: %w", err)
	default:
		return t.decode(localID, data)
	}

	err = t.db.QueryRowContext(ctx, `SELECT data FROM local_record WHERE local_id=?`, localID).Scan(&data)
	switch {
	case errors.Is(err, sql.ErrNoRows):
		break
	case err != nil:
		return record, fmt.Errorf("query local_record: %w", err)
	default:
		return t.decode(localID, data)
	}

	return record, errorx.NotExist
}

func (t *LocalTable[R]) scan(rows *sql.Rows, tableName string) ([]R, error) {
	var records []R
	for rows.Next() {
		var localID string
		var data []byte
		err := rows.Scan(&localID, &data)
		if err != nil {
			return nil, fmt.Errorf("scan %s: %w", tableName, err)
		}

		r, err := t.decode(localID, data)
		if err != nil {
			return nil, fmt.Errorf("decode: %w", err)
		}
		records = append(records, r)
	}
	return records, nil
}

func (t *LocalTable[R]) encode(localID string, r R) (data []byte, err error) {
	data, err = bytex.Marshal(r)
	if err != nil {
		data, err = json.Marshal(r)
		if err != nil {
			return
		}
	}

	return cryptox.Encrypt(data, t.password+localID)
}

func (t *LocalTable[R]) decode(localID string, data []byte) (record R, err error) {
	data, err = cryptox.Decrypt(data, t.password+localID)
	if err != nil {
		return
	}

	err = bytex.Unmarshal(data, &record)
	if err != nil {
		err = json.Unmarshal(data, &record)
	}
	return
}
