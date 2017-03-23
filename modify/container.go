package modify

import (
	"fmt"
	"github.com/go-scim/scimify/adt"
	"github.com/go-scim/scimify/filter"
	"github.com/go-scim/scimify/resource"
	"reflect"
)

type Container interface {
	Apply(paths adt.Queue, parentAttr *resource.Attribute, unit ModUnit, replaceFunc func(reflect.Value)) error
}

// No-op container to silently fail missing remove paths
type modNoOp int

func (m modNoOp) Apply(paths adt.Queue, parentAttr *resource.Attribute, unit ModUnit, replaceFunc func(reflect.Value)) error {
	return nil
}

// array container
type modArray []interface{}

func (m modArray) Apply(paths adt.Queue, parentAttr *resource.Attribute, unit ModUnit, replaceFunc func(reflect.Value)) error {
	if paths.Size() == 0 {
		panic("illegal state: nothing for array container to operate on.")
	}

	var attr *resource.Attribute
	predicate := func(index int) bool {
		return true
	}

	pNode := paths.Poll().(*adt.Node)
	switch pNode.Data.(filter.Token).Type {
	case filter.Relational, filter.Logical:
		predicate = func(index int) bool {
			return getPredicate(pNode)(pNode, func(path string) (reflect.Value, *resource.Attribute) {
				attr := parentAttr.GetAttribute(path)
				if attr == nil {
					return reflect.Value{}, nil
				}
				return reflect.ValueOf(m[index].(map[string]interface{})[attr.Assist.JSONName]), attr
			})
		}

		if paths.Size() == 0 {
			attr = parentAttr.Clone()
			attr.MultiValued = false
			switch unit.Op {
			case opReplace:
				m.replace(unit.Value, predicate)
			case opRemove:
				m.remove(predicate, replaceFunc)
			case opAdd:
				return &InvalidModificationError{"cannot add on filtered array"}
			default:
				return &InvalidModificationError{
					fmt.Sprintf("unsupported op %s", unit.Op),
				}
			}
		} else {
			pathTok := paths.Poll().(*adt.Node).Data.(filter.Token) // i.e. emails[<filter>].value
			if paths.Size() != 0 {
				return &InvalidPathError{
					path:   pathTok.Value,
					reason: "complex array is one level only.",
				}
			}
			attr = parentAttr.GetAttribute(pathTok.Value)
			if attr == nil {
				return &MissingAttributeForPathError{pathTok.Value}
			}
			switch unit.Op {
			case opAdd:
				m.addSub(unit.Value, attr, predicate)
			case opReplace:
				m.replaceSub(unit.Value, attr, predicate)
			case opRemove:
				m.removeSub(attr, predicate)
			}
		}

	default:
		if paths.Size() != 0 {
			return &InvalidPathError{reason: "complex array is one level only."}
		}

		attr = parentAttr.GetAttribute(pNode.Data.(filter.Token).Value)
		if attr == nil {
			return &MissingAttributeForPathError{pNode.Data.(filter.Token).Value}
		}

		switch unit.Op {
		case opAdd:
			m.addSub(unit.Value, attr, predicate)
		case opReplace:
			m.replaceSub(unit.Value, attr, predicate)
		case opRemove:
			m.removeSub(attr, predicate)
		}
	}

	return nil
}

func (m modArray) addSub(val interface{}, attr *resource.Attribute, predicate func(index int) bool) {
	arrVal := reflect.ValueOf(m)
	for i := 0; i < arrVal.Len(); i++ {
		if predicate(i) {
			mSub := arrVal.Index(i)
			if mSub.Kind() == reflect.Interface {
				mSub = mSub.Elem()
			}
			mSub.SetMapIndex(reflect.ValueOf(attr.Assist.JSONName), reflect.ValueOf(val))
		}
	}
}

func (m modArray) replaceSub(val interface{}, attr *resource.Attribute, predicate func(index int) bool) {
	arrVal := reflect.ValueOf(m)
	for i := 0; i < arrVal.Len(); i++ {
		if predicate(i) {
			mSub := arrVal.Index(i)
			if mSub.Kind() == reflect.Interface {
				mSub = mSub.Elem()
			}
			mSub.SetMapIndex(reflect.ValueOf(attr.Assist.JSONName), reflect.ValueOf(val))
		}
	}
}

func (m modArray) removeSub(attr *resource.Attribute, predicate func(index int) bool) {
	arrVal := reflect.ValueOf(m)
	for i := 0; i < arrVal.Len(); i++ {
		if predicate(i) {
			mSub := arrVal.Index(i)
			if mSub.Kind() == reflect.Interface {
				mSub = mSub.Elem()
			}
			mSub.SetMapIndex(reflect.ValueOf(attr.Assist.JSONName), reflect.Value{})
		}
	}
}

func (m modArray) replace(val interface{}, predicate func(index int) bool) {
	arrVal := reflect.ValueOf(m)
	for i := 0; i < arrVal.Len(); i++ {
		if predicate(i) {
			mSub := arrVal.Index(i)
			if mSub.Kind() == reflect.Interface {
				mSub = mSub.Elem()
			}
			mSub.Set(reflect.ValueOf(val))
		}
	}
}

func (m modArray) remove(predicate func(index int) bool, replaceFunc func(reflect.Value)) {
	newArr := make([]interface{}, 0)
	for i := 0; i < len(m); i++ {
		if !predicate(i) {
			newArr = append(newArr, m[i])
		}
	}
	replaceFunc(reflect.ValueOf(newArr))
}

// map container
type modMap map[string]interface{}

