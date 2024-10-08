package main

import (
	"bytes"
	"fmt"
	"reflect"
	"slices"
	"strings"
)

type RuntimeError struct {
	Reason string
}

func (e RuntimeError) Error() string {
	return fmt.Sprintf("runtime error: %s", e.Reason)
}

type RuntimeTypeError struct {
	Operator *Token
	Vals     []interface{}
}

func (e RuntimeTypeError) Error() string {
	var buffer bytes.Buffer
	buffer.WriteString("runtime error: invalid types for operator ")
	buffer.WriteString(e.Operator.Lexeme)
	buffer.WriteString("on values: ")
	for i, v := range e.Vals {
		buffer.WriteString(reflect.ValueOf(v).String())
		if i < len(e.Vals)-1 {
			buffer.WriteString(",")
		}
	}
	return buffer.String()
}

type RuntimeReturn struct {
	Value interface{}
}

func (e RuntimeReturn) Error() string {
	return "runtime error: return"
}

type RuntimeBreak struct{}

func (e RuntimeBreak) Error() string {
	return "runtime error: break"
}

type Environment struct {
	ParentEnv *Environment
	Bindings  map[string]interface{}
}

// CreateBinding creates a new binding in the current scope
func (e *Environment) CreateBinding(name string, val interface{}, update bool) bool {
	if _, ok := e.Bindings[name]; ok && !update {
		return false
	}
	e.Bindings[name] = val
	return true
}

func (e *Environment) FindBinding(name string, dist int) (interface{}, bool) {
	ancestorEnv := e
	for i := 0; i < dist; i++ {
		ancestorEnv = ancestorEnv.ParentEnv
	}
	val, ok := ancestorEnv.Bindings[name]
	return val, ok
}

func (e *Environment) UpdateBinding(name string, val interface{}, dist int) bool {
	ancestorEnv := e
	for i := 0; i < dist; i++ {
		ancestorEnv = ancestorEnv.ParentEnv
	}
	_, ok := ancestorEnv.Bindings[name]
	if !ok {
		return false
	}

	ancestorEnv.Bindings[name] = val
	return true
}

type Interpreter struct {
	// Number of scope hops between a variable usage and its declaration
	ScopeHops map[Expr]int
	// In Lox, runtime environment is a dynamic manifestation of static scope
	Globals *Environment
	CurrEnv *Environment
}

func MakeInterpreter() *Interpreter {
	globals := &Environment{
		Bindings: make(map[string]interface{}),
	}
	return &Interpreter{Globals: globals, CurrEnv: globals, ScopeHops: make(map[Expr]int)}
}

// Resolve tracks where a referenced variable is declared.
// This is possible since Lox uses static scope.
func (p *Interpreter) Resolve(expr Expr, dist int) {
	p.ScopeHops[expr] = dist
}

func (p *Interpreter) Evaluate(stmts []Stmt) error {
	for _, stmt := range stmts {
		if err := stmt.Accept(p); err != nil {
			return err
		}
	}
	return nil
}

func (p *Interpreter) EvaluateExpr(expr Expr) (interface{}, error) {
	return expr.Accept(p)
}

func (p *Interpreter) isTruthy(v interface{}) bool {
	t := reflect.ValueOf(v)
	if !t.IsValid() {
		return false
	}
	if t.Kind() == reflect.Bool {
		return v.(bool)
	}
	// TODO: not truthy for empty array and map
	return true
}

func (p *Interpreter) checkType(op *Token, v interface{}, expectedTypes []reflect.Kind) error {
	t := reflect.ValueOf(v)
	if slices.Contains(expectedTypes, t.Kind()) {
		return nil
	}
	return RuntimeTypeError{Operator: op, Vals: []interface{}{v}}
}

func (p *Interpreter) checkTypes(op *Token, vals []interface{}, expectedTypes []reflect.Kind) error {
	var err error

	for _, et := range expectedTypes {
		valid := true

		// Make sure all `vals` have the same types
		for _, v := range vals {
			if err = p.checkType(op, v, []reflect.Kind{et}); err != nil {
				valid = false
				break
			}
		}
		if valid {
			return nil
		}
	}
	return RuntimeTypeError{Operator: op, Vals: vals}
}

func (p *Interpreter) VisitIfStmt(stmt *IfStmt) error {
	var r interface{}
	var err error
	if r, err = stmt.Condition.Accept(p); err != nil {
		return err
	}
	if p.isTruthy(r) {
		return stmt.ThenBranch.Accept(p)
	}
	if stmt.ElseBranch != nil {
		return stmt.ElseBranch.Accept(p)
	}
	return nil
}

