package validation

import (
	"fmt"
	"github.com/go-scim/scimify/helper"
	"github.com/go-scim/scimify/resource"
	"reflect"
	"time"
)

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

type TypeMismatchError struct {
	T    reflect.Type
	Attr *resource.Attribute
}

func (tm *TypeMismatchError) Error() string {
	var expects string = ""
	switch tm.Attr.Type {
	case resource.String, resource.Binary, resource.Reference, resource.DateTime:
		expects = "string"
	case resource.Integer:
		expects = "integer"
	case resource.Decimal:
		expects = "decimal"
	case resource.Boolean:
		expects = "boolean"
	case resource.Complex:
		expects = "complex"
	}
	if tm.Attr.MultiValued {
		expects += " array"
	}
	return "type check expected type: " + expects + ", unsupported type: " + tm.T.String()
}

type FormatError struct {
	Attr   *resource.Attribute
	Format string
	Actual interface{}
}

func (fe *FormatError) Error() string {
	return fmt.Sprintf("required format of [%s] at %s, but got %v", fe.Format, fe.Attr.Assist.FullPath, fe.Actual)
}

type TypeUnsupportedError struct {
	T reflect.Type
}

func (ute *TypeUnsupportedError) Error() string {
	return fmt.Sprintf("type %s is not supported", ute.T.String())
}
