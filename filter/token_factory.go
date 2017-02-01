package filter

import (
	"strconv"
	"strings"
)

// Create a structured token based on a token string literal
func CreateToken(tok string) (Token, error) {
	switch strings.ToLower(tok) {
	case And:
		return Token{
			Value: And,
			Type:  Logical,
			Params: map[string]interface{}{
				Associativity: LeftAssociative,
				Precedence:    NormalPrecedence,
			},
		}, nil

	case Or:
		return Token{
			Value: Or,
			Type:  Logical,
			Params: map[string]interface{}{
				Associativity: LeftAssociative,
				Precedence:    NormalPrecedence - 1,
			},
		}, nil

	case Not:
		return Token{
			Value: Not,
			Type:  Logical,
			Params: map[string]interface{}{
				Associativity: RightAssociative,
				Precedence:    NormalPrecedence + 1,
			},
		}, nil

	case Eq:
		return Token{
			Value: Eq,
			Type:  Relational,
			Params: map[string]interface{}{
				Associativity: LeftAssociative,
				Precedence:    HighPrecedence,
				NumberOfArgs:  2,
			},
		}, nil

	case Ne:
		return Token{
			Value: Ne,
			Type:  Relational,
			Params: map[string]interface{}{
				Associativity: LeftAssociative,
				Precedence:    HighPrecedence,
				NumberOfArgs:  2,
			},
		}, nil

	case Sw:
		return Token{
			Value: Sw,
			Type:  Relational,
			Params: map[string]interface{}{
				Associativity: LeftAssociative,
				Precedence:    HighPrecedence,
				NumberOfArgs:  2,
			},
		}, nil

	case Ew:
		return Token{
			Value: Ew,
			Type:  Relational,
			Params: map[string]interface{}{
				Associativity: LeftAssociative,
				Precedence:    HighPrecedence,
				NumberOfArgs:  2,
			},
		}, nil

	case Co:
		return Token{
			Value: Co,
			Type:  Relational,
			Params: map[string]interface{}{
				Associativity: LeftAssociative,
				Precedence:    HighPrecedence,
				NumberOfArgs:  2,
			},
		}, nil

	case Pr:
		return Token{
			Value: Pr,
			Type:  Relational,
			Params: map[string]interface{}{
				Associativity: LeftAssociative,
				Precedence:    HighPrecedence,
				NumberOfArgs:  1,
			},
		}, nil

	case Gt:
		return Token{
			Value: Gt,
			Type:  Relational,
			Params: map[string]interface{}{
				Associativity: LeftAssociative,
				Precedence:    HighPrecedence,
				NumberOfArgs:  2,
			},
		}, nil

	case Ge:
		return Token{
			Value: Ge,
			Type:  Relational,
			Params: map[string]interface{}{
				Associativity: LeftAssociative,
				Precedence:    HighPrecedence,
				NumberOfArgs:  2,
			},
		}, nil

	case Lt:
		return Token{
			Value: Lt,
			Type:  Relational,
			Params: map[string]interface{}{
				Associativity: LeftAssociative,
				Precedence:    HighPrecedence,
				NumberOfArgs:  2,
			},
		}, nil

	case Le:
		return Token{
			Value: Le,
			Type:  Relational,
			Params: map[string]interface{}{
				Associativity: LeftAssociative,
				Precedence:    HighPrecedence,
				NumberOfArgs:  2,
			},
		}, nil

	case LeftParenthesis:
		return Token{
			Value:  LeftParenthesis,
			Type:   Parenthesis,
			Params: map[string]interface{}{},
		}, nil

	case RightParenthesis:
		return Token{
			Value:  RightParenthesis,
			Type:   Parenthesis,
			Params: map[string]interface{}{},
		}, nil

	default:
		if strings.HasPrefix(tok, "\"") && strings.HasSuffix(tok, "\"") {
			return Token{
				Value: tok,
				Type:  Constant,
				Params: map[string]interface{}{
					ConstantType: ConstString,
					ParsedValue:  tok[1 : len(tok)-1],
				},
			}, nil
		} else if b, err := strconv.ParseBool(tok); err == nil {
			return Token{
				Value: tok,
				Type:  Constant,
				Params: map[string]interface{}{
					ConstantType: ConstBool,
					ParsedValue:  b,
				},
			}, nil
		} else if i, err := strconv.ParseInt(tok, 10, 64); err == nil {
			return Token{
				Value: tok,
				Type:  Constant,
				Params: map[string]interface{}{
					ConstantType: ConstInteger,
					ParsedValue:  i,
				},
			}, nil
		} else if f, err := strconv.ParseFloat(tok, 64); err == nil {
			return Token{
				Value: tok,
				Type:  Constant,
				Params: map[string]interface{}{
					ConstantType: ConstDecimal,
					ParsedValue:  f,
				},
			}, nil
		} else {
			// TODO replace with scimError, need a standard error design
			return nil, nil
		}
	}
}
