package http_server

import (
	"github.com/betam/glb/lib/parser"
	"github.com/betam/glb/lib/try"
	"github.com/valyala/fasthttp"
)

func Parse[Destination any](body []byte, dest *Destination) *Destination {
	try.Catch(
		func() {
			parser.Parse(body, dest)
		},
		func(throwable error) {
			panic(NewError(fasthttp.StatusBadRequest, throwable.Error()))
		},
	)

	return dest
}
