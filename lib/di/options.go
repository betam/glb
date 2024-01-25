package di

import "reflect"

type ConstructorOptions interface {
	constructor()
}

type WireOptions interface {
	wire()
}

type DefineOptions interface {
	define()
}

func Defaults(v map[int]any) *defaultOptions {
	return &defaultOptions{v: v}
}

type defaultOptions struct {
	WireOptions
	DefineOptions
	ConstructorOptions
	v map[int]any
}

func For[Interface any]() *forOptions {
	i := reflect.TypeOf(new(Interface)).Elem()
	return &forOptions{t: name(i)}
}

type forOptions struct {
	WireOptions
	DefineOptions
	ConstructorOptions
	t string
}

func Alias[Interface any]() *aliasOptions {
	return &aliasOptions{t: reflect.TypeOf(new(Interface)).Elem()}
}

type aliasOptions struct {
	DefineOptions
	ConstructorOptions
	t reflect.Type
}

func Tag(tag any) *tagOptions {
	return &tagOptions{t: tag}
}

type tagOptions struct {
	WireOptions
	ConstructorOptions
	t any
}

func Fallback() *fallbackOptions {
	return &fallbackOptions{}
}

type fallbackOptions struct {
	WireOptions
	ConstructorOptions
}
