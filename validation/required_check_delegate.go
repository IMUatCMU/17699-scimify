package validation

import (
	"reflect"
	"github.com/go-scim/scimify/resource"
	"github.com/go-scim/scimify/helper"
	"fmt"
)

type requiredCheckDelegate struct {
	enforceReadOnlyAttributes	bool	// whether enforce required rules on readOnly attributes
}

func (rcd *requiredCheckDelegate) OnInvalidValue(rts *helper.ResourceTraversalState, v reflect.Value, attr *resource.Attribute) (abort bool) {
	if attr.Required {
		if attr.Mutability != resource.ReadOnly || rcd.enforceReadOnlyAttributes {
			rcd.error(&RequiredMissingError{attr})
		}
	}

	abort = true
	return
}

func (rcd *requiredCheckDelegate) OnValueIsInterface(rts *helper.ResourceTraversalState, v reflect.Value, attr *resource.Attribute) (abort bool) {
	abort = false
	return
}

func (rcd *requiredCheckDelegate) OnValueIsSlice(rts *helper.ResourceTraversalState, v reflect.Value, attr *resource.Attribute) (abort bool) {
	if v.Len() == 0 && attr.Required {
		if attr.Mutability != resource.ReadOnly || rcd.enforceReadOnlyAttributes {
			rcd.error(&RequiredUnassignedError{attr})
			abort = true
			return
		}
	}

	abort = false
	return
}

func (rcd *requiredCheckDelegate) OnValueIsArray(rts *helper.ResourceTraversalState, v reflect.Value, attr *resource.Attribute) (abort bool) {
	abort = rcd.OnValueIsSlice(rts, v, attr)
	return
}

func (rcd *requiredCheckDelegate) OnValueIsMap(rts *helper.ResourceTraversalState, v reflect.Value, attr *resource.Attribute) (abort bool) {
	if v.Len() == 0 && attr.Required {
		if attr.Mutability != resource.ReadOnly || rcd.enforceReadOnlyAttributes {
			rcd.error(&RequiredUnassignedError{attr})
			abort = true
			return
		}
	}

	abort = false
	return
}

func (rcd *requiredCheckDelegate) OnMapKeyIsNotString(rts *helper.ResourceTraversalState, v reflect.Value, attr *resource.Attribute) (abort bool) {
	abort = false
	return
}

func (rcd *requiredCheckDelegate) OnValueIsBool(rts *helper.ResourceTraversalState, v reflect.Value, attr *resource.Attribute) (abort bool) {
	abort = false
	return
}

func (rcd *requiredCheckDelegate) OnValueIsInt(rts *helper.ResourceTraversalState, v reflect.Value, attr *resource.Attribute) (abort bool) {
	abort = false
	return
}

func (rcd *requiredCheckDelegate) OnValueIsFloat(rts *helper.ResourceTraversalState, v reflect.Value, attr *resource.Attribute) (abort bool) {
	abort = false
	return
}

func (rcd *requiredCheckDelegate) OnValueIsString(rts *helper.ResourceTraversalState, v reflect.Value, attr *resource.Attribute) (abort bool) {
	if v.Len() == 0 && attr.Required {
		if attr.Mutability != resource.ReadOnly || rcd.enforceReadOnlyAttributes {
			rcd.error(&RequiredUnassignedError{attr})
			abort = true
			return
		}
	}

	abort = false
	return
}

func (rcd *requiredCheckDelegate) OnUnsupportedType(rts *helper.ResourceTraversalState, v reflect.Value, attr *resource.Attribute) (abort bool) {
	abort = false
	return
}

func (rcd *requiredCheckDelegate) error(err error) {
	panic(err)
}

type RequiredMissingError struct {
	Attr 	*resource.Attribute
}

func (rme *RequiredMissingError) Error() string {
	return fmt.Sprintf("Missing required attribute %s", rme.Attr.Assist.FullPath)
}

type RequiredUnassignedError struct {
	Attr 	*resource.Attribute
}

func (rue *RequiredUnassignedError) Error() string {
	return fmt.Sprintf("Attribute %s is unassigned", rue.Attr.Assist.FullPath)
}