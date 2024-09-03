package main

import (
	"bytes"
	"fmt"
)

// Resolver performs the semantic analysis that resolves a variable always to the same declaration,
// by calculating number of "hops" away the declared variable will be in the environment chain.
//
// An example chain of scopes: Global -> Block -> Class Decl -> Class Method
type Resolver struct {
	scopes         []map[string]bool
	intepreter     *Interpreter
	enclosingFunc  *FuncDeclStmt
	enclosingClass *ClassDeclStmt
}

type SemanticsError struct {
	Reason string
}

func (e SemanticsError) Error() string {
	var buf bytes.Buffer
	buf.WriteString("Semantics error: ")
	buf.WriteString(e.Reason)
	return buf.String()
}

func MakeResolver(interpreter *Interpreter) *Resolver {
	scopes := []map[string]bool{make(map[string]bool)}
	return &Resolver{scopes: scopes, intepreter: interpreter, enclosingFunc: nil}
}

func (r *Resolver) Resolve(stmts []Stmt) error {
	for _, s := range stmts {
		if err := s.Accept(r); err != nil {
			return err
		}
	}
	return nil
}

func (r *Resolver) beginScope() {
	r.scopes = append(r.scopes, make(map[string]bool))
}

func (r *Resolver) endScope() {
	r.scopes = r.scopes[0 : len(r.scopes)-1]
}

func (r *Resolver) declare(name string) bool {
	if _, declared := r.scopes[len(r.scopes)-1][name]; declared {
		return false
	}
	r.scopes[len(r.scopes)-1][name] = false
	return true
}

func (r *Resolver) define(name string) {
	r.scopes[len(r.scopes)-1][name] = true
}

func (r *Resolver) searchScopes(name string) (int, bool) {
	for i := len(r.scopes) - 1; i >= 0; i-- {
		defined, ok := r.scopes[i][name]
		if ok && defined {
			return len(r.scopes) - 1 - i, true
		}
	}
	return 0, false
}

func (r *Resolver) VisitVarDeclStmt(stmt *VarDeclStmt) error {
	// A variable declaration introduces a new binding in current scope
	if !r.declare(stmt.Name.Lexeme) {
		return &SemanticsError{fmt.Sprintf("redefining variable: %s", stmt.Name.Lexeme)}
	}
	if _, err := stmt.Initializer.Accept(r); err != nil {
		return err
	}
	r.define(stmt.Name.Lexeme)
	return nil
}

func (r *Resolver) VisitFunDeclStmt(stmt *FuncDeclStmt) error {
	// A function declaration introduces a new binding in the block/global level,
	// and creates a new scope for function body
	if !r.declare(stmt.Name.Lexeme) {
		return &SemanticsError{fmt.Sprintf("redefining function: %s", stmt.Name.Lexeme)}
	}
	r.define(stmt.Name.Lexeme)

	r.beginScope()
	lastEnclosingFunc := r.enclosingFunc
	r.enclosingFunc = stmt
	for _, param := range stmt.Params {
		r.declare(param.Lexeme)
		r.define(param.Lexeme)
	}
	if err := stmt.Body.Accept(r); err != nil {
		return err
	}
	r.enclosingFunc = lastEnclosingFunc
	r.endScope()
	return nil
}

func (r *Resolver) VisitClassDeclStmt(stmt *ClassDeclStmt) error {
	if !r.declare(stmt.Name.Lexeme) {
		return &SemanticsError{fmt.Sprintf("redefining class: %s", stmt.Name.Lexeme)}
	}
	r.define(stmt.Name.Lexeme)

	lastEnclosingClass := r.enclosingClass
	r.enclosingClass = stmt

	if stmt.SuperClass != nil {
		if _, err := stmt.SuperClass.Accept(r); err != nil {
			return err
		}

	}
	r.beginScope()
	r.define("super")
	r.define("this")
	for _, method := range stmt.Methods {
		if err := method.Accept(r); err != nil {
			return err
		}
	}
	r.endScope()

	r.enclosingClass = lastEnclosingClass
	return nil
}

func (r *Resolver) VisitInlineExprStmt(stmt *InlineExprStmt) error {
	if _, err := stmt.Child.Accept(r); err != nil {
		return err
	}
	return nil
}

func (r *Resolver) VisitPrintStmt(stmt *PrintStmt) error {
	if _, err := stmt.Child.Accept(r); err != nil {
		return err
	}
	return nil
}

func (r *Resolver) VisitBlockStmt(stmt *BlockStmt) error {
	// A block introduces a new scope
	r.beginScope()
	for _, s := range stmt.Stmts {
		if err := s.Accept(r); err != nil {
			return err
		}
	}
	r.endScope()
	return nil
}

func (r *Resolver) VisitIfStmt(stmt *IfStmt) error {
	if _, err := stmt.Condition.Accept(r); err != nil {
		return err
	}
	if err := stmt.ThenBranch.Accept(r); err != nil {
		return err
	}
	if stmt.ElseBranch != nil {
		if err := stmt.ElseBranch.Accept(r); err != nil {
			return err
		}
	}
	return nil
}

