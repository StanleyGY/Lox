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
	return ""
}

type Environment struct {
	ParentEnv *Environment
	Bindings  map[string]interface{}
}

// CreateBinding creates a new binding in the current scope
func (e *Environment) CreateBinding(name string, val interface{}) bool {
	if _, ok := e.Bindings[name]; ok {
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

// UpdateBinding searches for a binding, starting from the nearest scope
// and updates the value. If there's no binding in any scope, it returns
// false.
func (e *Environment) UpdateBinding(name string, val interface{}) bool {
	_, ok := e.Bindings[name]
	if ok {
		e.Bindings[name] = val
		return true
	}
	if e.ParentEnv == nil {
		return false
	}
	return e.ParentEnv.UpdateBinding(name, val)
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

	for _, c := range stmt.Stmts {
		if err := c.Accept(p); err != nil {
			return err
		}
	}

	p.CurrEnv = newEnv.ParentEnv
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
	if !p.CurrEnv.CreateBinding(stmt.Name.Lexeme, val) {
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
	if !p.CurrEnv.CreateBinding(stmt.Name.Lexeme, loxFunc) {
		return &RuntimeError{Reason: fmt.Sprintf("double declaration for function: %s", stmt.Name.Lexeme)}
	}
	return nil
}

func (p *Interpreter) VisitClassDeclStmt(stmt *ClassDeclStmt) error {
	var initializer *LoxFunction
	var methods []*LoxFunction
	for _, funcStmt := range stmt.Methods {
		m := &LoxFunction{Declaration: funcStmt, Closure: p.CurrEnv, IsInitializer: false}
		if funcStmt.Name.Lexeme == "init" {
			initializer = m
			m.IsInitializer = true
		}
		methods = append(methods, m)
	}

	klass := &LoxClass{Name: stmt.Name.Lexeme, Initializer: initializer, Methods: methods}
	if !p.CurrEnv.CreateBinding(stmt.Name.Lexeme, klass) {
		return &RuntimeError{Reason: fmt.Sprintf("double declaration for class: %s", stmt.Name.Lexeme)}
	}
	return nil
}

func (p *Interpreter) VisitReturnStmt(stmt *ReturnStmt) error {
	var value interface{}
	var err error
	if value, err = stmt.Value.Accept(p); err != nil {
		return err
	}
	// Return an error to unwind the call stack until reaching CallExpr
	return &RuntimeReturn{Value: value}
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
	if !p.CurrEnv.UpdateBinding(expr.Name.Lexeme, val) {
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
	if val, ok = loxInstance.Properties[expr.Property.Lexeme]; !ok {
		return nil, &RuntimeError{fmt.Sprintf("class %s does not have the field %s", loxInstance.Class, expr.Property.Lexeme)}
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
