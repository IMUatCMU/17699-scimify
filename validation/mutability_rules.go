package validation

import (
	"fmt"
	"github.com/go-scim/scimify/resource"
	"reflect"
	"sort"
)

// This mutability validator ignores caseExact directives
// and requirements for case insensitive keys - for another day
type mutabilityRulesValidator struct{}

func (v *mutabilityRulesValidator) Validate(r *resource.Resource, ctx *ValidatorContext) (bool, error) {
	ctx.RLock()
	schema := ctx.Data[Schema].(*resource.Schema)
	ref := ctx.Data[ReferenceResource]
	ctx.RUnlock()
	if nil == ref {
		return true, nil
	}

	for _, attr := range schema.Attributes {
		object, _ := getObjectByKey(r, attr.Assist.JSONName)
		reference, _ := getObjectByKey(ref, attr.Assist.JSONName)
		if ok, err := v.validate(object, reference, attr, ctx); !ok {
			return false, err
		}
	}
	return true, nil
}

func (v *mutabilityRulesValidator) validate(obj interface{}, ref interface{}, attr *resource.Attribute, context *ValidatorContext) (bool, error) {
	switch attr.Mutability {
	// Immutable attributes must have same value
	case resource.Immutable:
		if attr.IsUnassigned(ref) && v.bypassNilImmutable(context) {
			return true, nil
		} else if !v.equal(obj, ref, attr) {
			return false, &validationError{
				ViolationType: mutabilityCheck,
				FullPath:      attr.Assist.FullPath,
				Message:       fmt.Sprintf("immutable attribute [%s] has changed value", attr.Assist.FullPath),
			}
		}

	// ReadOnly attributes should have same value if provided in obj
	case resource.ReadOnly:
		if nil != obj && !v.equal(obj, ref, attr) {
			return false, &validationError{
				ViolationType: mutabilityCheck,
				FullPath:      attr.Assist.FullPath,
				Message:       fmt.Sprintf("readOnly attribute [%s] has changed value", attr.Assist.FullPath),
			}
		}
	}

	return true, nil
}

func (v *mutabilityRulesValidator) equal(obj interface{}, ref interface{}, attr *resource.Attribute) bool {
	if attr.MultiValued {
		objUnassigned, refUnassigned := attr.IsUnassigned(obj), attr.IsUnassigned(ref)
		if objUnassigned && refUnassigned {
			return true
		} else if objUnassigned != refUnassigned {
			return false
		}

		objVal, refVal := reflect.ValueOf(obj), reflect.ValueOf(ref)
		if objVal.Len() != refVal.Len() {
			return false
		}

		switch attr.Type {
		case resource.Complex:
			objWrap, refWrap := complexByKey{
				Attr:  attr,
				Slice: obj.([]interface{}),
			}, complexByKey{
				Attr:  attr,
				Slice: ref.([]interface{}),
			}
			sort.Sort(objWrap)
			sort.Sort(refWrap)
			return reflect.DeepEqual(objWrap.Slice, refWrap.Slice)

		case resource.String, resource.DateTime, resource.Reference, resource.Binary:
			sort.Strings(obj.([]string))
			sort.Strings(ref.([]string))
			return reflect.DeepEqual(obj, ref)

		case resource.Integer:
			sort.Ints(obj.([]int))
			sort.Ints(ref.([]int))
			return reflect.DeepEqual(obj, ref)

		case resource.Decimal:
			sort.Float64s(obj.([]float64))
			sort.Float64s(ref.([]float64))
			return reflect.DeepEqual(obj, ref)

		default:
			return false
		}
	} else {
		switch attr.Type {
		case resource.Complex:
			return reflect.DeepEqual(ref, obj)
		default:
			return ref == obj
		}
	}
}

func (v *mutabilityRulesValidator) bypassNilImmutable(ctx *ValidatorContext) bool {
	ctx.RLock()
	val := ctx.Data[IgnoreNilImmutable]
	ctx.RUnlock()
	if nil == val {
		return false
	} else {
		return val.(bool)
	}
}

type complexByKey struct {
	Slice []interface{}
	Attr  *resource.Attribute
}

func (b complexByKey) Len() int {
	return len(b.Slice)
}

func (b complexByKey) Swap(i, j int) {
	b.Slice[i], b.Slice[j] = b.Slice[j], b.Slice[i]
}

func (b complexByKey) Less(i, j int) bool {
	c1, c2 := b.Slice[i], b.Slice[j]
	k1, err := getObjectByKey(c1, b.Attr.Assist.ArrayIndexKey[0])
	k2, err := getObjectByKey(c2, b.Attr.Assist.ArrayIndexKey[0])
	if err != nil {
		return false
	}

	switch k1.(type) {
	case string:
		return k1.(string) < k2.(string)
	case int, int32, int64:
		return k1.(int64) < k2.(int64)
	case float32, float64:
		return k1.(float64) < k2.(float64)
	default:
		return false
	}
}
