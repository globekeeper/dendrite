package sqlutil

import (
	"context"
	"database/sql"

	"github.com/cockroachdb/cockroach-go/v2/crdb"
)

// DummyWriter implements sqlutil.Writer.
// The DummyWriter is designed to allow reuse of the sqlutil.Writer
// interface but, unlike ExclusiveWriter, it will not guarantee
// writer exclusivity. This is fine in PostgreSQL where overlapping
// transactions and writes are acceptable.
type DummyWriter struct {
}

// NewDummyWriter returns a new dummy writer.
func NewDummyWriter() Writer {
	return &DummyWriter{}
}

func (w *DummyWriter) Do(db *sql.DB, txn *sql.Tx, f func(txn *sql.Tx) error) error {
	if db != nil && txn == nil {
		return crdb.ExecuteTx(context.Background(), db, nil, f)
	} else {
		return f(txn)
	}
}