func (p *Interpreter) VisitWhileStmt(stmt *WhileStmt) error {
	var r interface{}
	var err error
	for {
		if r, err = stmt.Condition.Accept(p); err != nil {
			return err
		}
		if !p.isTruthy(r) {
			return nil
		}
		if err = stmt.Body.Accept(p); err != nil {
			if _, ok := err.(*RuntimeBreak); ok {
				return nil
			}
			return err
		}
	}
}

func (p *Interpreter) VisitBlockStmt(stmt *BlockStmt) error {
	newEnv := &Environment{
		ParentEnv: p.CurrEnv,
		Bindings:  make(map[string]interface{}),
	}
	p.CurrEnv = newEnv
	defer func() { p.CurrEnv = newEnv.ParentEnv }()

	for _, c := range stmt.Stmts {
		if err := c.Accept(p); err != nil {
			return err
		}
	}

	return nil
}

func (p *Interpreter) VisitInlineExprStmt(stmt *InlineExprStmt) error {
	_, err := stmt.Child.Accept(p)
	return err
}

func (p *Interpreter) VisitPrintStmt(stmt *PrintStmt) error {
	var val interface{}
	var err error

	if val, err = stmt.Child.Accept(p); err != nil {
		return err
	}
	fmt.Println(val)
	return nil
}

func (p *Interpreter) VisitVarDeclStmt(stmt *VarDeclStmt) error {
	var val interface{}
	var err error

	if val, err = stmt.Initializer.Accept(p); err != nil {
		return err
	}
	if !p.CurrEnv.CreateBinding(stmt.Name.Lexeme, val, false) {
		return &RuntimeError{Reason: fmt.Sprintf("double declaration for variable: %s", stmt.Name.Lexeme)}
	}
	return nil
}

func (p *Interpreter) VisitFunDeclStmt(stmt *FuncDeclStmt) error {
	// Things about closure:
	// For a function, closure is the runtime environment when the function is declared
	// For a method, closure is the runtime environment when the method is declared
	//					and a "this" property when the class is instantiated
	loxFunc := &LoxFunction{Declaration: stmt, Closure: p.CurrEnv}
	if !p.CurrEnv.CreateBinding(stmt.Name.Lexeme, loxFunc, false) {
		return &RuntimeError{Reason: fmt.Sprintf("double declaration for function: %s", stmt.Name.Lexeme)}
	}
	return nil
}

func (p *Interpreter) VisitClassDeclStmt(stmt *ClassDeclStmt) error {
	var initializer *LoxFunction
	var superClass *LoxClass

	// Handle inheritance
	if stmt.SuperClass != nil {
		val, ok := p.CurrEnv.FindBinding(stmt.SuperClass.Name.Lexeme, p.ScopeHops[stmt.SuperClass])
		if !ok {
			return &RuntimeError{fmt.Sprintf("super class is declared: %s", stmt.SuperClass.Name.Lexeme)}
		}
		if superClass, ok = val.(*LoxClass); !ok {
			return &RuntimeError{fmt.Sprintf("super class is not a class: %s", stmt.SuperClass.Name.Lexeme)}
		}
	}

	// Handle class methods
	methods := make(map[string]*LoxFunction)
	for _, funcStmt := range stmt.Methods {
		// With super class, the environment chain of a method when it's called looks like this:
		// Global -> Block -> Class Closure ("super", "this") -> Method env.
		// Unlike "this" which is specific to individual class instances,
		// "super" is shared across all instances. So we only bind once, and bind
		// when class is declared.
		closure := p.CurrEnv
		if stmt.SuperClass != nil {
			closure = &Environment{
				ParentEnv: closure,
				Bindings:  make(map[string]interface{}),
			}
			closure.CreateBinding("super", superClass, true)
		}

		m := &LoxFunction{Declaration: funcStmt, Closure: closure, IsInitializer: false}
		mName := funcStmt.Name.Lexeme
		methods[mName] = m

		if mName == "init" {
			initializer = m
			m.IsInitializer = true
		}
	}

	// Try create a binding for the class decl
	klass := &LoxClass{
		Name:        stmt.Name.Lexeme,
		SuperClass:  superClass,
		Initializer: initializer,
		Methods:     methods,
	}
	if !p.CurrEnv.CreateBinding(stmt.Name.Lexeme, klass, false) {
		return &RuntimeError{fmt.Sprintf("double declaration for class: %s", stmt.Name.Lexeme)}
	}
	return nil
}

