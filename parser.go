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

type ParseError struct {
    message string
}

func (pe ParseError) Error() string {
    return pe.message
}

func (ps *Parser) parse() ([]Stmt, error) {
    var statements []Stmt
    for !ps.finished() {
        stmt := ps.declaration()
        statements = append(statements, stmt)
    }
    return statements, nil
}

func (ps *Parser) declaration() Stmt {
    if ps.match(VAR) {
        stmt, err := ps.var_declaration()
        if err != nil {
            ps.synchronize()
            return nil
        }
        return stmt
    }
    if ps.match(FUN) {
        stmt, err := ps.function("function")
        if err != nil {
            ps.synchronize()
            return nil
        }
        return stmt
    }
    stmt, err := ps.statement()
    if err != nil {
        ps.synchronize()
        return nil
    }
    return stmt
}

func (ps *Parser) var_declaration() (Stmt, error) {
    name, err := ps.consume(IDENTIFIER, "Expect variable name")
    if err != nil {
        return nil, err
    }
    var initializer Expr
    if ps.match(EQUAL) {
        val, err := ps.expression()
        if err != nil {
            return nil, err
        }
        initializer = val
    }
    ps.consume(SEMICOLON, "Expect ';' after variable declaration")
    return Var{name, initializer}, nil
}

func (ps *Parser) function(kind string) (Func, error) {
    var new_function Func
    name, err := ps.consume(IDENTIFIER, "Expect " + kind + " name")
    if err != nil {
        return new_function, err
    }
    _, err = ps.consume(LEFT_PAREN, "Expect '(' after function " + kind + " name")
    if err != nil {
        return new_function, err
    }
    var parameters []Token
    if !ps.check(RIGHT_PAREN) {
        for {
            if len(parameters) >= 255 {
                ps.error(ps.peek(), "Cannot have more than 255 parameters")
            }
            param, err := ps.consume(IDENTIFIER, "Expect parameter name")
            if err != nil {
                return new_function, err
            }
            parameters = append(parameters, param)
            if !ps.match(COMMA) {
                break
            }
        }
    }
    _, err = ps.consume(RIGHT_PAREN, "Expect ')' after parameters")
    if err != nil {
        return new_function, err
    }
    _, err = ps.consume(LEFT_BRACE, "Expect '{' before " + kind + " body")
    if err != nil {
        return new_function, err
    }
    body, err := ps.block()
    if err != nil {
        return new_function, err
    }
    new_function.name = name
    new_function.params = parameters
    new_function.body = body
    return new_function, nil
}

func (ps *Parser) statement() (Stmt, error){
    if ps.match(FOR) {
        return ps.for_statement()
    }
    if ps.match(IF) {
        return ps.if_statement()
    }
    if ps.match(PRINT) {
        return ps.print_statement()
    }
    if ps.match(RETURN) {
        return ps.return_statement()
    }
    if ps.match(WHILE) {
        return ps.while_statement()
    }
    if ps.match(LEFT_BRACE) {
        stmts, err := ps.block()
        if err != nil {
            return nil, err
        }
        return Block{stmts}, nil
    }
    return ps.expression_statement()
}

func (ps *Parser) block() ([]Stmt, error) {
    var statements []Stmt
    for !ps.check(RIGHT_BRACE) && !ps.finished() {
        stmt := ps.declaration()
        statements = append(statements, stmt)
    }
    _, err := ps.consume(RIGHT_BRACE, "Expected '}' after block")
    if err != nil {
        return nil, err
    }
    return statements, nil
}

func (ps *Parser) print_statement() (Stmt, error) {
    value, err := ps.expression()
    if err != nil {
        return nil, err
    }
    _, err = ps.consume(SEMICOLON, "Expect ';' after value")
    if err != nil {
        return nil, err
    }
    return Print{value}, nil
}

func (ps *Parser) expression_statement() (Stmt, error) {
    expr, err := ps.expression()
    if err != nil {
        return nil, err
    }
    _, err = ps.consume(SEMICOLON, "Expect ';' after expression")
    if err != nil {
        return nil, err
    }
    return Expression{expr}, nil
}

