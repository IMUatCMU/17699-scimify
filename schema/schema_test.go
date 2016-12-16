package schema

import (
	"path/filepath"
	"testing"
	"github.com/stretchr/testify/assert"
)

func TestLoadSchema(t *testing.T) {
	if path, err := filepath.Abs("../schemas/common_schema.json"); err != nil {
		t.Fatal(err)
	} else if schema, err := LoadSchema(path); err != nil {
		t.Error(err)
	} else {
		assert.Contains(t, schema.Schemas, CommonUrn)
		assert.Equal(t, CommonUrn, schema.Id)
		assert.Equal(t, "Common Attributes", schema.Name)
		assert.Equal(t, 4, len(schema.Attributes))
		assertAttribute(t, &Attribute{
			Name:"schemas",
			Type:"reference",
			MultiValued:true,
			Required:true,
			CaseExact:true,
			Mutability:"readWrite",
			Returned:"default",
			Uniqueness:"none",
			ReferenceTypes:[]string{ "uri" },
			CanonicalValues:[]string{
				"urn:ietf:params:scim:schemas:core:2.0:User",
				"urn:ietf:params:scim:schemas:core:2.0:Group",
				"urn:ietf:params:scim:schemas:core:2.0:ResourceType",
				"urn:ietf:params:scim:schemas:core:2.0:ServiceProviderConfig",
				"urn:ietf:params:scim:schemas:core:2.0:Schema",
			},
			SubAttributes:make([]*Attribute, 0),
		}, schema.Attributes[0])
		assertAttribute(t, &Attribute{
			Name:"id",
			Type:"string",
			MultiValued:false,
			Required:true,
			CaseExact:true,
			Mutability:"readOnly",
			Returned:"always",
			Uniqueness:"global",
			ReferenceTypes:[]string{},
			CanonicalValues:[]string{},
			SubAttributes:make([]*Attribute, 0),
		}, schema.Attributes[1])
		assertAttribute(t, &Attribute{
			Name:"externalId",
			Type:"string",
			MultiValued:false,
			Required:false,
			CaseExact:true,
			Mutability:"readWrite",
			Returned:"default",
			Uniqueness:"none",
			ReferenceTypes:[]string{},
			CanonicalValues:[]string{},
			SubAttributes:make([]*Attribute, 0),
		}, schema.Attributes[2])
		assertAttribute(t, &Attribute{
			Name:"meta",
			Type:"complex",
			MultiValued:false,
			Required:false,
			CaseExact:false,
			Mutability:"readOnly",
			Returned:"default",
			Uniqueness:"none",
			ReferenceTypes:[]string{},
			CanonicalValues:[]string{},
			SubAttributes:[]*Attribute{
				{
					Name:"resourceType",
					Type:"string",
					MultiValued:false,
					Required:false,
					CaseExact:true,
					Mutability:"readOnly",
					Returned:"default",
					Uniqueness:"none",
					ReferenceTypes:[]string{},
					CanonicalValues:[]string{},
					SubAttributes:make([]*Attribute, 0),
				},
				{
					Name:"created",
					Type:"datetime",
					MultiValued:false,
					Required:false,
					CaseExact:true,
					Mutability:"readOnly",
					Returned:"default",
					Uniqueness:"none",
					ReferenceTypes:[]string{},
					CanonicalValues:[]string{},
					SubAttributes:make([]*Attribute, 0),
				},
				{
					Name:"lastModified",
					Type:"datetime",
					MultiValued:false,
					Required:false,
					CaseExact:true,
					Mutability:"readOnly",
					Returned:"default",
					Uniqueness:"none",
					ReferenceTypes:[]string{},
					CanonicalValues:[]string{},
					SubAttributes:make([]*Attribute, 0),
				},
				{
					Name:"location",
					Type:"reference",
					MultiValued:false,
					Required:false,
					CaseExact:true,
					Mutability:"readOnly",
					Returned:"default",
					Uniqueness:"none",
					ReferenceTypes:[]string{ "uri" },
					CanonicalValues:[]string{},
					SubAttributes:make([]*Attribute, 0),
				},
				{
					Name:"version",
					Type:"string",
					MultiValued:false,
					Required:false,
					CaseExact:true,
					Mutability:"readOnly",
					Returned:"default",
					Uniqueness:"none",
					ReferenceTypes:[]string{},
					CanonicalValues:[]string{},
					SubAttributes:make([]*Attribute, 0),
				},
			},
		}, schema.Attributes[3])
	}
}

func assertAttribute(t *testing.T, expected, actual *Attribute) {
	assert.NotNil(t, actual)
	assert.Equal(t, expected.Name, actual.Name)
	assert.Equal(t, expected.Type, actual.Type)
	assert.Equal(t, expected.MultiValued, actual.MultiValued)
	assert.Equal(t, expected.Required, actual.Required)
	assert.Equal(t, expected.CaseExact, actual.CaseExact)
	assert.Equal(t, expected.Mutability, actual.Mutability)
	assert.Equal(t, expected.Returned, actual.Returned)
	assert.Equal(t, expected.Uniqueness, actual.Uniqueness)
	assertStringArrayEquals(t, expected.CanonicalValues, actual.CanonicalValues)
	assertStringArrayEquals(t, expected.ReferenceTypes, actual.ReferenceTypes)
	assertAttributeArrayEquals(t, expected.SubAttributes, actual.SubAttributes)
}

func assertAttributeArrayEquals(t *testing.T, expected, actual []*Attribute) {
	if nil == expected {
		assert.Equal(t, 0, len(actual))
	} else {
		for i, expectedElem := range expected {
			assertAttribute(t, expectedElem, actual[i])
		}
	}
}

func assertStringArrayEquals(t *testing.T, expected, actual []string) {
	if nil == expected {
		assert.Equal(t, 0, len(actual))
	} else {
		for _, elem := range expected {
			assert.Contains(t, actual, elem)
		}
	}
}