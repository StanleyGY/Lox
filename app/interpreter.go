package main

import (
	"bytes"
	"fmt"
	"reflect"
	"slices"
	"strings"
)

type RuntimeTypeError struct {
	Operator *Token
	Vals     []interface{}
}

func (e RuntimeTypeError) Error() string {
	var buffer bytes.Buffer
	buffer.WriteString("invalid types for operator ")
	buffer.WriteString(e.Operator.Lexeme)
	buffer.WriteString(" : ")
	for i, v := range e.Vals {
		buffer.WriteString(reflect.ValueOf(v).String())
		if i < len(e.Vals)-1 {
			buffer.WriteString(",")
		}
	}
	return buffer.String()
}

type Environment struct {
	ParentEnv *Environment
	Bindings  map[string]interface{}
}

func (e Environment) FindBinding(name string) (interface{}, bool) {
	val, ok := e.Bindings[name]
	if ok {
		return val, true
	}
	if e.ParentEnv == nil {
		return nil, false
	}
	return e.ParentEnv.FindBinding(name)
}

type Interpreter struct {
	Env *Environment
}

func (p *Interpreter) createEnv() {
	if p.Env == nil {
		p.Env = &Environment{Bindings: make(map[string]interface{})}
	} else {
		newEnv := &Environment{
			ParentEnv: p.Env,
			Bindings:  make(map[string]interface{}),
		}
		p.Env = newEnv
	}
}

func (p *Interpreter) restoreEnv() {
	p.Env = p.Env.ParentEnv
}

func (p *Interpreter) Evaluate(stmts []Stmt) error {
	p.createEnv()
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
	return stmt.ElseBranch.Accept(p)
}

func (p *Interpreter) VisitBlockStmt(stmt *BlockStmt) error {
	p.createEnv()
	for _, c := range stmt.Stmts {
		if err := c.Accept(p); err != nil {
			return err
		}
	}
	p.restoreEnv()
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
	p.Env.Bindings[stmt.Name.Lexeme] = val
	return nil
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
	_, ok := p.Env.Bindings[expr.Name.Lexeme]
	if !ok {
		return nil, fmt.Errorf("assigns value to an undefined variable: %s", expr.Name.Lexeme)
	}

	p.Env.Bindings[expr.Name.Lexeme] = val
	return nil, nil
}

func (p *Interpreter) VisitLiteralExpr(expr *LiteralExpr) (interface{}, error) {
	return expr.Value, nil
}

func (p *Interpreter) VisitVariableExpr(expr *VariableExpr) (interface{}, error) {
	// Design choice: when reading the value of a variable, we can walk back
	// to the outer scopes
	val, ok := p.Env.FindBinding(expr.Name.Lexeme)
	if !ok {
		return nil, fmt.Errorf("reference an undefined variable: %s", expr.Name.Lexeme)
	}
	return val, nil
}
