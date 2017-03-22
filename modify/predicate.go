package modify

import (
	"github.com/go-scim/scimify/adt"
	"github.com/go-scim/scimify/filter"
	"github.com/go-scim/scimify/resource"
	"reflect"
	"strings"
)

type predicateFunc func(node *adt.Node, pfv pathForValueFunc) bool

type pathForValueFunc func(path string) (reflect.Value, *resource.Attribute)

func getPredicate(node *adt.Node) predicateFunc {
	tok := node.Data.(filter.Token)
	switch tok.Value {
	case filter.And:
		return andFunc
	case filter.Or:
		return orFunc
	case filter.Not:
		return notFunc
	case filter.Eq:
		return eqFunc
	case filter.Ne:
		return neFunc
	case filter.Sw:
		return swFunc
	case filter.Ew:
		return ewFunc
	case filter.Co:
		return coFunc
	case filter.Gt:
		return gtFunc
	case filter.Ge:
		return geFunc
	case filter.Lt:
		return ltFunc
	case filter.Le:
		return leFunc
	case filter.Pr:
		return prFunc
	default:
		return nil
	}
}

func andFunc(node *adt.Node, pfv pathForValueFunc) bool {
	lhs := getPredicate(node.Left)
	rhs := getPredicate(node.Right)
	return lhs(node.Left, pfv) && rhs(node.Right, pfv)
}

func orFunc(node *adt.Node, pfv pathForValueFunc) bool {
	lhs := getPredicate(node.Left)
	if lhs(node.Left, pfv) {
		return true
	} else {
		rhs := getPredicate(node.Right)
		return rhs(node.Right, pfv)
	}
}

func notFunc(node *adt.Node, pfv pathForValueFunc) bool {
	lhs := getPredicate(node.Left)
	return !lhs(node.Left, pfv)
}

func evaluate(node *adt.Node, pfv pathForValueFunc, f func(lhs, rhs reflect.Value, attr *resource.Attribute) bool) bool {
	v, attr := pfv(node.Left.Data.(filter.Token).Value)
	if attr == nil {
		return false
	}
	v0 := reflect.ValueOf(node.Right.Data.(filter.Token).Params[filter.ParsedValue])

	if attr.IsMultiValued() || attr.IsComplex() {
		return false
	}

	if !v.IsValid() || !v0.IsValid() {
		return false
	}

	if v.Kind() == reflect.Interface {
		v = v.Elem()
	}

	if v0.Kind() == reflect.Interface {
		v0 = v0.Elem()
	}

	return f(v, v0, attr)
}

func eqFunc(node *adt.Node, pfv pathForValueFunc) bool {
	return evaluate(node, pfv, func(v, v0 reflect.Value, attr *resource.Attribute) bool {
		if attr.Type == resource.Boolean {
			switch v.Kind() {
			case reflect.Bool:
				if v0.Kind() != reflect.Bool {
					return false
				} else {
					return v.Bool() == v0.Bool()
				}
			default:
				return false
			}
		} else {
			return compare(v, v0, attr, func(v, v0 string) bool {
				return v == v0
			}, func(v, v0 int64) bool {
				return v == v0
			}, func(v, v0 float64) bool {
				return v == v0
			})
		}
	})
}

func neFunc(node *adt.Node, pfv pathForValueFunc) bool {
	return !eqFunc(node, pfv)
}

func swFunc(node *adt.Node, pfv pathForValueFunc) bool {
	return evaluate(node, pfv, func(v, v0 reflect.Value, attr *resource.Attribute) bool {
		if v.Kind() != reflect.String || v0.Kind() != reflect.String {
			return false
		}

		if attr.CaseExact {
			return strings.HasPrefix(v.String(), v0.String())
		} else {
			return strings.HasPrefix(strings.ToLower(v.String()), strings.ToLower(v0.String()))
		}
	})
}

