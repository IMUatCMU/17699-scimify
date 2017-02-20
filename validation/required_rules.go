package validation

import (
	"fmt"
	"github.com/go-scim/scimify/resource"
	"reflect"
)

type requiredRulesValidator struct{}

func (v *requiredRulesValidator) Validate(r *resource.Resource, ctx *ValidatorContext) (bool, error) {
	ctx.RLock()
	schema := ctx.Data[Schema].(*resource.Schema)
	ctx.RUnlock()
	for _, attr := range schema.Attributes {
		object, _ := getObjectByKey(r, attr.Assist.JSONName)
		if ok, err := v.validate(object, attr, ctx); !ok {
			return false, err
		}
	}
	return true, nil
}

func (v *requiredRulesValidator) validate(object interface{}, attr *resource.Attribute, context *ValidatorContext) (bool, error) {
	if !attr.Required && (resource.Complex != attr.Type || nil == object) {
		return true, nil
	} else if attr.Required && nil == object {
		// required attribute with nil value could only be spared
		// if the attribute is readOnly and the validator is set to
		// not fail missing readOnly attributes
		if resource.ReadOnly != attr.Mutability || !v.shouldFailMissingReadOnly(context) {
			return false, &validationError{
				ViolationType: requiredMissing,
				FullPath:      attr.Assist.FullPath,
				Message:       v.formulateErrorMessage(attr),
			}
		}
	}

	if attr.MultiValued {
		if reflect.Slice != reflect.TypeOf(object).Kind() {
			return false, &validationError{
				ViolationType: requiredMissing,
				FullPath:      attr.Assist.FullPath,
				Message:       fmt.Sprintf("failed to check required rule for [%s]: not an array", attr.Assist.FullPath),
			}
		}

		slice := reflect.ValueOf(object)
		if attr.Required && slice.Len() == 0 {
			return false, &validationError{
				ViolationType: requiredMissing,
				FullPath:      attr.Assist.FullPath,
				Message:       fmt.Sprintf("required but unassigned array attribute at [%s]", attr.Assist.FullPath),
			}
		}

		clonedAttr := attr.Clone()
		clonedAttr.MultiValued = false
		if resource.Complex == attr.Type {
			for i := 0; i < slice.Len(); i++ {
				if ok, err := v.validate(slice.Index(i).Interface(), clonedAttr, context); !ok {
					return false, err
				}
			}
		}
	} else {

		if resource.Complex == attr.Type && nil != object {
			if m, ok := object.(map[string]interface{}); ok && len(m) == 0 {
				if attr.Required {
					return false, &validationError{
						ViolationType: requiredMissing,
						FullPath:      attr.Assist.FullPath,
						Message:       fmt.Sprintf("required but empty complex attribute value at [%s]", attr.Assist.FullPath),
					}
				} else {
					return true, nil
				}
			}

			// check a non-empty complex object, regardless of required attribute
			for _, subAttr := range attr.SubAttributes {
				subObject, _ := getObjectByKey(object, subAttr.Assist.JSONName)
				if ok, err := v.validate(subObject, subAttr, context); !ok {
					return false, err
				}
			}
		}
	}
	return true, nil
}

func (v *requiredRulesValidator) shouldFailMissingReadOnly(ctx *ValidatorContext) bool {
	ctx.RLock()
	opt := ctx.Data[FailReadOnlyRequired]
	ctx.RUnlock()

	if nil == opt {
		return false
	} else {
		return opt.(bool)
	}
}

func (v *requiredRulesValidator) formulateErrorMessage(attr *resource.Attribute) string {
	return fmt.Sprintf("missing reuqired attribute [%s]", attr.Assist.FullPath)
}
