package main

import (
	"fmt"
)

type RuntimeError struct {
    message string
    token Token
}

func (re RuntimeError) Error() string {
    return re.message
}

var global_env Environment = Environment{values: make(map[string]Value)}

func interpret(statements []Stmt) {
    curr_env := global_env
    for _, stmt := range statements {
        err := execute(stmt, curr_env)
        if err != nil {
            runtime_error(err.(RuntimeError))
            break
        }
    }
}

func execute(stmt Stmt, curr_env Environment) error {
    switch t := stmt.(type) {
        case Print:
            value, err := evaluate(t.expr, curr_env)
            if err != nil {
                return err
            }
            fmt.Println(stringify(value))
            return nil
        case Expression:
            value, err := evaluate(t.expr, curr_env)
            if err != nil {
                return err
            }
            if in_repl {
                fmt.Println(stringify(value))
            }
            return nil
        case Block:
            block_env := Environment{enclosing: &curr_env, values: make(map[string]Value)}
            err := execute_block(t.statements, block_env)
            if err != nil {
                return err
            }
            return nil
        case If:
            val, err := evaluate(t.condition, curr_env)
            if err != nil {
                return err
            }
            if is_truthy(val) {
                return execute(t.then_branch, curr_env)
            } else if t.else_branch != nil {
                return execute(t.else_branch, curr_env)
            }
            return nil
        case While:
            for {
                val, err := evaluate(t.condition, curr_env)
                if err!= nil {
                    return err
                }
                if !is_truthy(val) {
                    break
                }
                if err = execute(t.body, curr_env); err != nil {
                    return err
                }
            }
            return nil
        case Var:
            var value Value
            if t.initializer != nil {
                val, err := evaluate(t.initializer, curr_env)
                if err != nil {
                    return err
                }
                value = val
            }
            curr_env.define(t.name.lexeme, value)
            return nil
    }
    return RuntimeError{message: "Internal error, unknown statement type encountered"}
}

func execute_block(statements []Stmt, block_env Environment) error {
    for _, stmt := range statements {
        err := execute(stmt, block_env)
        if err != nil {
            break
        }
    }
    return nil
}

func evaluate(exp Expr, curr_env Environment) (Value, error) {
    switch t := exp.(type) {
        case Literal:
            return t.value, nil
        case Grouping:
            return evaluate(t.expression, curr_env)
        case Unary:
            right, r_err := evaluate(t.right, curr_env)
            if r_err != nil {
                return nil, r_err
            }
            switch t.operator.t_type {
                case MINUS:
                    v_err := valid_number_operand(t.operator, right)
                    if v_err != nil {
                        return nil, v_err
                    }
                    return -right.(float64), nil
                case BANG:
                    return !is_truthy(right), nil
            }
        case Variable:
            return curr_env.get(t.name)
        case Logical:
            left, err := evaluate(t.left, curr_env)
            if err != nil {
                return nil, err
            }
            if t.operator.t_type == OR {
                if is_truthy(left) {
                    return left, nil
                }
            } else if !is_truthy(left) {
                return left, nil
            }
            return evaluate(t.right, curr_env)
        case Assign:
            value, err := evaluate(t.value, curr_env)
            if err != nil {
                return nil, err
            }
            err = curr_env.assign(t.name, value)
            if err != nil {
                return nil, err
            }
            return value, nil
        case Binary:
            left, l_err := evaluate(t.left, curr_env)
            if l_err != nil {
                return nil, l_err
            }
            right, r_err := evaluate(t.right, curr_env)
            if r_err != nil {
                return nil, r_err
            }
            switch t.operator.t_type {
                case PLUS:
                    f_left, l_ok := left.(float64)
                    f_right, r_ok := right.(float64)
                    if l_ok && r_ok {
                        return f_left + f_right, nil
                    }
                    s_left, l_ok := left.(string)
                    s_right, r_ok := right.(string)
                    if l_ok && r_ok {
                        return s_left + s_right, nil
                    }
                    return nil, RuntimeError{"Operands must be two numbers or two strings", t.operator}
                case MINUS:
                    if err := valid_number_operands(t.operator, left, right); err != nil {
                        return nil, err
                    }
                    return left.(float64) - right.(float64), nil
                case SLASH:
                    if err := valid_number_operands(t.operator, left, right); err != nil {
                        return nil, err
                    }
                    return left.(float64) / right.(float64), nil
                case STAR:
                    if err := valid_number_operands(t.operator, left, right); err != nil {
                        return nil, err
                    }
                    return left.(float64) * right.(float64), nil
                case GREATER:
                    if err := valid_number_operands(t.operator, left, right); err != nil {
                        return nil, err
                    }
                    return left.(float64) > right.(float64), nil
                case GREATER_EQUAL:
                    if err := valid_number_operands(t.operator, left, right); err != nil {
                        return nil, err
                    }
                    return left.(float64) >= right.(float64), nil
                case LESS:
                    if err := valid_number_operands(t.operator, left, right); err != nil {
                        return nil, err
                    }
                    return left.(float64) < right.(float64), nil
                case LESS_EQUAL:
                    if err := valid_number_operands(t.operator, left, right); err != nil {
                        return nil, err
                    }
                    return left.(float64) <= right.(float64), nil
                case BANG_EQUAL:
                    return !is_equal(left, right), nil
                case EQUAL_EQUAL:
                    return is_equal(left, right), nil
            }
    }
    return nil, RuntimeError{message: "Internal error, unknown expr was passed in"}
}

func is_truthy(val Value) bool {
    if val == nil {
        return false
    }
    if b_val, ok := val.(bool); ok {
        return b_val
    }
    return true
}

func is_equal(val_one Value, val_two Value) bool {
    if val_one == nil && val_two == nil {
        return true
    }
    if val_one == nil {
        return false
    }
    return val_one == val_two
}

func valid_number_operand(operator Token, operand Value) error {
    if _, ok := operand.(float64); ok {
        return nil
    }
    return RuntimeError{"Operand must be a number", operator}
}

func valid_number_operands(operator Token, left Value, right Value) error {
    _, l_ok := left.(float64)
    _, r_ok := right.(float64)
    if l_ok && r_ok {
        return nil
    }
    return RuntimeError{"Operands must be two numbers or string", operator}
}

func stringify(value Value) string {
    if value == nil {
        return "nil"
    }
    return fmt.Sprintf("%v", value)
}
