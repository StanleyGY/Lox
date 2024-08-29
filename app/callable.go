package main

import "fmt"

type LoxCallable interface {
	Call(interpreter *Interpreter, args []Expr) (interface{}, error)
	// Closure() *Environment
	Arity() int
}

type LoxFunction struct {
	Closure     *Environment
	Declaration *FuncDeclStmt
}

func (f *LoxFunction) Call(interpreter *Interpreter, args []Expr) (interface{}, error) {
	var err error
	var ok bool

	// Create a new env for the function call
	env := &Environment{
		Bindings:  make(map[string]interface{}),
		ParentEnv: f.Closure,
	}

	// Copy arguments into current env
	for idx := range args {
		var param string
		var argv interface{}
		if argv, err = args[idx].Accept(interpreter); err != nil {
			return nil, err
		}
		param = f.Declaration.Params[idx].Lexeme
		env.CreateBinding(param, argv)
	}

	// Evaluate function body
	var returnVal *RuntimeReturn

	lastEnv := interpreter.CurrEnv
	interpreter.CurrEnv = env

	defer func() {
		interpreter.CurrEnv = lastEnv
	}()

	if err = f.Declaration.Body.Accept(interpreter); err != nil {
		if returnVal, ok = err.(*RuntimeReturn); ok {
			return returnVal.Value, nil
		}
		return nil, err
	}
	// In case there's no return value, return a nil to caller
	return nil, nil
}

// func (f *LoxFunction) Closure() *Environment {
// 	return f.Klosure
// }

func (f *LoxFunction) Arity() int {
	return len(f.Declaration.Params)
}

type LoxClass struct {
	Name    string
	Methods []*FuncDeclStmt
}

func (c *LoxClass) Call(interpreter *Interpreter, args []Expr) (interface{}, error) {
	// Create an instance of the class
	return MakeLoxClassInstance(c), nil
}

func (c *LoxClass) Arity() int {
	// For now, constructor doesn't take any parameter
	return 0
}

func (c LoxClass) String() string {
	return c.Name
}

type LoxClassInstance struct {
	Class      *LoxClass
	Properties map[string]interface{}
}

func MakeLoxClassInstance(klass *LoxClass) *LoxClassInstance {
	return &LoxClassInstance{
		Class:      klass,
		Properties: make(map[string]interface{}),
	}
}

func (i LoxClassInstance) String() string {
	return fmt.Sprintf("instance <%p> (class %s)", &i, i.Class.String())
}
