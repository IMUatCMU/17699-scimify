package validation

import (
	"github.com/go-scim/scimify/helper"
	"github.com/go-scim/scimify/resource"
)

type delegateValidator struct {
	Delegates  []interface{}
	Concurrent bool
}

func (v *delegateValidator) Validate(resource *resource.Resource, context *ValidatorContext) (ok bool, err error) {
	defer func() {
		if r := recover(); r != nil {
			ok = false
			err = &validationError{
				ViolationType: unknown,
				Message:       r.(error).Error(),
				FullPath:      "",
			}
		}
	}()

	processor := func(idx int, elem interface{}) (interface{}, error) {
		validator := elem.(Validator)
		ok, err := validator.Validate(resource, context)
		return ok, err
	}

	if v.Concurrent {
		_, err = helper.WalkSliceInParallel(v.Delegates, processor, &helper.DoNothingAggregator{})
	} else {
		_, err = helper.WalkSliceInSerial(v.Delegates, processor, &helper.DoNothingAggregator{})
	}

	if nil == err {
		ok = true
	} else {
		ok = false
	}
	return
}
