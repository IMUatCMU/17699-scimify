package processor

import (
	"github.com/go-scim/scimify/adt"
	"github.com/go-scim/scimify/helper"
	"github.com/go-scim/scimify/resource"
	"reflect"
	"strings"
)

type formatCaseProcessor struct{}

func (fcp *formatCaseProcessor) Process(ctx *ProcessorContext) (err error) {
	defer func() {
		if r := recover(); r != nil {
			switch r.(type) {
			case error:
				err = r.(error)
				return
			default:
				panic(r)
			}
		}
	}()

	r := fcp.getResource(ctx)
	schema := fcp.getSchema(ctx)
	delegate := &formatCaseDelegate{}

	helper.Traverse(r, schema, delegate)

	err = nil
	return
}

func (fcp *formatCaseProcessor) getResource(ctx *ProcessorContext) *resource.Resource {
	if ctx.Resource == nil {
		panic(&MissingContextValueError{"resource"})
	}
	return ctx.Resource
}

func (fcp *formatCaseProcessor) getSchema(ctx *ProcessorContext) *resource.Schema {
	if ctx.Schema == nil {
		panic(&MissingContextValueError{"schema"})
	}
	return ctx.Schema
}

type formatCaseDelegate struct{}

func (fcd *formatCaseDelegate) formatCase(ts *helper.TraverseState, v reflect.Value) {
	var (
		key       string
		path      string
		container reflect.Value
	)
	key = fcd.resolveCurrentKey(ts)
	path = fcd.resolveCurrentPath(ts)
	if ts.ContainerVal.Size() > 0 {
		container = ts.ContainerVal.Peek().(reflect.Value)
	}

	if len(path) > 0 {
		attr := ts.RootSchema.GetAttribute(path)
		if nil == attr {
			fcd.error(&NoDefinedAttributeError{path})
		}

		// key is the entry for the current map entry.
		// if we are on the top level, the key could (case insensitively) be the JSONName or FullPath (which is JSONName prepended with URN)
		// if we are below top level, the key should just (case insensitively) match the JSONName
		// hence, if the key fails to (case insensitively) match any of these two, we error out
		switch strings.ToLower(key) {
		case strings.ToLower(attr.Assist.JSONName):
			if key != attr.Assist.JSONName {
				keyVal := reflect.ValueOf(key)
				entryVal := container.MapIndex(keyVal)
				if entryVal.IsValid() {
					container.SetMapIndex(reflect.ValueOf(attr.Assist.JSONName), entryVal)
					container.SetMapIndex(keyVal, reflect.Value{})
				}
			}
		case strings.ToLower(attr.Assist.FullPath):
			if key != attr.Assist.FullPath {
				keyVal := reflect.ValueOf(key)
				entryVal := container.MapIndex(keyVal)
				container.SetMapIndex(reflect.ValueOf(attr.Assist.FullPath), entryVal)
				container.SetMapIndex(keyVal, reflect.Value{})
			}
		default:
			fcd.error(&AttributeMismatchWithKeyError{Key: key, Attr: attr})
		}
	}
}

func (fcd *formatCaseDelegate) OnValueIsInValid(ts *helper.TraverseState, v reflect.Value) (abort bool) {
	abort = true
	return
}

func (fcd *formatCaseDelegate) OnValueIsInterface(ts *helper.TraverseState, v reflect.Value) (abort bool) {
	abort = false
	return
}

func (fcd *formatCaseDelegate) OnValueIsArrayOrSlice(ts *helper.TraverseState, v reflect.Value) (abort bool) {
	fcd.formatCase(ts, v)
	abort = false
	return
}

func (fcd *formatCaseDelegate) OnValueIsBool(ts *helper.TraverseState, v reflect.Value) (abort bool) {
	fcd.formatCase(ts, v)
	abort = true
	return
}

func (fcd *formatCaseDelegate) OnValueIsInteger(ts *helper.TraverseState, v reflect.Value) (abort bool) {
	fcd.formatCase(ts, v)
	abort = true
	return
}

func (fcd *formatCaseDelegate) OnValueIsFloat(ts *helper.TraverseState, v reflect.Value) (abort bool) {
	fcd.formatCase(ts, v)
	abort = true
	return
}

func (fcd *formatCaseDelegate) OnValueIsString(ts *helper.TraverseState, v reflect.Value) (abort bool) {
	fcd.formatCase(ts, v)
	abort = true
	return
}

func (fcd *formatCaseDelegate) OnValueIsMap(ts *helper.TraverseState, v reflect.Value) (abort bool) {
	fcd.formatCase(ts, v)
	abort = false
	return
}

func (fcd *formatCaseDelegate) OnValueIsUnsupported(ts *helper.TraverseState, v reflect.Value) (abort bool) {
	fcd.error(&TypeUnsupportedError{T: v.Type()})
	abort = true
	return
}

func (fcd *formatCaseDelegate) resolveCurrentKey(ts *helper.TraverseState) string {
	key := ""

	if ts.ContainerKeys.Size() > 0 {
		key = ts.ContainerKeys.Peek().(adt.Stack).Peek().(reflect.Value).String()
	}

	return key
}

func (fcd *formatCaseDelegate) resolveCurrentPath(ts *helper.TraverseState) string {
	key := ""

	if ts.ContainerKeys.Size() > 0 {
		keyStack := ts.ContainerKeys.Peek().(adt.Stack).Clone()
		reverseStack := adt.NewStackWithoutLimit()
		for keyStack.Size() != 0 {
			reverseStack.Push(keyStack.Pop())
		}
		for reverseStack.Size() != 0 {
			if len(key) != 0 {
				key += "."
			}
			key += reverseStack.Pop().(reflect.Value).String()
		}
	}

	return key
}

func (tcd *formatCaseDelegate) error(err error) {
	panic(err)
}
