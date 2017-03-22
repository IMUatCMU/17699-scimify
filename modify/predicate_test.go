package modify

import (
	"github.com/go-scim/scimify/filter"
	"github.com/go-scim/scimify/helper"
	"github.com/go-scim/scimify/resource"
	"github.com/stretchr/testify/assert"
	"reflect"
	"testing"
)

func TestEvaluatePredicate(t *testing.T) {
	r, _, err := helper.LoadResource("../test_data/single_test_user_david.json")
	if err != nil {
		t.Fatal(err)
	}

	sch, _, err := helper.LoadSchema("../test_data/test_user_schema_all.json")
	if err != nil {
		t.Fatal(err)
	}
	sch.ConstructAttributeIndex()

	for _, test := range []struct {
		name          string
		predicateText string
		expectResult  bool
	}{
		{
			"string equality (true)",
			"userName eq \"david@example.com\"",
			true,
		},
		{
			"string equality (false)",
			"userName eq \"foobar\"",
			false,
		},
		{
			"string equality case ignore (true)",
			"userName eq \"DAVID@example.com\"",
			true,
		},
		{
			"string equality not equal",
			"userName ne \"foobar\"",
			true,
		},
		{
			"string sw",
			"userName sw \"david\"",
			true,
		},
		{
			"string sw case ignore",
			"userName sw \"DAVID\"",
			true,
		},
		{
			"string ew",
			"userName ew \".com\"",
			true,
		},
		{
			"string ew case ignore",
			"userName ew \".COM\"",
			true,
		},
		{
			"string co (true)",
			"userName co \"david\"",
			true,
		},
		{
			"string co (false)",
			"userName co \"ark\"",
			false,
		},
		{
			"groups pr (false)",
			"group pr",
			false,
		},
		{
			"emails pr",
			"emails pr",
			true,
		},
	} {
		tokens, err := filter.Tokenize(test.predicateText)
		if err != nil {
			t.Error(err)
		}

		node, err := filter.Parse(tokens)
		if err != nil {
			t.Error(err)
		}

		predicateFunc := getPredicate(node)
		assert.NotNil(t, predicateFunc)

		result := predicateFunc(node, func(path string) (reflect.Value, *resource.Attribute) {
			attr := sch.GetAttribute(path)
			if attr == nil {
				return reflect.Value{}, nil
			}
			return reflect.ValueOf(r.Data()[attr.Assist.JSONName]), attr
		})

		assert.Equal(t, test.expectResult, result)
	}
}
