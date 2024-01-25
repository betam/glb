package query

import (
	"encoding/json"
	"fmt"
	"github.com/betam/glb/lib/pointer"
	"github.com/betam/glb/lib/try"
	"reflect"
	"strings"
)

const (
	strategyAnd = "and"
	strategyOr  = "or"
)

type Expression interface {
	Add(...any) Expression
	Build() (string, *[]any)
	query(*[]any) (string, *[]any)
}

func And(expressions ...any) Expression {
	return pointer.Pointer(expression{strategyAnd, []any{}}).Add(expressions...)
}

func Or(expressions ...any) Expression {
	return pointer.Pointer(expression{strategyOr, []any{}}).Add(expressions...)
}

type expression struct {
	strategy string
	parts    []any
}

func (e *expression) UnmarshalJSON(bytes []byte) error {
	e.parts = []any{}
	var parser func(*expression, any) error
	parser = func(e *expression, data any) error {
		if part, ok := data.(map[string]any); ok {
			if mode, ok := part["mode"].(string); ok {
				e.strategy = mode
			}
			if part["conditions"] == nil {
				return fmt.Errorf("condition is required but missing")
			}
			for _, condition := range part["conditions"].([]any) {
				if _, ok := condition.(map[string]any); ok {
					nested := &expression{}
					err := parser(nested, condition)
					if err != nil {
						return err
					}
					e.Add(nested)
				} else {
					e.Add(condition.([]any)...)
				}
			}
		}
		return nil
	}
	data := map[string]any{}
	try.ThrowError(json.Unmarshal(bytes, &data))
	try.ThrowError(parser(e, data))
	return nil
}

func (e *expression) MarshalJSON() ([]byte, error) {
	data := map[string]any{
		"mode":       e.strategy,
		"conditions": e.parts,
	}
	return json.Marshal(data)
}

func (e *expression) Add(expressions ...any) Expression {
	// 0 — not initialized; 1 — Expression; 2 — operation
	mode := 0
	for _, p := range expressions {
		if _, ok := p.(Expression); ok {
			if mode == 2 {
				panic(fmt.Errorf("cannot mix Expressions and operations"))
			}
			mode = 1
		} else {
			mode = 2
		}
	}
	if mode == 2 {
		if len(expressions) != 3 {
			panic(fmt.Errorf("operations support only 3 argument"))
		}
		expressions = []any{expressions}
	}

	e.parts = append(e.parts, expressions...)
	return e
}

func (e *expression) Build() (string, *[]any) {
	return e.query(pointer.Pointer([]any{}))
}

func (e *expression) query(params *[]any) (string, *[]any) {
	var query []string
	var q string
	for _, part := range e.parts {
		if expr, ok := part.(Expression); ok {
			q, params = expr.query(params)
			query = append(query, q)
		} else {
			v := reflect.ValueOf(part)
			field := v.Index(0).Elem().String()
			operation := v.Index(1).Elem().String()
			value := v.Index(2).Interface()

			q, params = e.build(field, value, operation, params)
			query = append(query, q)
		}
	}

	if e.strategy == "" {
		e.strategy = strategyAnd
	}
	if e.strategy != strategyAnd && e.strategy != strategyOr {
		panic(fmt.Errorf("unsupported query mode: '%s'", e.strategy))
	}
	return fmt.Sprintf("(%s)", strings.Join(query, fmt.Sprintf(") %s (", e.strategy))), params
}

func (e *expression) build(field string, value any, operation string, params *[]any) (string, *[]any) {
	operationsList := map[string]string{
		"eq":      "=",
		"lt":      "<",
		"le":      "<=",
		"gt":      ">",
		"ge":      ">=",
		"ne":      "!=",
		"json_eq": "?",
		"json_in": "@>",
	}

	if op, ok := operationsList[operation]; ok {
		if value == nil || reflect.ValueOf(value).Kind() == reflect.Ptr && reflect.ValueOf(value).IsNil() {
			op = map[string]string{
				"=":  "is null",
				"!=": "is not null",
			}[op]
			return fmt.Sprintf("%s %s", field, op), params
		} else if v := reflect.ValueOf(value); v.Kind() == reflect.Slice {
			op = map[string]string{
				"=":  "in",
				"!=": "not in",
				"?":  "?|",
			}[op]
			list := make([]string, 0, v.Len())
			for i := 0; i < v.Len(); i++ {
				*params = append(*params, v.Index(i).Interface())
				list = append(list, fmt.Sprintf("$%d", len(*params)))
			}
			var in string
			if op == "?|" {
				in = fmt.Sprintf("array[%s]", strings.Join(list, ","))
			} else {
				in = fmt.Sprintf("(%s)", strings.Join(list, ","))
			}
			return fmt.Sprintf("%s %s %s", field, op, in), params
		} else {
			*params = append(*params, value)
			return fmt.Sprintf("%s %s $%d", field, op, len(*params)), params
		}
	}
	panic(fmt.Errorf("unsupported sql operation: '%s'", operation))
}
