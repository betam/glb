package err

import (
	"fmt"
)

func CheckClass(throwable error, class string) *Err {
	if err, ok := throwable.(Err); ok {
		if err.Class == class {
			return &err
		}
	}
	return nil
}

func NewErrf(class string, format string, args ...interface{}) Err {
	return Err{
		fmt.Sprintf(format, args...),
		class,
	}
}

type Err struct {
	message string
	Class   string
}

func (e Err) Error() string {
	return e.message
}
