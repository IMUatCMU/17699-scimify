package processor

import (
	"fmt"
	"github.com/go-scim/scimify/resource"
	"reflect"
)

// Error representing the scenario where a required parameter from context is not present.
type MissingContextValueError struct {
	Key string
}

func (mcv *MissingContextValueError) Error() string {
	return fmt.Sprintf("missing %s in context", mcv.Key)
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

// Error representing the scenario where a required attribute is missing (nil)
type RequiredMissingError struct {
	Attr *resource.Attribute
}

func (rme *RequiredMissingError) Error() string {
	return fmt.Sprintf("Missing required attribute %s", rme.Attr.Assist.FullPath)
}

// Error representing the scenario where a required attribute is unassigned (present but unassigned value, i.e. empty array)
type RequiredUnassignedError struct {
	Attr *resource.Attribute
}

func (rue *RequiredUnassignedError) Error() string {
	return fmt.Sprintf("Attribute %s is unassigned", rue.Attr.Assist.FullPath)
}

// Error representing the scenario where an immutable or read only attribute has its value changed (on update or patch)
type ValueChangedError struct {
	Attr *resource.Attribute
}

func (vce *ValueChangedError) Error() string {
	switch vce.Attr.Mutability {
	case resource.Immutable:
		return fmt.Sprintf("immutable attribute [%s] has changed value.", vce.Attr.Assist.FullPath)
	case resource.ReadOnly:
		return fmt.Sprintf("read only attribute [%s] has changed value.", vce.Attr.Assist.FullPath)
	default:
		return fmt.Sprintf("attribute [%s] has changed value.", vce.Attr.Assist.FullPath)
	}
}

type PrerequisiteFailedError struct {
	reporter    string
	requirement string
}

func (pfe *PrerequisiteFailedError) Error() string {
	return fmt.Sprintf("Prerequisite not met: %s requires %s", pfe.reporter, pfe.requirement)
}

// Error representing the scenario where there is no attribute defined at requested path
type NoDefinedAttributeError struct {
	Path string
}

func (nda *NoDefinedAttributeError) Error() string {
	return fmt.Sprintf("No attribute found for path: %s", nda.Path)
}

// Error representing the scenario where a found attribute does case insensitively match the map key
type AttributeMismatchWithKeyError struct {
	Key  string
	Attr *resource.Attribute
}

func (amk *AttributeMismatchWithKeyError) Error() string {
	return fmt.Sprintf("Attribute with path %s mismatches with entry key %s", amk.Attr.Assist.FullPath, amk.Key)
}

// Error representing the scenario where a type is not expected by the attribute (in a JSON serializer)
type UnexpectedTypeError struct {
	Type reflect.Type
	Attr *resource.Attribute
}

func (e *UnexpectedTypeError) Error() string {
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
	return "json expected type: " + expects + ", had type: " + e.Type.String()
}

// Error representing the scenario where a value cannot handled by the JSON serializer
type UnsupportedValueError struct {
	Value reflect.Value
	Str   string
}

func (e *UnsupportedValueError) Error() string {
	return "json: unsupported value: " + e.Str
}
