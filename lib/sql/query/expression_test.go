package query

import (
	"encoding/json"
	"github.com/betam/glb/lib/try"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestExpression(t *testing.T) {
	t.Run(
		"Serialize",
		func(t *testing.T) {
			q := And()
			compare := Or("id", "lt", 7.).Add(And("id", "eq", nil))
			payload := []byte(`{"mode":"or","conditions":[["id","lt",7],{"conditions":[["id","eq",null]]}]}`)
			try.ThrowError(json.Unmarshal(payload, &q))
			str, params := q.Build()
			compareStr, compareParams := compare.Build()
			assert.Equal(t, "(id < $1) or ((id is null))", str)
			assert.Equal(t, compareStr, str)
			assert.Equal(t, &[]any{7.}, params)
			assert.Equal(t, compareParams, params)
		},
	)
	t.Run(
		"Success",
		func(t *testing.T) {
			q := And("field", "eq", 1).
				Add(Or("field2", "eq", "str").Add("field3", "ne", nil)).
				Add("field4", "gt", 3.14).
				Add("field5", "json_eq", []int{7, 5, 4}).
				Add("field6", "eq", []string{"4", "fr"}).
				Add("f", "ne", []int{1})
			qJson := try.Throw(json.Marshal(q))
			assert.Equal(
				t,
				[]byte(`{"conditions":[["field","eq",1],{"conditions":[["field2","eq","str"],["field3","ne",null]],"mode":"or"},["field4","gt",3.14],["field5","json_eq",[7,5,4]],["field6","eq",["4","fr"]],["f","ne",[1]]],"mode":"and"}`),
				qJson,
			)

			sql, params := q.Build()
			assert.Equal(
				t,
				"(field = $1) and ((field2 = $2) or (field3 is not null)) and (field4 > $3) and (field5 ?| array[$4,$5,$6]) and (field6 in ($7,$8)) and (f not in ($9))",
				sql,
			)
			assert.Equal(t, &[]any{1, "str", 3.14, 7, 5, 4, "4", "fr", 1}, params)
		},
	)

	t.Run(
		"MixedOperations",
		func(t *testing.T) {
			assert.PanicsWithError(t, "cannot mix Expressions and operations", func() { And("field", "eq", 17, And("field", "ne", 6)) })
		},
	)

	t.Run(
		"WrongOperationsArgsCount",
		func(t *testing.T) {
			assert.PanicsWithError(t, "operations support only 3 argument", func() { And("field", "eq", 17, "fire", "ne", nil) })
		},
	)

	t.Run(
		"Unsupported operation",
		func(t *testing.T) {
			assert.PanicsWithError(t, "unsupported sql operation: 'some'", func() { And("field", "some", 17).Build() })
		},
	)
}
