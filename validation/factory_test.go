package validation

import (
	"github.com/go-scim/scimify/resource"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestValidator(t *testing.T) {
	validator := GetValidator()
	assert.NotNil(t, validator)

	for _, test := range []struct {
		name             string
		resourcePath     string
		contextDecorator func(*ValidatorContext) *ValidatorContext
		assertion        func(bool, error)
	}{
		{
			"correct resource",
			"../test_data/single_test_user_david.json",
			nil,
			func(ok bool, err error) {
				assert.True(t, ok)
				assert.Nil(t, err)
			},
		},
		{
			"bad_array_type_user",
			"../test_data/bad_array_type_user.json",
			nil,
			func(ok bool, err error) {
				assert.False(t, ok)
				assert.NotNil(t, err)
				assert.Equal(t, typeCheck, err.(*validationError).ViolationType)
				assert.Equal(t, "urn:ietf:params:scim:schemas:core:2.0:User:emails", err.(*validationError).FullPath)
			},
		},
		{
			"bad_bool_type_user",
			"../test_data/bad_bool_type_user.json",
			nil,
			func(ok bool, err error) {
				assert.False(t, ok)
				assert.NotNil(t, err)
				assert.Equal(t, typeCheck, err.(*validationError).ViolationType)
				assert.Equal(t, "urn:ietf:params:scim:schemas:core:2.0:User:active", err.(*validationError).FullPath)
			},
		},
		{
			"bad_complex_type_user",
			"../test_data/bad_complex_type_user.json",
			nil,
			func(ok bool, err error) {
				assert.False(t, ok)
				assert.NotNil(t, err)
				assert.Equal(t, typeCheck, err.(*validationError).ViolationType)
				assert.Equal(t, "urn:ietf:params:scim:schemas:core:2.0:User:name", err.(*validationError).FullPath)
			},
		},
		{
			"bad_datetime_format_user",
			"../test_data/bad_datetime_format_user.json",
			nil,
			func(ok bool, err error) {
				assert.False(t, ok)
				assert.NotNil(t, err)
				assert.Equal(t, typeCheck, err.(*validationError).ViolationType)
				assert.Equal(t, "meta.created", err.(*validationError).FullPath)
			},
		},
		{
			"bad_partial_array_type_user",
			"../test_data/bad_partial_array_type_user.json",
			nil,
			func(ok bool, err error) {
				assert.False(t, ok)
				assert.NotNil(t, err)
				assert.Equal(t, typeCheck, err.(*validationError).ViolationType)
				assert.Equal(t, "urn:ietf:params:scim:schemas:core:2.0:User:emails", err.(*validationError).FullPath)
			},
		},
		{
			"bad_string_type_user",
			"../test_data/bad_string_type_user.json",
			nil,
			func(ok bool, err error) {
				assert.False(t, ok)
				assert.NotNil(t, err)
				assert.Equal(t, typeCheck, err.(*validationError).ViolationType)
				assert.Equal(t, "urn:ietf:params:scim:schemas:core:2.0:User:displayName", err.(*validationError).FullPath)
			},
		},
		{
			"missing_schemas_user",
			"../test_data/missing_schemas_user.json",
			nil,
			func(ok bool, err error) {
				assert.False(t, ok)
				assert.NotNil(t, err)
				assert.Equal(t, requiredMissing, err.(*validationError).ViolationType)
				assert.Equal(t, "schemas", err.(*validationError).FullPath)
			},
		},
		{
			"missing_id_user",
			"../test_data/missing_id_user.json",
			func(ctx *ValidatorContext) *ValidatorContext {
				ctx.Lock()
				ctx.Data[FailReadOnlyRequired] = true
				ctx.Unlock()
				return ctx
			},
			func(ok bool, err error) {
				assert.False(t, ok)
				assert.NotNil(t, err)
				assert.Equal(t, requiredMissing, err.(*validationError).ViolationType)
				assert.Equal(t, "id", err.(*validationError).FullPath)
			},
		},
		{
			"missing_meta_user",
			"../test_data/missing_meta_user.json",
			func(ctx *ValidatorContext) *ValidatorContext {
				ctx.Lock()
				ctx.Data[FailReadOnlyRequired] = true
				ctx.Unlock()
				return ctx
			},
			func(ok bool, err error) {
				assert.True(t, ok)
				assert.Nil(t, err)
			},
		},
		{
			"user with no change",
			"../test_data/single_test_user_david.json",
			func(ctx *ValidatorContext) *ValidatorContext {
				data := loadTestDataFromJson(t, "../test_data/single_test_user_david.json")
				ref := resource.NewResourceFromMap(data)
				ctx.Lock()
				ctx.Data[ReferenceResource] = ref
				ctx.Unlock()
				return ctx
			},
			func(ok bool, err error) {
				assert.True(t, ok)
				assert.Nil(t, err)
			},
		},
		{
			"user with id changed",
			"../test_data/single_test_user_david.json",
			func(ctx *ValidatorContext) *ValidatorContext {
				data := loadTestDataFromJson(t, "../test_data/single_test_user_david.json")
				ref := resource.NewResourceFromMap(data)
				ref.Id = "foo"
				ctx.Lock()
				ctx.Data[ReferenceResource] = ref
				ctx.Unlock()
				return ctx
			},
			func(ok bool, err error) {
				assert.False(t, ok)
				assert.NotNil(t, err)
				assert.Equal(t, mutabilityCheck, err.(*validationError).ViolationType)
				assert.Equal(t, "id", err.(*validationError).FullPath)
			},
		},
		{
			"user with meta changed",
			"../test_data/single_test_user_david.json",
			func(ctx *ValidatorContext) *ValidatorContext {
				data := loadTestDataFromJson(t, "../test_data/single_test_user_david.json")
				ref := resource.NewResourceFromMap(data)
				ref.Meta.Version = "foo"
				ctx.Lock()
				ctx.Data[ReferenceResource] = ref
				ctx.Unlock()
				return ctx
			},
			func(ok bool, err error) {
				assert.False(t, ok)
				assert.NotNil(t, err)
				assert.Equal(t, mutabilityCheck, err.(*validationError).ViolationType)
				assert.Equal(t, "meta", err.(*validationError).FullPath)
			},
		},
		{
			"user with group changed",
			"../test_data/single_test_user_david.json",
			func(ctx *ValidatorContext) *ValidatorContext {
				data := loadTestDataFromJson(t, "../test_data/single_test_user_david.json")
				ref := resource.NewResourceFromMap(data)
				ref.Attributes["groups"] = []interface{}{
					map[string]interface{}{
						"value": "bar",
					},
				}
				ctx.Lock()
				ctx.Data[ReferenceResource] = ref
				ctx.Unlock()
				return ctx
			},
			func(ok bool, err error) {
				assert.True(t, ok) // defaults will overwrite the original nil
				assert.Nil(t, err)
			},
		},
	} {
		schema, err := loadSchema("../test_data/test_user_schema_all.json")
		if err != nil {
			t.Fatal(err)
		}

		data := loadTestDataFromJson(t, test.resourcePath)
		resource := resource.NewResourceFromMap(data)

		context := &ValidatorContext{
			Data: map[string]interface{}{
				Schema: schema,
			},
		}
		if nil != test.contextDecorator {
			context = test.contextDecorator(context)
		}
		ok, err := validator.Validate(resource, context)
		test.assertion(ok, err)
	}
}

func BenchmarkValidator(b *testing.B) {

}
