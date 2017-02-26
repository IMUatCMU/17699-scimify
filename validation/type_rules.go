package validation

import (
	"context"
	"encoding/base64"
	"fmt"
	"github.com/go-scim/scimify/resource"
	"reflect"
	"strings"
	"time"
)

// validator that enforces the type of the attributes
// don't force error if nothing is there, leave it to required rule
type typeRulesValidator struct{}

func (v *typeRulesValidator) Validate(r *resource.Resource, opt ValidationOptions, ctx context.Context) (bool, error) {
	schema, ok := ctx.Value(resource.CK_Schema).(*resource.Schema)
	if !ok {
		panic("missing required context parameter: CK_Schema")
	}

	for _, attr := range schema.Attributes {
		object, _ := getObjectByKey(r, attr.Assist.JSONName)
		if nil != object {
			if ok, err := v.validate(object, attr, opt, ctx); !ok {
				return false, err
			}
		}
	}
	return true, nil
}

func (v *typeRulesValidator) validate(object interface{}, attr *resource.Attribute, opt ValidationOptions, ctx context.Context) (bool, error) {
	if attr.MultiValued {
		clonedAttr := attr.Clone()
		clonedAttr.MultiValued = false

		switch reflect.TypeOf(object).Kind() {
		case reflect.Slice:
			slice := reflect.ValueOf(object)
			for i := 0; i < slice.Len(); i++ {
				if ok, err := v.validate(slice.Index(i).Interface(), clonedAttr, opt, ctx); !ok {
					return false, err
				}
			}
		default:
			return false, &validationError{
				ViolationType: typeCheck,
				FullPath:      attr.Assist.FullPath,
				Message:       fmt.Sprintf("invalid type at [%s], expects array of %s", attr.Assist.FullPath, strings.ToLower(attr.Type)),
			}
		}

	} else {
		switch attr.Type {
		case resource.String, resource.Reference, resource.DateTime, resource.Binary:
			if strVal := reflect.ValueOf(object); strVal.Kind() != reflect.String {
				return false, v.formatTypeError(attr)
			} else if s := strVal.String(); len(s) > 0 && resource.DateTime == attr.Type {
				if _, err := time.Parse("2006-01-02T15:04:05Z", s); err != nil {
					return false, &validationError{
						ViolationType: typeCheck,
						FullPath:      attr.Assist.FullPath,
						Message:       fmt.Sprintf("invalid datetime format at [%s]", attr.Assist.FullPath),
					}
				}
			} else if len(s) > 0 && resource.Binary == attr.Type {
				if _, err := base64.StdEncoding.DecodeString(s); err != nil {
					return false, &validationError{
						ViolationType: typeCheck,
						FullPath:      attr.Assist.FullPath,
						Message:       fmt.Sprintf("invalid base64 encoded data at [%s]", attr.Assist.FullPath),
					}
				}
			}

		case resource.Integer:
			if _, ok := object.(int64); !ok {
				return false, v.formatTypeError(attr)
			}

		case resource.Decimal:
			if _, ok := object.(float64); !ok {
				return false, v.formatTypeError(attr)
			}

		case resource.Boolean:
			if _, ok := object.(bool); !ok {
				return false, v.formatTypeError(attr)
			}

		case resource.Complex:
			if _, ok := object.(map[string]interface{}); !ok {
				return false, v.formatTypeError(attr)
			}
			for _, subAttr := range attr.SubAttributes {
				subObject, _ := getObjectByKey(object, subAttr.Assist.JSONName)
				if nil != subObject {
					if ok, err := v.validate(subObject, subAttr, opt, ctx); !ok {
						return false, err
					}
				}
			}
		}
	}

	return true, nil
}

func (v *typeRulesValidator) formatTypeError(attr *resource.Attribute) error {
	return &validationError{
		ViolationType: typeCheck,
		FullPath:      attr.Assist.FullPath,
		Message:       fmt.Sprintf("invalid type at [%s], expects %s", attr.Assist.FullPath, strings.ToLower(attr.Type)),
	}
}
