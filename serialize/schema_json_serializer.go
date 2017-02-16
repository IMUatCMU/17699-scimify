package serialize

import (
	"fmt"
	"github.com/go-scim/scimify/helper"
	"github.com/go-scim/scimify/resource"
	"strconv"
	"strings"
)

type SchemaJsonSerializerContext struct {
	InclusionPaths []string
	ExclusionPaths []string
	Schema         *resource.Schema
}

type SchemaJsonSerializer struct{}

func (s *SchemaJsonSerializer) Serialize(resource *resource.Resource, context interface{}) ([]byte, error) {
	data := resource.Attributes
	if len(resource.Schemas) > 0 {
		resource.Attributes["schemas"] = resource.Schemas
	}
	if len(resource.Id) > 0 {
		resource.Attributes["id"] = resource.Id
	}
	if len(resource.ExternalId) > 0 {
		resource.Attributes["externalId"] = resource.ExternalId
	}
	if resource.Meta != nil {
		meta := make(map[string]interface{})
		if len(resource.Meta.ResourceType) > 0 {
			meta["resourceType"] = resource.Meta.ResourceType
		}
		if len(resource.Meta.Created) > 0 {
			meta["created"] = resource.Meta.Created
		}
		if len(resource.Meta.LastModified) > 0 {
			meta["lastModified"] = resource.Meta.LastModified
		}
		if len(resource.Meta.Location) > 0 {
			meta["location"] = resource.Meta.Location
		}
		if len(resource.Meta.Version) > 0 {
			meta["version"] = resource.Meta.Version
		}
		resource.Attributes["meta"] = meta
	}

	json, err := serializeMap(data, context.(*SchemaJsonSerializerContext).Schema, context.(*SchemaJsonSerializerContext))
	if err != nil {
		return nil, err
	}
	return []byte(json), nil
}

// map aggregator
type jsonMapEntryAggregator struct {
	state []string
}

func (j *jsonMapEntryAggregator) Aggregate(key, input interface{}) {
	j.state = append(j.state, input.(string))
}
func (j *jsonMapEntryAggregator) Result() interface{} {
	return fmt.Sprintf("{%s}", strings.Join(j.state, ","))
}

// array aggregator
type jsonArrayElementAggregator struct {
	state []string
}

func (j *jsonArrayElementAggregator) Aggregate(key, input interface{}) {
	j.state = append(j.state, input.(string))
}
func (j *jsonArrayElementAggregator) Result() interface{} {
	return fmt.Sprintf("[%s]", strings.Join(j.state, ","))
}

func serializeMap(target map[string]interface{}, attrGuide resource.AttributeGetter, context *SchemaJsonSerializerContext) (string, error) {
	var processor func(key string, value interface{}) (interface{}, error)
	processor = func(key string, value interface{}) (interface{}, error) {
		attr := attrGuide.GetAttribute(key)
		if nil == attr {
			return "", fmt.Errorf("failed to obtain attribute for path segment '%s'.", key)
		}

		if attr.MultiValued {
			if array, ok := value.([]interface{}); !ok {
				return "", fmt.Errorf("%+v cannot be parsed as a slice", value)
			} else {
				arrayAggregator := &jsonArrayElementAggregator{state: make([]string, 0)}
				result, err := helper.WalkSliceInSerial(array, func(idx int, elem interface{}) (interface{}, error) {
					switch attr.Type {
					case resource.String, resource.DateTime, resource.Reference:
						if s, err := parseToString(elem); err != nil {
							return nil, err
						} else {
							return strconv.Quote(s), nil
						}

					case resource.Integer:
						if i, err := parseToInteger(value); err != nil {
							return "", err
						} else {
							return fmt.Sprintf("%d", i), nil
						}

					case resource.Decimal:
						if f, err := parseToFloat(value); err != nil {
							return "", err
						} else {
							return fmt.Sprintf("%f", f), nil
						}

					case resource.Boolean:
						if b, err := parseToBool(value); err != nil {
							return "", err
						} else {
							return fmt.Sprintf("%t", b), nil
						}

					case resource.Complex:
						clonedAttr := attr.Clone()
						clonedAttr.MultiValued = false
						if m, ok := elem.(map[string]interface{}); !ok {
							return "", fmt.Errorf("%+v cannot be parsed as a map 1 (%T)", value, value)
						} else {
							if subJson, err := serializeMap(m, clonedAttr, context); err != nil {
								return "", err
							} else {
								return subJson, nil
							}
						}

					default:
						return "", fmt.Errorf("serializer cannot handle type %s", attr.Type)
					}
				}, arrayAggregator)
				if err != nil {
					return "", err
				}
				return fmt.Sprintf("%s:%s", strconv.Quote(attr.Assist.JSONName), result.(string)), nil
			}
		} else {
			switch attr.Type {
			case resource.String, resource.DateTime, resource.Reference:
				if s, err := parseToString(value); err != nil {
					return "", err
				} else {
					return fmt.Sprintf("%s:%s", strconv.Quote(attr.Assist.JSONName), strconv.Quote(s)), nil
				}

			case resource.Integer:
				if i, err := parseToInteger(value); err != nil {
					return "", err
				} else {
					return fmt.Sprintf("%s:%d", strconv.Quote(attr.Assist.JSONName), i), nil
				}

			case resource.Decimal:
				if f, err := parseToFloat(value); err != nil {
					return "", err
				} else {
					return fmt.Sprintf("%s:%f", strconv.Quote(attr.Assist.JSONName), f), nil
				}

			case resource.Boolean:
				if b, err := parseToBool(value); err != nil {
					return "", err
				} else {
					return fmt.Sprintf("%s:%t", strconv.Quote(attr.Assist.JSONName), b), nil
				}

			case resource.Complex:
				if nil != value {
					if m, ok := value.(map[string]interface{}); !ok {
						return "", fmt.Errorf("%+v cannot be parsed as a map 2", value)
					} else {
						if subJson, err := serializeMap(m, attr, context); err != nil {
							return "", err
						} else {
							return fmt.Sprintf("%s:%s", strconv.Quote(attr.Assist.JSONName), subJson), nil
						}
					}
				}

			default:
				return "", fmt.Errorf("serializer cannot handle type %s", attr.Type)
			}
		}

		return "", nil
	}

	aggregator := &jsonMapEntryAggregator{state: make([]string, 0)}
	_, err := helper.WalkStringMapInParallel(target, processor, aggregator)
	if err != nil {
		return "", err
	}
	return aggregator.Result().(string), nil
}

func parseToString(value interface{}) (string, error) {
	if s, ok := value.(string); !ok {
		return "", fmt.Errorf("%+v cannot be serialized as a string", value)
	} else {
		return s, nil
	}
}

func parseToInteger(value interface{}) (int64, error) {
	if i, ok := value.(int64); !ok {
		return 0, fmt.Errorf("%+v cannot be serialized as an integer", value)
	} else {
		return i, nil
	}
}

func parseToFloat(value interface{}) (float64, error) {
	if f, ok := value.(float64); !ok {
		return 0.0, fmt.Errorf("%+v cannot be serialized as a float", value)
	} else {
		return f, nil
	}
}

func parseToBool(value interface{}) (bool, error) {
	if b, ok := value.(bool); !ok {
		return false, fmt.Errorf("%+v cannot be serialized as a boolean", value)
	} else {
		return b, nil
	}
}
