package helper

import (
	"fmt"
	"github.com/go-scim/scimify/adt"
	"github.com/go-scim/scimify/resource"
	"reflect"
)

type ResourceTraversalDelegate interface {
	OnInvalidValue(rts *ResourceTraversalState, v reflect.Value, attr *resource.Attribute) bool

	OnValueIsInterface(rts *ResourceTraversalState, v reflect.Value, attr *resource.Attribute) bool

	OnValueIsSlice(rts *ResourceTraversalState, v reflect.Value, attr *resource.Attribute) bool
	OnValueIsArray(rts *ResourceTraversalState, v reflect.Value, attr *resource.Attribute) bool

	OnValueIsMap(rts *ResourceTraversalState, v reflect.Value, attr *resource.Attribute) bool
	OnMapKeyIsNotString(rts *ResourceTraversalState, v reflect.Value, attr *resource.Attribute) bool

	OnValueIsBool(rts *ResourceTraversalState, v reflect.Value, attr *resource.Attribute) bool

	OnValueIsInt(rts *ResourceTraversalState, v reflect.Value, attr *resource.Attribute) bool

	OnValueIsFloat(rts *ResourceTraversalState, v reflect.Value, attr *resource.Attribute) bool

	OnValueIsString(rts *ResourceTraversalState, v reflect.Value, attr *resource.Attribute) bool

	OnUnsupportedType(rts *ResourceTraversalState, v reflect.Value, attr *resource.Attribute) bool
}

type ResourceTraversalState struct {
	RootSchema    *resource.Schema
	ContainerAttr adt.Stack
	ContainerVal  adt.Stack
	delegate      ResourceTraversalDelegate
}

func TraverseWithSchema(r *resource.Resource, sch *resource.Schema, dlg []ResourceTraversalDelegate) {
	state := &ResourceTraversalState{
		RootSchema:    sch,
		ContainerAttr: adt.NewStackWithoutLimit(),
		ContainerVal:  adt.NewStackWithoutLimit(),
		delegate:      &broadcastTraversalDelegate{delegates: dlg},
	}
	state.traverseWithReflection(reflect.ValueOf(r.Data()), sch.AsAttribute())
}

func (rts *ResourceTraversalState) traverseWithReflection(v reflect.Value, attr *resource.Attribute) {
	if !v.IsValid() {
		if abort := rts.delegate.OnInvalidValue(rts, v, attr); abort {
			return
		}
	}

	if v.Kind() == reflect.Interface {
		if abort := rts.delegate.OnValueIsInterface(rts, v, attr); abort {
			return
		}

		rts.traverseWithReflection(v.Elem(), attr)
		return
	}

	switch v.Kind() {
	case reflect.Slice, reflect.Array:
		if reflect.Slice == v.Kind() {
			if abort := rts.delegate.OnValueIsSlice(rts, v, attr); abort {
				return
			}
		} else if reflect.Array == v.Kind() {
			if abort := rts.delegate.OnValueIsArray(rts, v, attr); abort {
				return
			}
		}

		elemAttr := attr.Clone()
		elemAttr.MultiValued = false

		rts.ContainerAttr.Push(attr)
		rts.ContainerVal.Push(v)

		for i := 0; i < v.Len(); i++ {
			elem := v.Index(i)
			rts.traverseWithReflection(elem, elemAttr)
		}

		rts.ContainerAttr.Pop()
		rts.ContainerVal.Pop()

	case reflect.Map:
		if abort := rts.delegate.OnValueIsMap(rts, v, attr); abort {
			return
		}

		if v.Type().Key().Kind() != reflect.String {
			if abort := rts.delegate.OnMapKeyIsNotString(rts, v, attr); abort {
				return
			}
		}

		rts.ContainerAttr.Push(attr)
		rts.ContainerVal.Push(v)

		for _, subAttr := range attr.SubAttributes {
			rts.traverseWithReflection(v.MapIndex(reflect.ValueOf(subAttr.Assist.JSONName)), subAttr)
		}

		rts.ContainerAttr.Pop()
		rts.ContainerVal.Pop()

	case reflect.Bool:
		if abort := rts.delegate.OnValueIsBool(rts, v, attr); abort {
			return
		}

	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		if abort := rts.delegate.OnValueIsInt(rts, v, attr); abort {
			return
		}

	case reflect.Float32, reflect.Float64:
		if abort := rts.delegate.OnValueIsFloat(rts, v, attr); abort {
			return
		}

	case reflect.String:
		if abort := rts.delegate.OnValueIsString(rts, v, attr); abort {
			return
		}

	default:
		fmt.Println("type: ", v.Type())
		if abort := rts.delegate.OnUnsupportedType(rts, v, attr); abort {
			return
		}
	}
}

