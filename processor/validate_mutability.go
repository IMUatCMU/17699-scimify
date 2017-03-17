package processor

import (
	"github.com/go-scim/scimify/helper"
	"github.com/go-scim/scimify/resource"
	"reflect"
	"sync"
)

var (
	oneMutabilityValidation sync.Once
	mutabilityValidator     Processor
)

func MutabilityValidationProcessor() Processor {
	oneMutabilityValidation.Do(func() {
		mutabilityValidator = &mutabilityValidationProcessor{}
	})
	return mutabilityValidator
}

type mutabilityValidationProcessor struct{}

func (mvp *mutabilityValidationProcessor) Process(ctx *ProcessorContext) (err error) {
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

	r := mvp.getResource(ctx)
	schema := mvp.getSchema(ctx)
	ref := mvp.getReference(ctx)
	delegate := &mutabilityCheckState{PerformCopy: true, inArray: false} // TODO turn into configuration option

	helper.TraverseWithReference(r, ref, schema, []helper.TraversalWithReferenceDelegate{delegate})

	err = nil
	return
}

func (mvp *mutabilityValidationProcessor) getResource(ctx *ProcessorContext) *resource.Resource {
	if ctx.Resource == nil {
		panic(&MissingContextValueError{"resource"})
	}
	return ctx.Resource
}

func (mvp *mutabilityValidationProcessor) getSchema(ctx *ProcessorContext) *resource.Schema {
	if ctx.Schema == nil {
		panic(&MissingContextValueError{"schema"})
	}
	return ctx.Schema
}

func (mvp *mutabilityValidationProcessor) getReference(ctx *ProcessorContext) *resource.Resource {
	if ctx.Reference == nil {
		panic(&MissingContextValueError{"reference"})
	}
	return ctx.Reference
}

type mutabilityCheckState struct {
	PerformCopy bool
	inArray     bool
}

type invalidIndicator int

const (
	noneIsInvalid  = invalidIndicator(0)
	valueIsInvalid = invalidIndicator(1)
	refIsInvalid   = invalidIndicator(2)
	bothAreInvalid = invalidIndicator(3)
)

func (mcs *mutabilityCheckState) copyReferenceValueToResource(attr *resource.Attribute, v reflect.Value, ref reflect.Value, dts *helper.DoubleTraversalState) {
	if mcs.PerformCopy {
		container := dts.ContainerVal.Peek().(reflect.Value)
		container.SetMapIndex(reflect.ValueOf(attr.Assist.JSONName), ref)
	}
}

func (mcs *mutabilityCheckState) OnAtLeastOneIsInvalidValue(attr *resource.Attribute, v reflect.Value, ref reflect.Value, dts *helper.DoubleTraversalState) (abort bool) {
	invalidity := mcs.computeInvalidIndicator(v, ref)

	switch invalidity {
	case valueIsInvalid:
		switch attr.Mutability {
		case resource.Immutable:
			mcs.error(&ValueChangedError{attr})
		case resource.ReadOnly:
			mcs.copyReferenceValueToResource(attr, v, ref, dts)
		}

	case refIsInvalid:
		switch attr.Mutability {
		case resource.Immutable:
			break // immutable field not set, allow value to set it
		case resource.ReadOnly:
			mcs.copyReferenceValueToResource(attr, v, reflect.Value{}, dts)
		}

	case bothAreInvalid:
		break
	}

	abort = true // no need to continue as at least one is invalid
	return
}

func (mcs *mutabilityCheckState) OnTypeIsDifferent(attr *resource.Attribute, v reflect.Value, ref reflect.Value, dts *helper.DoubleTraversalState) (abort bool) {
	abort = true // just abort, type discrepancy should have been picked up by type_check rules
	return
}

func (mcs *mutabilityCheckState) OnTypeIsSame(attr *resource.Attribute, v reflect.Value, ref reflect.Value, dts *helper.DoubleTraversalState) (abort bool) {
	if mcs.inArray {
		return mcs.onTypeIsSameWhenInArray(attr, v, ref, dts)
	} else {
		return mcs.onTypeIsSameWhenNotInArray(attr, v, ref, dts)
	}
}

func (mcs *mutabilityCheckState) onTypeIsSameWhenInArray(attr *resource.Attribute, v reflect.Value, ref reflect.Value, dts *helper.DoubleTraversalState) (abort bool) {
	switch attr.Mutability {
	case resource.Immutable, resource.ReadOnly:
		switch attr.Type {
		case resource.Complex:
			if reflect.DeepEqual(
				v.MapIndex(reflect.ValueOf(attr.Assist.ArrayIndexKey)).Interface(),
				ref.MapIndex(reflect.ValueOf(attr.Assist.ArrayIndexKey)).Interface(),
			) {
				abort = false
			} else {
				abort = true
			}

		default:
			if reflect.DeepEqual(v.Interface(), ref.Interface()) {
				abort = false
			} else {
				abort = true
				mcs.error(&ValueChangedError{attr})
			}
		}

	default:
		abort = false
	}
	return
}

func (mcs *mutabilityCheckState) onTypeIsSameWhenNotInArray(attr *resource.Attribute, v reflect.Value, ref reflect.Value, dts *helper.DoubleTraversalState) (abort bool) {
	switch attr.Mutability {
	case resource.Immutable, resource.ReadOnly:
		if !reflect.DeepEqual(v.Interface(), ref.Interface()) {
			mcs.error(&ValueChangedError{attr})
		}
		abort = true
		return
	default:
		abort = false
		return
	}
}

func (mcs *mutabilityCheckState) OnElemCrossTraversalStart(attr *resource.Attribute, v reflect.Value, idx int, dts *helper.DoubleTraversalState) {
	mcs.inArray = true
}

func (mcs *mutabilityCheckState) OnElemCrossTraversalEnd(attr *resource.Attribute, v reflect.Value, idx int, dts *helper.DoubleTraversalState) {
	mcs.inArray = false
}

func (mcs *mutabilityCheckState) OnUnsupportedType(attr *resource.Attribute, v reflect.Value, ref reflect.Value, dts *helper.DoubleTraversalState) (abort bool) {
	abort = true // just abort, this error should have been handled by type_check
	return
}

func (mcs *mutabilityCheckState) computeInvalidIndicator(v reflect.Value, ref reflect.Value) invalidIndicator {
	switch {
	case !v.IsValid() && !ref.IsValid():
		return bothAreInvalid
	case !v.IsValid() && ref.IsValid():
		return valueIsInvalid
	case v.IsValid() && !ref.IsValid():
		return refIsInvalid
	default:
		return noneIsInvalid
	}
}

func (mcs *mutabilityCheckState) error(err error) {
	panic(err)
}
