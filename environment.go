package main

import "fmt"

type Environment struct {
    enclosing *Environment
    values map[string]Value
}

func (env *Environment) define(name string, value Value) {
    env.values[name] = value
}

func (env Environment) get(name Token) (Value, error) {
    if value, ok := env.values[name.lexeme]; ok {
        return value, nil
    }
    if env.enclosing != nil {
        value, err := env.enclosing.get(name)
        if err != nil {
            return nil, err
        }
        return value, nil
    }
    msg := fmt.Sprintf("Undefined variable '%v'.", name.lexeme)
    return nil, RuntimeError{message: msg, token: name}
}

func (env *Environment) assign(name Token, value Value) error {
    if _, ok := env.values[name.lexeme]; ok {
        env.values[name.lexeme] = value
        return nil
    }
    if env.enclosing != nil {
        err := env.enclosing.assign(name, value)
        if err != nil {
            return err
        }
        return nil
    }
    msg := fmt.Sprintf("Undefined variable '%v'.", name.lexeme)
    return RuntimeError{message: msg, token: name}
}
