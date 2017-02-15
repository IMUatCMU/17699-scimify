package persistence

import (
	"fmt"
	"github.com/go-scim/scimify/adt"
	"github.com/go-scim/scimify/filter"
	"github.com/go-scim/scimify/resource"
	"gopkg.in/mgo.v2/bson"
)

func TranspileToMongoQuery(root *adt.Node, schema *resource.Schema) (bson.M, error) {

	if token, ok := root.Data.(filter.Token); !ok {
		return nil, resource.CreateError(resource.InvalidFilter, "Encountered invalid root type during transpilation.")
	} else {
		var (
			left, right bson.M
			attribute   *resource.Attribute
			err         error
		)

		if filter.Logical == token.Type {
			if root.Left != nil {
				left, err = TranspileToMongoQuery(root.Left, schema)
				if err != nil {
					return nil, err
				}
			}

			if root.Right != nil {
				right, err = TranspileToMongoQuery(root.Right, schema)
				if err != nil {
					return nil, err
				}
			}
		} else if filter.Relational == token.Type {
			if token.Params[filter.NumberOfArgs].(int) > 0 && filter.Path != root.Left.Data.(filter.Token).Type {
				return nil, resource.CreateError(resource.InvalidFilter, "Expects path name on the left of the operator.")
			} else if token.Params[filter.NumberOfArgs].(int) > 1 && filter.Constant != root.Right.Data.(filter.Token).Type {
				return nil, resource.CreateError(resource.InvalidFilter, "Expects constant value on the right of the operator.")
			}

			pathName := root.Left.Data.(filter.Token).Value
			//if strings.HasPrefix(strings.ToLower(pathName), strings.ToLower(resource.UserUrn+":")) {
			//	pathName = pathName[len(resource.UserUrn+":"):]
			//} else if strings.HasPrefix(strings.ToLower(pathName), strings.ToLower(resource.GroupUrn+":")) {
			//	pathName = pathName[len(resource.GroupUrn+":"):]
			//}

			//attribute = schema.GetAttribute(strings.ToLower(pathName))
			attribute = schema.GetAttribute(pathName)
			if nil == attribute {
				return nil, resource.CreateError(resource.InvalidFilter, fmt.Sprintf("Unknown path name '%s'", pathName))
			} else if attribute.Type == resource.Complex && token.Value != filter.Pr {
				return nil, resource.CreateError(resource.InvalidFilter, "Cannot perform eq, ne, co, sw, ew, gt, ge, le or lt operation on complex attributes.")
			}

			if token.Value == filter.Ge ||
				token.Value == filter.Gt ||
				token.Value == filter.Lt ||
				token.Value == filter.Le {
				if attribute.Type == resource.Boolean || attribute.Type == resource.Binary {
					return nil, resource.CreateError(resource.InvalidFilter, "Cannot determine order on boolean or binary attributes.")
				}
			}
		}

		switch token.Value {
		case filter.And:
			return bson.M{
				"$and": []interface{}{left, right},
			}, nil

		case filter.Or:
			return bson.M{
				"$or": []interface{}{left, right},
			}, nil

		case filter.Not:
			return bson.M{
				"$nor": []interface{}{left},
			}, nil

		case filter.Eq:
			if attribute.CaseExact || attribute.Type != resource.String {
				return bson.M{
					attribute.Assist.Path: bson.M{
						"$eq": root.Right.Data.(filter.Token).Params[filter.ParsedValue],
					},
				}, nil
			} else {
				return bson.M{
					attribute.Assist.Path: bson.M{
						"$regex": bson.RegEx{
							Pattern: fmt.Sprintf("^%s$", root.Right.Data.(filter.Token).Params[filter.ParsedValue]),
							Options: "i",
						},
					},
				}, nil
			}

		case filter.Ne:
			if attribute.CaseExact || attribute.Type != resource.String {
				return bson.M{
					attribute.Assist.Path: bson.M{
						"ne": root.Right.Data.(filter.Token).Params[filter.ParsedValue],
					},
				}, nil
			} else {
				return bson.M{
					"$nor": []interface{}{
						bson.M{
							attribute.Assist.Path: bson.M{
								"$regex": bson.RegEx{
									Pattern: fmt.Sprintf("^%x$", root.Right.Data.(filter.Token).Params[filter.ParsedValue]),
									Options: "i",
								},
							},
						},
					},
				}, nil
			}

		case filter.Co:
			if attribute.Type != resource.String {
				return nil, resource.CreateError(resource.InvalidFilter, "Cannot use co operator on non-string attributes.")
			} else if parsedValue, ok := root.Right.Data.(filter.Token).Params[filter.ParsedValue].(string); !ok {
				return nil, resource.CreateError(resource.InvalidFilter, "Cannot use co operator with non-string value.")
			} else {
				if attribute.CaseExact {
					return bson.M{
						attribute.Assist.Path: bson.M{
							"$regex": bson.RegEx{
								Pattern: parsedValue,
							},
						},
					}, nil
				} else {
					return bson.M{
						attribute.Assist.Path: bson.M{
							"$regex": bson.RegEx{
								Pattern: parsedValue,
								Options: "i",
							},
						},
					}, nil
				}
			}

		case filter.Sw:
			if attribute.Type != resource.String {
				return nil, resource.CreateError(resource.InvalidFilter, "Cannot use sw operator on non-string attributes.")
			} else if parsedValue, ok := root.Right.Data.(filter.Token).Params[filter.ParsedValue].(string); !ok {
				return nil, resource.CreateError(resource.InvalidFilter, "Cannot use sw operator with non-string value.")
			} else {
				if attribute.CaseExact {
					return bson.M{
						attribute.Assist.Path: bson.M{
							"$regex": bson.RegEx{
								Pattern: "^" + parsedValue,
							},
						},
					}, nil
				} else {
					return bson.M{
						attribute.Assist.Path: bson.M{
							"$regex": bson.RegEx{
								Pattern: "^" + parsedValue,
								Options: "i",
							},
						},
					}, nil
				}
			}

		case filter.Ew:
			if attribute.Type != resource.String {
				return nil, resource.CreateError(resource.InvalidFilter, "Cannot use ew operator on non-string attributes.")
			} else if parsedValue, ok := root.Right.Data.(filter.Token).Params[filter.ParsedValue].(string); !ok {
				return nil, resource.CreateError(resource.InvalidFilter, "Cannot use ew operator with non-string value.")
			} else {
				if attribute.CaseExact {
					return bson.M{
						attribute.Assist.Path: bson.M{
							"$regex": bson.RegEx{
								Pattern: parsedValue + "$",
							},
						},
					}, nil
				} else {
					return bson.M{
						attribute.Assist.Path: bson.M{
							"$regex": bson.RegEx{
								Pattern: parsedValue + "$",
								Options: "i",
							},
						},
					}, nil
				}
			}

		case filter.Gt:
			return bson.M{
				attribute.Assist.Path: bson.M{
					"$gt": root.Right.Data.(filter.Token).Params[filter.ParsedValue],
				},
			}, nil

		case filter.Ge:
			return bson.M{
				attribute.Assist.Path: bson.M{
					"$gte": root.Right.Data.(filter.Token).Params[filter.ParsedValue],
				},
			}, nil

		case filter.Lt:
			return bson.M{
				attribute.Assist.Path: bson.M{
					"$lt": root.Right.Data.(filter.Token).Params[filter.ParsedValue],
				},
			}, nil

		case filter.Le:
			return bson.M{
				attribute.Assist.Path: bson.M{
					"$lte": root.Right.Data.(filter.Token).Params[filter.ParsedValue],
				},
			}, nil

		case filter.Pr:
			existsCriteria := bson.M{attribute.Assist.Path: bson.M{"$exists": true}}
			nullCriteria := bson.M{attribute.Assist.Path: bson.M{"$ne": nil}}
			emptyStringCriteria := bson.M{attribute.Assist.Path: bson.M{"$ne": ""}}
			emptyArrayCriteria := bson.M{attribute.Assist.Path: bson.M{"$not": bson.M{"$size": 0}}}
			emptyObjectCriteria := bson.M{attribute.Assist.Path: bson.M{"$ne": bson.M{}}}

			criterion := make([]interface{}, 0)
			criterion = append(criterion, existsCriteria, nullCriteria)
			if attribute.MultiValued {
				criterion = append(criterion, emptyArrayCriteria)
			} else {
				if resource.String == attribute.Type {
					criterion = append(criterion, emptyStringCriteria)
				} else if resource.Complex == attribute.Type {
					criterion = append(criterion, emptyObjectCriteria)
				}
			}

			return bson.M{"$and": criterion}, nil

		default:
			return nil, resource.CreateError(resource.InvalidFilter, fmt.Sprintf("Invalid operator %s", token.Value))
		}
	}
}
