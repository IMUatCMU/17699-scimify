package defaults

import (
	"github.com/go-scim/scimify/resource"
	"context"
	"strings"
)

type caseFormatValueDefaulter struct {}

func (d *caseFormatValueDefaulter) Default(r *resource.Resource, ctx context.Context) (bool, error) {
	schema, ok := ctx.Value(resource.CK_Schema).(*resource.Schema)
	if !ok {
		panic("missing required context parmaeter: CK_Schema")
	}

	r.Lock()
	defer r.Unlock()
	for _, attr := range schema.Attributes {
		d.performCaseCorrection(r.Attributes, attr, ctx)
	}
	return true, nil
}

func (d *caseFormatValueDefaulter) performCaseCorrection(container map[string]interface{}, attr *resource.Attribute, ctx context.Context) {
	var preCorrectionKey string
	if _, ok := container[attr.Assist.JSONName]; ok {
		return
	} else {
		for k := range container {
			if strings.ToLower(k) == strings.ToLower(attr.Assist.JSONName) {
				preCorrectionKey = k
				break
			}
		}
		if len(preCorrectionKey) == 0 {
			return
		}
	}

	switch attr.Type {
	case resource.Complex:
		if attr.MultiValued {
			if array, ok := container[preCorrectionKey].([]interface{}); !ok {
				return
			} else {
				for _, elem := range array {
					if subContainer, ok := elem.(map[string]interface{}); !ok {
						continue
					} else {
						for _, subAttr := range attr.SubAttributes {
							d.performCaseCorrection(subContainer, subAttr, ctx)
						}
					}
				}
			}
		} else {
			if subContainer, ok := container[preCorrectionKey].(map[string]interface{}); !ok {
				return
			} else {
				for _, subAttr := range attr.SubAttributes {
					d.performCaseCorrection(subContainer, subAttr, ctx)
				}
			}
		}
		fallthrough

	default:
		container[attr.Assist.JSONName] = container[preCorrectionKey]
		delete(container, preCorrectionKey)
	}
}
