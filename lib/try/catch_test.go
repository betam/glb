package try

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestThrow(t *testing.T) {
	t.Run(
		"ThrowError",
		func(t *testing.T) {
			assert.PanicsWithError(t, "oops", func() { ThrowError(errors.New("oops")) })
			assert.NotPanics(t, func() { ThrowError(nil) })
		},
	)

	t.Run(
		"Throw",
		func(t *testing.T) {
			value := 10
			err := errors.New("oops")
			assert.Equal(t, value, Throw(value, nil))
			assert.PanicsWithError(t, "oops", func() { Throw(value, err) })
		},
	)
}

func TestCatch(t *testing.T) {
	t.Run(
		"Throwable",
		func(t *testing.T) {
			try := func() {
				panic(errors.New("oops"))
			}
			catch := func(throwable error) {}
			assert.NotPanics(t, func() { Catch(try, nil) })
			assert.NotPanics(t, func() { Catch(try, catch) })
		},
	)

	t.Run(
		"NotThrowable",
		func(t *testing.T) {
			try := func() {
				panic("oops")
			}
			catch := func(throwable error) {}
			assert.PanicsWithValue(t, "oops", func() { Catch(try, nil) })
			assert.PanicsWithValue(t, "oops", func() { Catch(try, catch) })
		},
	)
}
