// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.15.0

package mrd

import (
	"context"
	"database/sql"
	"fmt"
)

type DBTX interface {
	ExecContext(context.Context, string, ...interface{}) (sql.Result, error)
	PrepareContext(context.Context, string) (*sql.Stmt, error)
	QueryContext(context.Context, string, ...interface{}) (*sql.Rows, error)
	QueryRowContext(context.Context, string, ...interface{}) *sql.Row
}

func New(db DBTX) *Queries {
	return &Queries{db: db}
}

func Prepare(ctx context.Context, db DBTX) (*Queries, error) {
	q := Queries{db: db}
	var err error
	if q.deleteMultiRoomVisibilityStmt, err = db.PrepareContext(ctx, deleteMultiRoomVisibility); err != nil {
		return nil, fmt.Errorf("error preparing query DeleteMultiRoomVisibility: %w", err)
	}
	if q.deleteMultiRoomVisibilityByExpireTSStmt, err = db.PrepareContext(ctx, deleteMultiRoomVisibilityByExpireTS); err != nil {
		return nil, fmt.Errorf("error preparing query DeleteMultiRoomVisibilityByExpireTS: %w", err)
	}
	if q.insertMultiRoomDataStmt, err = db.PrepareContext(ctx, insertMultiRoomData); err != nil {
		return nil, fmt.Errorf("error preparing query InsertMultiRoomData: %w", err)
	}
	if q.insertMultiRoomVisibilityStmt, err = db.PrepareContext(ctx, insertMultiRoomVisibility); err != nil {
		return nil, fmt.Errorf("error preparing query InsertMultiRoomVisibility: %w", err)
	}
	if q.selectMaxIdStmt, err = db.PrepareContext(ctx, selectMaxId); err != nil {
		return nil, fmt.Errorf("error preparing query SelectMaxId: %w", err)
	}
	if q.selectMultiRoomVisibilityRoomsStmt, err = db.PrepareContext(ctx, selectMultiRoomVisibilityRooms); err != nil {
		return nil, fmt.Errorf("error preparing query SelectMultiRoomVisibilityRooms: %w", err)
	}
	return &q, nil
}

func (q *Queries) Close() error {
	var err error
	if q.deleteMultiRoomVisibilityStmt != nil {
		if cerr := q.deleteMultiRoomVisibilityStmt.Close(); cerr != nil {
			err = fmt.Errorf("error closing deleteMultiRoomVisibilityStmt: %w", cerr)
		}
	}
	if q.deleteMultiRoomVisibilityByExpireTSStmt != nil {
		if cerr := q.deleteMultiRoomVisibilityByExpireTSStmt.Close(); cerr != nil {
			err = fmt.Errorf("error closing deleteMultiRoomVisibilityByExpireTSStmt: %w", cerr)
		}
	}
	if q.insertMultiRoomDataStmt != nil {
		if cerr := q.insertMultiRoomDataStmt.Close(); cerr != nil {
			err = fmt.Errorf("error closing insertMultiRoomDataStmt: %w", cerr)
		}
	}
	if q.insertMultiRoomVisibilityStmt != nil {
		if cerr := q.insertMultiRoomVisibilityStmt.Close(); cerr != nil {
			err = fmt.Errorf("error closing insertMultiRoomVisibilityStmt: %w", cerr)
		}
	}
	if q.selectMaxIdStmt != nil {
		if cerr := q.selectMaxIdStmt.Close(); cerr != nil {
			err = fmt.Errorf("error closing selectMaxIdStmt: %w", cerr)
		}
	}
	if q.selectMultiRoomVisibilityRoomsStmt != nil {
		if cerr := q.selectMultiRoomVisibilityRoomsStmt.Close(); cerr != nil {
			err = fmt.Errorf("error closing selectMultiRoomVisibilityRoomsStmt: %w", cerr)
		}
	}
	return err
}

func (q *Queries) exec(ctx context.Context, stmt *sql.Stmt, query string, args ...interface{}) (sql.Result, error) {
	switch {
	case stmt != nil && q.tx != nil:
		return q.tx.StmtContext(ctx, stmt).ExecContext(ctx, args...)
	case stmt != nil:
		return stmt.ExecContext(ctx, args...)
	default:
		return q.db.ExecContext(ctx, query, args...)
	}
}

func (q *Queries) query(ctx context.Context, stmt *sql.Stmt, query string, args ...interface{}) (*sql.Rows, error) {
	switch {
	case stmt != nil && q.tx != nil:
		return q.tx.StmtContext(ctx, stmt).QueryContext(ctx, args...)
	case stmt != nil:
		return stmt.QueryContext(ctx, args...)
	default:
		return q.db.QueryContext(ctx, query, args...)
	}
}

func (q *Queries) queryRow(ctx context.Context, stmt *sql.Stmt, query string, args ...interface{}) *sql.Row {
	switch {
	case stmt != nil && q.tx != nil:
		return q.tx.StmtContext(ctx, stmt).QueryRowContext(ctx, args...)
	case stmt != nil:
		return stmt.QueryRowContext(ctx, args...)
	default:
		return q.db.QueryRowContext(ctx, query, args...)
	}
}

type Queries struct {
	db                                      DBTX
	tx                                      *sql.Tx
	deleteMultiRoomVisibilityStmt           *sql.Stmt
	deleteMultiRoomVisibilityByExpireTSStmt *sql.Stmt
	insertMultiRoomDataStmt                 *sql.Stmt
	insertMultiRoomVisibilityStmt           *sql.Stmt
	selectMaxIdStmt                         *sql.Stmt
	selectMultiRoomVisibilityRoomsStmt      *sql.Stmt
}

func (q *Queries) WithTx(tx *sql.Tx) *Queries {
	return &Queries{
		db:                                      tx,
		tx:                                      tx,
		deleteMultiRoomVisibilityStmt:           q.deleteMultiRoomVisibilityStmt,
		deleteMultiRoomVisibilityByExpireTSStmt: q.deleteMultiRoomVisibilityByExpireTSStmt,
		insertMultiRoomDataStmt:                 q.insertMultiRoomDataStmt,
		insertMultiRoomVisibilityStmt:           q.insertMultiRoomVisibilityStmt,
		selectMaxIdStmt:                         q.selectMaxIdStmt,
		selectMultiRoomVisibilityRoomsStmt:      q.selectMultiRoomVisibilityRoomsStmt,
	}
}