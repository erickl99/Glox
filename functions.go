package main

import (
	"fmt"
	"time"
)

type LoxCallable interface {
    call(arguments []Value) Value
    arity() int
}

// Global functions

type Clock struct {}

func (cl Clock) call(globals Environment, arguments []Value) Value {
    return time.Now().UnixMilli()
}

func (cl Clock) arity() int {
    return 0
}

func (cl Clock) String() string {
    return "<native fn>"
}


type ToString struct {}

func (ts ToString) call(globals Environment, arguments []Value) Value {
    return fmt.Sprintf("%v", arguments[0])
}

func (ts ToString) arity() int {
    return 1
}

func (ts ToString) String() string {
    return "<native fn>"
}

// Type representing Lox functions
type LoxFunction struct {
    declaration Func
    closure *Environment
}

func (lf LoxFunction) call(arguments []Value) Value {
    func_env := Environment{lf.closure, make(map[string]Value)}
    for i := 0; i < len(lf.declaration.params); i++ {
        func_env.define(lf.declaration.params[i].lexeme, arguments[i])
    }
    err := execute_block(lf.declaration.body, &func_env)
    if err != nil {
        if return_val, ok := err.(ReturnVal); ok {
            return return_val.value
        }
        runtime_error(RuntimeError{message: err.Error()})
        return nil
    }
    return nil
}

func (lf LoxFunction) arity() int {
    return len(lf.declaration.params)
}

func (lf LoxFunction) String() string {
    return "<fn " + lf.declaration.name.lexeme + ">"
}
