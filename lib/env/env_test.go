package env

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNotExists(t *testing.T) {
	t.Run(
		"Value",
		func(t *testing.T) {
			assert.PanicsWithError(t, "no env found: TEST_NOT_EXISTS", func() { _ = NeedValue[string]("TEST_NOT_EXISTS") })
			assert.PanicsWithError(t, "no env found: TEST_NOT_EXISTS", func() { _ = NeedValue[int]("TEST_NOT_EXISTS") })
			assert.PanicsWithError(t, "no env found: TEST_NOT_EXISTS", func() { _ = NeedValue[float64]("TEST_NOT_EXISTS") })
			assert.PanicsWithError(t, "no env found: TEST_NOT_EXISTS", func() { _ = NeedValue[bool]("TEST_NOT_EXISTS") })

			assert.Equal(t, "just a string", Value("TEST_NOT_EXISTS", "just a string"))
			assert.Equal(t, 100, Value("TEST_NOT_EXISTS", 100))
			assert.Equal(t, 100.0, Value("TEST_NOT_EXISTS", 100.0))
			assert.Equal(t, true, Value("TEST_NOT_EXISTS", true))
			assert.Equal(t, false, Value("TEST_NOT_EXISTS", false))
		},
	)
	t.Run(
		"Array",
		func(t *testing.T) {
			assert.PanicsWithError(t, "no env found: TEST_NOT_EXISTS", func() { _ = NeedArray[string]("TEST_NOT_EXISTS") })
			assert.PanicsWithError(t, "no env found: TEST_NOT_EXISTS", func() { _ = NeedArray[int]("TEST_NOT_EXISTS") })
			assert.PanicsWithError(t, "no env found: TEST_NOT_EXISTS", func() { _ = NeedArray[float64]("TEST_NOT_EXISTS") })
			assert.PanicsWithError(t, "no env found: TEST_NOT_EXISTS", func() { _ = NeedArray[bool]("TEST_NOT_EXISTS") })

			assert.Equal(t, []string{"just", "a", "string"}, Array("TEST_NOT_EXISTS", []string{"just", "a", "string"}))
			assert.Equal(t, []int{100, 200, 300}, Array("TEST_NOT_EXISTS", []int{100, 200, 300}))
			assert.Equal(t, []float64{100.0, 200.0, 300.0}, Array("TEST_NOT_EXISTS", []float64{100, 200, 300}))
			assert.Equal(t, []bool{true, false, true}, Array("TEST_NOT_EXISTS", []bool{true, false, true}))
		},
	)
}

func TestEmpty(t *testing.T) {
	_ = os.Setenv("TEST_EMPTY", "")

	t.Run(
		"Value",
		func(t *testing.T) {
			assert.Equal(t, "", NeedValue[string]("TEST_EMPTY"))
			assert.PanicsWithError(
				t,
				`error during parse environment: TEST_EMPTY: strconv.ParseInt: parsing "": invalid syntax`,
				func() { _ = NeedValue[int]("TEST_EMPTY") },
			)
			assert.PanicsWithError(
				t,
				`error during parse environment: TEST_EMPTY: strconv.ParseFloat: parsing "": invalid syntax`,
				func() { _ = NeedValue[float64]("TEST_EMPTY") },
			)
			assert.PanicsWithError(
				t,
				`error during parse environment: TEST_EMPTY: non-bool expression: ""`,
				func() { _ = NeedValue[bool]("TEST_EMPTY") },
			)

			assert.Equal(t, "", Value("TEST_EMPTY", "just a string"))
			assert.Equal(t, 100, Value("TEST_EMPTY", 100))
			assert.Equal(t, 100.0, Value("TEST_EMPTY", 100.0))
			assert.Equal(t, true, Value("TEST_EMPTY", true))
		},
	)

	t.Run(
		"Array",
		func(t *testing.T) {
			assert.Equal(t, []string{""}, NeedArray[string]("TEST_EMPTY"))
			assert.PanicsWithError(
				t,
				`error during parse array environment: TEST_EMPTY: strconv.ParseInt: parsing "": invalid syntax`,
				func() { _ = NeedArray[int]("TEST_EMPTY") },
			)
			assert.PanicsWithError(
				t,
				`error during parse array environment: TEST_EMPTY: strconv.ParseFloat: parsing "": invalid syntax`,
				func() { _ = NeedArray[float64]("TEST_EMPTY") },
			)
			assert.PanicsWithError(
				t,
				`error during parse array environment: TEST_EMPTY: non-bool expression: ""`,
				func() { _ = NeedArray[bool]("TEST_EMPTY") },
			)

			assert.Equal(t, []string{""}, Array("TEST_EMPTY", []string{"just", "a", "string"}))
			assert.Equal(t, []int{100, 200, 300}, Array("TEST_EMPTY", []int{100, 200, 300}))
			assert.Equal(t, []float64{100.0, 200.0, 300.0}, Array("TEST_EMPTY", []float64{100, 200, 300}))
			assert.Equal(t, []bool{true, false, true}, Array("TEST_EMPTY", []bool{true, false, true}))
		},
	)
}

