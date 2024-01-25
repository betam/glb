package repository

import "context"

type Repository interface {
	Begin(ctx context.Context) (childCtx context.Context, commit func(), rollback func())
}
