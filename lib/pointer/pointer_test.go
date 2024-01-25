package pointer

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestPointer(t *testing.T) {
	t.Run(
		"Simple",
		func(t *testing.T) {
			value := 100
			assert.Equal(t, &value, Pointer(value))
		},
	)

	t.Run(
		"EmptyValue",
		func(t *testing.T) {
			var value int
			assert.Equal(t, &value, Pointer(value))
		},
	)
}

func TestValue(t *testing.T) {
	t.Run(
		"Simple",
		func(t *testing.T) {
			value := 100
			assert.Equal(t, value, Value(&value))
		},
	)

	t.Run(
		"NullValue",
		func(t *testing.T) {
			var value *string
			assert.Equal(t, "", Value(value))
		},
	)
}
