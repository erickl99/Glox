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
