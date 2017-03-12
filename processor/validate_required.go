package processor

import (
	"context"
	"github.com/go-scim/scimify/helper"
	"github.com/go-scim/scimify/resource"
	"reflect"
)

type requiredValidationProcessor struct{}

func (rvp *requiredValidationProcessor) Process(r *resource.Resource, ctx context.Context) (err error) {
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

	schema := getSchema(ctx, true)
	delegate := &requiredCheckDelegate{enforceReadOnlyAttributes: false} // TODO turn into configuration option

	helper.TraverseWithSchema(r, schema, []helper.ResourceTraversalDelegate{delegate})

	err = nil
	return
}

type requiredCheckDelegate struct {
	// "true" means validation will fail when "read only" attributes are missing
	// recommended to be relaxed to "false" as "read only copy" will cover this anyways
	enforceReadOnlyAttributes bool
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
