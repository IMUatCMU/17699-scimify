package helper

import (
	"testing"
	"strings"
	"fmt"
	"github.com/stretchr/testify/assert"
)

type testMapAggregator struct {
	state 	[]string
}
func (g *testMapAggregator) Aggregate(input interface{}) {
	if nil == g.state {
		g.state = make([]string, 0)
	}
	g.state = append(g.state, input.(string))
}
func (g *testMapAggregator) Result() interface{} {
	return "{" + strings.Join(g.state, ",") + "}"
}

type testSliceAggregator struct {
	state	[]string
}
func (g *testSliceAggregator) Aggregate(input interface{}) {
	if nil == g.state {
		g.state = make([]string, 0)
	}
	g.state = append(g.state, input.(string))
}
func (g *testSliceAggregator) Result() interface{} {
	return "[" + strings.Join(g.state, ",") + "]"
}

func TestWalkSliceInParallel(t *testing.T) {
	for _, test := range []struct{
		name 		string
		target 		[]interface{}
		processFunc	SliceElementProcessor
		aggregator	Aggregator
		assertion	func(interface{}, error)
	}{
		{
			"process simple slice",
			[]interface{}{"1", "2", "3"},
			func(_ int, elem interface{}) (interface{}, error) {
				return elem, nil
			},
			&testSliceAggregator{state:make([]string, 0)},
			func(result interface{}, err error) {
				assert.Nil(t, err)
				assert.True(t, result == "[1,2,3]" ||
					result == "[1,3,2]" ||
					result == "[2,1,3]" ||
					result == "[2,3,1]" ||
					result == "[3,1,2]" ||
					result == "[3,2,1]")
			},
		},
		{
			"process empty slice",
			[]interface{}{},
			func(_ int, elem interface{}) (interface{}, error) {
				return elem, nil
			},
			&testSliceAggregator{state:make([]string, 0)},
			func(result interface{}, err error) {
				assert.Nil(t, err)
				assert.Equal(t, "[]", result.(string))
			},
		},
	}{
		result, err := WalkSliceInParallel(test.target, test.processFunc, test.aggregator)
		test.assertion(result, err)
	}
}

func TestWalkStringMapInParallel(t *testing.T) {
	for _, test := range []struct{
		name 		string
		target		map[string]interface{}
		processFunc	MapEntryProcessor
		aggregator	Aggregator
		assertion	func(interface{}, error)
	}{
		{
			"process simple map",
			map[string]interface{}{
				"foo": "bar",
				"bar": 3,
			},
			func(key string, value interface{}) (interface{}, error) {
				return fmt.Sprintf("%s:%v", key, value), nil
			},
			&testMapAggregator{state:make([]string, 0)},
			func(result interface{}, err error) {
				assert.Nil(t, err)
				assert.True(t, result == "{foo:bar,bar:3}" || result == "{bar:3,foo:bar}")
			},
		},
		{
			"process empty map",
			map[string]interface{}{},
			func(key string, value interface{}) (interface{}, error) {
				return fmt.Sprintf("%s:%v", key, value), nil
			},
			&testMapAggregator{state:make([]string, 0)},
			func(result interface{}, err error) {
				assert.Nil(t, err)
				assert.Equal(t, "{}", result.(string))
			},
		},
	}{
		result, err := WalkStringMapInParallel(test.target, test.processFunc, test.aggregator)
		test.assertion(result, err)
	}
}
