package main

import "fmt"

type LoxCallable interface {
	Call(interpreter *Interpreter, args []Expr) (interface{}, error)
	// Closure() *Environment
	Arity() int
}

type LoxFunction struct {
	IsInitializer bool
	Closure       *Environment
	Declaration   *FuncDeclStmt
}

func (f *LoxFunction) Call(interpreter *Interpreter, args []Expr) (interface{}, error) {
	var err error

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
	var hasReturn bool

	lastEnv := interpreter.CurrEnv
	interpreter.CurrEnv = env
	defer func() { interpreter.CurrEnv = lastEnv }()

	if err = f.Declaration.Body.Accept(interpreter); err != nil {
		if returnVal, hasReturn = err.(*RuntimeReturn); !hasReturn {
			return nil, err
		}
	}
	if f.IsInitializer {
		// In case user calls the init() function explicitly,
		// force rewrite the return value of an initializer to the instance itself.
		instance, _ := f.Closure.FindBinding("this", 0)
		return instance, nil
	}
	if hasReturn {
		return returnVal.Value, nil
	}
	// In case there's no return value, return a nil to caller
	return nil, nil
}

func (f *LoxFunction) Arity() int {
	return len(f.Declaration.Params)
}

func (f *LoxFunction) Bind(name string, val interface{}) {
	f.Closure.UpdateBinding(name, val, true)
}

type LoxClass struct {
	Name        string
	SuperClass  *LoxClass
	Initializer *LoxFunction
	Methods     map[string]*LoxFunction
}

func (c *LoxClass) FindMethod(name string) (*LoxFunction, bool) {
	val, ok := c.Methods[name]
	return val, ok
}

func (c *LoxClass) Call(interpreter *Interpreter, args []Expr) (interface{}, error) {
	// Create axn instance of the class
	properties := make(map[string]interface{})

	instance := &LoxClassInstance{
		Class:      c,
		Properties: properties,
	}

	// Immediately call the user-defined constructor
	if c.Initializer != nil {
		c.Initializer.Bind("this", instance)
		c.Initializer.Call(interpreter, args)
	}
	return instance, nil
}

func (c *LoxClass) Arity() int {
	// The arity of class is determined by the number of arguments accepted in constructor
	if c.Initializer == nil {
		return 0
	}
	return c.Initializer.Arity()
}

func (c LoxClass) String() string {
	return c.Name
}

type LoxClassInstance struct {
	Class      *LoxClass
	Properties map[string]interface{}
}

func (i *LoxClassInstance) FindProperty(name string) (interface{}, bool) {
	var val interface{}
	var ok bool

	// Try get an instance property (owned by individual instance)
	if val, ok = i.Properties[name]; ok {
		return val, true
	}

	// Try get an instance class methods (shared by all class instances)
	if val, ok = i.Class.FindMethod(name); ok {
		f := val.(*LoxFunction)
		f.Bind("this", i)
		return val, true
	}

	// Try get from super class
	if i.Class.SuperClass != nil {
		if val, ok = i.Class.SuperClass.FindMethod(name); ok {
			return val, true
		}
	}
	return nil, false
}

func (i *LoxClassInstance) String() string {
	return fmt.Sprintf("%s object at <%p>", i.Class.String(), i)
}