func ewFunc(node *adt.Node, pfv pathForValueFunc) bool {
	return evaluate(node, pfv, func(v, v0 reflect.Value, attr *resource.Attribute) bool {
		if v.Kind() != reflect.String || v0.Kind() != reflect.String {
			return false
		}

		if attr.CaseExact {
			return strings.HasSuffix(v.String(), v0.String())
		} else {
			return strings.HasSuffix(strings.ToLower(v.String()), strings.ToLower(v0.String()))
		}
	})
}

func coFunc(node *adt.Node, pfv pathForValueFunc) bool {
	return evaluate(node, pfv, func(v, v0 reflect.Value, attr *resource.Attribute) bool {
		if v.Kind() != reflect.String || v0.Kind() != reflect.String {
			return false
		}

		if attr.CaseExact {
			return strings.Contains(v.String(), v0.String())
		} else {
			return strings.Contains(strings.ToLower(v.String()), strings.ToLower(v0.String()))
		}
	})
}

func compare(v, v0 reflect.Value, attr *resource.Attribute, strF func(v, v0 string) bool, intF func(v, v0 int64) bool, floatF func(v, v0 float64) bool) bool {
	switch attr.Type {
	case resource.String, resource.DateTime:
		if v.Kind() != reflect.String || v0.Kind() != reflect.String {
			return false
		} else {
			if attr.CaseExact {
				return strF(v.String(), v0.String())
			} else {
				return strF(strings.ToLower(v.String()), strings.ToLower(v0.String()))
			}
		}

	case resource.Integer:
		switch v.Kind() {
		case reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
			if v0.Kind() != reflect.Int64 {
				return false
			} else {
				return intF(v.Int(), v0.Int())
			}
		default:
			return false
		}

	case resource.Decimal:
		switch v.Kind() {
		case reflect.Float32, reflect.Float64:
			if v0.Kind() != reflect.Float64 {
				return false
			} else {
				return floatF(v.Float(), v0.Float())
			}
		default:
			return false
		}

	default:
		return false
	}
}

func gtFunc(node *adt.Node, pfv pathForValueFunc) bool {
	return evaluate(node, pfv, func(v, v0 reflect.Value, attr *resource.Attribute) bool {
		return compare(v, v0, attr, func(v, v0 string) bool {
			return v > v0
		}, func(v, v0 int64) bool {
			return v > v0
		}, func(v, v0 float64) bool {
			return v > v0
		})
	})
}

func geFunc(node *adt.Node, pfv pathForValueFunc) bool {
	return evaluate(node, pfv, func(v, v0 reflect.Value, attr *resource.Attribute) bool {
		return compare(v, v0, attr, func(v, v0 string) bool {
			return v >= v0
		}, func(v, v0 int64) bool {
			return v >= v0
		}, func(v, v0 float64) bool {
			return v >= v0
		})
	})
}

func ltFunc(node *adt.Node, pfv pathForValueFunc) bool {
	return evaluate(node, pfv, func(v, v0 reflect.Value, attr *resource.Attribute) bool {
		return compare(v, v0, attr, func(v, v0 string) bool {
			return v < v0
		}, func(v, v0 int64) bool {
			return v < v0
		}, func(v, v0 float64) bool {
			return v < v0
		})
	})
}

func leFunc(node *adt.Node, pfv pathForValueFunc) bool {
	return evaluate(node, pfv, func(v, v0 reflect.Value, attr *resource.Attribute) bool {
		return compare(v, v0, attr, func(v, v0 string) bool {
			return v <= v0
		}, func(v, v0 int64) bool {
			return v <= v0
		}, func(v, v0 float64) bool {
			return v <= v0
		})
	})
}

func prFunc(node *adt.Node, pfv pathForValueFunc) bool {
	v, attr := pfv(node.Left.Data.(filter.Token).Value)
	if attr == nil {
		return false
	}
	return attr.IsValueAssigned(v)
}
