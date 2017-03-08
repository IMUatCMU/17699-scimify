package validation

import (
	"fmt"
	"github.com/go-scim/scimify/resource"
	"strings"
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

// Get object by simple key, case insensitive
func getObjectByKey(target interface{}, key string) (interface{}, error) {
	if nil == target {
		return nil, fmt.Errorf("cannot get object by key '%s': target is nil", key)
	} else if m, ok := target.(map[string]interface{}); ok {
		for k, v := range m {
			if strings.ToLower(k) == strings.ToLower(key) {
				return v, nil
			}
		}
		return nil, fmt.Errorf("cannot get object by key '%s': missing key", key)
	} else if r, ok := target.(*resource.Resource); ok {
		for k, v := range r.Attributes {
			if strings.ToLower(k) == strings.ToLower(key) {
				return v, nil
			}
		}
		return nil, fmt.Errorf("cannot get object by key '%s': missing key", key)
	} else {
		return nil, fmt.Errorf("cannot get object by key '%s': unknown target type %T", key, target)
	}
}
