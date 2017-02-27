package defaults

import (
	"context"
	"github.com/go-scim/scimify/resource"
)

type copyReadOnlyValueDefaulter struct{}

type valueSetter func(key string)

func (d *copyReadOnlyValueDefaulter) Default(r *resource.Resource, ctx context.Context) (bool, error) {
	schema, ok := ctx.Value(resource.CK_Schema).(*resource.Schema)
	if !ok {
		panic("missing required context parmaeter: CK_Schema")
	}

	ref, ok := ctx.Value(resource.CK_Reference).(*resource.Resource)
	if !ok {
		panic("missing required context parameter: CK_Reference")
	}

	for _, attr := range schema.Attributes {
		d.copyReadOnly(r.Attributes, ref.Attributes, attr, ctx)
	}
	return true, nil
}

func (d *copyReadOnlyValueDefaulter) copyReadOnly(receiverContainer, referenceContainer map[string]interface{}, attr *resource.Attribute, ctx context.Context) {
	switch attr.Mutability {
	case resource.Immutable:
		return

	case resource.ReadOnly:
		if !attr.IsUnassigned(receiverContainer[attr.Assist.JSONName]) || !attr.IsUnassigned(referenceContainer[attr.Assist.JSONName]) {
			receiverContainer[attr.Assist.JSONName] = referenceContainer[attr.Assist.JSONName]
		}
		return

	case resource.WriteOnly, resource.ReadWrite:
		if resource.Complex == attr.Type {
			if attr.MultiValued {
				if receiverArray, ok := receiverContainer[attr.Assist.JSONName].([]interface{}); !ok {
					return
				} else if referenceArray, ok := referenceContainer[attr.Assist.JSONName].([]interface{}); !ok {
					return
				} else {
					for _, referenceElem := range referenceArray {
						subReferenceContainer, ok := referenceElem.(map[string]interface{})
						if !ok {
							continue
						}
						for _, receiverElem := range receiverArray {
							subReceiverContainer, ok := receiverElem.(map[string]interface{})
							if !ok {
								continue
							}
							if subReceiverContainer[attr.Assist.ArrayIndexKey[0]] == subReferenceContainer[attr.Assist.ArrayIndexKey[0]] {
								for _, subAttr := range attr.SubAttributes {
									d.copyReadOnly(subReceiverContainer, subReferenceContainer, subAttr, ctx)
								}
							}
						}
					}
				}
			} else {
				for _, subAttr := range attr.SubAttributes {
					if subReceiverContainer, ok := receiverContainer[attr.Assist.JSONName].(map[string]interface{}); !ok {
						return
					} else if subRefContainer, ok := referenceContainer[attr.Assist.JSONName].(map[string]interface{}); !ok {
						return
					} else {
						d.copyReadOnly(subReceiverContainer, subRefContainer, subAttr, ctx)
					}
				}
			}
		}
	}
}
