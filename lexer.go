package main

import (
	"strconv"
)

var keywords = map[string]TokenType{
	"and":    AND,
	"class":  CLASS,
	"else":   ELSE,
	"false":  FALSE,
	"for":    FOR,
	"fun":    FUN,
	"if":     IF,
	"nil":    NIL,
	"or":     OR,
	"print":  PRINT,
	"super":  SUPER,
	"this":   THIS,
	"true":   TRUE,
	"var":    VAR,
	"return": RETURN,
	"while":  WHILE,
}

type Lexer struct {
	source  string
	tokens  []Token
	start   int
	current int
	line    int
}

func NewLexer(source string) *Lexer {
	lexer := new(Lexer)
	lexer.source = source
	lexer.line = 1
	return lexer
}

func (lx *Lexer) scan_tokens() []Token {
	for !lx.finished() {
		lx.start = lx.current
		lx.scan_token()
	}
	lx.tokens = append(lx.tokens, Token{EOF, "", nil, lx.line})
	return lx.tokens
}

func (lx *Lexer) scan_token() {
	c := lx.advance()
	switch c {
	case '(':
		lx.add_token(LEFT_PAREN)
	case ')':
		lx.add_token(RIGHT_PAREN)
	case '{':
		lx.add_token(LEFT_BRACE)
	case '}':
		lx.add_token(RIGHT_BRACE)
	case ',':
		lx.add_token(COMMA)
	case '.':
		lx.add_token(DOT)
	case '-':
		lx.add_token(MINUS)
	case '+':
		lx.add_token(PLUS)
	case ';':
		lx.add_token(SEMICOLON)
	case '*':
		lx.add_token(STAR)
	case '!':
		var t_type TokenType
		if lx.matched('=') {
			t_type = BANG_EQUAL
		} else {
			t_type = BANG
		}
		lx.add_token(t_type)
	case '=':
		var t_type TokenType
		if lx.matched('=') {
			t_type = EQUAL_EQUAL
		} else {
			t_type = EQUAL
		}
		lx.add_token(t_type)
	case '<':
		var t_type TokenType
		if lx.matched('=') {
			t_type = LESS_EQUAL
		} else {
			t_type = LESS
		}
		lx.add_token(t_type)
	case '>':
		var t_type TokenType
		if lx.matched('=') {
			t_type = GREATER_EQUAL
		} else {
			t_type = GREATER
		}
		lx.add_token(t_type)
	case '/':
		if lx.matched('/') {
			for lx.peek() != '\n' && !lx.finished() {
				lx.advance()
			}
		} else {
			lx.add_token(SLASH)
		}
	case '"':
		lx.string()
	case ' ', '\r', '\t':
	case '\n':
		lx.line++
	default:
		if lx.is_digit(c) {
			lx.number()
		} else if lx.is_alpha(c) {
			lx.identifier()
		} else {
			line_error(lx.line, "Unexpected character.")
		}
	}
}

func (lx *Lexer) add_token(t_type TokenType) {
	lx.add_token_value(t_type, nil)
}

func (lx *Lexer) add_token_value(t_type TokenType, literal Value) {
	text := lx.source[lx.start:lx.current]
	lx.tokens = append(lx.tokens, Token{t_type, text, literal, lx.line})
}

func (lx *Lexer) advance() rune {
	rn := rune(lx.source[lx.current])
	lx.current++
	return rn
}

func (lx *Lexer) matched(expected rune) bool {
	if lx.finished() {
		return false
	}
	if rune(lx.source[lx.current]) != expected {
		return false
	}
	lx.current++
	return true
}

func (lx *Lexer) string() {
	for lx.peek() != '"' && !lx.finished() {
		if lx.peek() == '\n' {
			lx.line++
		}
		lx.advance()
	}

	if lx.finished() {
		line_error(lx.line, "Unterminated string")
		return
	}

	lx.advance()
	value := lx.source[lx.start+1 : lx.current-1]
	lx.add_token_value(STRING, value)
}

func (lx *Lexer) number() {
	for lx.is_digit(lx.peek()) {
		lx.advance()
	}

	if lx.peek() == '.' && lx.is_digit(lx.peek_next()) {
		lx.advance()
		for lx.is_digit(lx.peek()) {
			lx.advance()
		}
	}

	value, _ := strconv.ParseFloat(lx.source[lx.start:lx.current], 64)
	lx.add_token_value(NUMBER, value)
}

func (lx *Lexer) identifier() {
	for lx.is_alphanumeric(lx.peek()) {
		lx.advance()
	}
	text := lx.source[lx.start:lx.current]
	var t_type TokenType
	if kw, ok := keywords[text]; ok {
		t_type = kw
	} else {
		t_type = IDENTIFIER
	}
	lx.add_token(t_type)
}

func (lx Lexer) finished() bool {
	return lx.current >= len(lx.source)
}

func (lx Lexer) peek() rune {
	if lx.finished() {
		return '\000'
	}
	return rune(lx.source[lx.current])
}

func (lx Lexer) peek_next() rune {
	if lx.current+1 >= len(lx.source) {
		return '\000'
	}
	return rune(lx.source[lx.current+1])
}

func (lx Lexer) is_digit(rn rune) bool {
	return rn >= '0' && rn <= '9'
}

func (lx Lexer) is_alpha(rn rune) bool {
	return (rn >= 'A' && rn <= 'Z') || (rn >= 'a' && rn <= 'z') || (rn == '_')
}

func (lx Lexer) is_alphanumeric(rn rune) bool {
	return lx.is_alpha(rn) || lx.is_digit(rn)
}
