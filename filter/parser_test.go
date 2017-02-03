package filter

import (
	"github.com/go-scim/scimify/adt"
	"github.com/stretchr/testify/assert"
	"testing"
)

type parserTest struct {
	name      string
	filter    string
	assertion func(*testing.T, parserTest, *adt.Node)
}

func TestParse(t *testing.T) {
	for _, test := range []parserTest{
		{
			"test parse simple filter",
			"foo eq \"bar\"",
			func(t *testing.T, test parserTest, root *adt.Node) {
				assert.Equal(t, Eq, root.Data.(Token).Value)
				assert.Equal(t, "foo", root.Left.Data.(Token).Value)
				assert.Equal(t, "\"bar\"", root.Right.Data.(Token).Value)
			},
		},
		{
			"test logical operator",
			"foo eq \"bar\" or age gt 18",
			func(t *testing.T, test parserTest, root *adt.Node) {
				assert.Equal(t, Or, root.Data.(Token).Value)
				assert.Equal(t, Eq, root.Left.Data.(Token).Value)
				assert.Equal(t, "foo", root.Left.Left.Data.(Token).Value)
				assert.Equal(t, "\"bar\"", root.Left.Right.Data.(Token).Value)
				assert.Equal(t, Gt, root.Right.Data.(Token).Value)
				assert.Equal(t, "age", root.Right.Left.Data.(Token).Value)
				assert.Equal(t, "18", root.Right.Right.Data.(Token).Value)
			},
		},
		{
			"test not operator with parenthesis",
			"not (foo eq \"bar\") and age gt 18",
			func(t *testing.T, test parserTest, root *adt.Node) {
				assert.Equal(t, And, root.Data.(Token).Value)
				assert.Equal(t, Not, root.Left.Data.(Token).Value)
				assert.Equal(t, Eq, root.Left.Left.Data.(Token).Value)
				assert.Equal(t, "foo", root.Left.Left.Left.Data.(Token).Value)
				assert.Equal(t, "\"bar\"", root.Left.Left.Right.Data.(Token).Value)
				assert.Equal(t, Gt, root.Right.Data.(Token).Value)
				assert.Equal(t, "age", root.Right.Left.Data.(Token).Value)
				assert.Equal(t, "18", root.Right.Right.Data.(Token).Value)
			},
		},
		{
			"test nested filter",
			"email[type eq \"work\"]",
			func(t *testing.T, test parserTest, root *adt.Node) {
				assert.Equal(t, Eq, root.Data.(Token).Value)
				assert.Equal(t, "email.type", root.Left.Data.(Token).Value)
				assert.Equal(t, "\"work\"", root.Right.Data.(Token).Value)
			},
		},
		{
			"test fully qualified name",
			"urn:ietf:params:scim:schemas:core:2.0:User:username eq \"Ark\"",
			func(t *testing.T, test parserTest, root *adt.Node) {
				assert.Equal(t, Eq, root.Data.(Token).Value)
				assert.Equal(t, "urn:ietf:params:scim:schemas:core:2.0:User:username", root.Left.Data.(Token).Value)
				assert.Equal(t, "\"Ark\"", root.Right.Data.(Token).Value)
			},
		},
	} {
		tokens, err := Tokenize(test.filter)
		assert.Nil(t, err)
		root, err := Parse(tokens)
		assert.Nil(t, err)
		assert.NotNil(t, root)
		test.assertion(t, test, root)
	}
}
