package main

import "fmt"

type TokenType int

const (
    LEFT_PAREN TokenType = iota
    RIGHT_PAREN
    LEFT_BRACE
    RIGHT_BRACE
    COMMA
    DOT
    MINUS
    PLUS
    SEMICOLON
    SLASH
    STAR
    BANG
    BANG_EQUAL
    EQUAL
    EQUAL_EQUAL
    GREATER
    GREATER_EQUAL
    LESS
    LESS_EQUAL
    IDENTIFIER
    STRING
    NUMBER
    AND
    CLASS
    ELSE
    FALSE
    FUN
    FOR
    IF
    NIL
    OR
    PRINT
    RETURN
    SUPER
    THIS
    TRUE
    VAR
    WHILE
    EOF
)

type Value interface {}

type Token struct {
    t_type TokenType
    lexeme string
    literal Value
    line int
}

func (t Token) String() string {
    return fmt.Sprintf("%v %s %v", t.t_type, t.lexeme, t.literal)
}
