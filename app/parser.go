package main

import (
	"errors"
	"fmt"
)

/*
Use right-associative notations:

	program        → statement* EOF
	statement      → block | exprStmt | printStmt | varDeclStmt | ifStmt
	block 		   → "{" statement* "}"

	exprStmt       → expression ";"
	printStmt      → "print" expression ";"
	varDeclStmt    → "var" IDENTIFIER ("=" EXPRESSION)? ";"
	ifStmt		   → "if" "(" expression ")" statement ( "else" statement )?
	whileStmt      → "while" "(" expression ")" statement

	expression     → assignment

	assignment     → lvalue "=" assignment) | logic_or
	logic_or	   → logic_and ( "or" logic_and )*
	logic_and      → equality ( "and" equality )*

	equality       → comparison ( ( "!=" | "==" ) comparison )*
	comparison     → term ( ( ">" | ">=" | "<" | "<=" ) term )*
	term           → factor ( ( "-" | "+" ) factor )*
	factor         → unary ( ( "/" | "*" ) unary )*
	unary          → ( "!" | "-" ) unary | primary
	primary        → NUMBER | STRING | "true" | "false" | "nil" | "(" expression ")" | IDENTIFIER
*/

type Parser interface {
	Parse(tokens []*Token) []Stmt
}

type RDParser struct {
	tokens  []*Token
	currIdx int
}

var (
	errorStmtMissingSemiColon = errors.New("statement missing semicolon")
)

func (p *RDParser) Parse(tokens []*Token) ([]Stmt, error) {
	var stmts []Stmt
	var err error

	p.tokens = tokens
	p.currIdx = 0

	for p.currIdx < len(p.tokens) {
		var stmt Stmt

		if p.match(EOF) {
			// TODO: what's the use for EOF?
			break
		}
		if stmt, err = p.statement(); err != nil {
			return nil, err
		}
		stmts = append(stmts, stmt)
	}
	return stmts, nil
}

func (p *RDParser) peek() *Token {
	if p.currIdx >= len(p.tokens) {
		return nil
	}
	return p.tokens[p.currIdx]
}

func (p *RDParser) previous() *Token {
	if p.currIdx == 0 {
		return nil
	}
	return p.tokens[p.currIdx-1]
}

func (p *RDParser) match(tokenType int) bool {
	token := p.peek()
	return token != nil && token.Type == tokenType
}

func (p *RDParser) advance() {
	p.currIdx += 1
}

func (p *RDParser) hasNext() bool {
	return p.currIdx < len(p.tokens)
}

// advanceIfMatch checks if current token matches any of `tokenTypes`.
// If yes, return the current token, and advance to the next token
func (p *RDParser) advanceIfMatch(tokenTypes ...int) bool {
	for _, t := range tokenTypes {
		if p.match(t) {
			p.advance()
			return true
		}
	}
	return false
}

func (p *RDParser) statement() (Stmt, error) {
	if p.advanceIfMatch(Print) {
		return p.printStatement()
	}
	if p.advanceIfMatch(Var) {
		return p.varDeclStatement()
	}
	if p.advanceIfMatch(LeftBrace) {
		return p.blockStatement()
	}
	if p.advanceIfMatch(If) {
		return p.ifStatement()
	}
	if p.advanceIfMatch(While) {
		return p.whileStatement()
	}
	return p.expressionStatement()
}

func (p *RDParser) printStatement() (Stmt, error) {
	var expr Expr
	var err error

	if expr, err = p.expression(); err != nil {
		return nil, err
	}
	if !p.advanceIfMatch(SemiColon) {
		return nil, errorStmtMissingSemiColon
	}
	return &PrintStmt{Child: expr}, nil
}

func (p *RDParser) expressionStatement() (Stmt, error) {
	var expr Expr
	var err error

	if expr, err = p.expression(); err != nil {
		return nil, err
	}
	if !p.advanceIfMatch(SemiColon) {
		return nil, errorStmtMissingSemiColon
	}
	return &InlineExprStmt{Child: expr}, nil
}

