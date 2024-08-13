package main

import "errors"

/*
Use right-associative notations:
expression     → equality ;
equality       → comparison ( ( "!=" | "==" ) comparison )* ;
comparison     → term ( ( ">" | ">=" | "<" | "<=" ) term )* ;
term           → factor ( ( "-" | "+" ) factor )* ;
factor         → unary ( ( "/" | "*" ) unary )* ;
unary          → ( "!" | "-" ) unary | primary ;
primary        → NUMBER | STRING | "true" | "false" | "nil" | "(" expression ")" ;
*/
type Parser interface {
	Parse(tokens []*Token) Expr
}

type RDParser struct {
	tokens  []*Token
	currIdx int
}

func (p *RDParser) Parse(tokens []*Token) (Expr, error) {
	var expr Expr
	var err error

	p.tokens = tokens
	p.currIdx = 0
	if expr, err = p.expression(); err != nil {
		return nil, err
	}
	if p.currIdx < len(p.tokens) {
		return nil, errors.New("unused tokens")
	}
	return expr, nil
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

func (p *RDParser) expression() (Expr, error) {
	return p.equality()
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
	return nil, errors.New("expect a valid expression")
}