func (m modMap) add(val interface{}, attr *resource.Attribute) {
	mapVal := reflect.ValueOf(m)
	keyVal := reflect.ValueOf(attr.Assist.JSONName)
	dataVal := reflect.ValueOf(val)

	entryVal := mapVal.MapIndex(keyVal)
	if !entryVal.IsValid() {
		entryVal = attr.ZeroValue()
	}

	if entryVal.Kind() == reflect.Interface {
		entryVal = entryVal.Elem()
	}

	switch {
	// array array
	case attr.IsMultiValued():
		switch dataVal.Kind() {
		case reflect.Slice, reflect.Array:
			entryVal = reflect.AppendSlice(entryVal, dataVal)
		default:
			entryVal = reflect.Append(entryVal, dataVal)
		}

	// complex object
	case attr.IsComplex() && !attr.IsMultiValued():
		for _, dataKeyVal := range dataVal.MapKeys() {
			entryVal.SetMapIndex(dataKeyVal, dataVal.MapIndex(dataKeyVal))
		}

	// simple field
	case !attr.IsComplex() && !attr.IsMultiValued():
		entryVal = dataVal
	}

	mapVal.SetMapIndex(keyVal, entryVal)
}

func (m modMap) replace(val interface{}, attr *resource.Attribute) {
	mapVal := reflect.ValueOf(m)
	keyVal := reflect.ValueOf(attr.Assist.JSONName)
	dataVal := reflect.ValueOf(val)
	mapVal.SetMapIndex(keyVal, dataVal)
}

func (m modMap) remove(attr *resource.Attribute) {
	mapVal := reflect.ValueOf(m)
	keyVal := reflect.ValueOf(attr.Assist.JSONName)
	mapVal.SetMapIndex(keyVal, reflect.Value{})
}

func (m modMap) Apply(paths adt.Queue, parentAttr *resource.Attribute, unit ModUnit, replaceFunc func(reflect.Value)) error {
	switch paths.Size() {
	case 0:
		val := reflect.ValueOf(unit.Value)
		if val.Kind() != reflect.Map {
			return &InvalidModificationError{
				reason: "implicit path modification must use map value",
			}
		}

		for _, kVal := range val.MapKeys() {
			nextTok, err := filter.CreateToken(kVal.String())
			if err != nil {
				return err
			}

			nextPaths := adt.NewQueue(1)
			nextPaths.Offer(adt.NewNode(nextTok))

			nextModUnit := ModUnit{
				Op:    unit.Op,
				Path:  kVal.String(),
				Value: val.MapIndex(kVal).Interface(),
			}

			err = m.Apply(nextPaths, parentAttr, nextModUnit, nil)
			if err != nil {
				return err
			}
		}

	case 1:
		p := paths.Poll().(*adt.Node).Data.(filter.Token)
		if p.Type != filter.Path {
			return &InvalidPathError{
				path:   p.Value,
				reason: fmt.Sprintf("unexpected token type '%s'", p.Type),
			}
		}

		attr := parentAttr.GetAttribute(p.Value)
		if attr == nil {
			return &MissingAttributeForPathError{p.Value}
		}

		switch unit.Op {
		case opAdd:
			m.add(unit.Value, attr)
		case opReplace:
			m.replace(unit.Value, attr)
		case opRemove:
			m.remove(attr)
		default:
			return &InvalidModificationError{
				fmt.Sprintf("unsupported op %s", unit.Op),
			}
		}

	default:
		p := paths.Poll().(*adt.Node).Data.(filter.Token)
		if p.Type != filter.Path {
			return &InvalidPathError{
				path:   p.Value,
				reason: fmt.Sprintf("unexpected token type '%s'", p.Type),
			}
		}

		nextCon, nextPAttr, err := m.nextContainer(p, parentAttr, unit)
		if err != nil {
			return err
		}
		return nextCon.Apply(paths, nextPAttr, unit, func(val reflect.Value) {
			mVal := reflect.ValueOf(m)
			mVal.SetMapIndex(reflect.ValueOf(nextPAttr.Assist.JSONName), val)
		})
	}

	return nil
}

func (m modMap) nextContainer(p filter.Token, parentAttr *resource.Attribute, unit ModUnit) (Container, *resource.Attribute, error) {
	attr := parentAttr.GetAttribute(p.Value)
	if attr == nil {
		return nil, nil, &MissingAttributeForPathError{p.Value}
	}

	cVal := reflect.ValueOf(m)
	kVal := reflect.ValueOf(attr.Assist.JSONName)
	val := cVal.MapIndex(kVal)
	if !val.IsValid() {
		switch unit.Op {
		case opAdd:
			cVal.SetMapIndex(kVal, attr.ZeroValue())
		case opRemove:
			return modNoOp(0), nil, nil
		default:
			return nil, nil, &InvalidPathError{
				path:   p.Value,
				reason: fmt.Sprintf("no value is found at path (component) %s", p.Value),
			}
		}
	}

	if reflect.Interface == val.Kind() {
		val = val.Elem()
	}

	switch val.Kind() {
	case reflect.Map:
		if attr.IsMultiValued() || !attr.IsComplex() {
			return nil, nil, &InvalidPathError{
				path:   p.Value,
				reason: "unexpected complex object",
			}
		}
		return modMap(val.Interface().(map[string]interface{})), attr, nil

	case reflect.Array, reflect.Slice:
		if !attr.IsMultiValued() {
			return nil, nil, &InvalidPathError{
				path:   p.Value,
				reason: "unexpected array",
			}
		}
		return modArray(val.Interface().([]interface{})), attr, nil

	default:
		return nil, nil, &InvalidPathError{
			path:   p.Value,
			reason: fmt.Sprintf("terminal data type %s encountered in the middle", val.Type()),
		}
	}
}
