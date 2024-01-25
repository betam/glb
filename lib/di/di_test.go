package di

import (
	"fmt"
	"reflect"
	"runtime"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/betam/glb/lib/try"
)

type WireForInterface interface {
	Dep() WireInterface
}

type WireInterface interface{}
type WireInterfaceAlias interface{}

type wireInterfaceOne struct {
	WireInterface
	WireInterfaceAlias
}
type wireInterfaceTwo struct {
	WireInterface
	WireInterfaceAlias
}

func NewWireInterfaceOne() *wireInterfaceOne { return &wireInterfaceOne{} }
func NewWireInterfaceTwo() *wireInterfaceTwo { return &wireInterfaceTwo{} }

type wireForInterface struct {
	dep WireInterface
}

func (w *wireForInterface) Dep() WireInterface {
	return w.dep
}

func NewWireForInterface(wireInterface WireInterface) *wireForInterface {
	return &wireForInterface{dep: wireInterface}
}

func TestWire(t *testing.T) {
	t.Run("WrongConstructor", func(t *testing.T) {
		container = make(map[string]*ioc)
		instances = make(map[string]*ioc)
		tags = make(map[any][]*tag)
		assert.NotPanics(t, func() {
			try.Catch(
				func() {
					Wire[WireInterface](new(WireInterface))
					panic(fmt.Errorf("unexpected"))
				},
				func(throwable error) {
					assert.ErrorIs(t, throwable, ErrConstructorNotFunction)
				},
			)
		})
		assert.Len(t, container, 0)
		assert.Len(t, instances, 0)
		assert.Len(t, tags, 0)
	})

	t.Run("Wire", func(t *testing.T) {
		container = make(map[string]*ioc)
		instances = make(map[string]*ioc)
		tags = make(map[any][]*tag)
		Wire[WireInterface](NewWireInterfaceOne)
		keyIface := encrypt(name(reflect.TypeOf(new(WireInterface)).Elem()), "")
		keyConstruct := name(reflect.TypeOf(wireInterfaceOne{}))
		assert.Len(t, container, 1)
		assert.Len(t, instances, 1)
		assert.Len(t, tags, 0)
		assert.Contains(t, container, keyIface)
		assert.Contains(t, instances, keyConstruct)
		assert.Nil(t, container[keyIface].instance)
		assert.Nil(t, container[keyIface].defaults)
		assert.Equal(t, runtime.FuncForPC(reflect.ValueOf(NewWireInterfaceOne).Pointer()).Name(), runtime.FuncForPC(reflect.ValueOf(container[keyIface].constructor).Pointer()).Name())
		assert.Equal(t, container[keyIface], instances[keyConstruct])
	})
	t.Run("For", func(t *testing.T) {
		container = make(map[string]*ioc)
		instances = make(map[string]*ioc)
		tags = make(map[any][]*tag)
		Wire[WireInterface](NewWireInterfaceOne)
		Wire[WireInterface](NewWireInterfaceTwo, For[WireForInterface]())
		keyIface := encrypt(name(reflect.TypeOf(new(WireInterface)).Elem()), name(reflect.TypeOf(new(WireForInterface)).Elem()))
		keyConstruct := name(reflect.TypeOf(wireInterfaceTwo{}))
		assert.Len(t, container, 2)
		assert.Len(t, instances, 2)
		assert.Len(t, tags, 0)
		assert.Contains(t, container, keyIface)
		assert.Contains(t, instances, keyConstruct)
		assert.Nil(t, container[keyIface].instance)
		assert.Nil(t, container[keyIface].defaults)
		assert.Equal(t, runtime.FuncForPC(reflect.ValueOf(NewWireInterfaceTwo).Pointer()).Name(), runtime.FuncForPC(reflect.ValueOf(container[keyIface].constructor).Pointer()).Name())
		assert.Equal(t, container[keyIface], instances[keyConstruct])
	})
}