func (ps *Parser) if_statement() (Stmt, error) {
    _, err := ps.consume(LEFT_PAREN, "Expect  '(' after if")
    if err != nil {
        return nil, err
    }
    cond, err := ps.expression()
    if err != nil {
        return nil, err
    }
    _, err = ps.consume(RIGHT_PAREN, "Expect  ')' after condition")
    if err != nil {
        return nil, err
    }
    then_branch, err := ps.statement()
    if err != nil {
        return nil, err
    }
    var else_branch Stmt = nil
    if ps.match(ELSE) {
        else_branch, err = ps.statement()
        if err != nil {
            return nil, err
        }
    }
    return If{cond,then_branch, else_branch}, nil
}

func (ps *Parser) return_statement() (Stmt, error) {
    keyword := ps.previous()
    var value Expr
    var err error
    if !ps.check(SEMICOLON) {
        value, err = ps.expression()
        if err != nil {
            return nil, err
        }
    }
    _, err = ps.consume(SEMICOLON, "Expect semicolon after return value")
    return Return{keyword, value}, nil
}

func (ps *Parser) while_statement() (Stmt, error) {
    _, err := ps.consume(LEFT_PAREN, "Expect '(' after 'while'")
    if err != nil {
        return nil, err
    }
    expr, err := ps.expression()
    if err != nil {
        return nil, err
    }
    _, err = ps.consume(RIGHT_PAREN, "Expect ')' after condition")
    if err != nil {
        return nil, err
    }
    body, err := ps.statement()
    if err != nil {
        return nil, err
    }
    return While{expr, body}, nil
}

func (ps *Parser) for_statement() (Stmt, error) {
    _, err := ps.consume(LEFT_PAREN, "Expect '(' after 'for'")
    if err != nil {
        return nil, err
    }
    var initializer Stmt
    if ps.match(SEMICOLON) {
        initializer = nil
    } else if ps.match(VAR) {
        initializer, err = ps.var_declaration()
    } else {
        initializer, err = ps.expression_statement()
    }
    if err != nil {
        return nil, err
    }
    var condition Expr = nil
    if !ps.match(SEMICOLON) {
        condition, err = ps.expression()
        if err != nil {
            return nil, err
        }
    }
    _, err = ps.consume(SEMICOLON, "Expect ';' after loop condition")
    if err != nil {
        return nil, err
    }
    var increment Expr = nil
    if !ps.check(RIGHT_PAREN) {
        increment, err = ps.expression()
        if err != nil {
            return nil, err
        }
    }
    _, err = ps.consume(RIGHT_PAREN, "Expect ')' after for clauses")
    if err != nil {
        return nil, err
    }
    body, err := ps.statement()
    if err != nil {
        return nil, err
    }
    if increment != nil  {
        stmts := []Stmt{body, Expression{increment}}
        body = Block{stmts}
    }
    if condition == nil {
        condition = Literal{true}
    }
    body = While{condition, body}
    if initializer != nil {
        stmts := []Stmt{initializer, body}
        body = Block{stmts}
    }
    return body, nil
}

func (ps *Parser) expression() (Expr, error) {
    return ps.assignment()
}

func (ps *Parser) assignment() (Expr, error) {
    expr, err := ps.or()
    if err != nil {
        return nil, err
    }
    if ps.match(EQUAL) { 
        equals := ps.previous()
        value, err := ps.assignment()
        if err != nil {
            return nil, err
        }
        if assignee, ok := expr.(Variable); ok {
            name := assignee.name
            return Assign{name, value}, nil
        }
        ps.error(equals, "Invalid assignment target")
    }
    return expr, nil
}

func (ps *Parser) or() (Expr, error) {
    expr, err := ps.and()
    if err != nil {
        return nil, err
    }
    for ps.match(OR) {
        op := ps.previous()
        right, err := ps.and()
        if err != nil {
            return nil, err
        }
        expr = Logical{expr, op, right}
    }
    return expr, nil
}

