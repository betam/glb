package cmd

import (
	"context"
	"github.com/betam/glb/lib/sdk"
	"strings"
)

func NewContextWithCommand(name string, args []string) context.Context {
	return sdk.NewContextWithSession(
		context.Background(),
		sdk.Session{
			Uri:    name,
			Method: "CLI",
			Body:   strings.Join(args, " "),
		},
	)
}
