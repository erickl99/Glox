package main

import (
	"errors"
	"fmt"
	"strings"
)

type Parser struct {
    tokens []Token
    current int
}

func (ps *Parser) parse() (Expr, error) {
    return ps.expression()
}

func (ps *Parser) expression() (Expr, error) {
    return ps.equality()
}

func (ps *Parser) equality() (Expr, error) {
    expr, _ := ps.comparison()
    for ps.match(BANG_EQUAL, EQUAL_EQUAL) {
        op := ps.previous()
        right, _ := ps.comparison()
        //fmt.Println("Matched binary in equality")
        expr = Binary{expr, op, right}
    }
    return expr, nil
}

func (ps *Parser) comparison() (Expr, error) {
    expr, _ := ps.term()
    for ps.match(GREATER, GREATER_EQUAL, LESS, LESS_EQUAL) {
        op := ps.previous()
        right, _ := ps.term()
        //fmt.Println("Matched binary in comp")
        expr = Binary{expr, op, right}
    }
    return expr, nil
}

func (ps *Parser) term() (Expr, error) {
    expr, _ := ps.factor()

    for ps.match(MINUS, PLUS) {
        op := ps.previous()
        right, _ := ps.factor()
        //fmt.Println("Matched binary in term")
        expr = Binary{expr, op, right}
    }

    return expr, nil
}

func (ps *Parser) factor() (Expr, error) {
    expr, _ := ps.unary()

    for ps.match(SLASH, STAR) {
        op := ps.previous()
        right, _ := ps.unary()
        //fmt.Println("Matched binary in factor")
        expr = Binary{expr, op, right}
    }

    return expr, nil
}

func (ps *Parser) unary() (Expr, error) {
    if ps.match(BANG, MINUS) {
        op := ps.previous()
        right, _ := ps.unary()
        //fmt.Println("Matched unary")
        return Unary{op, right}, nil
    }
    return ps.primary()
}

func (ps *Parser) primary() (Expr, error) {
    if ps.match(FALSE) {
        //fmt.Println("Matched false")
        return Literal{false}, nil
    }
    if ps.match(TRUE) {
        //fmt.Println("Matched true")
        return Literal{true}, nil
    }
    if ps.match(NIL) {
        //fmt.Println("Matched a nil")
        return Literal{nil}, nil
    }

    if ps.match(NUMBER, STRING) {
        //fmt.Println("Matched a number or string")
        return Literal{ps.previous().literal}, nil
    }

    if ps.match(LEFT_PAREN) {
        //fmt.Println("Matched a left paren")
        expr, _ := ps.expression()
        ps.consume(RIGHT_PAREN, "Expect ')' after expression")
        return Grouping{expr}, nil
    }
    //fmt.Println("Matched a left paren")
    return Binary{}, errors.New("Expect expression")
}

func (ps *Parser) match(t_types ...TokenType) bool {
    for _, tt := range t_types {
        if ps.check(tt) {
            ps.advance()
            return true
        }
    }
    return false
}

func (ps Parser) check(t_type TokenType) bool {
    if ps.finished() {
        return false
    }
    return ps.peek().t_type == t_type
}

func (ps *Parser) advance() Token {
    if !ps.finished() {
        ps.current++
    }
    return ps.previous()
}

func (ps *Parser) consume(t_type TokenType, message string) (Token, error) {
    if ps.check(t_type) {
        return ps.advance(), nil
    }

    return Token{}, ps.error(ps.peek(), message)
}

func (ps Parser) error(token Token, message string) error {
    token_error(token, message)
    return errors.New("")
}

func (ps *Parser) synchronize() {
    ps.advance()

    for !ps.finished() {
        if ps.previous().t_type == SEMICOLON {
            return
        }
        switch ps.peek().t_type {
            case CLASS:
                fallthrough
            case FUN:
                fallthrough
            case VAR:
                fallthrough
            case FOR:
                fallthrough
            case IF:
                fallthrough
            case WHILE:
                fallthrough
            case PRINT:
                fallthrough
            case RETURN:
                return
        }
        ps.advance()
    }
}

func (ps Parser) finished() bool {
    return ps.peek().t_type == EOF
}

func (ps Parser) peek() Token {
    return ps.tokens[ps.current]
}

func (ps Parser) previous() Token {
    return ps.tokens[ps.current - 1]
}

func print(expr Expr) string {
    var ast string
    switch t := expr.(type) {
        case Binary:
            ast = parenthesize(t.operator.lexeme, t.left, t.right)
        case Grouping:
            ast = parenthesize("group", t.expression)
        case Literal:
            if t.value == nil {
                ast = "nil"
            } else {
                ast = fmt.Sprintf("%v", t.value)
            }
        case Unary:
            ast = parenthesize(t.operator.lexeme, t.right)
    }
    return ast
}

func parenthesize(name string, exprs ...Expr) string {
    var builder strings.Builder
    builder.WriteString("(")
    builder.WriteString(name)
    for _, expr := range exprs {
        builder.WriteString(" ")
        builder.WriteString(print(expr))
    }
    builder.WriteString(")")
    return builder.String()
}