package query

import (
	"context"
	"database/sql"

	"github.com/betam/glb/lib/try"
	"github.com/sirupsen/logrus"
)

type selectHandler func(ctx context.Context, dest any, query string, args ...any) error

type execHandler func(ctx context.Context, query string, args ...any) (sql.Result, error)

type execNamedHandler func(ctx context.Context, query string, arg any) (sql.Result, error)

func Query[Result any](ctx context.Context, handler selectHandler, builder Builder, dest ...*Result) (result Result) {
	query, args := builder.Returning("*").Build()

	//logrus.Trace(query, *args)

	err := handler(ctx, &result, query, *args...)
	if err != nil {
		logrus.Warn(query, *args)
		try.ThrowError(err)
	}
	if len(dest) == 1 {
		*dest[0] = result
	}
	return
}

func Exec(ctx context.Context, handler execHandler, builder Builder) int {
	query, args := builder.Build()

	result, err := handler(ctx, query, *args...)
	if err != nil {
		logrus.Warn(query, *args)
		try.ThrowError(err)
	}
	return int(try.Throw(result.RowsAffected()))
}

func ExecNamed(ctx context.Context, handler execNamedHandler, builder Builder) int {
	query, args := builder.named(true).Build()

	result, err := handler(ctx, query, *args)
	if err != nil {
		logrus.Warn(query, *args)
		try.ThrowError(err)
	}
	return int(try.Throw(result.RowsAffected()))
}
