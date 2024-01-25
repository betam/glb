package http_server

import (
	"encoding/json"
	"github.com/betam/glb/lib/sdk"
	"runtime/debug"

	"github.com/sirupsen/logrus"
	"github.com/valyala/fasthttp"

	"github.com/betam/glb/lib/try"
)

type Middleware func(ctx *fasthttp.RequestCtx, next func(ctx *fasthttp.RequestCtx) Response) Response

func Handler(handler func(ctx *fasthttp.RequestCtx) Response, middlewares ...Middleware) func(ctx *fasthttp.RequestCtx) {
	return func(ctx *fasthttp.RequestCtx) {
		ctx = NewContextWithSession(ctx)
		session := sdk.SessionFromContext(ctx)
		ctx.Response.Header.Set(sdk.HttpCorrelationIdHeader, session.CorrelationId)
		try.Catch(
			func() {
				chain := handler
				for i := len(middlewares) - 1; i >= 0; i-- {
					chain = func(mw Middleware, chain func(ctx *fasthttp.RequestCtx) Response) func(ctx *fasthttp.RequestCtx) Response {
						return func(ctx *fasthttp.RequestCtx) Response {
							return mw(ctx, chain)
						}
					}(middlewares[i], chain)
				}
				result := chain(ctx)
				if result != nil {
					ctx.Response.Header.Add("Content-Type", result.Type())
					ctx.SetStatusCode(int(result.Code()))
					ctx.SetBodyString(result.Content())
				}
				logrus.Infof("%s %s [%d]", ctx.Method(), ctx.Path(), 200)
			},
			func(throwable error) {
				code := 500
				isNotHttpError := true
				if err, ok := throwable.(Error); ok {
					code = err.Code
					isNotHttpError = false
				}
				ctx.SetStatusCode(code)
				ctx.SetContentType("application/json")
				message := map[string]any{
					"message": throwable.Error(),
					"code":    code,
				}
				ctx.SetBody(try.Throw(json.Marshal(message)))
				logrus.Infof("%s %s [%d]", ctx.Method(), ctx.Path(), code)
				logrus.Tracef(string(ctx.PostBody()))
				if isNotHttpError {
					logrus.Tracef("stacktrace from panic: \n" + string(debug.Stack()))
				}
				if code < 500 {
					return
				}
				logrus.WithContext(ctx).Error(throwable)

			},
		)
	}
}
