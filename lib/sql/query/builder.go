package query

import (
	"fmt"
	"reflect"
	"strings"

	"github.com/betam/glb/lib/list"
	"github.com/betam/glb/lib/pointer"
)

const (
	modeSelect uint8 = iota
	modeDelete
	modeInsert
	modeUpdate

	doNotSort = "doNotSort"
)

const mainTblAlias = "maintbl"

type Builder interface {
	SubTable(Builder) Builder
	Where(Expression) Builder
	Page(int, int) Builder
	Sort(sort ...string) Builder
	NotSort() Builder
	Window(window string) Builder
	Build() (string, *[]any)
	Returning(fields ...string) Builder
	Group(fields ...string) Builder
	named(value bool) Builder
	parameters(*[]any) Builder
}

type InsertBuilder interface {
	Builder
	Values(...any) InsertBuilder
	Conflict(conflict string, excluded ...string) InsertBuilder
}

type SelectBuilder interface {
	Builder
	Join(mode string, table string, alias string, on string) SelectBuilder
}

type TableBuilder interface {
	Builder
	Select(fields ...string) SelectBuilder
	Delete() Builder
	Insert(fields ...string) InsertBuilder
	Update(values map[string]any) Builder
}

func NewBuilder(table ...string) TableBuilder {
	b := &builder{}
	if len(table) == 1 {
		b.table = table[0]
	} else if len(table) > 1 {
		panic(fmt.Errorf("unexpected argument count: expected 0 or 1 given %d", len(table)))
	}
	return b
}

func (b *builder) Select(fields ...string) SelectBuilder {
	b.queryMode = modeSelect
	b.fields = fields
	return b
}

func (b *builder) Delete() Builder {
	b.queryMode = modeDelete
	return b
}

func (b *builder) Insert(fields ...string) InsertBuilder {
	b.queryMode = modeInsert
	b.fields = fields
	return b
}

func (b *builder) Update(values map[string]any) Builder {
	b.queryMode = modeUpdate
	b.updates = values
	return b
}

type join struct {
	mode  string
	table string
	alias string
	on    string
}

type builder struct {
	queryMode      uint8
	fields         []string
	table          string
	subSelect      Builder
	inserts        [][]any
	updates        map[string]any
	where          Expression
	join           []*join
	window         string
	page           int
	count          int
	sort           []string
	conflictKey    string
	conflictFields []string
	returning      []string
	namedMode      bool
	params         *[]any
	group          []string
}

func (b *builder) Group(fields ...string) Builder {
	b.group = fields
	return b
}

func (b *builder) Returning(fields ...string) Builder {
	b.returning = fields
	return b
}

func (b *builder) named(value bool) Builder {
	b.namedMode = value
	return b
}

func (b *builder) parameters(params *[]any) Builder {
	b.params = params
	return b
}

func (b *builder) Where(expression Expression) Builder {
	b.where = expression
	return b
}

func (b *builder) SubTable(subtable Builder) Builder {
	switch b.queryMode {
	case modeSelect:
		if b.table != "" {
			panic(fmt.Errorf("cannot use both table and subquery"))
		}
		b.subSelect = subtable
	case modeInsert:
		if b.inserts != nil {
			panic(fmt.Errorf("cannot use both values and subquery"))
		}
		b.subSelect = subtable
	}

	return b
}

func (b *builder) Page(page, count int) Builder {
	b.page = page
	b.count = count
	return b
}

func (b *builder) Sort(sort ...string) Builder {
	b.sort = sort
	return b
}

func (b *builder) NotSort() Builder {
	b.sort = make([]string, 0)
	b.sort = append(b.sort, doNotSort)
	return b
}

func (b *builder) Window(window string) Builder {
	b.window = window
	return b
}

