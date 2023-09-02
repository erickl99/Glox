package main

type Stmt interface {
    saccept()
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
