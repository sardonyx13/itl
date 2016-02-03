package token

import "strconv"

// Token is the set of lexical tokens of the Go programming language.
type Token int

// The list of tokens.
const (
	// Special tokens
	ILLEGAL Token = iota
	EOF
	COMMENT

	literal_beg
	// Identifiers and basic type literals
	// (these tokens stand for classes of literals)
	IDENT  // main
	INT    // 12345
	FLOAT  // 123.45
	literal_end

	operator_beg
	// Operators and delimiters
	ADD // +
	SUB // -
	MUL // *
	QUO // /

	AND     // AND
	OR      // OR
	NOT     // NOT

	EQL    // =
	LSS    // <
	GTR    // >
	ASSIGN // :=

	NEQ      // <>
	LEQ      // <=
	GEQ      // >=
	
	XABOVE
	GOINGUP  
	TURNSUP
	WHEN
	XBELOW
	GOINGDOWN
	TURNSDOWN
	IF

	LPAREN // (
	LBRACK // [
	COMMA  // ,
	PERIOD // .

	RPAREN    // )
	RBRACK    // ]
	SEMICOLON // ;
	COLON     // :
	operator_end
)

var tokens = [...]string{
	ILLEGAL: "ILLEGAL",

	EOF:     "EOF",
	COMMENT: "COMMENT",

	IDENT:  "IDENT",
	INT:    "INT",
	FLOAT:  "FLOAT",

	ADD: "+",
	SUB: "-",
	MUL: "*",
	QUO: "/",

	AND: "AND",
	OR:  "OR",
	NOT: "NOT",

	EQL:  "=",
	LSS:  "<",
	GTR:  ">",
	ASSIGN: ":=",

	NEQ: "<>",
	LEQ: "<=",
	GEQ: ">=",
	
	XABOVE:    "XABOVE",
	GOINGUP:   "GOINGUP",  
	TURNSUP:   "TURNSUP",
	WHEN:      "WHEN",
	XBELOW:    "XBELOW",
	GOINGDOWN: "GOINGDOWN",
	TURNSDOWN: "TURNSDOWN",
	IF:        "IF",

	LPAREN: "(",
	LBRACK: "[",
	COMMA: ",",
	PERIOD: ".",

	RPAREN: ")",
	RBRACK: "]",
	SEMICOLON: ";",
	COLON: ":",
}

// String returns the string corresponding to the token tok.
// For operators, delimiters, and keywords the string is the actual
// token character sequence (e.g., for the token ADD, the string is
// "+"). For all other tokens the string corresponds to the token
// constant name (e.g. for the token IDENT, the string is "IDENT").
//
func (tok Token) String() string {
	s := ""
	if 0 <= tok && tok < Token(len(tokens)) {
		s = tokens[tok]
	}
	if s == "" {
		s = "token(" + strconv.Itoa(int(tok)) + ")"
	}
	return s
}

// A set of constants for precedence-based expression parsing.
// Non-operators have lowest precedence, followed by operators
// starting with precedence 1 up to unary operators. The highest
// precedence serves as "catch-all" precedence for selector,
// indexing, and other operator and delimiter tokens.
//
const (
	LowestPrec  = 0 // non-operators
	UnaryPrec   = 6
	HighestPrec = 7
)

// Precedence returns the operator precedence of the binary
// operator op. If op is not a binary operator, the result
// is LowestPrecedence.
//
func (op Token) Precedence() int {
	switch op {
	case LOR:
		return 1
	case LAND:
		return 2
	case EQL, NEQ, LSS, LEQ, GTR, GEQ:
		return 3
	case ADD, SUB, OR, XOR:
		return 4
	case MUL, QUO, REM, SHL, SHR, AND, AND_NOT:
		return 5
	}
	return LowestPrec
}

var keywords map[string]Token

func init() {
	keywords = make(map[string]Token)
	for i := keyword_beg + 1; i < keyword_end; i++ {
		keywords[tokens[i]] = i
	}
}

// Lookup maps an identifier to its keyword token or IDENT (if not a keyword).
//
func Lookup(ident string) Token {
	if tok, is_keyword := keywords[ident]; is_keyword {
		return tok
	}
	return IDENT
}

// Predicates

// IsLiteral returns true for tokens corresponding to identifiers
// and basic type literals; it returns false otherwise.
//
func (tok Token) IsLiteral() bool { return literal_beg < tok && tok < literal_end }

// IsOperator returns true for tokens corresponding to operators and
// delimiters; it returns false otherwise.
//
func (tok Token) IsOperator() bool { return operator_beg < tok && tok < operator_end }
