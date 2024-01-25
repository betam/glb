package sql

import (
	"context"

	"github.com/jmoiron/sqlx"

	baseRepository "github.com/betam/glb/lib/repository"
)

type Repository interface {
	baseRepository.Repository
	Transaction(ctx context.Context) (context.Context, *sqlx.Tx, func(), func())
}

func NewRepository(connection Db) *repository {
	return &repository{DB: connection}
}

type repository struct {
	DB Db
}

func (r *repository) Begin(ctx context.Context) (context.Context, func(), func()) {
	childCtx, _, commit, rollback := NewContextWithTransaction(ctx, r.DB)
	return childCtx, commit, rollback
}

func (r *repository) Transaction(ctx context.Context) (context.Context, *sqlx.Tx, func(), func()) {
	childCtx, tx, commit, rollback := NewContextWithTransaction(ctx, r.DB)
	return childCtx, tx, commit, rollback
}
