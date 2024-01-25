package sql

import (
	"context"
	"fmt"

	"github.com/jmoiron/sqlx"

	"github.com/betam/glb/lib/try"
)

type transactionKey struct{}

type transaction struct {
	*sqlx.Tx
	level  int
	isDone bool
}

func NewContextWithTransaction(ctx context.Context, connection Db) (childCtx context.Context, tx *sqlx.Tx, commit func(), rollback func()) {
	if existsTx, ok := ctx.Value(transactionKey{}).(*transaction); ok && !existsTx.isDone {
		childCtx = ctx
		existsTx.level++
		isDone := false
		savepoint := fmt.Sprintf("level_%d", existsTx.level)
		_ = try.Throw(existsTx.ExecContext(ctx, fmt.Sprintf("savepoint %s", savepoint)))
		commit = func() {
			if !isDone {
				_ = try.Throw(existsTx.ExecContext(ctx, fmt.Sprintf("release savepoint %s", savepoint)))
			}
			isDone = true
		}
		rollback = func() {
			if !isDone {
				_ = try.Throw(existsTx.ExecContext(ctx, fmt.Sprintf("rollback to savepoint %s", savepoint)))
			}
			isDone = true
		}
		tx = existsTx.Tx
	} else {
		existsTx = &transaction{Tx: try.Throw(connection.Connect().BeginTxx(ctx, nil))}
		childCtx = context.WithValue(ctx, transactionKey{}, existsTx)
		commit = func() {
			try.ThrowError(existsTx.Tx.Commit())
			existsTx.isDone = true
		}
		rollback = func() {
			if !existsTx.isDone {
				try.ThrowError(existsTx.Tx.Rollback())
			}
			existsTx.isDone = true
		}
		tx = existsTx.Tx
	}

	return childCtx, tx, commit, rollback
}
