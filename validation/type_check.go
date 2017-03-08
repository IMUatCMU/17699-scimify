package validation

import (
	"github.com/go-scim/scimify/resource"
	"context"
	"reflect"
	"fmt"
)

type typeCheckValidator struct{}

func (v *typeCheckValidator) Validate(r *resource.Resource, opt ValidationOptions, ctx context.Context) (bool, error) {
	schema, ok := ctx.Value(resource.CK_Schema).(*resource.Schema)
	if !ok {
		panic("missing required context parameter: CK_Schema")
	}

	if err := v.validate(r.Data(), typeCheckOpts{}, schema.AsAttribute()); err != nil {
		return false, err
	}

	return true, nil
}

type typeCheckFunc func(v reflect.Value, opts typeCheckOpts, attr *resource.Attribute)

func (v *typeCheckValidator) validate(val interface{}, opts typeCheckOpts, attr *resource.Attribute) (err error) {
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
	v.reflectValue(reflect.ValueOf(val), opts, attr)
	return nil
}

func (v *typeCheckValidator) reflectValue(val reflect.Value, opts typeCheckOpts, attr *resource.Attribute) {
	v.valueChecker(val, attr)(val, opts, attr)
}

func (v *typeCheckValidator) valueChecker(val reflect.Value, attr *resource.Attribute) typeCheckFunc {
	if !val.IsValid() {
		return v.skipValueChecker		// invalid value means there's no value there (possibly be a missing key), it's fine for type check, just skip this one.
	}
	return v.newValueChecker(val.Type(), attr)
}

func (v *typeCheckValidator) newValueChecker(t reflect.Type, attr *resource.Attribute) typeCheckFunc {
	if t.Kind() == reflect.Interface {
		return v.interfaceValueChecker
	}

	multiValued, dataType := attr.MultiValued, attr.Type
	if multiValued {
		switch t.Kind() {
		case reflect.Slice:
			return v.sliceValueChecker
		case reflect.Array:
			return v.arrayValueChecker
		default:
			return v.unsupportedType
		}
	} else {
		switch dataType {
		case type_bool:
			return v.boolValueChecker
		case type_int:
			return v.intValueChecker
		case type_float:
			return v.floatValueChecker
		// TODO extend this
		case type_string, type_ref, type_binary, type_datetime:
			return v.stringValueChecker
		case type_complex:
			return v.complexValueChecker
		default:
			return v.unsupportedType
		}
	}
}

func (v *typeCheckValidator) arrayValueChecker(val reflect.Value, opt typeCheckOpts, attr *resource.Attribute) {
	elemAttr := attr.Clone()
	elemAttr.MultiValued = false

	n := val.Len()
	for i := 0; i < n; i++ {
		elem := val.Index(i)
		v.valueChecker(elem, elemAttr)(elem, opt, elemAttr)
	}
}

func (v *typeCheckValidator) sliceValueChecker(val reflect.Value, opt typeCheckOpts, attr *resource.Attribute) {
	v.arrayValueChecker(val, opt, attr)
}

func (v *typeCheckValidator) complexValueChecker(val reflect.Value, opt typeCheckOpts, attr *resource.Attribute) {
	switch val.Kind() {
	case reflect.Map:
		if val.Type().Key().Kind() != reflect.String {
			v.unsupportedType(val, opt, attr)
		}
		keyAttrs := make([]*resource.Attribute, 0, len(attr.SubAttributes))
		for _, subAttr := range attr.SubAttributes {
			keyAttrs = append(keyAttrs, subAttr)
		}
		for _, subAttr := range keyAttrs {
			val := val.MapIndex(reflect.ValueOf(subAttr.Assist.JSONName))
			v.valueChecker(val, subAttr)(val, opt, subAttr)
		}
		return
	default:
		v.unsupportedType(val, opt, attr)
	}
}

func (v *typeCheckValidator) stringValueChecker(val reflect.Value, opt typeCheckOpts, attr *resource.Attribute) {
	switch val.Kind() {
	case reflect.String:
		return
	default:
		v.unsupportedType(val, opt, attr)
	}
}

func (v *typeCheckValidator) floatValueChecker(val reflect.Value, opt typeCheckOpts, attr *resource.Attribute) {
	switch val.Kind() {
	case reflect.Float32, reflect.Float64:
		return
	default:
		v.unsupportedType(val, opt, attr)
	}
}

func (v *typeCheckValidator) intValueChecker(val reflect.Value, opt typeCheckOpts, attr *resource.Attribute) {
	switch val.Kind() {
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return
	default:
		v.unsupportedType(val, opt, attr)
	}
}

func (v *typeCheckValidator) boolValueChecker(val reflect.Value, opt typeCheckOpts, attr *resource.Attribute) {
	switch val.Kind() {
	case reflect.Bool:
		return
	default:
		v.unsupportedType(val, opt, attr)
	}
}

func (v *typeCheckValidator) interfaceValueChecker(val reflect.Value, opts typeCheckOpts, attr *resource.Attribute) {
	if val.IsNil() {
		return
	}
	v.reflectValue(val.Elem(), opts, attr)
}

func (v *typeCheckValidator) unsupportedType(val reflect.Value, _ typeCheckOpts, attr *resource.Attribute) {
	v.error(&UnsupportedTypeError{val.Type(), attr})
}

func (v *typeCheckValidator) invalidValueChecker(val reflect.Value, opts typeCheckOpts, attr *resource.Attribute) {
	v.error(fmt.Errorf("Value at [%s] is not valid.", attr.Assist.FullPath))
}

func (v *typeCheckValidator) skipValueChecker(val reflect.Value, opts typeCheckOpts, attr *resource.Attribute) {
	return
}

func (v *typeCheckValidator) error(e error) {
	panic(e)
}

type typeCheckOpts struct {}

type UnsupportedTypeError struct {
	Type reflect.Type
	Attr *resource.Attribute
}

func (e *UnsupportedTypeError) Error() string {
	var expects string = ""
	switch e.Attr.Type {
	case resource.String, resource.Binary, resource.Reference, resource.DateTime:
		expects = "string"
	case resource.Integer:
		expects = "integer"
	case resource.Decimal:
		expects = "decimal"
	case resource.Boolean:
		expects = "boolean"
	case resource.Complex:
		expects = "complex"
	}
	if e.Attr.MultiValued {
		expects += " array"
	}
	return "type check expected type: " + expects + ", unsupported type: " + e.Type.String()
}