package validation

import (
	"context"
	"encoding/json"
	"github.com/go-scim/scimify/resource"
	"github.com/stretchr/testify/assert"
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"
	"github.com/go-scim/scimify/helper"
)

type typesCheckValidatorTest struct {
	name         string
	resourcePath string
	assertion    func(bool, error)
}

func BenchmarkRulesValidator_Validate(b *testing.B) {
	validator := &typeRulesValidator{}
	schema, _, err := helper.LoadSchema("../test_data/test_user_schema_all.json")
	if err != nil {
		b.Fatal(err)
	}
	schema.ConstructAttributeIndex()
	r, _, err := helper.LoadResource("../test_data/single_test_user_david.json")
	if err != nil {
		b.Fatal(err)
	}
	opt := ValidationOptions{UnassignedImmutableIsIgnored: false, ReadOnlyIsMandatory: false}
	ctx := context.WithValue(context.Background(), resource.CK_Schema, schema)

	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			_, err := validator.Validate(r, opt, ctx)
			if err != nil {
				b.Fatal(err)
			}
		}
	})
}

func TestTypeRulesValidator_Validate(t *testing.T) {
	validator := &typeRulesValidator{}
	for _, test := range []typesCheckValidatorTest{
		{
			"test valid resource",
			"../test_data/single_test_user_david.json",
			func(ok bool, err error) {
				assert.True(t, ok)
				assert.Nil(t, err)
			},
		},
		//{
		//	"test string type has number",
		//	"../test_data/bad_string_type_user.json",
		//	func(ok bool, err error) {
		//		assert.False(t, ok)
		//		assert.NotNil(t, err)
		//		assert.Equal(t, "urn:ietf:params:scim:schemas:core:2.0:User:displayName", err.(*validationError).FullPath)
		//	},
		//},
		//{
		//	"test invalid datetime format",
		//	"../test_data/bad_datetime_format_user.json",
		//	func(ok bool, err error) {
		//		assert.False(t, ok)
		//		assert.NotNil(t, err)
		//		assert.Equal(t, "meta.created", err.(*validationError).FullPath)
		//	},
		//},
		//{
		//	"test bool type has string",
		//	"../test_data/bad_bool_type_user.json",
		//	func(ok bool, err error) {
		//		assert.False(t, ok)
		//		assert.NotNil(t, err)
		//		assert.Equal(t, "urn:ietf:params:scim:schemas:core:2.0:User:active", err.(*validationError).FullPath)
		//	},
		//},
		//{
		//	"test array type has string",
		//	"../test_data/bad_array_type_user.json",
		//	func(ok bool, err error) {
		//		assert.False(t, ok)
		//		assert.NotNil(t, err)
		//		assert.Equal(t, "urn:ietf:params:scim:schemas:core:2.0:User:emails", err.(*validationError).FullPath)
		//	},
		//},
		//{
		//	"test complex type has string",
		//	"../test_data/bad_complex_type_user.json",
		//	func(ok bool, err error) {
		//		assert.False(t, ok)
		//		assert.NotNil(t, err)
		//		assert.Equal(t, "urn:ietf:params:scim:schemas:core:2.0:User:name", err.(*validationError).FullPath)
		//	},
		//},
		//{
		//	"test bad partial array type",
		//	"../test_data/bad_partial_array_type_user.json",
		//	func(ok bool, err error) {
		//		assert.False(t, ok)
		//		assert.NotNil(t, err)
		//		assert.Equal(t, "urn:ietf:params:scim:schemas:core:2.0:User:emails", err.(*validationError).FullPath)
		//	},
		//},
	} {
		// prepare schema
		schema := &resource.Schema{
			Schemas:    []string{resource.SchemaUrn},
			Id:         resource.UserUrn,
			Name:       "User schema",
			Attributes: make([]*resource.Attribute, 0),
		}
		coreSchema, err := loadSchema("../schemas/common_schema.json")
		userSchema, err := loadSchema("../schemas/user_schema.json")
		if err != nil {
			t.Fatal(err)
		}
		schema.MergeWith(coreSchema, userSchema)
		schema.ConstructAttributeIndex()

		// prepare test resource
		resourceData := loadTestDataFromJson(t, test.resourcePath)
		r := resource.NewResourceFromMap(resourceData)
		opt := ValidationOptions{UnassignedImmutableIsIgnored: false, ReadOnlyIsMandatory: false}
		ctx := context.WithValue(context.Background(), resource.CK_Schema, schema)

		ok, err := validator.Validate(r, opt, ctx)
		test.assertion(ok, err)
	}
}

func loadTestDataFromJson(t *testing.T, filePath string) map[string]interface{} {
	path, err := filepath.Abs(filePath)
	if err != nil {
		t.Fatal(err)
		return nil
	}

	file, err := os.Open(path)
	if err != nil {
		t.Fatal(err)
		return nil
	}
	defer file.Close()

	fileBytes, err := ioutil.ReadAll(file)
	if err != nil {
		t.Fatal(err)
		return nil
	}

	data := make(map[string]interface{})
	err = json.Unmarshal(fileBytes, &data)
	if err != nil {
		t.Fatal(err)
		return nil
	}

	return data
}

func loadSchema(filePath string) (*resource.Schema, error) {
	if path, err := filepath.Abs(filePath); err != nil {
		return nil, err
	} else if schema, err := resource.LoadSchema(path); err != nil {
		return nil, err
	} else {
		return schema, nil
	}
}
