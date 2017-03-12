package processor

import (
	"github.com/go-scim/scimify/resource"
	"fmt"
	"reflect"
)

// Error representing the scenario where a required parameter from context is not present.
type MissingContextValueError struct {
	Key 	resource.ContextKey
}

func (mcv *MissingContextValueError) Error() string {
	switch mcv.Key {
	case resource.CK_Schema:
		return "missing schema in context"
	case resource.CK_Reference:
		return "missing reference resource in context"
	case resource.CK_ResourceType:
		return "missing resource type in context"
	case resource.CK_ResourceTypeURI:
		return "missing resource type uri in context"
	default:
		return fmt.Sprintf("missing attribute with index %d in context", mcv.Key)
	}
}

// Error representing the scenario where a certain type is not what was expected by the defined attribute.
type TypeMismatchError struct {
	T    reflect.Type
	Attr *resource.Attribute
}

func (tm *TypeMismatchError) Error() string {
	var expects string = ""
	switch tm.Attr.Type {
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
	if tm.Attr.MultiValued {
		expects += " array"
	}
	return "type check expected type: " + expects + ", unsupported type: " + tm.T.String()
}

// Error representing the scenario where a data format is not what was expected by the defined attribute.
type FormatError struct {
	Attr   *resource.Attribute
	Format string
	Actual interface{}
}

func (fe *FormatError) Error() string {
	return fmt.Sprintf("required format of [%s] at %s, but got %v", fe.Format, fe.Attr.Assist.FullPath, fe.Actual)
}

// Error representing the scenario where a type is not supported for processing at all.
type TypeUnsupportedError struct {
	T reflect.Type
}

func (ute *TypeUnsupportedError) Error() string {
	return fmt.Sprintf("type %s is not supported", ute.T.String())
}