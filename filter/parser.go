package filter

import (
	"fmt"
	"github.com/go-scim/scimify/adt"
	"github.com/go-scim/scimify/resource"
)

func Parse(tokens []Token) (*adt.Node, error) {
	nestedCount := 0
	for _, tok := range tokens {
		if tok.Type == NestedPath {
			nestedCount += len(tok.Params[NestedTokens].([]Token))
		}
	}

	sy := &shuntingYard{
		input:    adt.NewQueue(len(tokens) + nestedCount),
		operator: adt.NewStack(len(tokens) + nestedCount),
		output:   adt.NewStack(len(tokens) + nestedCount),
	}

	return sy.run(tokens)
}

type shuntingYard struct {
	input    adt.Queue
	operator adt.Stack
	output   adt.Stack
}

func (sy *shuntingYard) run(tokens []Token) (*adt.Node, error) {
	for _, tok := range tokens {
		if tok.Type != NestedPath {
			sy.input.Offer(tok)
		} else {
			nestedTokens := tok.Params[NestedTokens].([]Token)
			for _, nestedTok := range nestedTokens {
				if nestedTok.Type != Path {
					sy.input.Offer(nestedTok)
				} else {
					sy.input.Offer(Token{
						Value:  tok.Value + "." + nestedTok.Value,
						Type:   nestedTok.Type,
						Params: nestedTok.Params,
					})
				}
			}
		}
	}

	for sy.input.Size() > 0 {
		tok := sy.input.Poll().(Token)

		switch {
		case tok.Type == Path || tok.Type == Constant:
			if err := sy.pushToOutput(tok); err != nil {
				return nil, err
			}

		case tok.Type == Relational || tok.Type == Logical:
			//o1 := tok.(token.OperatorToken)
			for {
				if peek := sy.operator.Peek(); peek == nil {
					break
				} else if tok2 := peek.(Token); tok2.Type != Relational && tok2.Type != Logical {
					break
				} else {
					//o2 := peek.(token.OperatorToken)
					if tok.Params[Associativity] == LeftAssociative &&
						tok.Params[Precedence].(int) <= tok2.Params[Precedence].(int) {
						if err := sy.pushToOutput(sy.operator.Pop().(Token)); err != nil {
							return nil, err
						}
					} else if tok.Params[Associativity] == RightAssociative &&
						tok.Params[Precedence].(int) < tok2.Params[Precedence].(int) {
						if err := sy.pushToOutput(sy.operator.Pop().(Token)); err != nil {
							return nil, err
						}
					} else {
						break
					}
				}
			}
			sy.operator.Push(tok)

		case tok.Type == Parenthesis && tok.Value == LeftParenthesis:
			sy.operator.Push(tok)

		case tok.Type == Parenthesis && tok.Value == RightParenthesis:
			for {
				if peek := sy.operator.Peek(); peek == nil {
					return nil, errParenMismatch
				} else if peek.(Token).Type == Parenthesis && peek.(Token).Value == LeftParenthesis {
					sy.operator.Pop()
					break
				} else {
					if err := sy.pushToOutput(sy.operator.Pop().(Token)); err != nil {
						return nil, err
					}
				}
			}
		}
	}

	for sy.operator.Size() > 0 {
		if peek := sy.operator.Peek(); peek != nil && peek.(Token).Type == Parenthesis {
			return nil, errParenMismatch
		} else {
			if err := sy.pushToOutput(sy.operator.Pop().(Token)); err != nil {
				return nil, err
			}
		}
	}

	return sy.output.Pop().(*adt.Node), nil
}

func (sy *shuntingYard) pushToOutput(tok Token) error {
	node := adt.NewNode(tok)
	if tok.Type != Constant && tok.Type != Path {
		numOfArgs := tok.Params[NumberOfArgs].(int)
		if numOfArgs == 1 {
			if arg := sy.output.Pop(); arg == nil {
				return errInsufficientArg(tok)
			} else {
				node.Left = arg.(*adt.Node)
			}
		} else if numOfArgs == 2 {
			arg2 := sy.output.Pop()
			arg1 := sy.output.Pop()

			if arg1 == nil || arg2 == nil {
				return errInsufficientArg(tok)
			} else {
				node.Left = arg1.(*adt.Node)
				node.Right = arg2.(*adt.Node)
			}
		} else {
			return errTooManyArguments
		}
	}
	sy.output.Push(node)
	return nil
}

var (
	errParenMismatch   = resource.CreateError(resource.InvalidFilter, "Mismatched parenthesis.")
	errInsufficientArg = func(tok Token) error {
		return resource.CreateError(
			resource.InvalidFilter, fmt.Errorf("Operator [%s] has insufficient arguments", tok.Value))
	}
	errTooManyArguments = resource.CreateError(resource.InvalidFilter, "Arguments more than 2 is not supported.")
)