func TestDefine(t *testing.T) {
	t.Run("WrongConstructor", func(t *testing.T) {
		container = make(map[string]*ioc)
		instances = make(map[string]*ioc)
		tags = make(map[any][]*tag)
		assert.NotPanics(t, func() {
			try.Catch(
				func() {
					Define(new(WireInterface))
					panic(fmt.Errorf("unexpected"))
				},
				func(throwable error) {
					assert.ErrorIs(t, throwable, ErrConstructorNotFunction)
				},
			)
		})
		assert.Len(t, container, 0)
		assert.Len(t, instances, 0)
		assert.Len(t, tags, 0)
	})
	t.Run("Simple", func(t *testing.T) {
		container = make(map[string]*ioc)
		instances = make(map[string]*ioc)
		tags = make(map[any][]*tag)
		instanceKey := Define(NewWireInterfaceOne, Alias[WireInterface]())
		keyIface := encrypt(name(reflect.TypeOf(new(WireInterface)).Elem()), "")
		keyConstruct := name(reflect.TypeOf(wireInterfaceOne{}))
		assert.Len(t, container, 1)
		assert.Len(t, instances, 1)
		assert.Len(t, tags, 0)
		assert.Contains(t, container, keyIface)
		assert.Contains(t, instances, keyConstruct)
		assert.Nil(t, container[keyIface].instance)
		assert.Nil(t, container[keyIface].defaults)
		assert.Equal(t, runtime.FuncForPC(reflect.ValueOf(NewWireInterfaceOne).Pointer()).Name(), runtime.FuncForPC(reflect.ValueOf(container[keyIface].constructor).Pointer()).Name())
		assert.Equal(t, container[keyIface], instances[keyConstruct])
		assert.Equal(t, keyConstruct, instanceKey)
	})

	t.Run("DefineByWrongKey", func(t *testing.T) {
		container = make(map[string]*ioc)
		instances = make(map[string]*ioc)
		tags = make(map[any][]*tag)
		instanceKey := Define(NewWireInterfaceOne)
		assert.NotPanics(t, func() {
			try.Catch(
				func() {
					Define(instanceKey+":::new", Alias[WireInterfaceAlias]())
					panic(fmt.Errorf("unexpected"))
				},
				func(throwable error) {
					assert.ErrorIs(t, throwable, ErrKeyNotFound)
				},
			)
		})
		assert.Len(t, container, 0)
		assert.Len(t, instances, 1)
		assert.Len(t, tags, 0)
	})

	t.Run("DefineByKey", func(t *testing.T) {
		container = make(map[string]*ioc)
		instances = make(map[string]*ioc)
		tags = make(map[any][]*tag)
		instanceKey := Define(NewWireInterfaceOne)
		Define(instanceKey, Alias[WireInterfaceAlias]())
		keyIface := encrypt(name(reflect.TypeOf(new(WireInterfaceAlias)).Elem()), "")
		keyConstruct := name(reflect.TypeOf(wireInterfaceOne{}))
		assert.Len(t, container, 1)
		assert.Len(t, instances, 1)
		assert.Len(t, tags, 0)
		assert.Contains(t, container, keyIface)
		assert.Contains(t, instances, keyConstruct)
		assert.Nil(t, container[keyIface].instance)
		assert.Nil(t, container[keyIface].defaults)
		assert.Equal(t, runtime.FuncForPC(reflect.ValueOf(NewWireInterfaceOne).Pointer()).Name(), runtime.FuncForPC(reflect.ValueOf(container[keyIface].constructor).Pointer()).Name())
		assert.Equal(t, container[keyIface], instances[keyConstruct])
		assert.Equal(t, keyConstruct, instanceKey)
	})
}

type NewInterface interface{}
type newInterface struct{ NewInterface }

func NewInterfaceConstructor() *newInterface {
	return &newInterface{}
}

func TestNew(t *testing.T) {
	t.Run("NonWired", func(t *testing.T) {
		container = make(map[string]*ioc)
		instances = make(map[string]*ioc)
		tags = make(map[any][]*tag)
		assert.NotPanics(t, func() {
			try.Catch(
				func() {
					New[NewInterface]()
					panic(fmt.Errorf("unexpected"))
				},
				func(throwable error) {
					assert.ErrorIs(t, throwable, ErrNotWired)
				},
			)
		})
	})

	t.Run("New", func(t *testing.T) {
		container = make(map[string]*ioc)
		instances = make(map[string]*ioc)
		tags = make(map[any][]*tag)
		Wire[NewInterface](NewInterfaceConstructor)
		i := New[NewInterface]()
		assert.Implements(t, (*NewInterface)(nil), i)
	})

	t.Run("Define", func(t *testing.T) {
		container = make(map[string]*ioc)
		instances = make(map[string]*ioc)
		tags = make(map[any][]*tag)
		Define(NewInterfaceConstructor, Alias[NewInterface]())
		i := New[NewInterface]()
		assert.Implements(t, (*NewInterface)(nil), i)
	})

	t.Run("For", func(t *testing.T) {
		Wire[WireInterface](func() WireInterface { return &wireInterfaceTwo{} }, For[WireForInterface]())
		Wire[WireForInterface](NewWireForInterface)
		f := New[WireForInterface]()
		assert.Implements(t, (*WireForInterface)(nil), f)
		assert.Implements(t, (*WireInterface)(nil), f.Dep())
		_, casted := f.Dep().(*wireInterfaceTwo)
		assert.True(t, casted)
	})
}

type CircularOneInterface interface{}
type circularOneInterface struct {
	CircularOneInterface
	two CircularTwoInterface
}

func NewCircularOne(twoInterface CircularTwoInterface) *circularOneInterface {
	return &circularOneInterface{two: twoInterface}
}

type CircularTwoInterface interface{}
type circularTwoInterface struct {
	CircularTwoInterface
	three CircularThreeInterface
}

func NewCircularTwo(threeInterface CircularThreeInterface) *circularTwoInterface {
	return &circularTwoInterface{three: threeInterface}
}

type CircularThreeInterface interface{}
type circularThreeInterface struct {
	CircularThreeInterface
	one CircularOneInterface
}

func NewCircularThree(oneInterface CircularOneInterface) *circularThreeInterface {
	return &circularThreeInterface{one: oneInterface}
}

func TestCircular(t *testing.T) {
	Wire[CircularOneInterface](NewCircularOne)
	Wire[CircularTwoInterface](NewCircularTwo)
	Wire[CircularThreeInterface](NewCircularThree)
	assert.NotPanics(t, func() {
		try.Catch(
			func() {
				New[CircularOneInterface]()
				panic(fmt.Errorf("unexpected"))
			},
			func(throwable error) {
				assert.ErrorIs(t, throwable, ErrCircularDependencies)
			},
		)
	})
}
