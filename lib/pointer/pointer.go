package pointer

func Pointer[T any](value T) *T {
	return &value
}

func Value[T any](value *T) T {
	if value == nil {
		return *new(T)
	}
	return *value
}
