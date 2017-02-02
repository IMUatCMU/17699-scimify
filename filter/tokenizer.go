package filter

import (
	"fmt"
	"github.com/go-scim/scimify/resource"
)

func Tokenize(filter string) ([]Token, error) {
	t := &tokenizer{
		textMode:  false,
		remaining: make([]rune, 0, len(filter)),
		buffer:    make([]rune, 0),
		tokens:    make([]Token, 0),
	}
	for _, r := range filter {
		t.remaining = append(t.remaining, r)
	}

	if err := t.tokenize(); err != nil {
		return nil, err
	} else {
		return t.tokens, nil
	}
}

type tokenizer struct {
	textMode     bool    // treat everything as text
	remaining    []rune  // the remaining runes to be tokenized
	buffer       []rune  // buffer for the runes to be converted to the next token
	tokens       []Token // tokens
	bracketLevel int     // matching count for brackets
	parenLevel   int     // matching count for parenthesis
}

func (t *tokenizer) tokenize() error {
	for len(t.remaining) > 0 {
		r := t.getAndDropTopRune()
		switch r {
		case spaceRune:
			if t.textMode {
				t.addToBuffer(r)
			} else {
				t.addBufferToTokens()
			}

		case quoteRune:
			t.addToBuffer(r)
			t.textMode = !t.textMode

		case leftBracketRune:
			t.addToBuffer(r)
			t.textMode = true
			t.bracketLevel++
			if t.bracketLevel > 1 {
				return resource.CreateError(
					resource.InvalidFilter,
					"Only one level of nested filter is allowed.")
			}

		case rightBracketRune:
			t.addToBuffer(r)
			t.textMode = false
			t.bracketLevel--

		case leftParenRune:
			t.addToTokens(r)
			t.parenLevel++

		case rightParenRune:
			t.addToTokens(r)
			t.parenLevel--

		case commaRune:
			t.addToTokens(r)

		default:
			t.addToBuffer(r)
		}
	}
	t.addBufferToTokens()

	switch {
	case t.bracketLevel > 0:
		return resource.CreateError(
			resource.InvalidFilter,
			"Failed to parse filter: mismatched brackets")
	case t.parenLevel > 0:
		return resource.CreateError(
			resource.InvalidFilter,
			"Failed to parse filter: mismatched parenthesis")
	default:
		return nil
	}
}

func (t *tokenizer) getAndDropTopRune() rune {
	r := t.remaining[0]
	t.remaining = t.remaining[1:]
	return r
}

func (t *tokenizer) addToBuffer(r rune) {
	t.buffer = append(t.buffer, r)
}

func (t *tokenizer) addToTokens(r rune) error {
	if len(t.buffer) > 0 {
		t.addBufferToTokens()
	}
	if tok, err := CreateToken(fmt.Sprintf("%c", r)); err != nil {
		return err
	} else {
		t.tokens = append(t.tokens, tok)
		return nil
	}
}

func (t *tokenizer) addBufferToTokens() error {
	if len(t.buffer) > 0 {
		if tok, err := CreateToken(string(t.buffer)); err != nil {
			return err
		} else {
			t.tokens = append(t.tokens, tok)
			t.buffer = make([]rune, 0)
			return nil
		}
	} else {
		return resource.CreateError(
			resource.InvalidFilter,
			"Failed to parse filter: unexpected filter content.")
	}
}

// Internally used constant
const (
	spaceRune        = ' '
	quoteRune        = '"'
	commaRune        = ','
	leftBracketRune  = '['
	rightBracketRune = ']'
	leftParenRune    = '('
	rightParenRune   = ')'
)