func (p *RDParser) varDeclStatement() (Stmt, error) {
	var initializer Expr
	var name *Token
	var err error

	if !p.advanceIfMatch(Identifier) {
		return nil, errors.New("variable declaration missing identifier")
	}
	name = p.previous()

	if p.advanceIfMatch(Equal) {
		if initializer, err = p.expression(); err != nil {
			return nil, err
		}
	}
	if !p.advanceIfMatch(SemiColon) {
		return nil, errorStmtMissingSemiColon
	}
	return &VarDeclStmt{Name: name, Initializer: initializer}, nil
}

func (p *RDParser) blockStatement() (Stmt, error) {
	var stmts []Stmt
	for p.hasNext() && !p.match(RightBrace) {
		var stmt Stmt
		var err error

		if stmt, err = p.statement(); err != nil {
			return nil, err
		}
		stmts = append(stmts, stmt)
	}
	if !p.advanceIfMatch(RightBrace) {
		return nil, errors.New("block missing \"}\"")
	}
	return &BlockStmt{Stmts: stmts}, nil
}

func (p *RDParser) ifStatement() (Stmt, error) {
	var condition Expr
	var thenBranch Stmt
	var elseBranch Stmt
	var err error

	if !p.advanceIfMatch(LeftParen) {
		return nil, errors.New("if statement missing left parenthesis")
	}
	if condition, err = p.expression(); err != nil {
		return nil, err
	}
	if !p.advanceIfMatch(RightParen) {
		return nil, errors.New("if statement missing right parenthesis")
	}
	if thenBranch, err = p.statement(); err != nil {
		return nil, err
	}

	// In case of the "dangling else" problem (i.e. if A if B else C), the "else" statement
	// is bounded to the nearest "if" statement
	if p.advanceIfMatch(Else) {
		if elseBranch, err = p.statement(); err != nil {
			return nil, err
		}
	}
	return &IfStmt{Condition: condition, ThenBranch: thenBranch, ElseBranch: elseBranch}, nil
}

func (p *RDParser) whileStatement() (Stmt, error) {
	var condition Expr
	var body Stmt
	var err error
	if !p.advanceIfMatch(LeftParen) {
		return nil, errors.New("while statement missing left parenthesis")
	}
	if condition, err = p.expression(); err != nil {
		return nil, err
	}
	if !p.advanceIfMatch(RightParen) {
		return nil, errors.New("while statement missing right parenthesis")
	}
	if body, err = p.statement(); err != nil {
		return nil, err
	}
	return &WhileStmt{Condition: condition, Body: body}, nil
}

func (p *RDParser) expression() (Expr, error) {
	return p.assignment()
}

func (p *RDParser) assignment() (Expr, error) {
	var name Expr
	var value Expr
	var err error

	// If an "assignment" rule is satisfied by an assignment expression,
	// The LHS of this expr is a l-value expr that evaluates to the storage location.
	// The RHS of this expr is a r-value expr that evaluates to a value.
	// A l-value expr happens to satisfy "logicOr" rule, so we can use "logicOr" rule to parse it
	// and filter out the well-defined variants.
	if name, err = p.logicOr(); err != nil {
		return nil, err
	}

	// This is an "logicOr" expr
	if !p.advanceIfMatch(Equal) {
		return name, nil
	}

	// This expr should start with an identifier followed by "logicOr" expr
	if value, err = p.assignment(); err != nil {
		return nil, err
	}

	// TODO: support more l-value expr
	varName, ok := name.(*VariableExpr)
	if !ok {
		return nil, fmt.Errorf("invalid assignment target")
	}
	return &AssignExpr{Name: varName.Name, Value: value}, nil
}