type broadcastTraversalDelegate struct {
	delegates []ResourceTraversalDelegate
}

func (bd *broadcastTraversalDelegate) OnInvalidValue(rts *ResourceTraversalState, v reflect.Value, attr *resource.Attribute) bool {
	for _, d := range bd.delegates {
		if abort := d.OnInvalidValue(rts, v, attr); abort {
			return abort
		}
	}
	return false
}

func (bd *broadcastTraversalDelegate) OnValueIsInterface(rts *ResourceTraversalState, v reflect.Value, attr *resource.Attribute) bool {
	for _, d := range bd.delegates {
		if abort := d.OnValueIsInterface(rts, v, attr); abort {
			return abort
		}
	}
	return false
}

func (bd *broadcastTraversalDelegate) OnValueIsSlice(rts *ResourceTraversalState, v reflect.Value, attr *resource.Attribute) bool {
	for _, d := range bd.delegates {
		if abort := d.OnValueIsSlice(rts, v, attr); abort {
			return abort
		}
	}
	return false
}

func (bd *broadcastTraversalDelegate) OnValueIsArray(rts *ResourceTraversalState, v reflect.Value, attr *resource.Attribute) bool {
	for _, d := range bd.delegates {
		if abort := d.OnValueIsArray(rts, v, attr); abort {
			return abort
		}
	}
	return false
}

func (bd *broadcastTraversalDelegate) OnValueIsMap(rts *ResourceTraversalState, v reflect.Value, attr *resource.Attribute) bool {
	for _, d := range bd.delegates {
		if abort := d.OnValueIsMap(rts, v, attr); abort {
			return abort
		}
	}
	return false
}

func (bd *broadcastTraversalDelegate) OnMapKeyIsNotString(rts *ResourceTraversalState, v reflect.Value, attr *resource.Attribute) bool {
	for _, d := range bd.delegates {
		if abort := d.OnMapKeyIsNotString(rts, v, attr); abort {
			return abort
		}
	}
	return false
}

func (bd *broadcastTraversalDelegate) OnValueIsBool(rts *ResourceTraversalState, v reflect.Value, attr *resource.Attribute) bool {
	for _, d := range bd.delegates {
		if abort := d.OnValueIsBool(rts, v, attr); abort {
			return abort
		}
	}
	return false
}

func (bd *broadcastTraversalDelegate) OnValueIsInt(rts *ResourceTraversalState, v reflect.Value, attr *resource.Attribute) bool {
	for _, d := range bd.delegates {
		if abort := d.OnValueIsInt(rts, v, attr); abort {
			return abort
		}
	}
	return false
}

func (bd *broadcastTraversalDelegate) OnValueIsFloat(rts *ResourceTraversalState, v reflect.Value, attr *resource.Attribute) bool {
	for _, d := range bd.delegates {
		if abort := d.OnValueIsFloat(rts, v, attr); abort {
			return abort
		}
	}
	return false
}

func (bd *broadcastTraversalDelegate) OnValueIsString(rts *ResourceTraversalState, v reflect.Value, attr *resource.Attribute) bool {
	for _, d := range bd.delegates {
		if abort := d.OnValueIsString(rts, v, attr); abort {
			return abort
		}
	}
	return false
}

func (bd *broadcastTraversalDelegate) OnUnsupportedType(rts *ResourceTraversalState, v reflect.Value, attr *resource.Attribute) bool {
	for _, d := range bd.delegates {
		if abort := d.OnUnsupportedType(rts, v, attr); abort {
			return abort
		}
	}
	return false
}
