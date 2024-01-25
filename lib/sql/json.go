package sql

import (
	"database/sql/driver"
	"encoding/json"
	"fmt"
)

type Json[T any] struct {
	value T
}

func NewJson[T any](value T) *Json[T] {
	return &Json[T]{value}
}

func (f *Json[T]) Value() (driver.Value, error) {
	var v T
	switch any(v).(type) {
	case int, int8, int16, int32, int64, uint, uint8, uint16, uint32, uint64, float32, float64, bool:
		return driver.DefaultParameterConverter.ConvertValue(f.value)
	default:
		return json.Marshal(f.value)
	}
}

func (f *Json[T]) Scan(value any) error {
	var v []byte
	switch value.(type) {
	case int, int8, int16, int32, int64, uint, uint8, uint16, uint32, uint64, float32, float64, bool:
		v = []byte(fmt.Sprintf("%v", value))
	default:
		v = value.([]byte)
	}
	return json.Unmarshal(v, &f.value)
}

func (f *Json[T]) MarshalJSON() ([]byte, error) {
	return json.Marshal(f.value)
}

func (f *Json[T]) UnmarshalJSON(value []byte) error {
	return json.Unmarshal(value, &f.value)
}

func (f *Json[T]) Unwrap() T {
	return f.value
}
