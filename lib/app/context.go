package app

import (
	"context"
	"github.com/betam/glb/lib/sdk"
	"strings"
)

func NewContextWithCancelWithCommand(name string, args []string) (context.Context, context.CancelFunc) {

	ctx, cancel := context.WithCancel(context.Background())

	return sdk.NewContextWithSession(
		ctx,
		sdk.Session{
			Uri:    name,
			Method: "CLI",
			Body:   strings.Join(args, " "),
		},
	), cancel
}
