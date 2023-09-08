package main

import (
	"fmt"
	"os"
)

type FunctionType int

const (
	NONE FunctionType = iota
	FUNCTION
)


type Stack []map[string]bool

func (st Stack) empty() bool {
	return len(st) == 0
}

func (st Stack) peek() (map[string]bool, bool) {
	if st.empty() {
		return nil, false
	}
	return st[len(st)-1], true
}

func (st *Stack) pop() (map[string]bool, bool) {
	if st.empty() {
		return nil, false
	}
	entry, _ := st.peek()
	(*st) = (*st)[:len(*st)-1]
	return entry, true
}

func (st *Stack) push(entry map[string]bool) {
	(*st) = append((*st), entry)
}

var init_scopes *Stack
var curr_function FunctionType = NONE

func resolve(statements []Stmt) {
	init_scopes = new(Stack)
	resolve_stmts(statements, init_scopes)
}
func resolve_stmts(statements []Stmt, scopes *Stack) {
	for _, stmt := range statements {
		resolve_stmt(stmt, scopes)
	}
}

func resolve_stmt(stmt Stmt, scopes *Stack) {
	switch t := stmt.(type) {
	case Block:
		begin_scope(scopes)
		resolve_stmts(t.statements, scopes)
		end_scope(scopes)
		return
	case Expression:
		resolve_expr(t.expr, scopes)
		return
	case Func:
		declare(t.name, scopes)
		define(t.name, scopes)
		resolve_func(t, scopes, FUNCTION)
		return
	case If:
		resolve_expr(t.condition, scopes)
		resolve_stmt(t.then_branch, scopes)
		if t.else_branch != nil {
			resolve_stmt(t.else_branch, scopes)
		}
		return
	case Print:
		resolve_expr(t.expr, scopes)
		return
	case Return:
		if curr_function == NONE {
			token_error(t.keyword, "Can't return from top level routine")
		}
		if t.value != nil {
			resolve_expr(t.value, scopes)
		}
		return
	case Var:
		declare(t.name, scopes)
		if t.initializer != nil {
			resolve_expr(t.initializer, scopes)
		}
		define(t.name, scopes)
		return
	case While:
		resolve_expr(t.condition, scopes)
		resolve_stmt(t.body, scopes)
		return
	}
	fmt.Fprintf(os.Stderr, "Internal error, encountered unkown statement type: %v", stmt)
	panic(69)
}

func resolve_expr(expr Expr, scopes *Stack) {
	switch t := expr.(type) {
	case Assign:
		resolve_expr(t.value, scopes)
		resolve_local(t, t.name, scopes)
		return
	case Binary:
		resolve_expr(t.left, scopes)
		resolve_expr(t.right, scopes)
		return
	case Call:
		resolve_expr(t.callee, scopes)
		for _, arg := range t.arguments {
			resolve_expr(arg, scopes)
		}
		return
	case Grouping:
		resolve_expr(t.expression, scopes)
		return
	case Literal:
		return
	case Logical:
		resolve_expr(t.left, scopes)
		resolve_expr(t.right, scopes)
		return
	case Unary:
		resolve_expr(t.right, scopes)
		return
	case Variable:
		if scope, ok := scopes.peek(); ok {
			if resolved, ok := scope[t.name.lexeme]; ok && !resolved {
				token_error(t.name, "Can't read local variable in its own initializer")
			}
		}
		resolve_local(t, t.name, scopes)
		return
	}
	fmt.Fprintf(os.Stderr, "Internal error, encountered unkown expression type: %v", expr)
	panic(69)
}

func resolve_func(function Func, scopes *Stack, f_type FunctionType) {
	enclosing_function := curr_function
	curr_function = f_type
	begin_scope(scopes)
	for _, param := range function.params {
		declare(param, scopes)
		define(param, scopes)
	}
	resolve_stmts(function.body, scopes)
	end_scope(scopes)
	curr_function = enclosing_function
}

func resolve_local(expr Expr, name Token, scopes *Stack) {
	for i := len(*scopes) - 1; i > -1; i-- {
		if _, ok := (*scopes)[i][name.lexeme]; ok {
			set_scope(expr, len(*scopes) - i - 1)
			return
		}
	}
}

func declare(name Token, scopes *Stack) {
	if scopes.empty() {
		return
	}
	scope, _ := scopes.peek()
	if _, ok := scope[name.lexeme]; ok {
		token_error(name, "Already a variable with this name in this scope")
	}
	scope[name.lexeme] = false
}

func define(name Token, scopes *Stack) {
	if scopes.empty() {
		return
	}
	scope, _ := scopes.peek()
	scope[name.lexeme] = true
}

func begin_scope(scopes *Stack) {
	scopes.push(make(map[string]bool))
}

func end_scope(scopes *Stack) {
	scopes.pop()
}
