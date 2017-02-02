package filter

type Token struct {
	Value  string                 // Face value of the token
	Type   string                 // Type of the token
	Params map[string]interface{} // Extra parameters or metadata about this token
}

// Token types
const (
	Relational  = "rel"
	Logical     = "logic"
	Parenthesis = "paren"
	Constant    = "const"
	Path        = "path"
	NestedPath  = "nest_path"
)

// Parameter keys
const (
	Associativity = "assoc"
	Precedence    = "prec"
	ConstantType  = "const_type"
	NumberOfArgs  = "num_args"
	ParsedValue   = "parsed_val"
	NestedTokens  = "nest_tok"
)

// Parameter values
const (
	LeftAssociative  = "assoc_left"
	RightAssociative = "assoc_right"
	HighPrecedence   = 1000
	NormalPrecedence = 100
	LowPrecedence    = 10
	ConstString      = "str"
	ConstBool        = "bool"
	ConstInteger     = "int"
	ConstDecimal     = "dec"
)

// Stock token values
const (
	And              = "and"
	Or               = "or"
	Not              = "not"
	Eq               = "eq"
	Ne               = "ne"
	Sw               = "sw"
	Ew               = "ew"
	Co               = "co"
	Pr               = "pr"
	Gt               = "gt"
	Ge               = "ge"
	Lt               = "lt"
	Le               = "le"
	LeftParenthesis  = "("
	RightParenthesis = ")"
)
