package try

func Catch(try func(), catch func(throwable error)) {
	defer func() {
		if exception := recover(); exception != nil {
			if throwable, ok := exception.(error); ok {
				if catch != nil {
					catch(throwable)
				}
				return
			}
			panic(exception)
		}
	}()

	try()
}

func Throw[T any](result T, err error) T {
	ThrowError(err)
	return result
}

func ThrowError(err error) {
	if err != nil {
		panic(err)
	}
}

/*type try struct {
	try     func()
	catch   func(Throwable)
	finally func()
}

func Try(t func()) *try {
	return &try{try: t}
}

func (t *try) Catch(catch func(throwable Throwable)) *try {
	t.catch = catch
	return t
}

func (t *try) Finally(finally func()) *try {
	t.finally = finally
	return t
}

func (t *try) Run() {
	if t.finally != nil {
		defer t.finally()
	}
	defer func() {
		if exception := recover(); exception != nil {
			if throwable, ok := exception.(Throwable); ok {
				if t.catch != nil {
					t.catch(throwable)
					return
				}
			}
			panic(exception)
		}
	}()

	t.try()
}
*/
