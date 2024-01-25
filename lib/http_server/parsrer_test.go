package http_server

import (
	"github.com/stretchr/testify/assert"
	"github.com/valyala/fasthttp"
	"testing"
	"time"
)

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
			assert.PanicsWithValue(
				t,
				NewError(fasthttp.StatusBadRequest, "cannot parse payload: json: cannot unmarshal string into Go struct field .one of type int"),
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
			assert.PanicsWithValue(
				t,
				NewError(fasthttp.StatusBadRequest, "cannot parse payload: field 'three' is required but missing"),
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
}
