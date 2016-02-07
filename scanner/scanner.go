package scanner

import (
	"fmt"
	"unicode"
	"unicode/utf8"
)

type ErrorHandler func(pos int, msg string)

type Scanner struct {
	src []byte       // source
	err ErrorHandler // error reporting; or nil

	// scanning state
	ch       rune // current character
	offset   int  // character offset
	rdOffset int  // reading offset (position after current character)

	ErrorCount int // number of errors encountered
}

const bom = 0xFEFF // byte order mark, only permitted as very first character

func (s *Scanner) Init(src []byte, err ErrorHandler) {
	s.src = src
	s.err = err
	s.ch = ' '
	s.offset = 0
	s.rdOffset = 0
	s.ErrorCount = 0

	s.next()
	if s.ch == bom {
		s.next()
	}
}

func (s *Scanner) Scan() (pos int, tok Token, lit string) {
	s.skipWhitespace()

	lit = ""
	pos = s.offset
	tok = ILLEGAL

	switch ch := s.ch; {
	case isLetter(ch):
		lit = s.scanIdentifier()
		tok = Lookup(lit)
	case isDigit(ch):
		tok, lit = s.scanNumber(false)
	default:
		s.next() // always make progress
		switch ch {
		case -1:
			tok = EOF
		case ',':
			tok = COMMA
		case '(':
			tok = LPAREN
		case ')':
			tok = RPAREN
		case '[':
			tok = LBRACK
		case ']':
			tok = RBRACK
		case '+':
			tok = ADD
		case '-':
			tok = SUB
		case '*':
			tok = MUL
		case '/':
			tok = QUO
		case '=':
			tok = EQL
		case ':':
			if '=' == s.ch {
				tok = ASSIGN
				s.next()
			}
		case '.':
			if '0' <= s.ch && s.ch <= '9' {
				tok, lit = s.scanNumber(true)
			}
		case '<':
			tok = s.switch3(LSS, LEQ, NEQ)
		case '>':
			tok = s.switch2(GTR, GEQ)
		case ';':
			tok = SEMICOLON
		default:
			// next reports unexpected BOMs - don't repeat
			if ch != bom {
				s.error(pos, fmt.Sprintf("illegal character %#U", ch))
			}
		}
	}
	return
}

func (s *Scanner) scanIdentifier() string {
	offs := s.offset
	for isLetter(s.ch) || isDigit(s.ch) {
		s.next()
	}
	return string(s.src[offs:s.offset])
}

func (s *Scanner) scanNumber(seenDecimalPoint bool) (tok Token, lit string) {
	offs := s.offset
	tok = INT

	if seenDecimalPoint {
		offs--
	}

	for isDigit(s.ch) {
		s.next()
	}

	if s.ch == '.' && seenDecimalPoint == false {
		s.next()
		seenDecimalPoint = true
		for isDigit(s.ch) {
			s.next()
		}
	}

	if seenDecimalPoint {
		tok = FLOAT
	}

	return tok, string(s.src[offs:s.offset])
}

func (s *Scanner) switch2(tok0, tok1 Token) Token {
	if s.ch == '=' {
		s.next()
		return tok1
	}
	return tok0
}

func (s *Scanner) switch3(tok0, tok1, tok3 Token) Token {
	switch s.ch {
	case '=':
		s.next()
		return tok1
	case '>':
		s.next()
		return tok3
	}
	return tok0
}

func (s *Scanner) next() {
	if s.rdOffset < len(s.src) {
		s.offset = s.rdOffset
		r, w := rune(s.src[s.rdOffset]), 1
		switch {
		case r == 0:
			s.error(s.offset, "illegal character NULL")
		case r >= 0x80:
			// not ASCII
			r, w = utf8.DecodeRune(s.src[s.rdOffset:])
			if r == utf8.RuneError && w == 1 {
				s.error(s.offset, "illegal UTF-8 encoding")
			} else if r == bom && s.offset > 0 {
				s.error(s.offset, "illegal byte order mark")
			}
		}
		s.rdOffset += w
		s.ch = r
	} else {
		s.offset = len(s.src)
		s.ch = -1 // eof
	}
}

func (s *Scanner) skipWhitespace() {
	for s.ch == ' ' || s.ch == '\t' || s.ch == '\n' || s.ch == '\r' {
		s.next()
	}
}

func (s *Scanner) error(offs int, msg string) {
	if s.err != nil {
		s.err(offs, msg)
	}
	s.ErrorCount++
}

func isLetter(ch rune) bool {
	return 'a' <= ch && ch <= 'z' || 'A' <= ch && ch <= 'Z' || ch == '_' || ch >= 0x80 && unicode.IsLetter(ch)
}

func isDigit(ch rune) bool {
	return '0' <= ch && ch <= '9' || ch >= 0x80 && unicode.IsDigit(ch)
}
