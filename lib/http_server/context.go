package http_server

import (
	"github.com/betam/glb/lib/sdk"
	"github.com/valyala/fasthttp"
)

func NewContextWithSession(ctx *fasthttp.RequestCtx) *fasthttp.RequestCtx {
	ctx.SetUserValue(sdk.SessionKey{}, sdk.SetupSession(sdk.Session{
		Uri:           ctx.Request.URI().String(),
		Method:        string(ctx.Request.Header.Method()),
		Body:          string(ctx.Request.Body()),
		CorrelationId: string(ctx.Request.Header.Peek(sdk.HttpCorrelationIdHeader)),
		UserAgent:     string(ctx.Request.Header.Peek(sdk.HttpUserAgentHeader)),
		UserLanguage:  string(ctx.Request.Header.Peek(sdk.HttpUserLanguageHeader)),
		SourceIp:      string(ctx.Request.Header.Peek(sdk.HttpSourceIpHeader)),
	}))
	return ctx
}
