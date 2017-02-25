package validation

import (
	"context"
	"fmt"
	"github.com/go-scim/scimify/resource"
	"reflect"
)

type requiredRulesValidator struct{}

func (v *requiredRulesValidator) Validate(r *resource.Resource, opt ValidationOptions, ctx context.Context) (bool, error) {
	schema, ok := ctx.Value(resource.CK_Schema).(*resource.Schema)
	if !ok {
		panic("missing required context parameter: CK_Schema")
	}

	for _, attr := range schema.Attributes {
		object, _ := getObjectByKey(r, attr.Assist.JSONName)
		if ok, err := v.validate(object, attr, opt, ctx); !ok {
			return false, err
		}
	}
	return true, nil
}

func (v *requiredRulesValidator) validate(object interface{}, attr *resource.Attribute, opt ValidationOptions, ctx context.Context) (bool, error) {
	if attr.Required && attr.IsUnassigned(object) {
		switch {
		case resource.ReadOnly == attr.Mutability && !opt.ReadOnlyIsMandatory:
			return true, nil
		default:
			return false, &validationError{
				ViolationType: requiredMissing,
				FullPath:      attr.Assist.FullPath,
				Message:       v.formulateErrorMessage(attr),
			}
		}
	}

	if attr.MultiValued {
		if attr.IsUnassigned(object) {
			return true, nil
		} else if reflect.Slice != reflect.TypeOf(object).Kind() {
			return false, &validationError{
				ViolationType: requiredMissing,
				FullPath:      attr.Assist.FullPath,
				Message:       fmt.Sprintf("failed to check required rule for [%s]: not an array", attr.Assist.FullPath),
			}
		}

		slice := reflect.ValueOf(object)
		clonedAttr := attr.Clone()
		clonedAttr.MultiValued = false
		if resource.Complex == attr.Type {
			for i := 0; i < slice.Len(); i++ {
				if ok, err := v.validate(slice.Index(i).Interface(), clonedAttr, opt, ctx); !ok {
					return false, err
				}
			}
		}
	} else if resource.Complex == attr.Type {
		switch object.(type) {
		case map[string]interface{}, *resource.Meta:
		default:
			return false, &validationError{
				ViolationType: requiredMissing,
				FullPath:      attr.Assist.FullPath,
				Message:       fmt.Sprintf("failed to check required rule for [%s]: not a complex object", attr.Assist.FullPath),
			}
		}

		for _, subAttr := range attr.SubAttributes {
			subObject, _ := getObjectByKey(object, subAttr.Assist.JSONName)
			if ok, err := v.validate(subObject, subAttr, opt, ctx); !ok {
				return false, err
			}
		}
	}
	return true, nil
}

func (v *requiredRulesValidator) formulateErrorMessage(attr *resource.Attribute) string {
	return fmt.Sprintf("missing reuqired attribute [%s]", attr.Assist.FullPath)
}
