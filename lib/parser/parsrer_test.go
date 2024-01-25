package parser

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func TestValidate(t *testing.T) {
	t.Run(
		"Single",
		func(t *testing.T) {
			var d struct {
				A []struct {
					S string `json:"s,required"`
				} `json:"a,required"`
				B *struct {
					S string `json:"s,required"`
				} `json:"b,required"`
				C *[]struct {
					S string `json:"s,required"`
				} `json:"c,required"`
				D struct {
					S string `json:"s,required"`
				} `json:"d,required"`
				E *[]*struct {
					S string `json:"s,required"`
				} `json:"e,required"`
				F string     `json:"f,required"`
				G int        `json:"g,required"`
				T *time.Time `json:"t,required"`
			}

			assert.NotPanics(
				t, func() {
					validate([]byte(`{"a":[],"b":{"s":"true"},"c":[],"d":{"s":"true"},"e":[],"f":"","g":0,"t":"2022-12-12"}`), &d)
				},
			)
			assert.NotPanics(
				t, func() {
					validate([]byte(`{"a":[{"s":"true"}],"b":{"s":"true"},"c":[],"d":{"s":"true"},"e":[],"f":"","g":0,"t":"2022-12-12"}`), &d)
				},
			)
			assert.NotPanics(
				t, func() {
					validate([]byte(`{"a":[{"s":"true"}],"b":{"s":"true"},"c":[{"s":"true"}],"d":{"s":"true"},"e":[],"f":"","g":0,"t":"2022-12-12"}`), &d)
				},
			)
			assert.NotPanics(
				t, func() {
					validate([]byte(`{"a":[{"s":"true"}],"b":{"s":"true"},"c":[{"s":"true"}],"d":{"s":"true"},"e":[{"s":"true"}],"f":"","g":0,"t":"2022-12-12"}`), &d)
				},
			)
			assert.Panics(
				t, func() {
					validate([]byte(`{"a":[{}],"b":{"s":"true"},"c":[{"s":"true"}],"d":{"s":"true"},"e":[{"s":"true"}],"f":"","g":0,"t":"2022-12-12"}`), &d)
				},
			)
			assert.Panics(
				t, func() {
					validate([]byte(`{"a":[{"s":"true"}],"b":{},"c":[{"s":"true"}],"d":{"s":"true"},"e":[{"s":"true"}],"f":"","g":0,"t":"2022-12-12"}`), &d)
				},
			)
			assert.Panics(
				t, func() {
					validate([]byte(`{"a":[{"s":"true"}],"b":{"s":"true"},"c":[{}],"d":{"s":"true"},"e":[{"s":"true"}],"f":"","g":0,"t":"2022-12-12"}`), &d)
				},
			)
			assert.Panics(
				t, func() {
					validate([]byte(`{"a":[{"s":"true"}],"b":{"s":"true"},"c":[{"s":"true"}],"d":{},"e":[{"s":"true"}],"f":"","g":0,"t":"2022-12-12"}`), &d)
				},
			)
			assert.Panics(
				t, func() {
					validate([]byte(`{"a":[{"s":"true"}],"b":{"s":"true"},"c":[{"s":"true"}],"d":{"s":"true"},"e":[{}],"f":"","g":0,"t":"2022-12-12"}`), &d)
				},
			)
			assert.Panics(
				t, func() {
					validate([]byte(`{"a":[{"s":"true"}],"b":{"s":"true"},"c":[{"s":"true"}],"d":{"s":"true"},"e":[{"s":"true"}],"g":0,"t":"2022-12-12"}`), &d)
				},
			)
			assert.Panics(
				t, func() {
					validate([]byte(`{"a":[{"s":"true"}],"b":{"s":"true"},"c":[{"s":"true"}],"d":{"s":"true"},"e":[{"s":"true"}],"f":"","t":"2022-12-12"}`), &d)
				},
			)
			assert.Panics(
				t, func() {
					validate([]byte(`{"a":[{"s":"true"}],"b":{"s":"true"},"c":[{"s":"true"}],"d":{"s":"true"},"e":[{"s":"true"}],"f":"","g":0}`), &d)
				},
			)
		},
	)

	t.Run(
		"List",
		func(t *testing.T) {
			var d []struct {
				A []struct {
					S string `json:"s,required"`
				} `json:"a,required"`
				B *struct {
					S string `json:"s,required"`
				} `json:"b,required"`
				C *[]struct {
					S string `json:"s,required"`
				} `json:"c,required"`
				D struct {
					S string `json:"s,required"`
				} `json:"d,required"`
				E *[]*struct {
					S string `json:"s,required"`
				} `json:"e,required"`
				F string     `json:"f,required"`
				G int        `json:"g,required"`
				T *time.Time `json:"t,required"`
			}

			assert.NotPanics(
				t, func() {
					validate([]byte(`[{"a":[],"b":{"s":"true"},"c":[],"d":{"s":"true"},"e":[],"f":"","g":0,"t":"2022-12-12"}]`), &d)
				},
			)
			assert.NotPanics(
				t, func() {
					validate([]byte(`[{"a":[{"s":"true"}],"b":{"s":"true"},"c":[],"d":{"s":"true"},"e":[],"f":"","g":0,"t":"2022-12-12"}]`), &d)
				},
			)
			assert.NotPanics(
				t, func() {
					validate([]byte(`[{"a":[{"s":"true"}],"b":{"s":"true"},"c":[{"s":"true"}],"d":{"s":"true"},"e":[],"f":"","g":0,"t":"2022-12-12"}]`), &d)
				},
			)
			assert.NotPanics(
				t, func() {
					validate([]byte(`[{"a":[{"s":"true"}],"b":{"s":"true"},"c":[{"s":"true"}],"d":{"s":"true"},"e":[{"s":"true"}],"f":"","g":0,"t":"2022-12-12"}]`), &d)
				},
			)
			assert.Panics(
				t, func() {
					validate([]byte(`[{"a":[{}],"b":{"s":"true"},"c":[{"s":"true"}],"d":{"s":"true"},"e":[{"s":"true"}],"f":"","g":0,"t":"2022-12-12"}]`), &d)
				},
			)
			assert.Panics(
				t, func() {
					validate([]byte(`[{"a":[],"b":null,"c":[{"s":"true"}],"d":{"s":"true"},"e":[{"s":"true"}],"f":"","g":0,"t":"2022-12-12"}]`), &d)
				},
			)
			assert.Panics(
				t, func() {
					validate([]byte(`[{"a":[],"c":[{"s":"true"}],"d":{"s":"true"},"e":[{"s":"true"}],"f":"","g":0,"t":"2022-12-12"}]`), &d)
				},
			)
			assert.Panics(
				t, func() {
					validate([]byte(`[{"b":{"s":"true"},"c":[{"s":"true"}],"d":{"s":"true"},"e":[{"s":"true"}],"f":"","g":0,"t":"2022-12-12"}]`), &d)
				},
			)
			assert.Panics(
				t, func() {
					validate([]byte(`[{"a":null,"b":{"s":"true"},"c":[{"s":"true"}],"d":{"s":"true"},"e":[{"s":"true"}],"f":"","g":0,"t":"2022-12-12"}]`), &d)
				},
			)
			assert.Panics(
				t, func() {
					validate([]byte(`[{"a":123,"b":{"s":"true"},"c":[{"s":"true"}],"d":{"s":"true"},"e":[{"s":"true"}],"f":"","g":0,"t":"2022-12-12"}]`), &d)
				},
			)
			assert.Panics(
				t, func() {
					validate([]byte(`[{"a":[{"s":"true"}],"b":{},"c":[{"s":"true"}],"d":{"s":"true"},"e":[{"s":"true"}],"f":"","g":0,"t":"2022-12-12"}]`), &d)
				},
			)
			assert.Panics(
				t, func() {
					validate([]byte(`[{"a":[{"s":"true"}],"b":{"s":"true"},"c":[{}],"d":{"s":"true"},"e":[{"s":"true"}],"f":"","g":0,"t":"2022-12-12"}]`), &d)
				},
			)
			assert.Panics(
				t, func() {
					validate([]byte(`[{"a":[{"s":"true"}],"b":{"s":"true"},"c":[{"s":"true"}],"d":{},"e":[{"s":"true"}],"f":"","g":0,"t":"2022-12-12"}]`), &d)
				},
			)
			assert.Panics(
				t, func() {
					validate([]byte(`[{"a":[{"s":"true"}],"b":{"s":"true"},"c":[{"s":"true"}],"d":{"s":"true"},"e":[{}],"f":"","g":0,"t":"2022-12-12"}]`), &d)
				},
			)
			assert.Panics(
				t, func() {
					validate([]byte(`[{"a":[{"s":"true"}],"b":{"s":"true"},"c":[{"s":"true"}],"d":{"s":"true"},"e":[{"s":"true"}],"g":0,"t":"2022-12-12"}]`), &d)
				},
			)
			assert.Panics(
				t, func() {
					validate([]byte(`[{"a":[{"s":"true"}],"b":{"s":"true"},"c":[{"s":"true"}],"d":{"s":"true"},"e":[{"s":"true"}],"f":"","t":"2022-12-12"}]`), &d)
				},
			)
			assert.Panics(
				t, func() {
					validate([]byte(`[{"a":[{"s":"true"}],"b":{"s":"true"},"c":[{"s":"true"}],"d":{"s":"true"},"e":[{"s":"true"}],"f":"","g":0}]`), &d)
				},
			)
		},
	)
}

