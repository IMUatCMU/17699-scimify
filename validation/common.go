package validation

import (
	"fmt"
	"github.com/go-scim/scimify/resource"
	"strings"
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
		switch strings.ToLower(key) {
		case "schemas":
			return r.Schemas, nil
		case "id":
			return r.Id, nil
		case "externalid":
			return r.ExternalId, nil
		case "meta":
			return r.Meta, nil
		default:
			for k, v := range r.Attributes {
				if strings.ToLower(k) == strings.ToLower(key) {
					return v, nil
				}
			}
			return nil, fmt.Errorf("cannot get object by key '%s': missing key", key)
		}
	} else if meta, ok := target.(*resource.Meta); ok {
		switch strings.ToLower(key) {
		case "resourcetype":
			return meta.ResourceType, nil
		case "location":
			return meta.Location, nil
		case "version":
			return meta.Version, nil
		case "created":
			return meta.Created, nil
		case "lastmodified":
			return meta.LastModified, nil
		default:
			return nil, fmt.Errorf("cannot get object by key '%s': missing key", key)
		}
	} else {
		return nil, fmt.Errorf("cannot get object by key '%s': unknown target type %T", key, target)
	}
}