func (p *RDParser) logicOr() (Expr, error) {
	var left Expr
	var right Expr
	var err error

	if left, err = p.logicAnd(); err != nil {
		return nil, err
	}
	for p.advanceIfMatch(Or) {
		op := p.previous()
		if right, err = p.logicAnd(); err != nil {
			return nil, err
		}
		left = &LogicExpr{Operator: op, Left: left, Right: right}
	}
	return left, nil
}

func (p *RDParser) logicAnd() (Expr, error) {
	var left Expr
	var right Expr
	var err error

	if left, err = p.equality(); err != nil {
		return nil, err
	}
	for p.advanceIfMatch(And) {
		op := p.previous()
		if right, err = p.equality(); err != nil {
			return nil, err
		}
		left = &LogicExpr{Operator: op, Left: left, Right: right}
	}
	return left, nil
}

func (p *RDParser) equality() (Expr, error) {
	var left Expr
	var right Expr
	var err error

	if left, err = p.comparison(); err != nil {
		return nil, err
	}
	for p.advanceIfMatch(BangEqual, EqualEqual) {
		op := p.previous()
		if right, err = p.comparison(); err != nil {
			return nil, err
		}
		left = &BinaryExpr{Operator: op, Left: left, Right: right}
	}
	return left, nil
}

func (p *RDParser) comparison() (Expr, error) {
	var left Expr
	var right Expr
	var err error

	if left, err = p.term(); err != nil {
		return nil, err
	}
	for p.advanceIfMatch(Greater, GreaterEqual, Less, LessEqual) {
		op := p.previous()
		if right, err = p.term(); err != nil {
			return nil, err
		}
		left = &BinaryExpr{Operator: op, Left: left, Right: right}
	}
	return left, nil
}

func (p *RDParser) term() (Expr, error) {
	var left Expr
	var right Expr
	var err error

	if left, err = p.factor(); err != nil {
		return nil, err
	}
	for p.advanceIfMatch(Plus, Minus) {
		op := p.previous()
		if right, err = p.factor(); err != nil {
			return nil, err
		}
		left = &BinaryExpr{Operator: op, Left: left, Right: right}
	}
	return left, nil
}

func (p *RDParser) factor() (Expr, error) {
	var left Expr
	var right Expr
	var err error

	if left, err = p.unary(); err != nil {
		return nil, err
	}
	for p.advanceIfMatch(Slash, Star) {
		op := p.previous()
		if right, err = p.unary(); err != nil {
			return nil, err
		}
		left = &BinaryExpr{Operator: op, Left: left, Right: right}
	}
	return left, nil
}

func (p *RDParser) unary() (Expr, error) {
	var right Expr
	var err error

	if p.advanceIfMatch(Bang, Minus) {
		op := p.previous()
		if right, err = p.unary(); err != nil {
			return nil, err
		}
		return &UnaryExpr{Operator: op, Right: right}, nil
	}
	return p.primary()
}

func (p *RDParser) primary() (Expr, error) {
	var expr Expr
	var err error

	if p.advanceIfMatch(False) {
		return &LiteralExpr{Value: false}, nil
	}
	if p.advanceIfMatch(True) {
		return &LiteralExpr{Value: true}, nil
	}
	if p.advanceIfMatch(Nil) {
		return &LiteralExpr{Value: nil}, nil
	}
	if p.advanceIfMatch(String, Number) {
		return &LiteralExpr{Value: p.previous().Literal}, nil
	}
	if p.advanceIfMatch(LeftParen) {
		if expr, err = p.expression(); err != nil {
			return nil, err
		}
		if !p.match(RightParen) {
			return nil, errors.New("missing right parenthesis")
		}
		p.advance()
		return &GroupingExpr{Child: expr}, nil
	}
	if p.advanceIfMatch(Identifier) {
		return &VariableExpr{Name: p.previous()}, nil
	}
	return nil, fmt.Errorf("expect a valid expression: %s", p.peek().Lexeme)
}
