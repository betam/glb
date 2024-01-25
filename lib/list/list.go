package list

func Shift[Element any](list *[]Element) *Element {
	if list == nil {
		return nil
	}
	if len(*list) == 0 {
		return nil
	}
	value := (*list)[0]
	if len(*list) > 1 {
		*list = (*list)[1:]
	} else {
		list = &[]Element{}
	}
	return &value
}

func Pop[Element any](list *[]Element) *Element {
	if list == nil {
		return nil
	}
	if len(*list) == 0 {
		return nil
	}
	size := len(*list)
	value := (*list)[size-1]
	if len(*list) > 1 {
		*list = (*list)[:size-1]
	} else {
		list = &[]Element{}
	}
	return &value
}

func Push[Element any](list *[]Element, item Element) *[]Element {
	*list = append(*list, item)
	return list
}

func Unshift[Element any](list *[]Element, item Element) *[]Element {
	l := make([]Element, 0, len(*list)+1)
	l = append(l, item)
	l = append(l, *list...)
	*list = l
	return list
}

func Foreach[Element any](list []Element, callback func(item Element)) {
	for _, value := range list {
		callback(value)
	}
}

func Map[Element any, Result any](list []Element, callback func(item Element) Result) []Result {
	e := make([]Result, 0, len(list))
	for _, item := range list {
		e = append(e, callback(item))
	}

	return e
}

func Filter[Element any](list []Element, callback func(item Element) bool) []Element {
	result := make([]Element, 0, len(list))
	for _, item := range list {
		if callback(item) {
			result = append(result, item)
		}
	}

	return result
}

func Reduce[Element, Result any](list []Element, callback func(result Result, item Element) Result, result Result) Result {
	for _, item := range list {
		result = callback(result, item)
	}

	return result
}

func Keys[Key comparable, Value any](m map[Key]Value) []Key {
	keys := make([]Key, 0, len(m))
	for key := range m {
		keys = append(keys, key)
	}
	return keys
}

func Values[Key comparable, Value any](m map[Key]Value) []Value {
	values := make([]Value, 0, len(m))
	for _, value := range m {
		values = append(values, value)
	}
	return values
}

func KeysValues[Key comparable, Value any](m map[Key]Value) ([]Key, []Value) {
	keys := make([]Key, 0, len(m))
	values := make([]Value, 0, len(m))
	for key, value := range m {
		keys = append(keys, key)
		values = append(values, value)
	}
	return keys, values
}

func Contains[Element comparable](item Element, array []Element) bool {
	for _, el := range array {
		if item == el {
			return true
		}
	}
	return false
}
