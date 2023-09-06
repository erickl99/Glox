package main

import (
	"fmt"
	"strconv"
	"strings"
)

type Environment struct {
	enclosing *Environment
	values    map[string]Value
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

func (env Environment) get_at(distance int, name string) Value {
	curr_env := &env
	for i := 0; i < distance; i++ {
		curr_env = curr_env.enclosing
	}
	return curr_env.values[name]
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

func (env *Environment) assign_at(distance int, name Token, value Value) {
	curr_env := env
	for i := 0; i < distance; i++ {
		curr_env = curr_env.enclosing
	}
	curr_env.assign(name, value)
}

func (env Environment) String() string {
	curr := &env
	var result strings.Builder
	var i = 1
	for curr != nil {
		result.WriteString("Level " + strconv.Itoa(i) + ":\n")
		for k, v := range curr.values {
			entry := fmt.Sprintf("%s : %v\n", k, v)
			result.WriteString(entry)
		}
		curr = curr.enclosing
		i++
		if i > 100 {
			fmt.Println("In an infinite loop!")
			panic(1)
		}
	}
	return result.String()
}
