package http_server

import (
	"encoding/json"
	"github.com/betam/glb/lib/try"
)

type Response interface {
	Content() string
	Type() string
	Code() uint
}

type response struct {
	contentType    string
	content        any
	code           uint
	doNotSerialize bool
}

func (r *response) Content() string {
	if r.doNotSerialize {
		return *(r.content).(*string)
	}
	return string(try.Throw(json.Marshal(r.content)))
}

func (r *response) Type() string {
	return r.contentType
}

func (r *response) Code() uint {
	return r.code
}

func NewResponse[T any](code uint, content *T) Response {
	return rawResponse(code, "text/plain", content)
}

func NewEmptyResponse(code uint) Response {
	resp := "[]"
	r := rawResponse(code, "application/json", &resp)
	r.doNotSerialize = true
	return r
}

func NewJsonResponse[T any](code uint, content *T) Response {
	return rawResponse(code, "application/json", content)
}

func NewSerializedJsonResponse[T any](code uint, content *T) Response {
	r := rawResponse(code, "application/json", content)
	r.doNotSerialize = true
	return r
}

func rawResponse(code uint, contentType string, content any) *response {
	return &response{
		code:        code,
		contentType: contentType,
		content:     content,
	}
}
