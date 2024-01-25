package http_server

func NewError(code int, message string) Error {
	return Error{message, code}
}

type Error struct {
	message string
	Code    int
}

func (e Error) Error() string {
	return e.message
}
