package http_server

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/valyala/fasthttp"

	"github.com/betam/glb/lib/pointer"
)

func TestHandler(t *testing.T) {
	t.Run(
		"Success",
		func(t *testing.T) {
			ctx := &fasthttp.RequestCtx{}

			var success bool
			handler := func(ctx *fasthttp.RequestCtx) Response {
				success = true
				return NewResponse[int](207, pointer.Pointer(777))
			}

			assert.NotPanics(t, func() { Handler(handler)(ctx) })
			assert.Equal(t, []byte("text/plain"), ctx.Response.Header.Peek("Content-Type"))
			assert.Equal(t, 207, ctx.Response.StatusCode())
			assert.Equal(t, []byte("777"), ctx.Response.Body())
			assert.True(t, success)
		},
	)

	t.Run(
		"SuccessEmpty",
		func(t *testing.T) {
			ctx := &fasthttp.RequestCtx{}

			var success bool
			handler := func(ctx *fasthttp.RequestCtx) Response {
				success = true
				return NewJsonResponse[int](204, nil)
			}

			assert.NotPanics(t, func() { Handler(handler)(ctx) })
			assert.Equal(t, []byte("application/json"), ctx.Response.Header.Peek("Content-Type"))
			assert.Equal(t, 204, ctx.Response.StatusCode())
			assert.Equal(t, []byte("null"), ctx.Response.Body())
			assert.True(t, success)
		},
	)

	t.Run(
		"HttpError",
		func(t *testing.T) {
			ctx := &fasthttp.RequestCtx{}

			err := NewError(415, "something went wrong")
			handler := func(ctx *fasthttp.RequestCtx) Response {
				panic(err)
			}

			assert.NotPanics(t, func() { Handler(handler)(ctx) })
			assert.Equal(t, []byte("application/json"), ctx.Response.Header.Peek("Content-Type"))
			assert.Equal(t, 415, ctx.Response.StatusCode())
			assert.Equal(t, []byte(`{"code":415,"message":"something went wrong"}`), ctx.Response.Body())
		},
	)

	t.Run(
		"InternalError",
		func(t *testing.T) {
			ctx := &fasthttp.RequestCtx{}

			err := fmt.Errorf("something went wrong")
			handler := func(ctx *fasthttp.RequestCtx) Response {
				panic(err)
			}

			assert.NotPanics(t, func() { Handler(handler)(ctx) })
			assert.Equal(t, []byte("application/json"), ctx.Response.Header.Peek("Content-Type"))
			assert.Equal(t, 500, ctx.Response.StatusCode())
			assert.Equal(t, []byte(`{"code":500,"message":"something went wrong"}`), ctx.Response.Body())
		},
	)
}

func TestMiddleware(t *testing.T) {
	var before []int
	var after []int

	var mw1 Middleware = func(ctx *fasthttp.RequestCtx, next func(ctx *fasthttp.RequestCtx) Response) Response {
		before = append(before, 1)
		result := next(ctx)
		after = append(after, 1)
		return result
	}

	var mw2 Middleware = func(ctx *fasthttp.RequestCtx, next func(ctx *fasthttp.RequestCtx) Response) Response {
		before = append(before, 2)
		result := next(ctx)
		after = append(after, 2)
		return result
	}

	handler := func(ctx *fasthttp.RequestCtx) Response {
		return NewJsonResponse(200, pointer.Pointer("ok"))
	}

	ctx := fasthttp.RequestCtx{}
	Handler(handler, mw1, mw2)(&ctx)

	assert.Equal(t, []int{1, 2}, before)
	assert.Equal(t, []int{2, 1}, after)
	assert.Equal(t, []byte(`"ok"`), ctx.Response.Body())
}