func (ps *Parser) and() (Expr, error) {
    expr, err := ps.equality()
    if err != nil {
        return nil, err
    }
    for ps.match(AND) {
        op := ps.previous()
        right, err := ps.equality()
        if err != nil {
            return nil, err
        }
        expr = Logical{expr, op, right}
    }
    return expr, nil
}

func (ps *Parser) equality() (Expr, error) {
    expr, err := ps.comparison()
    if err != nil {
        return nil, err
    }
    for ps.match(BANG_EQUAL, EQUAL_EQUAL) {
        op := ps.previous()
        right, err := ps.comparison()
        if err != nil {
            return nil, err
        }
        //fmt.Println("Matched binary in equality")
        expr = Binary{expr, op, right}
    }
    return expr, nil
}

func (ps *Parser) comparison() (Expr, error) {
    expr, err := ps.term()
    if err != nil {
        return nil, err
    }
    for ps.match(GREATER, GREATER_EQUAL, LESS, LESS_EQUAL) {
        op := ps.previous()
        right, err := ps.term()
        if err != nil {
            return nil, err
        }
        expr = Binary{expr, op, right}
    }
    return expr, nil
}

func (ps *Parser) term() (Expr, error) {
    expr, err := ps.factor()
    if err != nil {
        return nil, err
    }
    for ps.match(MINUS, PLUS) {
        op := ps.previous()
        right, err := ps.factor()
        if err != nil {
            return nil, err
        }
        //fmt.Println("Matched binary in term")
        expr = Binary{expr, op, right}
    }

    return expr, nil
}

func (ps *Parser) factor() (Expr, error) {
    expr, err := ps.unary()
    if err != nil {
        return nil, err
    }
    for ps.match(SLASH, STAR) {
        op := ps.previous()
        right, err := ps.unary()
        if err != nil {
            return nil, err
        }
        //fmt.Println("Matched binary in factor")
        expr = Binary{expr, op, right}
    }

    return expr, nil
}

func (ps *Parser) unary() (Expr, error) {
    if ps.match(BANG, MINUS) {
        op := ps.previous()
        right, err := ps.unary()
        if err != nil {
            return nil, err
        }
        //fmt.Println("Matched unary")
        return Unary{op, right}, nil
    }
    return ps.call()
}

func (ps *Parser) call() (Expr, error) {
    expr, err := ps.primary()
    if err != nil {
        return nil, err
    }
    for {
        if ps.match(LEFT_PAREN) {
            expr, err = ps.finish_call(expr)
            if err != nil {
                return nil, err
            }
        } else {
            break
        }
    }
    return expr, nil
}

func (ps *Parser) finish_call(callee Expr) (Expr, error) {
    var args []Expr
    if !ps.check(RIGHT_PAREN) {
        for {
            expr, err := ps.expression()
            if err != nil {
                return nil, err
            }
            if len(args) >= 255 {
                ps.error(ps.peek(), "Cannot have more than 255 arguments to a function")
            }
            args = append(args, expr)
            if !ps.match(COMMA) {
                break
            }
        }
    }
    paren, err := ps.consume(RIGHT_PAREN, "Expect ')' after arguments")
    if err != nil {
        return nil, err
    }
    return Call{callee, paren, args}, nil
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
        return Literal{value: nil}, nil
    }
    if ps.match(IDENTIFIER) {
        return Variable{ps.previous()}, nil
    }

    if ps.match(NUMBER, STRING) {
        return Literal{ps.previous().literal}, nil
    }

    if ps.match(LEFT_PAREN) {
        //fmt.Println("Matched a left paren")
        expr, err := ps.expression()
        if err != nil {
            return nil, err
        }
        _, err = ps.consume(RIGHT_PAREN, "Expect ')' after expression")
        if err != nil {
            return nil, err
        }
        return Grouping{expr}, nil
    }
    //fmt.Println("Matched a left paren")
    return nil, errors.New("Expect expression")
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