func TestString(t *testing.T) {
	_ = os.Setenv("TEST_STRING", "test")
	_ = os.Setenv("TEST_ARRAY_ONE_ELEMENT", "1.0.0")
	_ = os.Setenv("TEST_ARRAY", "one,two,three")

	t.Run(
		"Value",
		func(t *testing.T) {
			assert.Equal(t, "test", NeedValue[string]("TEST_STRING"))

			assert.Equal(t, "test", Value("TEST_STRING", "just a string"))
		},
	)

	t.Run(
		"Array",
		func(t *testing.T) {
			assert.Equal(t, []string{NeedValue[string]("TEST_ARRAY_ONE_ELEMENT")}, NeedArray[string]("TEST_ARRAY_ONE_ELEMENT"))
			assert.Equal(t, []string{"one", "two", "three"}, NeedArray[string]("TEST_ARRAY"))

			assert.Equal(t, []string{NeedValue[string]("TEST_ARRAY_ONE_ELEMENT")}, Array("TEST_ARRAY_ONE_ELEMENT", []string{"default", "array"}))
			assert.Equal(t, []string{"one", "two", "three"}, Array("TEST_ARRAY", []string{"default", "array"}))
		},
	)
}

func TestInt(t *testing.T) {
	_ = os.Setenv("TEST_STRING", "1.0")
	_ = os.Setenv("TEST_INT", "1")
	_ = os.Setenv("TEST_UINT", "-1")
	_ = os.Setenv("TEST_ARRAY_ONE_ELEMENT", "1")
	_ = os.Setenv("TEST_ARRAY", "1 ,2, 3")
	_ = os.Setenv("TEST_ARRAY_NOT_INT", "1,2d,3")

	t.Run(
		"Value",
		func(t *testing.T) {
			assert.PanicsWithError(
				t,
				`error during parse environment: TEST_STRING: strconv.ParseInt: parsing "1.0": invalid syntax`,
				func() { _ = NeedValue[int]("TEST_STRING") },
			)
			assert.Equal(t, 1, NeedValue[int]("TEST_INT"))

			assert.Equal(t, 100, Value("TEST_STRING", 100))
			assert.Equal(t, 1, Value("TEST_INT", 100))
			assert.Equal(t, int64(1), Value("TEST_INT", int64(100)))
			assert.Equal(t, int32(1), Value("TEST_INT", int32(100)))
			assert.Equal(t, int16(1), Value("TEST_INT", int16(100)))
			assert.Equal(t, int8(1), Value("TEST_INT", int8(100)))
			assert.Equal(t, uint(1), Value("TEST_INT", uint(100)))
			assert.Equal(t, uint64(1), Value("TEST_INT", uint64(100)))
			assert.Equal(t, uint32(1), Value("TEST_INT", uint32(100)))
			assert.Equal(t, uint16(1), Value("TEST_INT", uint16(100)))
			assert.Equal(t, uint8(1), Value("TEST_INT", uint8(100)))
			assert.Equal(t, uint(100), Value("TEST_UINT", uint(100)))
		},
	)

	t.Run(
		"Array",
		func(t *testing.T) {
			assert.PanicsWithError(
				t,
				`error during parse array environment: TEST_STRING: strconv.ParseInt: parsing "1.0": invalid syntax`,
				func() { _ = NeedArray[int]("TEST_STRING") },
			)
			assert.PanicsWithError(
				t,
				`error during parse array environment: TEST_ARRAY_NOT_INT: strconv.ParseInt: parsing "2d": invalid syntax`,
				func() { _ = NeedArray[int]("TEST_ARRAY_NOT_INT") },
			)
			assert.Equal(t, []int{NeedValue[int]("TEST_ARRAY_ONE_ELEMENT")}, NeedArray[int]("TEST_ARRAY_ONE_ELEMENT"))
			assert.Equal(t, []int{1, 2, 3}, NeedArray[int]("TEST_ARRAY"))

			assert.Equal(t, []int{100, 200, 300}, Array("TEST_STRING", []int{100, 200, 300}))
			assert.Equal(t, []int{100, 200, 300}, Array("TEST_ARRAY_NOT_INT", []int{100, 200, 300}))
			assert.Equal(t, []int{NeedValue[int]("TEST_ARRAY_ONE_ELEMENT")}, Array("TEST_ARRAY_ONE_ELEMENT", []int{1000, 2000}))
			assert.Equal(t, []int{1, 2, 3}, Array("TEST_ARRAY", []int{100, 200, 300}))
		},
	)
}

