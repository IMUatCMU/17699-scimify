package helper

import (
	"github.com/go-scim/scimify/adt"
	"github.com/go-scim/scimify/resource"
	"reflect"
)

type TraversalWithReferenceDelegate interface {
	OnAtLeastOneIsInvalidValue(attr *resource.Attribute, v reflect.Value, ref reflect.Value, dts *DoubleTraversalState) (abort bool)

	OnTypeIsDifferent(attr *resource.Attribute, v reflect.Value, ref reflect.Value, dts *DoubleTraversalState) (abort bool)

	OnTypeIsSame(attr *resource.Attribute, v reflect.Value, ref reflect.Value, dts *DoubleTraversalState) (abort bool)

	OnElemCrossTraversalStart(attr *resource.Attribute, v reflect.Value, idx int, dts *DoubleTraversalState)

	OnElemCrossTraversalEnd(attr *resource.Attribute, v reflect.Value, idx int, dts *DoubleTraversalState)

	OnUnsupportedType(attr *resource.Attribute, v reflect.Value, ref reflect.Value, dts *DoubleTraversalState) (abort bool)
}

type DoubleTraversalState struct {
	RootSchema    *resource.Schema
	ContainerAttr adt.Stack
	ContainerVal  adt.Stack
	ContainerRef  adt.Stack
	delegate      TraversalWithReferenceDelegate
}

func TraverseWithReference(r *resource.Resource, ref *resource.Resource, sch *resource.Schema, dlg []TraversalWithReferenceDelegate) {
	state := &DoubleTraversalState{
		RootSchema:    sch,
		ContainerAttr: adt.NewStackWithoutLimit(),
		ContainerVal:  adt.NewStackWithoutLimit(),
		ContainerRef:  adt.NewStackWithoutLimit(),
		delegate:      &broadcastDoubleTraversalDelegate{delegates: dlg},
	}
	state.traverseWithReflection(reflect.ValueOf(r.Data()), reflect.ValueOf(ref.Data()), sch.AsAttribute())
}

func (dts *DoubleTraversalState) traverseWithReflection(v reflect.Value, w reflect.Value, attr *resource.Attribute) {
	if !v.IsValid() || !w.IsValid() {
		if abort := dts.delegate.OnAtLeastOneIsInvalidValue(attr, v, w, dts); abort {
			return
		}
	}

	if v.Kind() == reflect.Interface {
		dts.traverseWithReflection(v.Elem(), w, attr)
		return
	} else if w.Kind() == reflect.Interface {
		dts.traverseWithReflection(v, w.Elem(), attr)
		return
	}

	if v.Kind() != w.Kind() {
		if abort := dts.delegate.OnTypeIsDifferent(attr, v, w, dts); abort {
			return
		}
	} else {
		switch v.Kind() {
		case reflect.Slice, reflect.Array:
			if abort := dts.delegate.OnTypeIsSame(attr, v, w, dts); abort {
				return
			}

			elemAttr := attr.Clone()
			elemAttr.MultiValued = false

			dts.ContainerAttr.Push(attr)
			dts.ContainerVal.Push(v)
			dts.ContainerRef.Push(w)
			for i := 0; i < v.Len(); i++ {
				dts.delegate.OnElemCrossTraversalStart(attr, v, i, dts)
				for j := 0; j < w.Len(); j++ {
					dts.traverseWithReflection(v.Index(i), w.Index(j), elemAttr)
				}
				dts.delegate.OnElemCrossTraversalEnd(attr, v, i, dts)
			}
			dts.ContainerAttr.Pop()
			dts.ContainerVal.Pop()
			dts.ContainerRef.Pop()

		case reflect.Map:
			if abort := dts.delegate.OnTypeIsSame(attr, v, w, dts); abort {
				return
			}

			dts.ContainerAttr.Push(attr)
			dts.ContainerVal.Push(v)
			dts.ContainerRef.Push(w)

			for _, subAttr := range attr.SubAttributes {
				dts.traverseWithReflection(
					v.MapIndex(reflect.ValueOf(subAttr.Assist.JSONName)),
					w.MapIndex(reflect.ValueOf(subAttr.Assist.JSONName)),
					subAttr,
				)
			}

			dts.ContainerAttr.Pop()
			dts.ContainerVal.Pop()
			dts.ContainerRef.Pop()

		case reflect.Bool, reflect.Float32, reflect.Float64, reflect.String,
			reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
			if abort := dts.delegate.OnTypeIsSame(attr, v, w, dts); abort {
				return
			}
		default:
			if abort := dts.delegate.OnUnsupportedType(attr, v, w, dts); abort {
				return
			}
		}
	}
}

type broadcastDoubleTraversalDelegate struct {
	delegates []TraversalWithReferenceDelegate
}

func (d *broadcastDoubleTraversalDelegate) OnAtLeastOneIsInvalidValue(attr *resource.Attribute, v reflect.Value, ref reflect.Value, dts *DoubleTraversalState) bool {
	for _, delegate := range d.delegates {
		if abort := delegate.OnAtLeastOneIsInvalidValue(attr, v, ref, dts); abort {
			return abort
		}
	}
	return false
}

func (d *broadcastDoubleTraversalDelegate) OnTypeIsDifferent(attr *resource.Attribute, v reflect.Value, ref reflect.Value, dts *DoubleTraversalState) bool {
	for _, delegate := range d.delegates {
		if abort := delegate.OnTypeIsDifferent(attr, v, ref, dts); abort {
			return abort
		}
	}
	return false
}

func (d *broadcastDoubleTraversalDelegate) OnTypeIsSame(attr *resource.Attribute, v reflect.Value, ref reflect.Value, dts *DoubleTraversalState) bool {
	for _, delegate := range d.delegates {
		if abort := delegate.OnTypeIsSame(attr, v, ref, dts); abort {
			return abort
		}
	}
	return false
}

func (d *broadcastDoubleTraversalDelegate) OnElemCrossTraversalStart(attr *resource.Attribute, v reflect.Value, idx int, dts *DoubleTraversalState) {
	for _, delegate := range d.delegates {
		delegate.OnElemCrossTraversalStart(attr, v, idx, dts)
	}
}

func (d *broadcastDoubleTraversalDelegate) OnElemCrossTraversalEnd(attr *resource.Attribute, v reflect.Value, idx int, dts *DoubleTraversalState) {
	for _, delegate := range d.delegates {
		delegate.OnElemCrossTraversalEnd(attr, v, idx, dts)
	}
}

func (d *broadcastDoubleTraversalDelegate) OnUnsupportedType(attr *resource.Attribute, v reflect.Value, ref reflect.Value, dts *DoubleTraversalState) (abort bool) {
	for _, delegate := range d.delegates {
		if abort := delegate.OnUnsupportedType(attr, v, ref, dts); abort {
			return abort
		}
	}
	return false
}
