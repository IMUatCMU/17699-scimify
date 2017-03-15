package processor

import (
	"crypto/sha1"
	"encoding/base64"
	"fmt"
	"github.com/go-scim/scimify/resource"
	"time"
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

func get(ctx *ProcessorContext, key AName, panicIfAbsent bool, defaultValue interface{}) interface{} {
	if val, ok := ctx.MiscArgs[key]; !ok {
		if panicIfAbsent {
			panic(&MissingContextValueError{key})
		} else {
			return defaultValue
		}
	} else {
		return val
	}
}

func getString(ctx *ProcessorContext, key AName, panicIfAbsent bool, defaultValue string) string {
	if val, ok := ctx.MiscArgs[key].(string); !ok {
		if panicIfAbsent {
			panic(&MissingContextValueError{key})
		} else {
			return defaultValue
		}
	} else {
		return val
	}
}

func getInt(ctx *ProcessorContext, key AName, panicIfAbsent bool, defaultValue int) int {
	if val, ok := ctx.MiscArgs[key].(int); !ok {
		if panicIfAbsent {
			panic(&MissingContextValueError{key})
		} else {
			return defaultValue
		}
	} else {
		return val
	}
}

func getBool(ctx *ProcessorContext, key AName, panicIfAbsent bool, defaultValue bool) bool {
	if val, ok := ctx.MiscArgs[key].(bool); !ok {
		if panicIfAbsent {
			panic(&MissingContextValueError{key})
		} else {
			return defaultValue
		}
	} else {
		return val
	}
}

func getSchema(ctx *ProcessorContext, panicIfAbsent bool) *resource.Schema {
	if nil == ctx.Schema && panicIfAbsent {
		panic(&MissingContextValueError{ArgSchema})
	}
	return ctx.Schema
}

func getReference(ctx *ProcessorContext, panicIfAbsent bool) *resource.Resource {
	if nil == ctx.Reference && panicIfAbsent {
		panic(&MissingContextValueError{ArgReference})
	}
	return ctx.Reference
}

func getResource(ctx *ProcessorContext, panicIfAbsent bool) *resource.Resource {
	if nil == ctx.Resource && panicIfAbsent {
		panic(&MissingContextValueError{ArgResource})
	}
	return ctx.Resource
}

func getError(ctx *ProcessorContext, panicIfAbsent bool) error {
	if e, ok := ctx.MiscArgs[ArgError].(error); !ok {
		if panicIfAbsent {
			panic(&MissingContextValueError{ArgError})
		} else {
			return nil
		}
	} else {
		return e
	}
}

func getCurrentTime() string {
	return time.Now().Format("2006-01-02T15:04:05Z")
}

func generateNewVersion(id string) string {
	hash := sha1.New()
	now := time.Now().Format(time.RFC3339Nano)
	hash.Write([]byte(id))
	hash.Write([]byte(now))
	return fmt.Sprintf("W/\"%s\"", base64.StdEncoding.EncodeToString(hash.Sum(nil)))
}