func (r *Resolver) VisitWhileStmt(stmt *WhileStmt) error {
	if _, err := stmt.Condition.Accept(r); err != nil {
		return err
	}
	if err := stmt.Body.Accept(r); err != nil {
		return err
	}
	return nil
}

func (r *Resolver) VisitReturnStmt(stmt *ReturnStmt) error {
	if r.enclosingFunc == nil {
		return &SemanticsError{"return must be inside of a function"}
	}

	if stmt.Value != nil {
		if r.enclosingFunc.Name.Lexeme == "init" {
			// A init function declared in class should just be a return without value in code
			return &SemanticsError{"class initializer should return nothing"}
		}
		if _, err := stmt.Value.Accept(r); err != nil {
			return err
		}
	}
	return nil
}

func (r *Resolver) VisitBinaryExpr(expr *BinaryExpr) (interface{}, error) {
	if _, err := expr.Left.Accept(r); err != nil {
		return nil, err
	}
	if _, err := expr.Right.Accept(r); err != nil {
		return nil, err
	}
	return nil, nil
}

func (r *Resolver) VisitUnaryExpr(expr *UnaryExpr) (interface{}, error) {
	if _, err := expr.Right.Accept(r); err != nil {
		return nil, err
	}
	return nil, nil
}

func (r *Resolver) VisitLogicalExpr(expr *LogicExpr) (interface{}, error) {
	if _, err := expr.Left.Accept(r); err != nil {
		return nil, err
	}
	if _, err := expr.Right.Accept(r); err != nil {
		return nil, err
	}
	return nil, nil
}

func (r *Resolver) VisitGroupingExpr(expr *GroupingExpr) (interface{}, error) {
	if _, err := expr.Child.Accept(r); err != nil {
		return nil, err
	}
	return nil, nil
}

func (r *Resolver) VisitLiteralExpr(expr *LiteralExpr) (interface{}, error) {
	return nil, nil
}

func (r *Resolver) VisitAssignExpr(expr *AssignExpr) (interface{}, error) {
	if _, err := expr.Value.Accept(r); err != nil {
		return nil, err
	}

	dist, defined := r.searchScopes(expr.Name.Lexeme)
	if !defined {
		return nil, &SemanticsError{Reason: fmt.Sprintf("undefined variable: %s", expr.Name.Lexeme)}
	}
	r.intepreter.Resolve(expr, dist)
	return nil, nil
}

func (r *Resolver) VisitCallExpr(expr *CallExpr) (interface{}, error) {
	if _, err := expr.Callee.Accept(r); err != nil {
		return nil, err
	}
	for _, argv := range expr.Arguments {
		if _, err := argv.Accept(r); err != nil {
			return nil, err
		}
	}
	return nil, nil
}

func (r *Resolver) VisitGetPropertyExpr(expr *GetPropertyExpr) (interface{}, error) {
	return expr.Object.Accept(r)
}

func (r *Resolver) VisitSetPropertyExpr(expr *SetPropertyExpr) (interface{}, error) {
	if _, err := expr.Object.Accept(r); err != nil {
		return nil, err
	}
	if _, err := expr.Value.Accept(r); err != nil {
		return nil, err
	}
	return nil, nil
}

func (r *Resolver) VisitVariableExpr(expr *VariableExpr) (interface{}, error) {
	if defined, declared := r.scopes[len(r.scopes)-1][expr.Name.Lexeme]; declared && !defined {
		return nil, &SemanticsError{fmt.Sprintf("variable referencing itself in its own initializer: %s", expr.Name.Lexeme)}
	}

	// Start from the innermost till the global scope, look for a matching name
	dist, defined := r.searchScopes(expr.Name.Lexeme)
	if !defined {
		return nil, &SemanticsError{fmt.Sprintf("undefined variable: %s", expr.Name.Lexeme)}
	}
	r.intepreter.Resolve(expr, dist)
	return nil, nil
}

func (r *Resolver) VisitThisExpr(expr *ThisExpr) (interface{}, error) {
	dist, defined := r.searchScopes("this")
	if !defined {
		return nil, &SemanticsError{Reason: "ed \"this\""}
	}
	r.intepreter.Resolve(expr, dist)
	return nil, nil
}

func (r *Resolver) VisitSuperExpr(expr *SuperExpr) (interface{}, error) {
	// Resolve `super` as if it were a variable
	dist, defined := r.searchScopes("super")
	if !defined {
		return nil, &SemanticsError{"unbounded \"super\""}
	}
	// Check if this class has a super class
	if r.enclosingClass.SuperClass == nil {
		return nil, &SemanticsError{"calling super on a class that doesn't have a super class"}
	}
	r.intepreter.Resolve(expr, dist)
	return nil, nil
}
