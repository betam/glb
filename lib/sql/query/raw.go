package query

type raw struct {
	expression string
}

func Raw(expression string) *raw {
	return &raw{expression}
}