func TestParse(t *testing.T) {
	t.Run(
		"Success",
		func(t *testing.T) {
			a := struct {
				One int    `json:"one"`
				Two string `json:"two"`
			}{}
			Parse([]byte(`{"one":1,"two":"second"}`), &a)
			assert.Equal(t, 1, a.One)
			assert.Equal(t, "second", a.Two)
		},
	)

	t.Run(
		"Error",
		func(t *testing.T) {
			a := struct {
				One int    `json:"one"`
				Two string `json:"two"`
			}{}
			assert.PanicsWithError(
				t,
				fmt.Sprintf("%v: json: cannot unmarshal string into Go struct field .one of type int", ErrCannotParse),
				func() { Parse([]byte(`{"one":"1","two":"second"}`), &a) },
			)
		},
	)

	t.Run(
		"ErrorValidation",
		func(t *testing.T) {
			a := struct {
				One   int    `json:"one"`
				Two   string `json:"two"`
				Three int    `json:"three,required"`
			}{}
			assert.PanicsWithError(
				t,
				"cannot parse payload: field 'three' is required but missing",
				func() { Parse([]byte(`{"one":1,"two":"second"}`), &a) },
			)
		},
	)

	t.Run(
		"ParseDate",
		func(t *testing.T) {
			var date time.Time
			Parse([]byte(`"2022-11-12T12:12:12Z"`), &date)
			assert.Equal(t, "2022-11-12 12:12:12", date.Format("2006-01-02 15:04:05"))
		},
	)

	t.Run(
		"EmptyDest",
		func(t *testing.T) {
			type tst struct {
				One int    `json:"one"`
				Two string `json:"two"`
			}
			a := Parse([]byte(`{"one":1,"two":"second"}`), new(tst))
			assert.Equal(t, 1, a.One)
			assert.Equal(t, "second", a.Two)
		},
	)

	t.Run(
		"EmptyDest",
		func(t *testing.T) {
			type tst struct {
				One int    `json:"one"`
				Two string `json:"two"`
			}
			a := Parse[tst]([]byte(`{"one":1,"two":"second"}`), nil)
			assert.Equal(t, 1, a.One)
			assert.Equal(t, "second", a.Two)
		},
	)

	t.Run(
		"ToString",
		func(t *testing.T) {
			a := *Parse[string]([]byte(`"{\"t\":true}"`), nil)
			t.Log(a)
			assert.NotNil(t, a)
			b := *Parse[map[string]any]([]byte(a), nil)
			t.Log(b)
			assert.Equal(t, true, b["t"])
		},
	)
}
