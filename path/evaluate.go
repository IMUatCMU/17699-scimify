package path

import (
	"github.com/go-scim/scimify/resource"
	"github.com/go-scim/scimify/adt"
	"reflect"
	"github.com/go-scim/scimify/filter"
	"strings"
)

type PathEvaluator struct {}

type EvaluationNotes struct {

}

func (pe *PathEvaluator) Evaluate(r *resource.Resource, paths adt.Queue, sch *resource.Schema) (reflect.Value, *EvaluationNotes, error) {
	notes := &EvaluationNotes{}
	v := pe.stepIn(reflect.ValueOf(r.Data()), paths, sch.AsAttribute(), notes)
	return v, notes, nil
}

func (pe *PathEvaluator) stepIn(v reflect.Value, paths adt.Queue, attr *resource.Attribute, notes *EvaluationNotes) reflect.Value {
	if paths.Size() == 0 {
		return v
	}

	if !v.IsValid() {
		return reflect.Value{}
	}

	if v.Kind() == reflect.Interface {
		return pe.stepIn(v.Elem(), paths, attr, notes)
	}

	pathNode := paths.Peek().(*adt.Node)
	switch pathNode.Data.(filter.Token).Type  {
	case filter.Path:
		if pathAttr := pe.findSubAttrThatMatchesPathToken(pathNode.Data.(filter.Token), attr); pathAttr == nil {
			pe.error("TODO err, no such path")
		} else {
			switch v.Kind() {
			case reflect.Array, reflect.Slice:
				resolvedVals := make([]interface{}, 0)
				for i := 0; i < v.Len(); i++ {
					clonedPath := paths.Clone()
					val := pe.stepIn(v.Index(i), clonedPath, attr, notes)
					if val.IsValid() {
						resolvedVals = append(resolvedVals, val.Interface())
					}
				}
				paths.Poll()
				return reflect.ValueOf(resolvedVals)

			case reflect.Map:
				paths.Poll()
				return pe.stepIn(v.MapIndex(reflect.ValueOf(pathAttr.Assist.JSONName)), paths, pathAttr, notes)

			default:
				pe.error("TODO err, cannot step in")
			}
		}

	case filter.Logical, filter.Relational:
		switch v.Kind() {
		case reflect.Array, reflect.Slice:
			matchedVals := make([]interface{}, 0)
			for i := 0; i < v.Len(); i++ {
				if pe.evaluatePredicate(v.Index(i), pathNode, attr, notes) {
					matchedVals = append(matchedVals, v.Index(i).Interface())
				}
			}
			paths.Poll()
			return pe.stepIn(reflect.ValueOf(matchedVals), paths, attr, notes)

		default:
			pe.error("TODO err, invalid path, cannot evaluate")
		}
		pe.error("not implemented yet")
	default:
		pe.error("TODO err, invalid path node type")
	}

	return reflect.Value{}
}

func (pe *PathEvaluator) evaluatePredicate(v reflect.Value, root *adt.Node, attr *resource.Attribute, notes *EvaluationNotes) bool {
	if !v.IsValid() {
		return false
	}

	if v.Kind() == reflect.Interface {
		return pe.evaluatePredicate(v.Elem(), root, attr, notes)
	}

	switch root.Data.(filter.Token).Value {
	case filter.And:
		return pe.evaluateAnd(v, root, attr, notes)

	case filter.Or:
		return pe.evaluateOr(v, root, attr, notes)

	case filter.Not:
		return pe.evaluateNot(v, root, attr, notes)

	case filter.Eq:
		return pe.evaluateEq(v, root, attr, notes)

	case filter.Ne:
		return pe.evaluateNe(v, root, attr, notes)

	case filter.Sw:
		return pe.evaluateSw(v, root, attr, notes)

	case filter.Ew:
		return pe.evaluateEw(v, root, attr, notes)

	case filter.Co:
		return pe.evaluateCo(v, root, attr, notes)

	case filter.Gt:
		return pe.evaluateGt(v, root, attr, notes)

	case filter.Ge:
		return pe.evaluateGe(v, root, attr, notes)

	case filter.Lt:
		return pe.evaluateLt(v, root, attr, notes)

	case filter.Le:
		return pe.evaluateLe(v, root, attr, notes)

	case filter.Pr:
		return pe.evaluatePr(v, root, attr, notes)

	default:
		return false
	}
}

