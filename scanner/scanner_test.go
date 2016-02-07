package scanner

import (
	"testing"
)

const whitespace = "  \t  \n\n\n" // to separate tokens

type testCase struct {
	tok Token
	lit string
}

var testTokens = [...]testCase{

	{IDENT, "variable"},
	{INT, "123"},
	{FLOAT, "123.456"},
	{FLOAT, ".123"},
	{FLOAT, "123."},

	{ADD, "+"},
	{SUB, "-"},
	{MUL, "*"},
	{QUO, "/"},

	{AND, "AND"},
	{OR, "OR"},
	{NOT, "NOT"},

	{EQL, "="},
	{LSS, "<"},
	{GTR, ">"},
	{ASSIGN, ":="},

	{LPAREN, "("},

	{LBRACK, "["},
	{COMMA, ","},

	{RPAREN, ")"},
	{RBRACK, "]"},

	{SEMICOLON, ";"},

	{NEQ, "<>"},
	{LEQ, "<="},
	{GEQ, ">="},

	{IF, "IF"},
}

var source = func() []byte {
	var src []byte
	for _, i := range testTokens {
		src = append(src, i.lit...)
		src = append(src, whitespace...)
	}
	return src
}()

func TestInit(t *testing.T) {
	var s Scanner

	src1 := "test"
	s.Init([]byte(src1), nil)

	if s.ErrorCount != 0 {
		t.Errorf("found %d errors", s.ErrorCount)
	}
}

func TestScan(t *testing.T) {

	var s Scanner

	s.Init(source, nil)

	for _, i := range testTokens {
		_, tok, _ := s.Scan()
		if tok != i.tok {
			t.Errorf("return %s , expected %s", tok, i.tok)
		}
	}
}
