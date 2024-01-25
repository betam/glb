package http_server

import (
	"encoding/json"
	"github.com/betam/glb/lib/try"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestResponse(t *testing.T) {
	t.Run(
		"TextResponse",
		func(t *testing.T) {
			data := struct {
				Field string
				Code  int
			}{
				Field: "string",
				Code:  100,
			}
			jsonData := try.Throw(json.Marshal(data))
			var resp Response
			assert.NotPanics(t, func() { resp = NewResponse(200, &data) })
			assert.Equal(t, uint(200), resp.Code())
			assert.Equal(t, "text/plain", resp.Type())
			assert.Equal(t, jsonData, []byte(resp.Content()))
			assert.False(t, resp.(*response).doNotSerialize)
			assert.Equal(t, &data, resp.(*response).content)
		},
	)

	t.Run(
		"JsonResponse",
		func(t *testing.T) {
			data := struct {
				Field string
				Code  int
			}{
				Field: "string",
				Code:  100,
			}
			jsonData := try.Throw(json.Marshal(data))
			var resp Response
			assert.NotPanics(t, func() { resp = NewJsonResponse(200, &data) })
			assert.Equal(t, uint(200), resp.Code())
			assert.Equal(t, "application/json", resp.Type())
			assert.Equal(t, jsonData, []byte(resp.Content()))
			assert.False(t, resp.(*response).doNotSerialize)
			assert.Equal(t, &data, resp.(*response).content)
		},
	)

	t.Run(
		"JsonSerializedResponse",
		func(t *testing.T) {
			data := struct {
				Field string
				Code  int
			}{
				Field: "string",
				Code:  100,
			}
			jsonData := string(try.Throw(json.Marshal(data)))
			var resp Response
			assert.NotPanics(t, func() { resp = NewSerializedJsonResponse(200, &jsonData) })
			assert.Equal(t, uint(200), resp.Code())
			assert.Equal(t, "application/json", resp.Type())
			assert.Equal(t, jsonData, resp.Content())
			assert.True(t, resp.(*response).doNotSerialize)
			assert.Equal(t, &jsonData, resp.(*response).content)
		},
	)
}
