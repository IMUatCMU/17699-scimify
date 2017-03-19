package path

import (
	"testing"
	"github.com/go-scim/scimify/helper"
)

func TestPathEvaluator_Evaluate(t *testing.T) {
	pe := &PathEvaluator{}

	r, _, err := helper.LoadResource("../test_data/single_test_user_david.json")
	if err != nil {
		t.Fatal(err)
	}

	sch, _, err := helper.LoadSchema("../test_data/test_user_schema_all.json")
	if err != nil {
		t.Fatal(err)
	}

	q, err := Tokenize("emails[primary pr and value co \"home\"]")
	if err != nil {
		t.Fatal(err)
	}

	v, _, err := pe.Evaluate(r, q, sch)
	if err != nil {
		t.Fatal(err)
	}

	t.Logf("%+v", v.Interface())
}
