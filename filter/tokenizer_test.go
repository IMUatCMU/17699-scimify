package filter

import "testing"
import (
	"github.com/go-scim/scimify/resource"
	"github.com/stretchr/testify/assert"
)

type tokenizerTest struct {
	testName         string
	filter           string
	expected         []string
	additionalAssert func(*testing.T, tokenizerTest, []Token)
}

type tokenizerInvalidTest struct {
	testName        string
	filter          string
	expectedErrType string
}

func TestTokenize(t *testing.T) {
	for _, test := range []tokenizerTest{
		{
			"test simple filter",
			"age gt 10",
			[]string{"age", "gt", "10"},
			func(t *testing.T, test tokenizerTest, tokens []Token) {
				assert.Equal(t, Path, tokens[0].Type, "[%s] token #1 should be a path token", test.testName)
				assert.Equal(t, Relational, tokens[1].Type, "[%s] token #2 should be a relational token", test.testName)
				assert.Equal(t, Constant, tokens[2].Type, "[%s] token #3 should be a constant token", test.testName)
				assert.Equal(t, int64(10), tokens[2].Params[ParsedValue], "[%s] token #4 should be int64 value of 10", test.testName)
			},
		},
		{
			"test normal filter",
			"(age gt 10) and (name eq \"David\")",
			[]string{"(", "age", "gt", "10", ")", "and", "(", "name", "eq", "\"David\"", ")"},
			func(t *testing.T, test tokenizerTest, tokens []Token) {
				assert.Equal(t, Parenthesis, tokens[0].Type, "[%s] token #1 should be a parenthesis token", test.testName)
				assert.Equal(t, Parenthesis, tokens[4].Type, "[%s] token #5 should be a parenthesis token", test.testName)
				assert.Equal(t, Logical, tokens[5].Type, "[%s] token #6 should be a logical token", test.testName)
				assert.Equal(t, "David", tokens[9].Params[ParsedValue], "[%s] token #10 should be string value of 'David'", test.testName)
			},
		},
		{
			"test nested filter",
			"email[type eq \"work\"] and name sw \"D\"",
			[]string{"email", "and", "name", "sw", "\"D\""},
			func(t *testing.T, test tokenizerTest, tokens []Token) {
				assert.Equal(t, NestedPath, tokens[0].Type, "[%s] token #1 should be a nested path token", test.testName)
				nestedTokens := tokens[0].Params[NestedTokens].([]Token)
				assert.NotNil(t, nestedTokens, "[%s] token #1 should have non-nil nested tokens parsed", test.testName)
				assert.Equal(t, Path, nestedTokens[0].Type, "[%s] nested token #1 should be a path token", test.testName)
				assert.Equal(t, Relational, nestedTokens[1].Type, "[%s] nested token #2 should be a relational token", test.testName)
				assert.Equal(t, Constant, nestedTokens[2].Type, "[%s] nested token #3 should be a constant token", test.testName)
				assert.Equal(t, "type", nestedTokens[0].Value)
				assert.Equal(t, "eq", nestedTokens[1].Value)
				assert.Equal(t, "\"work\"", nestedTokens[2].Value)
			},
		},
	} {
		tokens, err := Tokenize(test.filter)
		assert.Nil(t, err)
		assert.Equal(t, len(test.expected), len(tokens), "there should be %d tokens", len(test.expected))
		for i, tok := range tokens {
			assert.Equal(t, test.expected[i], tok.Value)
		}
		if test.additionalAssert != nil {
			test.additionalAssert(t, test, tokens)
		}
	}
}

func TestTokenizeInvalid(t *testing.T) {
	for _, test := range []tokenizerInvalidTest{
		{
			"test more than one level of nested path",
			"foo[bar[type eq \"work\"]] and name sw \"D\"",
			resource.InvalidFilter,
		},
		{
			"test unbalanced brackets",
			"foo[bar[type eq \"work\"] and name sw \"D\"",
			resource.InvalidFilter,
		},
		{
			"test unbalanced parenthesis",
			"( foo eq bar",
			resource.InvalidFilter,
		},
	} {
		_, err := Tokenize(test.filter)
		assert.NotNil(t, err, "[%s] expects an error", test.testName)
		assert.Equal(t, test.expectedErrType, err.(resource.Error).ScimType, "[%s] expected scimType to be %s", test.testName, test.expectedErrType)
	}
}
