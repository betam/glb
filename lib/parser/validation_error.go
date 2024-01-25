package parser

type ValidationError struct {
	m string
}

func (v ValidationError) Error() string {
	return v.m
}