func (pe *PathEvaluator) evaluatePr(v reflect.Value, root *adt.Node, attr *resource.Attribute, notes *EvaluationNotes) bool {
	pathAttr := pe.findSubAttrThatMatchesPathToken(root.Left.Data.(filter.Token), attr)
	if pathAttr == nil {
		pe.error("TODO err, invalid path (3):" + root.Left.Data.(filter.Token).Value)
	}

	val := v.MapIndex(reflect.ValueOf(pathAttr.Assist.JSONName))
	return pathAttr.IsValueAssigned(val)
}

func (pe *PathEvaluator) evaluateLe(v reflect.Value, root *adt.Node, attr *resource.Attribute, notes *EvaluationNotes) bool {
	return pe.evaluateComparable(v, root, attr, notes, func(lhs, rhs string) bool {
		return lhs <= rhs
	}, func (lhs, rhs int64) bool {
		return lhs <= rhs
	}, func (lhs, rhs float64) bool {
		return lhs <= rhs
	})
}

func (pe *PathEvaluator) evaluateLt(v reflect.Value, root *adt.Node, attr *resource.Attribute, notes *EvaluationNotes) bool {
	return pe.evaluateComparable(v, root, attr, notes, func(lhs, rhs string) bool {
		return lhs < rhs
	}, func (lhs, rhs int64) bool {
		return lhs < rhs
	}, func (lhs, rhs float64) bool {
		return lhs < rhs
	})
}

func (pe *PathEvaluator) evaluateGe(v reflect.Value, root *adt.Node, attr *resource.Attribute, notes *EvaluationNotes) bool {
	return pe.evaluateComparable(v, root, attr, notes, func(lhs, rhs string) bool {
		return lhs >= rhs
	}, func (lhs, rhs int64) bool {
		return lhs >= rhs
	}, func (lhs, rhs float64) bool {
		return lhs >= rhs
	})
}

func (pe *PathEvaluator) evaluateGt(v reflect.Value, root *adt.Node, attr *resource.Attribute, notes *EvaluationNotes) bool {
	return pe.evaluateComparable(v, root, attr, notes, func(lhs, rhs string) bool {
		return lhs > rhs
	}, func (lhs, rhs int64) bool {
		return lhs > rhs
	}, func (lhs, rhs float64) bool {
		return lhs > rhs
	})
}

func (pe *PathEvaluator) evaluateCo(v reflect.Value, root *adt.Node, attr *resource.Attribute, notes *EvaluationNotes) bool {
	return pe.evaluateTwoStrings(v, root, attr, notes, func(lhs, rhs string) bool {
		return strings.Contains(lhs, rhs)
	})
}

func (pe *PathEvaluator) evaluateEw(v reflect.Value, root *adt.Node, attr *resource.Attribute, notes *EvaluationNotes) bool {
	return pe.evaluateTwoStrings(v, root, attr, notes, func(lhs, rhs string) bool {
		return strings.HasSuffix(lhs, rhs)
	})
}

func (pe *PathEvaluator) evaluateSw(v reflect.Value, root *adt.Node, attr *resource.Attribute, notes *EvaluationNotes) bool {
	return pe.evaluateTwoStrings(v, root, attr, notes, func(lhs, rhs string) bool {
		return strings.HasPrefix(lhs, rhs)
	})
}

func (pe *PathEvaluator) evaluateNe(v reflect.Value, root *adt.Node, attr *resource.Attribute, notes *EvaluationNotes) bool {
	pathAttr := pe.findSubAttrThatMatchesPathToken(root.Left.Data.(filter.Token), attr)
	if pathAttr == nil {
		pe.error("TODO err, invalid path (3):" + root.Left.Data.(filter.Token).Value)
	}

	val := v.MapIndex(reflect.ValueOf(pathAttr.Assist.JSONName))
	if val.IsValid() {
		// TODO case sensitivity
		return !reflect.DeepEqual(val.Interface(), root.Right.Data.(filter.Token).Params[filter.ParsedValue])
	} else {
		return true
	}
}

func (pe *PathEvaluator) evaluateEq(v reflect.Value, root *adt.Node, attr *resource.Attribute, notes *EvaluationNotes) bool {
	pathAttr := pe.findSubAttrThatMatchesPathToken(root.Left.Data.(filter.Token), attr)
	if pathAttr == nil {
		pe.error("TODO err, invalid path (2):" + root.Left.Data.(filter.Token).Value)
	}

	val := v.MapIndex(reflect.ValueOf(pathAttr.Assist.JSONName))
	if val.IsValid() {
		// TODO case sensitivity
		return reflect.DeepEqual(val.Interface(), root.Right.Data.(filter.Token).Params[filter.ParsedValue])
	} else {
		return false
	}
}

