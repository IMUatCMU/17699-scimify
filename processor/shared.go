package processor

import (
	"context"
	"github.com/go-scim/scimify/resource"
)

const (
	type_bool     = resource.Boolean
	type_int      = resource.Integer
	type_float    = resource.Decimal
	type_string   = resource.String
	type_ref      = resource.Reference
	type_binary   = resource.Binary
	type_datetime = resource.DateTime
	type_complex  = resource.Complex
)

func getSchema(ctx context.Context, panicIfAbsent bool) *resource.Schema {
	if schema, ok := ctx.Value(resource.CK_Schema).(*resource.Schema); !ok {
		if panicIfAbsent {
			panic(&MissingContextValueError{resource.CK_Schema})
		} else {
			return nil
		}
	} else {
		return schema
	}
}

func getReference(ctx context.Context, panicIfAbsent bool) *resource.Resource {
	if ref, ok := ctx.Value(resource.CK_Reference).(*resource.Resource); !ok {
		if panicIfAbsent {
			panic(&MissingContextValueError{resource.CK_Reference})
		} else {
			return nil
		}
	} else {
		return ref
	}
}