func (p *Interpreter) VisitReturnStmt(stmt *ReturnStmt) error {
	var value interface{}
	var err error

	if stmt.Value != nil {
		if value, err = stmt.Value.Accept(p); err != nil {
			return err
		}
	}
	// Return an error to unwind the call stack until reaching CallExpr
	return &RuntimeReturn{Value: value}
}

func (p *Interpreter) VisitBreakStmt(stmt *BreakStmt) error {
	// Return an error to unwind the call stack until reaching WhileStmt
	return &RuntimeBreak{}
}

func (p *Interpreter) VisitBinaryExpr(expr *BinaryExpr) (interface{}, error) {
	var leftVal interface{}
	var rightVal interface{}
	var err error

	if leftVal, err = expr.Left.Accept(p); err != nil {
		return nil, err
	}
	if rightVal, err = expr.Right.Accept(p); err != nil {
		return nil, err
	}

	switch expr.Operator.Type {
	case Plus:
		if err = p.checkTypes(
			expr.Operator,
			[]interface{}{leftVal, rightVal},
			[]reflect.Kind{reflect.String, reflect.Float64},
		); err != nil {
			return nil, err
		}
	case Minus:
		fallthrough
	case Star:
		fallthrough
	case Slash:
		fallthrough
	case Greater:
		fallthrough
	case GreaterEqual:
		fallthrough
	case Less:
		fallthrough
	case LessEqual:
		if err = p.checkTypes(
			expr.Operator,
			[]interface{}{leftVal, rightVal},
			[]reflect.Kind{reflect.Float64},
		); err != nil {
			return nil, err
		}
	}

	switch expr.Operator.Type {
	case Plus:
		if reflect.TypeOf(leftVal).Kind() == reflect.String {
			return strings.Join([]string{leftVal.(string), rightVal.(string)}, ""), nil
		} else {
			return leftVal.(float64) + (rightVal).(float64), nil
		}
	case Minus:
		return leftVal.(float64) - (rightVal).(float64), nil
	case Star:
		return leftVal.(float64) * (rightVal).(float64), nil
	case Slash:
		return leftVal.(float64) / (rightVal).(float64), nil
	case Greater:
		return leftVal.(float64) > rightVal.(float64), nil
	case GreaterEqual:
		return leftVal.(float64) >= rightVal.(float64), nil
	case Less:
		return leftVal.(float64) < rightVal.(float64), nil
	case LessEqual:
		return leftVal.(float64) <= rightVal.(float64), nil
	case BangEqual:
		return !reflect.DeepEqual(leftVal, rightVal), nil
	case EqualEqual:
		return reflect.DeepEqual(leftVal, rightVal), nil
	}

	return nil, RuntimeTypeError{Operator: expr.Operator, Vals: []interface{}{leftVal, rightVal}}
}

func (p *Interpreter) VisitLogicalExpr(expr *LogicExpr) (interface{}, error) {
	var leftVal interface{}
	var err error

	if leftVal, err = expr.Left.Accept(p); err != nil {
		return nil, err
	}

	// Logical operator will return a value that guarantees
	// the truthness of this operator
	switch expr.Operator.Type {
	case And:
		if !p.isTruthy(leftVal) {
			return false, nil
		}
		return expr.Right.Accept(p)
	case Or:
		if p.isTruthy(leftVal) {
			return leftVal, nil
		}
		return expr.Right.Accept(p)
	}

	return nil, RuntimeTypeError{Operator: expr.Operator, Vals: []interface{}{leftVal}}
}

func (p *Interpreter) VisitUnaryExpr(expr *UnaryExpr) (interface{}, error) {
	var rightVal interface{}
	var err error

	if rightVal, err = expr.Right.Accept(p); err != nil {
		return nil, err
	}

	switch expr.Operator.Type {
	case Minus:
		if err = p.checkType(expr.Operator, rightVal, []reflect.Kind{reflect.Float64}); err != nil {
			return nil, err
		}
		return -rightVal.(float64), nil
	case Bang:
		return !p.isTruthy(rightVal), nil
	}

	return nil, RuntimeTypeError{Operator: expr.Operator, Vals: []interface{}{rightVal}}
}

func (p *Interpreter) VisitGroupingExpr(expr *GroupingExpr) (interface{}, error) {
	return expr.Child.Accept(p)
}

