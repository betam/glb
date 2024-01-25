package sdk

import (
	"context"

	"github.com/google/uuid"

	"github.com/betam/glb/lib/try"
)

const (
	HttpCorrelationIdHeader = "X-Correlation-Id"
	HttpUserAgentHeader     = "User-Agent"
	HttpUserLanguageHeader  = "Accept-Language"
	HttpSourceIpHeader      = "X-Real-Ip"
)

type SessionKey struct{}

type Session struct {
	Uri           string
	Method        string
	Body          string
	CorrelationId string
	UserAgent     string
	UserLanguage  string
	SourceIp      string
}

func SetupSession(session Session) *Session {
	if session.CorrelationId == "" {
		session.CorrelationId = try.Throw(uuid.NewRandom()).String()
	}
	return &session
}

func NewContextWithSession(ctx context.Context, session Session) context.Context {
	ctx = context.WithValue(ctx, SessionKey{}, SetupSession(session))
	return ctx
}

func SessionFromContext(ctx context.Context) *Session {
	s, _ := ctx.Value(SessionKey{}).(*Session)
	return s
}
