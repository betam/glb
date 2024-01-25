package query

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestBuilder(t *testing.T) {
	t.Run(
		"SimpleErors",
		func(t *testing.T) {
			assert.PanicsWithError(t, "no table specified", func() { NewBuilder().Build() })
			assert.PanicsWithError(t, "unexpected argument count: expected 0 or 1 given 2", func() { NewBuilder("table", "test") })
		},
	)

	t.Run(
		"Select",
		func(t *testing.T) {
			b := NewBuilder("test")
			q, params := b.Build()
			assert.Equal(t, "select * from test order by id asc", q)
			assert.Empty(t, params)

			q, params = b.Select("f", "s").Build()
			assert.Equal(t, "select f,s from test order by id asc", q)
			assert.Empty(t, params)

			q, params = b.Sort("m asc", "l desc").Build()
			assert.Equal(t, "select f,s from test order by m asc,l desc", q)
			assert.Empty(t, params)

			q, params = b.Sort().Build()
			assert.Equal(t, "select f,s from test order by id asc", q)
			assert.Empty(t, params)

			q, params = b.Page(1, 0).Build()
			assert.Equal(t, "select f,s from test order by id asc offset 1", q)
			assert.Empty(t, params)

			q, params = b.Page(1, 10).Build()
			assert.Equal(t, "select f,s from test order by id asc offset 10 limit 10", q)
			assert.Empty(t, params)

			q, params = b.Page(0, 0).Build()
			assert.Equal(t, "select f,s from test order by id asc", q)
			assert.Empty(t, params)

			q, params = b.Page(0, 10).Build()
			assert.Equal(t, "select f,s from test order by id asc limit 10", q)
			assert.Empty(t, params)

			q, params = b.Page(0, 0).Build()
			assert.Equal(t, "select f,s from test order by id asc", q)
			assert.Empty(t, params)

			q, params = b.Returning("a", "b", "c").Build()
			assert.Equal(t, "select f,s from test order by id asc", q)
			assert.Empty(t, params)

			q, params = b.Window("s as (partition by date)").Build()
			assert.Equal(t, "select f,s from test window s as (partition by date) order by id asc", q)
			assert.Empty(t, params)

			q, params = b.Window("").Build()
			assert.Equal(t, "select f,s from test order by id asc", q)
			assert.Empty(t, params)

			q, params = b.Where(And("f", "eq", 7)).Build()
			assert.Equal(t, "select f,s from test where (f = $1) order by id asc", q)
			assert.Equal(t, &[]any{7}, params)

			b = NewBuilder("test")
			q, params = b.Where(And("f", "eq", 7)).NotSort().Build()
			assert.Equal(t, "select * from test where (f = $1)", q)
			assert.Equal(t, &[]any{7}, params)
		},
	)

	t.Run(
		"Select Join",
		func(t *testing.T) {
			b := NewBuilder("test")

			q, params := b.Select().Join("LEFT OUTER", "other_table", "ot", "ot.id = mainTbl.id").Build()
			assert.Equal(t, "select maintbl.* from test AS maintbl LEFT OUTER JOIN other_table AS ot ON ot.id = mainTbl.id order by maintbl.id asc", q)
			assert.Empty(t, params)

			b = NewBuilder("test")
			q, params = b.Select("id as ma, ot.id as ba").Join("LEFT OUTER", "other_table", "ot", "ot.id = mainTbl.id").Join("LEFT OUTER", "other_table2", "ot2", "ot.id = ot2.id").Build()
			assert.Equal(t, "select id as ma, ot.id as ba from test AS maintbl LEFT OUTER JOIN other_table AS ot ON ot.id = mainTbl.id LEFT OUTER JOIN other_table2 AS ot2 ON ot.id = ot2.id order by maintbl.id asc", q)
			assert.Empty(t, params)

		},
	)

	t.Run(
		"Insert",
		func(t *testing.T) {

		},
	)

	t.Run(
		"Update",
		func(t *testing.T) {

		},
	)

	t.Run(
		"Subquery",
		func(t *testing.T) {
			sub := NewBuilder("test").Where(And("f", "eq", "r"))
			builderWithTable := NewBuilder("test")
			selectBuilder := NewBuilder()

			assert.PanicsWithError(t, "cannot use both table and subquery", func() { builderWithTable.SubTable(sub) })

			q, params := selectBuilder.SubTable(sub).Where(And("m", "eq", []int{1, 2})).Build()
			assert.Equal(t, "select * from (select * from test where (f = $1) order by id asc) s where (m in ($2,$3)) order by id asc", q)
			assert.Equal(t, &[]any{"r", 1, 2}, params)

			subInsert := NewBuilder("test").Where(And("f", "eq", "r"))
			insertBuilder := NewBuilder("tost").Insert("a", "b").SubTable(subInsert)
			q, params = insertBuilder.Build()
			assert.Equal(t, "insert into tost (a, b) select * from test where (f = $1) order by id asc", q)
			assert.Equal(t, &[]any{"r"}, params)
		},
	)

}
