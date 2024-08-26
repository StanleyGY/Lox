package main

import (
	"bytes"
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

type Interpreter struct {
}

func (p *Interpreter) Evaluate(expr Expr) (interface{}, error) {
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

func (p *Interpreter) VisitLiteralExpr(expr *LiteralExpr) (interface{}, error) {
	return expr.Value, nil
}
