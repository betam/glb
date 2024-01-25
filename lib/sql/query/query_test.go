package query

import (
	"context"
	"database/sql"
	"fmt"
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"
)

type driverResult struct {
}

func (dr driverResult) LastInsertId() (int64, error) {
	return 1, nil
}

func (dr driverResult) RowsAffected() (int64, error) {
	return 2, nil
}

func TestQuery(t *testing.T) {
	t.Run(
		"Success",
		func(t *testing.T) {
			handler := func(ctx context.Context, dest any, query string, args ...any) error {
				assert.Equal(t, "select * from test where (f = $1) order by id asc", query)
				assert.Equal(t, []any{17}, args)
				assert.Equal(t, reflect.Ptr, reflect.TypeOf(dest).Kind())
				assert.Equal(t, reflect.Int, reflect.TypeOf(dest).Elem().Kind())

				*dest.(*int) = 1
				return nil
			}

			result := Query[int](context.Background(), handler, NewBuilder("test").Select("*").Where(And("f", "eq", 17)))
			assert.Equal(t, 1, result)
		},
	)

	t.Run(
		"Insert",
		func(t *testing.T) {
			handler := func(ctx context.Context, dest any, query string, args ...any) error {
				assert.Equal(t, "insert into test (a, f) values ($1, $2) returning *", query)
				assert.Equal(t, []any{1, 2}, args)
				assert.Equal(t, reflect.Ptr, reflect.TypeOf(dest).Kind())
				assert.Equal(t, reflect.Int, reflect.TypeOf(dest).Elem().Kind())

				*dest.(*int) = 1
				return nil
			}

			result := Query[int](context.Background(), handler, NewBuilder("test").Insert("a", "f").Values(1, 2))
			assert.Equal(t, 1, result)
		},
	)

	t.Run(
		"SuccessDestLink",
		func(t *testing.T) {
			handler := func(ctx context.Context, dest any, query string, args ...any) error {
				assert.Equal(t, "select * from test where (f = $1) order by id asc", query)
				assert.Equal(t, []any{17}, args)
				assert.Equal(t, reflect.Ptr, reflect.TypeOf(dest).Kind())
				assert.Equal(t, reflect.Int, reflect.TypeOf(dest).Elem().Kind())

				*dest.(*int) = 1
				return nil
			}

			var link int
			result := Query[int](context.Background(), handler, NewBuilder("test").Select("*").Where(And("f", "eq", 17)), &link)
			assert.Equal(t, 1, result)
			assert.Equal(t, link, result)
		},
	)

	t.Run(
		"InternalError",
		func(t *testing.T) {
			handler := func(ctx context.Context, dest any, query string, args ...any) error {
				panic(fmt.Errorf("something went wrong"))
			}

			assert.PanicsWithError(t, "something went wrong", func() {
				Query[int](context.Background(), handler, NewBuilder("test").Select("*").Where(And("f", "eq", 17)))
			})
		},
	)

	t.Run(
		"ErrorOnQuery",
		func(t *testing.T) {
			handler := func(ctx context.Context, dest any, query string, args ...any) error {
				return fmt.Errorf("something went wrong")
			}

			assert.PanicsWithError(t, "something went wrong", func() {
				Query[int](context.Background(), handler, NewBuilder("test").Select("*").Where(And("f", "eq", 17)))
			})
		},
	)
}

func TestExec(t *testing.T) {
	t.Run(
		"Success",
		func(t *testing.T) {
			handler := func(ctx context.Context, query string, args ...any) (sql.Result, error) {
				assert.Equal(t, "delete from test where (f = $1)", query)
				assert.Equal(t, []any{17}, args)

				return &driverResult{}, nil
			}

			result := Exec(context.Background(), handler, NewBuilder("test").Delete().Where(And("f", "eq", 17)))
			assert.Equal(t, 2, result)
		},
	)

	t.Run(
		"ErrorOnExec",
		func(t *testing.T) {

			handler := func(ctx context.Context, query string, args ...any) (sql.Result, error) {
				return nil, fmt.Errorf("something went wrong")
			}

			assert.PanicsWithError(t, "something went wrong", func() { Exec(context.Background(), handler, NewBuilder("test").Select("*").Where(And("f", "eq", 17))) })
		},
	)
}

func TestExecNamed(t *testing.T) {
	t.Run(
		"SuccessSingle",
		func(t *testing.T) {
			single := struct {
				Field int
			}{
				Field: 111,
			}
			handler := func(ctx context.Context, query string, arg any) (sql.Result, error) {
				assert.Equal(t, "insert into test (field) values (:field)", query)
				assert.Equal(t, []any{single}, arg)

				return &driverResult{}, nil
			}

			result := ExecNamed(context.Background(), handler, NewBuilder("test").Insert("field").Values(single))
			assert.Equal(t, 2, result)
		},
	)

	t.Run(
		"SuccessMulti",
		func(t *testing.T) {
			list := []struct {
				Field int
			}{
				{
					Field: 111,
				},
				{
					Field: 222,
				},
			}
			handler := func(ctx context.Context, query string, arg any) (sql.Result, error) {
				assert.Equal(t, "insert into test (field) values (:field)", query)
				assert.Equal(t, []any{list[0], list[1]}, arg)

				return &driverResult{}, nil
			}

			result := ExecNamed(context.Background(), handler, NewBuilder("test").Insert("field").Values(list))
			assert.Equal(t, 2, result)
		},
	)

	t.Run(
		"ErrorOnExec",
		func(t *testing.T) {

			handler := func(ctx context.Context, query string, arg any) (sql.Result, error) {
				return nil, fmt.Errorf("something went wrong")
			}

			assert.PanicsWithError(t, "something went wrong", func() {
				ExecNamed(context.Background(), handler, NewBuilder("test").Select("*").Where(And("f", "eq", 17)))
			})
		},
	)
}
