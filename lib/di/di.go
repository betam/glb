package di

import (
	"crypto/sha1"
	"encoding/hex"
	"fmt"
	"io"
	"reflect"

	"github.com/betam/glb/lib/list"
	"github.com/betam/glb/lib/pointer"
	"github.com/betam/glb/lib/try"
)

var (
	ErrConstructorNotFunction = fmt.Errorf("unsupported type for constructor")
	ErrDefaultNotFound        = fmt.Errorf("default property is not found")
	ErrDefaultForWired        = fmt.Errorf("got default value for wired interface")
	ErrMismatchedTypes        = fmt.Errorf("expected and given types are different")
	ErrNotWired               = fmt.Errorf("missed wired constructor")
	ErrKeyNotFound            = fmt.Errorf("cannot find instance by key")
	ErrCircularDependencies   = fmt.Errorf("circular dependencies")
	ErrCloseFailed            = fmt.Errorf("fail to stop a service")
)

var (
	container = make(map[string]*ioc)
	instances = make(map[string]*ioc)
	tags      = make(map[any][]*tag)
)

type ioc struct {
	instance     *reflect.Value
	constructor  any
	defaults     map[int]any
	dependencies []*ioc
}

type tag struct {
	t   reflect.Type
	ioc *ioc
}

func Define(constructor any, options ...DefineOptions) string {
	opts := list.Map(options, func(item DefineOptions) ConstructorOptions {
		return item.(ConstructorOptions)
	})
	return wire(constructor, opts...)
}

func Wire[Interface any](constructor any, options ...WireOptions) {
	opts := list.Map(options, func(item WireOptions) ConstructorOptions {
		return item.(ConstructorOptions)
	})
	opts = append(opts, &aliasOptions{t: reflect.TypeOf(new(Interface)).Elem()})
	wire(constructor, opts...)
}

func wire(constructor any, options ...ConstructorOptions) string {
	f := reflect.TypeOf(constructor)
	if f.Kind() != reflect.Func && f.Kind() != reflect.String {
		panic(fmt.Errorf("%w: '%s'", ErrConstructorNotFunction, f.String()))
	}
	var def map[int]any
	var target string
	var aliases []reflect.Type
	var tagList []any
	var fallbackOnly bool
	for _, o := range options {
		switch option := o.(type) {
		case *defaultOptions:
			def = verifyDefaults(f, option.v)
		case *forOptions:
			target = option.t
		case *aliasOptions:
			aliases = append(aliases, option.t)
		case *tagOptions:
			tagList = append(tagList, option.t)
		case *fallbackOptions:
			fallbackOnly = true
		}
	}
	var IoC *ioc
	var instanceKey string
	if f.Kind() == reflect.String {
		instanceKey = constructor.(string)
		if IoC = instances[instanceKey]; IoC == nil {
			panic(fmt.Errorf("%w: '%s'", ErrKeyNotFound, constructor))
		}
	} else {
		IoC = &ioc{
			constructor: constructor,
			defaults:    def,
		}
	}
	for _, alias := range aliases {
		hash := encrypt(name(alias), target)
		if _, found := container[hash]; !found || !fallbackOnly {
			container[hash] = IoC
		}
		for _, label := range tagList {
			singleTag := &tag{
				t:   alias,
				ioc: IoC,
			}
			if _, ok := tags[label]; !ok {
				tags[label] = []*tag{}
			}
			tags[label] = append(tags[label], singleTag)
		}
	}

	if instanceKey == "" && f.NumOut() > 0 {
		instanceKey = name(deref(f.Out(0)))
		instances[instanceKey] = IoC
	}
	return instanceKey
}

func verifyDefaults(f reflect.Type, v map[int]any) map[int]any {
	for idx := 0; idx < f.NumMethod(); idx++ {
		_, ok := v[idx]
		if ok && f.In(idx).Kind() == reflect.Interface {
			panic(fmt.Errorf("%w: '%d'", ErrDefaultForWired, idx))
		}
		if ok && f.In(idx).Kind() != reflect.TypeOf(v[idx]).Kind() {
			panic(fmt.Errorf("%w: '%s', '%s'", ErrMismatchedTypes, name(f.In(idx)), name(reflect.TypeOf(v[idx]))))
		}
		if !ok && f.In(idx).Kind() != reflect.Interface {
			panic(fmt.Errorf("%w: '%d'", ErrDefaultNotFound, idx))
		}
	}
	return v
}

