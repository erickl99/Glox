package main

type LoxClass struct {
	name       string
	superclass *LoxClass
	methods    map[string]LoxFunction
}

func (lc LoxClass) find_method(name string) (LoxFunction, bool) {
	if md, ok := lc.methods[name]; ok {
		return md, true
	}
	if lc.superclass != nil {
		return lc.superclass.find_method(name)
	}
	return LoxFunction{}, false
}

func (lc LoxClass) call(arguments []Value) Value {
	instance := LoxInstance{lc, make(map[string]Value)}
	if initializer, ok := lc.find_method("init"); ok {
		initializer.bind(instance).call(arguments)
	}
	return instance
}

func (lc LoxClass) arity() int {
	if initializer, ok := lc.find_method("init"); ok {
		return initializer.arity()
	}
	return 0
}

func (lc LoxClass) String() string {
	return lc.name
}

type LoxInstance struct {
	klass  LoxClass
	fields map[string]Value
}

func (li LoxInstance) get(name Token) (Value, error) {
	if val, ok := li.fields[name.lexeme]; ok {
		return val, nil
	} else if md, ok := li.klass.find_method(name.lexeme); ok {
		return md.bind(li), nil
	} else {
		return nil, RuntimeError{"Undefined property '" + name.lexeme + "'.", name}
	}
}

func (li *LoxInstance) set(name Token, val Value) {
	li.fields[name.lexeme] = val
}

func (li LoxInstance) String() string {
	return li.klass.name + " instance"
}
