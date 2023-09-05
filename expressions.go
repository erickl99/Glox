package main

type Expr interface {
    accept()
}

type Assign struct {
    name Token
    value Expr
}

func (as Assign) accept() {
}

type Binary struct {
    left Expr
    operator Token
    right Expr
}

func (bn Binary) accept() {
}

type Call struct {
    callee Expr
    paren Token
    arguments []Expr
}

func (ca Call) accept() {
}

type Get struct {
    object Expr
    name Token
}

func (gt Get) accept() {
}

type Grouping struct {
    expression Expr
}

func (gp Grouping) accept() {
}

type Literal struct {
    value Value
}

func (lt Literal) accept() {
}

type Logical struct {
    left Expr
    operator Token
    right Expr
}

func (lg Logical) accept() {
}

type Set struct {
    object Expr
    name Token
    value Expr
}

func (st Set) accept() {
}

type Super struct {
    keyword Token
    method Token
}

func (sp Super) accept() {
}

type This struct {
    keyword Token
}

func (th This) accept() {
}

type Unary struct {
    operator Token
    right Expr
}

func (un Unary) accept() {
}

type Variable struct {
    name Token
}

func (vr Variable) accept() {
}
