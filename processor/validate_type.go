package processor

import (
	"github.com/go-scim/scimify/helper"
	"github.com/go-scim/scimify/resource"
	"reflect"
	"sync"
	"time"
)

var (
	oneTypeValidation sync.Once
	typeValidator     Processor
)

func TypeValidationProcessor() Processor {
	oneTypeValidation.Do(func() {
		typeValidator = &typeValidationProcessor{}
	})
	return typeValidator
}

type typeValidationProcessor struct{}

func (tvp *typeValidationProcessor) Process(ctx *ProcessorContext) (err error) {
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

	r := tvp.getResource(ctx)
	schema := tvp.getSchema(ctx)
	delegate := &typeCheckDelegate{}

	helper.TraverseWithSchema(r, schema, []helper.ResourceTraversalDelegate{delegate})

	err = nil
	return
}

func (tvp *typeValidationProcessor) getResource(ctx *ProcessorContext) *resource.Resource {
	if ctx.Resource == nil {
		panic(&MissingContextValueError{"resource"})
	}
	return ctx.Resource
}

func (tvp *typeValidationProcessor) getSchema(ctx *ProcessorContext) *resource.Schema {
	if ctx.Schema == nil {
		panic(&MissingContextValueError{"schema"})
	}
	return ctx.Schema
}

type typeCheckDelegate struct{}

func (tcd *typeCheckDelegate) OnInvalidValue(rts *helper.ResourceTraversalState, v reflect.Value, attr *resource.Attribute) (abort bool) {
	abort = true // type checking does not care if something is nil, just abort
	return
}

func (tcd *typeCheckDelegate) OnValueIsInterface(rts *helper.ResourceTraversalState, v reflect.Value, attr *resource.Attribute) (abort bool) {
	abort = false // continue type check for interface values
	return
}

func (tcd *typeCheckDelegate) OnValueIsSlice(rts *helper.ResourceTraversalState, v reflect.Value, attr *resource.Attribute) (abort bool) {
	if attr.MultiValued {
		abort = false
	} else {
		abort = true
		tcd.error(&TypeMismatchError{v.Type(), attr})
	}
	return
}

func (tcd *typeCheckDelegate) OnValueIsArray(rts *helper.ResourceTraversalState, v reflect.Value, attr *resource.Attribute) (abort bool) {
	abort = tcd.OnValueIsSlice(rts, v, attr)
	return
}

func (tcd *typeCheckDelegate) OnValueIsMap(rts *helper.ResourceTraversalState, v reflect.Value, attr *resource.Attribute) (abort bool) {
	if !attr.MultiValued && type_complex == attr.Type {
		abort = false
	} else {
		abort = true
		tcd.error(&TypeMismatchError{v.Type(), attr})
	}
	return
}

func (tcd *typeCheckDelegate) OnMapKeyIsNotString(rts *helper.ResourceTraversalState, v reflect.Value, attr *resource.Attribute) (abort bool) {
	abort = true
	tcd.error(&TypeMismatchError{v.Type(), &resource.Attribute{Type: type_string}})
	return
}

func (tcd *typeCheckDelegate) OnValueIsBool(rts *helper.ResourceTraversalState, v reflect.Value, attr *resource.Attribute) (abort bool) {
	if !attr.MultiValued && type_bool == attr.Type {
		abort = false
	} else {
		abort = true
		tcd.error(&TypeMismatchError{v.Type(), attr})
	}
	return
}

func (tcd *typeCheckDelegate) OnValueIsInt(rts *helper.ResourceTraversalState, v reflect.Value, attr *resource.Attribute) (abort bool) {
	if !attr.MultiValued && type_int == attr.Type {
		abort = false
	} else {
		abort = true
		tcd.error(&TypeMismatchError{v.Type(), attr})
	}
	return
}

func (tcd *typeCheckDelegate) OnValueIsFloat(rts *helper.ResourceTraversalState, v reflect.Value, attr *resource.Attribute) (abort bool) {
	if !attr.MultiValued && type_float == attr.Type {
		abort = false
	} else {
		abort = true
		tcd.error(&TypeMismatchError{v.Type(), attr})
	}
	return
}

func (tcd *typeCheckDelegate) OnValueIsString(rts *helper.ResourceTraversalState, v reflect.Value, attr *resource.Attribute) (abort bool) {
	if attr.MultiValued {
		abort = true
		tcd.error(&TypeMismatchError{v.Type(), attr})
		return
	}

	switch attr.Type {
	case type_string, type_ref, type_binary:
		abort = false
	case type_datetime:
		datetimeFormat := "2006-01-02T15:04:05Z"
		if _, err := time.Parse(datetimeFormat, v.Interface().(string)); err != nil {
			abort = true
			tcd.error(&FormatError{attr, datetimeFormat, v.Interface()})
		} else {
			abort = false
		}
	default:
		abort = true
		tcd.error(&TypeMismatchError{v.Type(), attr})
	}
	return
}

func (tcd *typeCheckDelegate) OnUnsupportedType(rts *helper.ResourceTraversalState, v reflect.Value, attr *resource.Attribute) (abort bool) {
	abort = true
	tcd.error(&TypeUnsupportedError{v.Type()})
	return
}

func (tcd *typeCheckDelegate) error(err error) {
	panic(err)
}