func TestFloat(t *testing.T) {
	_ = os.Setenv("TEST_STRING", "1.0.0")
	_ = os.Setenv("TEST_FLOAT", "1")
	_ = os.Setenv("TEST_ARRAY_ONE_ELEMENT", "1.3")
	_ = os.Setenv("TEST_ARRAY", "1 ,2, 3")
	_ = os.Setenv("TEST_ARRAY_NOT_FLOAT", "1,2d,3")

	t.Run(
		"Value",
		func(t *testing.T) {
			assert.PanicsWithError(
				t,
				`error during parse environment: TEST_STRING: strconv.ParseFloat: parsing "1.0.0": invalid syntax`,
				func() { _ = NeedValue[float64]("TEST_STRING") },
			)
			assert.Equal(t, 1.0, NeedValue[float64]("TEST_FLOAT"))

			assert.Equal(t, 100.0, Value("TEST_STRING", 100.0))
			assert.Equal(t, 1.0, Value("TEST_FLOAT", 100.0))
			assert.Equal(t, float32(1), Value("TEST_FLOAT", float32(100.0)))
		},
	)

	t.Run(
		"Array",
		func(t *testing.T) {
			assert.PanicsWithError(
				t,
				`error during parse array environment: TEST_STRING: strconv.ParseFloat: parsing "1.0.0": invalid syntax`,
				func() { _ = NeedArray[float64]("TEST_STRING") },
			)
			assert.PanicsWithError(
				t,
				`error during parse array environment: TEST_ARRAY_NOT_FLOAT: strconv.ParseFloat: parsing "2d": invalid syntax`,
				func() { _ = NeedArray[float64]("TEST_ARRAY_NOT_FLOAT") },
			)
			assert.Equal(t, []float64{NeedValue[float64]("TEST_ARRAY_ONE_ELEMENT")}, NeedArray[float64]("TEST_ARRAY_ONE_ELEMENT"))
			assert.Equal(t, []float64{1, 2, 3}, NeedArray[float64]("TEST_ARRAY"))

			assert.Equal(t, []float64{100, 200, 300}, Array("TEST_STRING", []float64{100, 200, 300}))
			assert.Equal(t, []float64{100, 200, 300}, Array("TEST_ARRAY_NOT_INT", []float64{100, 200, 300}))
			assert.Equal(t, []float64{NeedValue[float64]("TEST_ARRAY_ONE_ELEMENT")}, Array("TEST_ARRAY_ONE_ELEMENT", []float64{1000, 2000}))
			assert.Equal(t, []float64{1, 2, 3}, Array("TEST_ARRAY", []float64{100, 200, 300}))
		},
	)
}

func TestBool(t *testing.T) {
	_ = os.Setenv("TEST_STRING", "1.0.0")
	_ = os.Setenv("TEST_BOOL_FALSE", "false")
	_ = os.Setenv("TEST_BOOL_TRUE", "true")
	_ = os.Setenv("TEST_BOOL_ONE", "1")
	_ = os.Setenv("TEST_ARRAY_ONE_ELEMENT", "true")
	_ = os.Setenv("TEST_ARRAY", "1, false ,true, 0")
	_ = os.Setenv("TEST_ARRAY_NOT_BOOL", "1,true,3")

	t.Run(
		"Value",
		func(t *testing.T) {
			assert.PanicsWithError(
				t,
				`error during parse environment: TEST_STRING: non-bool expression: "1.0.0"`,
				func() { _ = NeedValue[bool]("TEST_STRING") },
			)
			assert.Equal(t, false, NeedValue[bool]("TEST_BOOL_FALSE"))
			assert.Equal(t, true, NeedValue[bool]("TEST_BOOL_TRUE"))
			assert.Equal(t, true, NeedValue[bool]("TEST_BOOL_ONE"))

			assert.Equal(t, true, Value("TEST_EMPTY", true))
			assert.Equal(t, false, Value("TEST_STRING", false))
			assert.Equal(t, false, Value("TEST_BOOL_FALSE", true))
			assert.Equal(t, true, Value("TEST_BOOL_TRUE", false))
			assert.Equal(t, true, Value("TEST_BOOL_ONE", false))
		},
	)

	t.Run(
		"Array",
		func(t *testing.T) {
			assert.PanicsWithError(
				t,
				`error during parse array environment: TEST_STRING: non-bool expression: "1.0.0"`,
				func() { _ = NeedArray[bool]("TEST_STRING") },
			)
			assert.PanicsWithError(
				t,
				`error during parse array environment: TEST_ARRAY_NOT_BOOL: non-bool expression: "3"`,
				func() { _ = NeedArray[bool]("TEST_ARRAY_NOT_BOOL") },
			)
			assert.Equal(t, []bool{NeedValue[bool]("TEST_ARRAY_ONE_ELEMENT")}, NeedArray[bool]("TEST_ARRAY_ONE_ELEMENT"))
			assert.Equal(t, []bool{true, false, true, false}, NeedArray[bool]("TEST_ARRAY"))

			assert.Equal(t, []bool{true, false}, Array("TEST_STRING", []bool{true, false}))
			assert.Equal(t, []bool{false, true, false}, Array("TEST_ARRAY_NOT_BOOL", []bool{false, true, false}))
			assert.Equal(t, []bool{NeedValue[bool]("TEST_ARRAY_ONE_ELEMENT")}, Array("TEST_ARRAY_ONE_ELEMENT", []bool{true, true, false}))
			assert.Equal(t, []bool{true, false, true, false}, Array("TEST_ARRAY", []bool{false, true, false, true}))
		},
	)
}
