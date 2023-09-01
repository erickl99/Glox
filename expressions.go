package main

type Expr interface {
    accept() Value
}

type Assign struct {
    name Token
    value Expr
}

func (as Assign) accept() Value {
    return 0
}

type Binary struct {
    left Expr
    operator Token
    right Expr
}

func (bn Binary) accept() Value {
    return 0
}

type Call struct {
    calle Expr
    paren Token
    arguemnts []Expr
}

func (ca Call) accept() Value {
    return 0
}

type Get struct {
    object Expr
    name Token
}

func (gt Get) accept() Value {
    return 0
}

type Grouping struct {
    expression Expr
}

func (gp Grouping) accept() Value {
    return 0
}

type Literal struct {
    value Value
}

func (lt Literal) accept() Value {
    return 0
}

type Logical struct {
    left Expr
    operator Token
    right Expr
}

func (lg Logical) accept() Value {
    return 0
}

type Set struct {
    object Expr
    name Token
    value Expr
}

func (st Set) accept() Value {
    return 0
}

type Super struct {
    keyword Token
    method Token
}

func (sp Super) accept() Value {
    return 0
}

type This struct {
    keyword Token
}

func (th This) accept() Value {
    return 0
}

type Unary struct {
    operator Token
    right Expr
}

func (un Unary) accept() Value {
    return 0
}

type Variable struct {
    name Token
}

func (vr Variable) accept() Value {
    return 0
}
