package path

import (
	"testing"
	"github.com/go-scim/scimify/adt"
	"github.com/stretchr/testify/assert"
	"github.com/go-scim/scimify/filter"
)

func TestTokenize(t *testing.T) {
	for _, test := range []struct{
		name 		string
		path 		string
		assertion 	func(adt.Queue, error)
	}{
		{
			"simple path",
			"emails",
			func(q adt.Queue, err error) {
				assert.Nil(t, err)
				assert.Equal(t, 1, q.Size())
				assert.Equal(t, "emails", q.Peek().(*adt.Node).Data.(filter.Token).Value)
			},
		},
		{
			"simple path x2",
			"name.familyName",
			func(q adt.Queue, err error) {
				assert.Nil(t, err)
				assert.Equal(t, 2, q.Size())
				assert.Equal(t, "name", q.Poll().(*adt.Node).Data.(filter.Token).Value)
				assert.Equal(t, "familyName", q.Poll().(*adt.Node).Data.(filter.Token).Value)
			},
		},
		{
			"simple with filter",
			"emails[type eq \"work\"]",
			func(q adt.Queue, err error) {
				assert.Nil(t, err)
				assert.Equal(t, 2, q.Size())
				assert.Equal(t, "emails", q.Poll().(*adt.Node).Data.(filter.Token).Value)
				tok := q.Poll().(*adt.Node).Data.(filter.Token)
				assert.Equal(t, "eq", tok.Value)
				assert.Equal(t, filter.Relational, tok.Type)
			},
		},
		{
			"simple with filter and another simple",
			"groups[type eq \"internal\"].displayName",
			func(q adt.Queue, err error) {
				assert.Nil(t, err)
				assert.Equal(t, 3, q.Size())
				assert.Equal(t, "groups", q.Poll().(*adt.Node).Data.(filter.Token).Value)
				tok := q.Poll().(*adt.Node).Data.(filter.Token)
				assert.Equal(t, "eq", tok.Value)
				assert.Equal(t, filter.Relational, tok.Type)
				assert.Equal(t, "displayName", q.Poll().(*adt.Node).Data.(filter.Token).Value)
			},
		},
	}{
		q, err := Tokenize(test.path)
		test.assertion(q, err)
	}
}
