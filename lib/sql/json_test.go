package sql

import (
	"context"
	"database/sql"
	"encoding/json"
	"testing"

	"github.com/jmoiron/sqlx"
	"github.com/stretchr/testify/assert"

	"github.com/betam/glb/lib/try"
)

func TestJson(t *testing.T) {
	mockDB, mock, err := sqlmock.New()
	assert.Nil(t, err)
	defer func(mockDB *sql.DB) { _ = mockDB.Close() }(mockDB)

	dbMock := sqlx.NewDb(mockDB, "sqlmock")

	type J struct {
		Name   string  `json:"name"`
		Age    int     `json:"age"`
		Weight float64 `json:"weight"`
	}

	t.Run(
		"GetMeta",
		func(t *testing.T) {
			type R struct {
				Id    int            `json:"id" db:"id"`
				Value *Json[J]       `json:"value" db:"value"`
				Int   *Json[int]     `json:"int"`
				Float *Json[float64] `json:"float"`
				Bool  *Json[bool]    `json:"bool"`
			}
			var result R

			var data = []byte(`{"name":"me","age":18,"weight":106}`)
			var document = []byte(`{"id":1,"value":{"name":"me","age":18,"weight":106},"int":150,"float":3.14,"bool":true}`)

			var fromString Json[J]
			err = json.Unmarshal(data, &fromString)
			assert.Nil(t, err)
			structValue := J{
				Name:   "me",
				Age:    18,
				Weight: 106,
			}
			var fromStruct = NewJson(structValue)
			assert.Equal(t, fromStruct, &fromString)
			assert.Equal(t, structValue, fromString.Unwrap())

			mock.ExpectQuery("select").WillReturnRows(sqlmock.NewRows([]string{"id", "value", "int", "float", "bool"}).AddRow(1, data, 150, 3.14, true))
			assert.NotPanics(t, func() {
				try.ThrowError(dbMock.GetContext(context.Background(), &result, "select * from somewhere where condition=?", 16))
			})
			assert.Equal(
				t,
				R{
					Id:    1,
					Value: &fromString,
					Int:   NewJson(150),
					Float: NewJson(3.14),
					Bool:  NewJson(true),
				},
				result,
			)

			jsonData, err := json.Marshal(result)
			assert.Nil(t, err)
			assert.Equal(t, document, jsonData, string(jsonData), result)
		},
	)

	t.Run(
		"Insert",
		func(t *testing.T) {
			var data = J{
				Name:   "me",
				Age:    18,
				Weight: 106,
			}

			mock.ExpectExec("insert").WillReturnResult(sqlmock.NewResult(1, 1))
			assert.NotPanics(
				t, func() {
					try.Throw(dbMock.ExecContext(context.Background(), "insert into somewhere values (?, ?)", 1, NewJson(data), NewJson(150), NewJson(3.14), NewJson(true)))
				},
			)
		},
	)
}
