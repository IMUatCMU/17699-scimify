package helper

import (
	"github.com/go-scim/scimify/adt"
	"github.com/go-scim/scimify/resource"
	"reflect"
)

type TraverseDelegate interface {
	OnValueIsInValid(ts *TraverseState, v reflect.Value) (abort bool)

	OnValueIsInterface(ts *TraverseState, v reflect.Value) (abort bool)

	OnValueIsArrayOrSlice(ts *TraverseState, v reflect.Value) (abort bool)

	OnValueIsBool(ts *TraverseState, v reflect.Value) (abort bool)

	OnValueIsInteger(ts *TraverseState, v reflect.Value) (abort bool)

	OnValueIsFloat(ts *TraverseState, v reflect.Value) (abort bool)

	OnValueIsString(ts *TraverseState, v reflect.Value) (abort bool)

	OnValueIsMap(ts *TraverseState, v reflect.Value) (abort bool)

	OnValueIsUnsupported(ts *TraverseState, v reflect.Value) (abort bool)
}

func Traverse(r *resource.Resource, sch *resource.Schema, dlg TraverseDelegate) {
	state := &TraverseState{
		Resource:      r,
		RootSchema:    sch,
		ContainerKeys: adt.NewStackWithoutLimit(),
		ContainerVal:  adt.NewStackWithoutLimit(),
		delegate:      dlg,
	}
	state.traverseWithReflection(reflect.ValueOf(r.Data()))
}

type TraverseState struct {
	Resource      *resource.Resource
	RootSchema    *resource.Schema
	ContainerKeys adt.Stack
	ContainerVal  adt.Stack
	delegate      TraverseDelegate
}

func (ts *TraverseState) pushKey(k reflect.Value) {
	var keyStack adt.Stack
	if top, ok := ts.ContainerKeys.Peek().(adt.Stack); !ok {
		keyStack = adt.NewStackWithoutLimit()
	} else {
		keyStack = top.Clone()
	}
	keyStack.Push(k)
	ts.ContainerKeys.Push(keyStack)
}

func (ts *TraverseState) popKey() {
	ts.ContainerKeys.Pop()
}

func (ts *TraverseState) traverseWithReflection(v reflect.Value) {
	if !v.IsValid() {
		if abort := ts.delegate.OnValueIsInValid(ts, v); abort {
			return
		}
	}

	if v.Kind() == reflect.Interface {
		if abort := ts.delegate.OnValueIsInterface(ts, v); abort {
			return
		}

		ts.traverseWithReflection(v.Elem())
		return
	}

	switch v.Kind() {
	case reflect.Bool:
		if abort := ts.delegate.OnValueIsBool(ts, v); abort {
			return
		}
	case reflect.String:
		if abort := ts.delegate.OnValueIsString(ts, v); abort {
			return
		}
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		if abort := ts.delegate.OnValueIsInteger(ts, v); abort {
			return
		}
	case reflect.Float32, reflect.Float64:
		if abort := ts.delegate.OnValueIsFloat(ts, v); abort {
			return
		}
	case reflect.Array, reflect.Slice:
		if abort := ts.delegate.OnValueIsArrayOrSlice(ts, v); abort {
			return
		}
		for i := 0; i < v.Len(); i++ {
			ts.traverseWithReflection(v.Index(i))
		}
	case reflect.Map:
		if abort := ts.delegate.OnValueIsMap(ts, v); abort {
			return
		}

		ts.ContainerVal.Push(v)
		for _, k := range v.MapKeys() {
			ts.pushKey(k)
			ts.traverseWithReflection(v.MapIndex(k))
			ts.popKey()
		}
		ts.ContainerVal.Pop()
	default:
		if abort := ts.delegate.OnValueIsUnsupported(ts, v); abort {
			return
		}
	}
}
