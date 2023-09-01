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

func interpret(expr Expr) {
    value, err := evaluate(expr)
    if err != nil {
        runtime_error(err.(RuntimeError))
        return
    }
    fmt.Println(stringify(value))
}

func evaluate(exp Expr) (Value, error) {
    switch t := exp.(type) {
        case Literal:
            return t.value, nil
        case Grouping:
            return evaluate(t.expression)
        case Unary:
            right, r_err := evaluate(t.right)
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
        case Binary:
            left, l_err := evaluate(t.left)
            if l_err != nil {
                return nil, l_err
            }
            right, r_err := evaluate(t.right)
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
