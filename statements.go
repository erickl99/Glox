package main

type Stmt interface {
    saccept()
}

type Block struct {
    statements []Stmt
}

func (bl Block) saccept() {
}

type Expression struct {
    expr Expr
}

func (ex Expression) saccept() {
}

type Print struct {
    expr Expr
}

func (pr Print) saccept() {
}

type Var struct {
    name Token
    initializer Expr
}

func (vr Var) saccept() {
}

type If struct {
    condition Expr
    then_branch Stmt
    else_branch Stmt
}

func (iff If) saccept() {
}

type While struct {
    condition Expr
    body Stmt
}

func (wh While) saccept() {
}