func (p *Interpreter) VisitAssignExpr(expr *AssignExpr) (interface{}, error) {
	var val interface{}
	var err error

	if val, err = expr.Value.Accept(p); err != nil {
		return nil, err
	}

	// Design choice: the variable must be defined in the current scope
	// before assigning another value
	if !p.CurrEnv.UpdateBinding(expr.Name.Lexeme, val, p.ScopeHops[expr]) {
		return nil, &RuntimeError{Reason: fmt.Sprintf("assigns value to an undefined variable: %s", expr.Name.Lexeme)}
	}
	return nil, nil
}

func (p *Interpreter) VisitCallExpr(expr *CallExpr) (interface{}, error) {
	var err error
	var ok bool

	// Look up the call binding (i.e. function / class constructor)
	var callee interface{}
	var callable LoxCallable

	if callee, err = expr.Callee.Accept(p); err != nil {
		return nil, err
	}
	if callable, ok = callee.(LoxCallable); !ok {
		return nil, &RuntimeError{Reason: "not a function declaration"}
	}

	// Validate arity
	if callable.Arity() != len(expr.Arguments) {
		return nil, &RuntimeError{Reason: "function call supplies incorrect number of parameters"}
	}

	// Call function
	return callable.Call(p, expr.Arguments)
}

func (p *Interpreter) VisitGetPropertyExpr(expr *GetPropertyExpr) (interface{}, error) {
	var object interface{}
	var loxInstance *LoxClassInstance
	var ok bool
	var err error

	// Get lox class instance
	if object, err = expr.Object.Accept(p); err != nil {
		return nil, err
	}
	if loxInstance, ok = object.(*LoxClassInstance); !ok {
		return nil, &RuntimeError{"cannot convert to a LoxClass instance"}
	}

	// Access property from the Lox class instance
	var val interface{}
	name := expr.Property.Lexeme

	if val, ok = loxInstance.FindProperty(name); !ok {
		return nil, &RuntimeError{fmt.Sprintf("class %s does not have the field %s", loxInstance.Class, name)}
	}
	return val, nil
}

func (p *Interpreter) VisitSetPropertyExpr(expr *SetPropertyExpr) (interface{}, error) {
	var object interface{}
	var loxInstance *LoxClassInstance
	var val interface{}
	var ok bool
	var err error

	// Get lox class instance
	if object, err = expr.Object.Accept(p); err != nil {
		return nil, err
	}
	if loxInstance, ok = object.(*LoxClassInstance); !ok {
		return nil, &RuntimeError{"cannot convert to a LoxClass instance"}
	}

	// Evaluate value
	if val, err = expr.Value.Accept(p); err != nil {
		return nil, err
	}

	loxInstance.Properties[expr.Property.Lexeme] = val
	return val, nil
}

func (p *Interpreter) VisitLiteralExpr(expr *LiteralExpr) (interface{}, error) {
	return expr.Value, nil
}

func (p *Interpreter) VisitVariableExpr(expr *VariableExpr) (interface{}, error) {
	// Design choice: when searching for the value of a variable,
	// it can be traced back to the outer scopes
	val, ok := p.CurrEnv.FindBinding(expr.Name.Lexeme, p.ScopeHops[expr])
	if !ok {
		return nil, &RuntimeError{Reason: fmt.Sprintf("reference an undefined variable: %s", expr.Name.Lexeme)}
	}
	return val, nil
}

func (p *Interpreter) VisitThisExpr(expr *ThisExpr) (interface{}, error) {
	val, ok := p.CurrEnv.FindBinding("this", p.ScopeHops[expr])
	if !ok {
		return nil, &RuntimeError{Reason: "reference an unbounded \"this\""}
	}
	return val, nil
}

func (p *Interpreter) VisitSuperExpr(expr *SuperExpr) (interface{}, error) {
	// When `SuperExpr` is evaluated, we must be in a declaration body of a method
	val, ok := p.CurrEnv.FindBinding("super", p.ScopeHops[expr])
	if !ok {
		return nil, &RuntimeError{Reason: "reference an unbounded \"super\""}
	}
	superClass, ok := val.(*LoxClass)
	if !ok {
		return nil, &RuntimeError{Reason: "super does not refer to a class"}
	}
	method, ok := superClass.FindMethod(expr.Property.Lexeme)
	if !ok {
		return nil, &RuntimeError{fmt.Sprintf("super class does not have this method: %s", expr.Property.Lexeme)}
	}

	// When a super method is called, it's still binded to the current instance
	thisInstance, _ := p.CurrEnv.FindBinding("this", p.ScopeHops[expr]-1)
	method.Closure.CreateBinding("this", thisInstance, true)
	return method, nil
}