func (pe *PathEvaluator) evaluateNot(v reflect.Value, root *adt.Node, attr *resource.Attribute, notes *EvaluationNotes) bool {
	return !pe.evaluatePredicate(v, root.Left, attr, notes)
}

func (pe *PathEvaluator) evaluateOr(v reflect.Value, root *adt.Node, attr *resource.Attribute, notes *EvaluationNotes) bool {
	lhs := pe.evaluatePredicate(v, root.Left, attr, notes)
	if lhs {
		return true
	} else {
		return pe.evaluatePredicate(v, root.Right, attr, notes)
	}
}

func (pe *PathEvaluator) evaluateAnd(v reflect.Value, root *adt.Node, attr *resource.Attribute, notes *EvaluationNotes) bool {
	return pe.evaluatePredicate(v, root.Left, attr, notes) && pe.evaluatePredicate(v, root.Right, attr, notes)
}

func (pe *PathEvaluator) evaluateTwoStrings(v reflect.Value, root *adt.Node, attr *resource.Attribute, notes *EvaluationNotes, comparator func(lhs, rhs string) bool) bool {
	constantTok := root.Right.Data.(filter.Token)
	if constantTok.Params[filter.ConstantType] != filter.ConstString {
		return false
	}

	pathAttr := pe.findSubAttrThatMatchesPathToken(root.Left.Data.(filter.Token), attr)
	if pathAttr == nil {
		pe.error("TODO err, invalid path (3):" + root.Left.Data.(filter.Token).Value)
	}

	val := v.MapIndex(reflect.ValueOf(pathAttr.Assist.JSONName))

	if val.IsValid() {
		switch val.Kind() {
		case reflect.Interface:
			val = val.Elem()
			fallthrough

		case reflect.String:
			if pathAttr.CaseExact {
				return comparator(
					val.String(),
					constantTok.Params[filter.ParsedValue].(string),
				)
			} else {
				return comparator(
					strings.ToLower(val.String()),
					strings.ToLower(constantTok.Params[filter.ParsedValue].(string)),
				)
			}
		default:
			return false
		}
	} else {
		return false
	}
}

func (pe *PathEvaluator) evaluateComparable(v reflect.Value, root *adt.Node, attr *resource.Attribute, notes *EvaluationNotes, stringComparator func(lhs, rhs string) bool, intComparator func(lhs, rhs int64) bool, floatComparator func(lhs, rhs float64) bool) bool {
	constantTok := root.Right.Data.(filter.Token)
	pathAttr := pe.findSubAttrThatMatchesPathToken(root.Left.Data.(filter.Token), attr)
	if pathAttr == nil {
		pe.error("TODO err, invalid path (3):" + root.Left.Data.(filter.Token).Value)
	}

	val := v.MapIndex(reflect.ValueOf(pathAttr.Assist.JSONName))
	switch val.Kind() {
	case reflect.Interface:
		val = val.Elem()
		fallthrough

	case reflect.String:
		if constantTok.Params[filter.ConstantType] == filter.ConstString {
			if pathAttr.CaseExact {
				return stringComparator(val.String(), constantTok.Params[filter.ParsedValue].(string))
			} else {
				return stringComparator(strings.ToLower(val.String()), strings.ToLower(constantTok.Params[filter.ParsedValue].(string)))
			}
		} else {
			return false
		}
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		if constantTok.Params[filter.ConstantType] == filter.ConstInteger {
			return intComparator(val.Int(), constantTok.Params[filter.ParsedValue].(int64))
		} else {
			return false
		}
	case reflect.Float32, reflect.Float64:
		if constantTok.Params[filter.ConstantType] == filter.ConstDecimal {
			return floatComparator(val.Float(), constantTok.Params[filter.ParsedValue].(float64))
		} else {
			return false
		}
	default:
		return false
	}
}

func (pe *PathEvaluator) findSubAttrThatMatchesPathToken(tok filter.Token, attr *resource.Attribute) *resource.Attribute {
	for _, subAttr := range attr.SubAttributes {
		if strings.ToLower(subAttr.Assist.JSONName) == strings.ToLower(tok.Value) {
			return subAttr
		}
	}
	return nil
}

func (pe *PathEvaluator) error(err interface{}) {
	panic(err)
}