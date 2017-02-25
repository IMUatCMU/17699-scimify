package validation

import (
	"context"
	"github.com/go-scim/scimify/resource"
	"github.com/stretchr/testify/assert"
	"testing"
)

type mutabilityRulesTest struct {
	name         string
	resourcePath string
	assertion    func(bool, error)
}

func TestMutabilityRulesValidator_Validate(t *testing.T) {
	validator := &mutabilityRulesValidator{}

	for _, test := range []mutabilityRulesTest{
		{
			"test same resource passes",
			"../test_data/single_test_user_david.json",
			func(ok bool, err error) {
				assert.True(t, ok)
				assert.Nil(t, err)
			},
		},
		{
			"test changes readOnly string attribute",
			"../test_data/changes_readonly_string.json",
			func(ok bool, err error) {
				assert.False(t, ok)
				assert.NotNil(t, err)
				assert.Equal(t, "id", err.(*validationError).FullPath)
			},
		},
		{
			"test changes readOnly complex attribute",
			"../test_data/changes_readonly_complex.json",
			func(ok bool, err error) {
				assert.False(t, ok)
				assert.NotNil(t, err)
				assert.Equal(t, "meta", err.(*validationError).FullPath)
			},
		},
	} {
		schema, err := loadSchema("../test_data/test_user_schema_all.json")
		if err != nil {
			t.Fatal(err)
		}

		referenceData := loadTestDataFromJson(t, "../test_data/single_test_user_david.json")
		reference := resource.NewResourceFromMap(referenceData)

		resourceData := loadTestDataFromJson(t, test.resourcePath)
		r := resource.NewResourceFromMap(resourceData)

		opt := ValidationOptions{UnassignedImmutableIsIgnored: false, ReadOnlyIsMandatory: false}

		ctx := context.WithValue(context.Background(), resource.CK_Schema, schema)
		ctx = context.WithValue(ctx, resource.CK_Reference, reference)

		ok, err := validator.Validate(r, opt, ctx)
		test.assertion(ok, err)
	}
}