func Tags[Interface any](label any) func() ([]Interface, []*ioc) {
	return func() (result []Interface, deps []*ioc) {
		if _, found := tags[label]; found {
			for _, singleTag := range tags[label] {
				IoC := instance(singleTag.ioc, singleTag.t)
				result = append(result, IoC.instance.Interface().(Interface))
				deps = append(deps, IoC)
			}
		}
		return result, deps
	}
}

func NewWithCloser[Interface any]() (Interface, func()) {
	i := reflect.TypeOf(new(Interface)).Elem()
	f := constructor(i)
	IoC := instance(f, i)

	var liner func(IoC *ioc) (result []*ioc)
	liner = func(IoC *ioc) (result []*ioc) {
		for _, dep := range IoC.dependencies {
			result = append(result, dep)
		}
		for _, dep := range IoC.dependencies {
			result = append(result, liner(dep)...)
		}
		return
	}

	done := map[*ioc]bool{}
	deps := append([]*ioc{IoC}, liner(IoC)...)
	for idx := len(deps) - 1; idx >= 0; idx-- {
		if _, closed := done[deps[idx]]; closed {
			deps[idx] = nil
		}
		done[deps[idx]] = true
	}

	closer := func(deps []*ioc) func() {
		return func() {
			for idx := range deps {
				if deps[idx] == nil {
					continue
				}
				if _, ok := deps[idx].instance.Interface().(io.Closer); ok {
					try.Catch(
						func() {
							deps[idx].instance.MethodByName("Close").Call([]reflect.Value{})
						},
						func(throwable error) {
							panic(fmt.Errorf("%w: %s", ErrCloseFailed, deps[idx].instance.String()))
						},
					)
				}
			}
		}
	}

	return IoC.instance.Interface().(Interface), closer(deps)
}

func New[Interface any]() Interface {
	i := reflect.TypeOf(new(Interface)).Elem()
	f := constructor(i)
	IoC := instance(f, i)
	return IoC.instance.Interface().(Interface)
}

func constructor(need reflect.Type, required ...reflect.Type) *ioc {
	if len(required) > 0 {
		hash := encrypt(name(need), name(required[0]))
		if IoC, found := container[hash]; found {
			return IoC
		}
	}

	hash := encrypt(name(need), "")
	if IoC, found := container[hash]; found {
		return IoC
	}

	panic(fmt.Errorf("%w: '%s'", ErrNotWired, name(need)))
}

func instance(IoC *ioc, need ...reflect.Type) *ioc {
	if IoC.instance != nil {
		return IoC
	}
	fType := reflect.TypeOf(IoC.constructor)

	var args []reflect.Value
	var dependencies []*ioc
	for idx := 0; idx < fType.NumIn(); idx++ {
		if field := deref(fType.In(idx)); field.Kind() == reflect.Interface {
			for _, dep := range need {
				if name(dep) == name(field) {
					panic(fmt.Errorf("%w: '%s', '%s'", ErrCircularDependencies, name(field), name(dep)))
				}
			}
			fConstructor := constructor(field, need...)
			fNeed := append([]reflect.Type{field}, need...)
			fIoC := instance(fConstructor, fNeed...)
			args = append(args, *fIoC.instance)
			dependencies = append(dependencies, fIoC)
		} else {
			arg := IoC.defaults[idx]
			if field.Kind() == reflect.Slice && field.Elem().Kind() == reflect.Interface && reflect.TypeOf(IoC.defaults[idx]).Kind() == reflect.Func {
				tagSlice := reflect.ValueOf(arg).Call([]reflect.Value{})
				arg = tagSlice[0].Interface()
				dependencies = append(dependencies, tagSlice[1].Interface().([]*ioc)...)
			}
			value := reflect.ValueOf(arg)
			if value.Kind() != reflect.Invalid {
				args = append(args, value)
			}
		}

	}
	if fType.NumIn() != len(args) {
		panic(fmt.Errorf("constructor for %s needs %d args, %d given", name(need[0]), fType.NumIn(), len(args)))
	}
	if result := reflect.ValueOf(IoC.constructor).Call(args); len(result) > 0 {
		IoC.instance = &result[0]
	} else {
		IoC.instance = pointer.Pointer(reflect.New(need[0]))
	}
	IoC.dependencies = dependencies
	return IoC

}

func deref(t reflect.Type) reflect.Type {
	for t.Kind() == reflect.Ptr {
		t = t.Elem()
	}
	return t
}

func encrypt(need, target string) string {
	hasher := sha1.New()
	hasher.Write([]byte(fmt.Sprintf("%s::%s", need, target)))
	return hex.EncodeToString(hasher.Sum(nil))
}

func name(i reflect.Type) string {
	return fmt.Sprintf("%s.%s", i.PkgPath(), i.Name())
}