func (b *builder) Build() (string, *[]any) {
	query := ""
	if b.params == nil {
		b.params = pointer.Pointer([]any{})
	}

	switch b.queryMode {
	case modeSelect:

		asterics := "*"
		if len(b.join) != 0 {
			asterics = fmt.Sprintf("%s.*", mainTblAlias)
		}

		if b.fields == nil || len(b.fields) == 0 {
			b.fields = append(b.fields, asterics)
		}

		query = fmt.Sprintf("select %s from", strings.Join(b.fields, ","))
	case modeDelete:
		query = "delete from"
	case modeInsert:
		query = "insert into"
	case modeUpdate:
		query = "update"
	default:
		panic(fmt.Errorf("unsupported mode '%d'", b.queryMode))
	}

	if b.table == "" && b.subSelect == nil {
		panic(fmt.Errorf("no table specified"))
	}
	if b.table != "" {
		query = fmt.Sprintf("%s %s", query, b.table)
		if len(b.join) != 0 {
			query = fmt.Sprintf("%s AS %s", query, mainTblAlias)
		}
	}
	if b.queryMode == modeSelect && b.subSelect != nil {
		subquery, params := b.subSelect.Build()
		query = fmt.Sprintf("%s (%s) s", query, subquery)
		*b.params = append(*b.params, *params...)
	}

	switch b.queryMode {
	case modeInsert:
		query = fmt.Sprintf("%s (%s)", query, strings.Join(b.fields, ", "))
		var inserts []string
		if b.namedMode {
			for _, field := range b.fields {
				inserts = append(inserts, ":"+field)
			}
			b.params = pointer.Pointer(b.inserts[0])
			query = fmt.Sprintf("%s values (%s)", query, strings.Join(inserts, ", "))
		} else if b.subSelect != nil {
			subquery, params := b.subSelect.Build()
			query = fmt.Sprintf("%s %s", query, subquery)
			*b.params = append(*b.params, *params...)
		} else {
			for _, value := range b.inserts {
				if len(value) != len(b.fields) {
					panic(fmt.Errorf("wrong insert value (count of fields not match count of values)"))
				}
				var line []string
				for idx := range value {
					if r, ok := value[idx].(*raw); ok {
						line = append(line, r.expression)
					} else {
						*b.params = append(*b.params, value[idx])
						line = append(line, fmt.Sprintf("$%d", len(*b.params)))
					}
				}
				inserts = append(inserts, strings.Join(line, ", "))
			}
			query = fmt.Sprintf("%s values (%s)", query, strings.Join(inserts, "), ("))
		}
	case modeUpdate:
		var updates []string
		for field, value := range b.updates {
			if r, ok := value.(*raw); ok {
				updates = append(updates, fmt.Sprintf("%s=%s", field, r.expression))
			} else {
				*b.params = append(*b.params, value)
				updates = append(updates, fmt.Sprintf("%s=$%d", field, len(*b.params)))
			}
		}
		query = fmt.Sprintf("%s set %s", query, strings.Join(updates, ", "))
	}

	if len(b.join) > 0 && b.queryMode == modeSelect {
		for _, j := range b.join {
			appendJoin := fmt.Sprintf("%s JOIN %s AS %s ON %s", j.mode, j.table, j.alias, j.on)
			query = fmt.Sprintf("%s %s", query, appendJoin)
		}
	}

	if b.where != nil && b.queryMode != modeInsert {
		where, _ := b.where.query(b.params)
		query = fmt.Sprintf("%s where %s", query, where)
	}

	if b.window != "" {
		query = fmt.Sprintf("%s window %s", query, b.window)
	}

	if b.queryMode == modeSelect && len(b.group) > 0 {
		query = fmt.Sprintf("%s group by %s", query, strings.Join(b.group, ","))
	}

	if b.queryMode == modeSelect {
		if len(b.sort) == 1 && b.sort[0] == doNotSort {
			//do nothing with query
		} else {
			if b.sort == nil || len(b.sort) == 0 {
				sortid := "id"
				if len(b.join) != 0 {
					sortid = fmt.Sprintf("%s.id", mainTblAlias)
				}
				b.sort = []string{fmt.Sprintf("%s asc", sortid)}
			}
			query = fmt.Sprintf("%s order by %s", query, strings.Join(b.sort, ","))
		}
	}

	if b.page != 0 {
		offset := b.page
		if b.count != 0 {
			offset *= b.count
		}
		query = fmt.Sprintf("%s offset %d", query, offset)
	}
	if b.count != 0 {
		query = fmt.Sprintf("%s limit %d", query, b.count)
	}

	if b.queryMode != modeSelect && b.conflictKey != "" {
		if b.conflictFields == nil || len(b.conflictFields) == 0 {
			query = fmt.Sprintf("%s on conflict (%s) do nothing", query, b.conflictKey)
		} else {
			conflicts := list.Map(b.conflictFields, func(field string) string { return fmt.Sprintf("%s=excluded.%s", field, field) })
			query = fmt.Sprintf("%s on conflict (%s) do update set %s", query, b.conflictKey, strings.Join(conflicts, ", "))
		}
	}

	if b.returning != nil && len(b.returning) > 0 && b.queryMode != modeSelect {
		query = fmt.Sprintf("%s returning %s", query, strings.Join(b.returning, ","))
	}

	return query, b.params
}

func (b *builder) Values(values ...any) InsertBuilder {
	// 0 — uninitialized; 1 — [][]any; 2 — []any
	mode := 0
	var inserts [][]any
	for _, value := range values {
		if value != nil && reflect.TypeOf(value).Kind() == reflect.Slice {
			if mode != 0 && mode != 1 {
				panic(fmt.Errorf("cannot mix values, structs and slices"))
			}
			mode = 1
			v := reflect.ValueOf(value)
			var item []any
			for i := 0; i < v.Len(); i++ {
				item = append(item, v.Index(i).Interface())
			}
			inserts = append(inserts, item)
		} else {
			if mode != 0 && mode != 2 {
				panic(fmt.Errorf("cannot mix values, structs and slices"))
			}
			mode = 2
		}
	}
	if mode == 2 {
		inserts = append(inserts, values)
	}
	b.inserts = inserts
	return b
}

func (b *builder) Conflict(conflict string, excluded ...string) InsertBuilder {
	b.conflictKey = conflict
	b.conflictFields = excluded
	return b
}

func (b *builder) Join(mode string, table string, alias string, on string) SelectBuilder {

	j := &join{
		mode:  mode,
		table: table,
		alias: alias,
		on:    on,
	}
	b.join = append(b.join, j)
	return b
}
