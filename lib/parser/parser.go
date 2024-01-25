package parser

import (
	"encoding/json"
	"fmt"
	"reflect"
	"strings"
	"time"

	"github.com/betam/glb/lib/try"
)

var ErrCannotParse = fmt.Errorf("cannot parse payload")

func Parse[Destination any](body []byte, dest *Destination) *Destination {
	if dest == nil {
		dest = new(Destination)
	}
	try.Catch(
		func() {
			try.ThrowError(json.Unmarshal(body, dest))
			validate(body, dest)
		},
		func(throwable error) {
			panic(fmt.Errorf("%w: %v", ErrCannotParse, throwable))
		},
	)

	return dest
}

func validate[Destination any](body []byte, dest *Destination) {
	t := reflect.TypeOf(*dest)
	var value any
	if _, ok := getStructElem(t); ok {
		var v map[string]any
		try.ThrowError(json.Unmarshal(body, &v))
		value = v
	} else if t.Kind() == reflect.Slice {
		if _, ok := getStructElem(t.Elem()); ok {
			var v []map[string]any
			try.ThrowError(json.Unmarshal(body, &v))
			value = v
		}
	}
	try.Catch(
		func() {
			analyze(value, t)
		},
		func(throwable error) {
			if err, ok := throwable.(ValidationError); ok {
				panic(fmt.Errorf("field '%s' is required but missing", err))
			}
			panic(throwable)
		},
	)
}

func getStructElem(t reflect.Type) (reflect.Type, bool) {
	for {
		if t.Kind() == reflect.Ptr {
			t = t.Elem()
		} else {
			break
		}
	}

	if t.Kind() == reflect.Struct && !t.ConvertibleTo(reflect.TypeOf(time.Time{})) {
		return t, true
	}
	return nil, false
}

func analyze(data any, t reflect.Type) {
	for {
		if t.Kind() == reflect.Ptr {
			t = t.Elem()
		} else {
			break
		}
	}

	if tVal, ok := getStructElem(t); ok {
		analyzeStruct(data, tVal)
	} else if t.Kind() == reflect.Slice {
		if tVal, ok := getStructElem(t.Elem()); ok {
			analyzeSlice(data, tVal)
		}
	}
}

func analyzeSlice(data any, t reflect.Type) {
	v := reflect.ValueOf(data)
	if v.Kind() == reflect.Invalid {
		return
	}
	for i := 0; i < v.Len(); i++ {
		try.Catch(
			func() {
				analyzeStruct(v.Index(i).Interface().(map[string]any), t)
			},
			func(throwable error) {
				if err, ok := throwable.(ValidationError); ok {
					panic(ValidationError{fmt.Sprintf("[%d]->%v", i, err)})
				}
				panic(throwable)
			},
		)
	}
}

func analyzeStruct(data any, t reflect.Type) {
	for i := 0; i < t.NumField(); i++ {
		tag := t.Field(i).Tag.Get("json")
		tagName := strings.Split(tag, ",")[0]
		if value, ok := data.(map[string]any)[tagName]; strings.Contains(tag, ",required") && (!ok || value == nil) {
			panic(ValidationError{tagName})
		}
		try.Catch(
			func() {
				if data.(map[string]any)[tagName] == nil {
					return
				}
				analyze(data.(map[string]any)[tagName], t.Field(i).Type)
			},
			func(throwable error) {
				if err, ok := throwable.(ValidationError); ok {
					panic(ValidationError{fmt.Sprintf("%s->%v", tagName, err)})
				}
				panic(throwable)
			},
		)
	}
}
