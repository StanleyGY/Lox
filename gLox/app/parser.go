package main

import (
	"bytes"
)

/*
Use right-associative notations:

	program        → declaration* EOF

	declaration    → classDecl | funDecl | varDecl | statement
	classDecl      → "class" IDENTIFIER ( "<" IDENTIFIER )? "{" function* "}"
	varDecl    	   → "var" IDENTIFIER ( "=" EXPRESSION )? ";"
	funDecl        → "fun" function
	function       → IDENTIFIER "(" parameters? ")" block
	parameters     → IDENTIFIER ( "," IDENTIFIER )*

	statement      → block | exprStmt | printStmt | ifStmt | forStmt | whileStmt | returnStmt | breakStmt
	block 		   → "{" declaration* "}"
	exprStmt       → expression ";"
	printStmt      → "print" expression ";"
	ifStmt		   → "if" "(" expression ")" statement ( "else" statement )?
	forStmt        → "for" "(" ( varDecl | exprStmt | ";" ) | expression? ";" expression? ")" statement
	whileStmt      → "while" "(" expression ")" statement
	returnStmt     → "return" expression? ";"
	breakStmt      → "break" ";"

	expression     → assignment
	assignment     → ( call "." )? IDENTIFIER "=" assignment | logic_or
	logic_or	   → logic_and ( "or" logic_and )*
	logic_and      → equality ( "and" equality )*
	equality       → comparison (( "!=" | "==" ) comparison )*
	comparison     → term (( ">" | ">=" | "<" | "<=" ) term )*
	term           → factor (( "-" | "+" ) factor )*
	factor         → unary (( "/" | "*" ) unary )*
	unary          → (( "!" | "-" ) unary) | call
	call           → primary ( "(" arguments? ")" | "." IDENTIFIER )*
	arguments      → expression ( "," expression )*
	primary        → NUMBER | STRING | "true" | "false" | "nil" | "(" expression ")" | IDENTIFIER | "super" "." IDENTIFIER
*/

const (
	MaxNumFunCallArguments = 255
)

type Parser interface {
	Parse(tokens []*Token) []Stmt
}

type RDParser struct {
	tokens  []*Token
	currIdx int
}

type ParsingError struct {
	Reason   string
	TokenIdx int
	Tokens   []*Token
}

func (e ParsingError) Error() string {
	var buf bytes.Buffer
	buf.WriteString("Parsing failed around lines: \n")
	for i := max(0, e.TokenIdx-10); i <= min(e.TokenIdx+5, len(e.Tokens)-1); i++ {
		buf.WriteString(e.Tokens[i].Lexeme)
		buf.WriteString(" ")
	}
	buf.WriteString("\n")
	for i := max(0, e.TokenIdx-10); i <= min(e.TokenIdx+5, len(e.Tokens)-1); i++ {
		for j := 0; j < len(e.Tokens[i].Lexeme); j++ {
			if i == e.TokenIdx && j <= len(e.Tokens[i].Lexeme)/2 {
				buf.WriteString("^")
			} else {
				buf.WriteString(" ")
			}
		}
		buf.WriteString(" ")
	}
	buf.WriteString("\nFailed reason: ")
	buf.WriteString(e.Reason)
	return buf.String()
}

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
		if stmt, err = p.declaration(); err != nil {
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

func (p *RDParser) emitParsingError(reason string) error {
	return &ParsingError{Reason: reason, TokenIdx: p.currIdx, Tokens: p.tokens}
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

func (p *RDParser) declaration() (Stmt, error) {
	if p.match(Var) {
		return p.varDecl()
	}
	if p.match(Fun) {
		return p.funDecl()
	}
	if p.match(Class) {
		return p.classDecl()
	}
	return p.statement()
}

func (p *RDParser) varDecl() (Stmt, error) {
	var initializer Expr
	var name *Token
	var err error

	if !p.advanceIfMatch(Var) {
		return nil, p.emitParsingError("variable declaration missing \"var\" keyword")
	}
	if !p.advanceIfMatch(Identifier) {
		return nil, p.emitParsingError("variable declaration missing identifier")
	}
	name = p.previous()

	if p.advanceIfMatch(Equal) {
		if initializer, err = p.expression(); err != nil {
			return nil, err
		}
	}
	if !p.advanceIfMatch(SemiColon) {
		return nil, p.emitParsingError("missing \";\"")
	}
	return &VarDeclStmt{Name: name, Initializer: initializer}, nil
}

func (p *RDParser) funDecl() (Stmt, error) {
	if !p.advanceIfMatch(Fun) {
		return nil, p.emitParsingError("func declaration missing \"fun\" keyword")
	}
	return p.function()
}

func (p *RDParser) function() (*FuncDeclStmt, error) {
	var name *Token
	var parameters []*Token
	var body Stmt
	var err error

	// Match function signatures
	if !p.advanceIfMatch(Identifier) {
		return nil, p.emitParsingError("func declaration missing name")
	}
	name = p.previous()

	if !p.advanceIfMatch(LeftParen) {
		return nil, p.emitParsingError("func declaration missing \"(\"")
	}
	if parameters, err = p.parameters(); err != nil {
		return nil, err
	}
	if len(parameters) > MaxNumFunCallArguments {
		return nil, p.emitParsingError("func declaration argument list too long")
	}
	if !p.advanceIfMatch(RightParen) {
		return nil, p.emitParsingError("func declaration missing \")\"")
	}

	// Match function implementation
	if body, err = p.blockStmt(); err != nil {
		return nil, err
	}
	return &FuncDeclStmt{Name: name, Params: parameters, Body: body}, nil
}

func (p *RDParser) classDecl() (Stmt, error) {
	var name *Token
	var superClass *VariableExpr
	var methods []*FuncDeclStmt
	var err error

	if !p.advanceIfMatch(Class) {
		return nil, p.emitParsingError("class declaration missing \"class\" keyword")
	}
	if !p.advanceIfMatch(Identifier) {
		return nil, p.emitParsingError("class declaration missing name")
	}
	name = p.previous()

	if p.advanceIfMatch(Less) {
		if !p.advanceIfMatch(Identifier) {
			return nil, p.emitParsingError("class declaration missing super class")
		}
		// Wrap this additionally in an expr so semantic analysis
		// can be done on this identifier
		superClass = &VariableExpr{p.previous()}
	}
	if !p.advanceIfMatch(LeftBrace) {
		return nil, p.emitParsingError("class declaration missing \"{\"")
	}
	for !p.advanceIfMatch(RightBrace) {
		var m *FuncDeclStmt
		if m, err = p.function(); err != nil {
			return nil, err
		}
		methods = append(methods, m)
	}
	return &ClassDeclStmt{Name: name, SuperClass: superClass, Methods: methods}, nil
}

func (p *RDParser) parameters() ([]*Token, error) {
	var params []*Token

	// Function with no parameters
	if !p.match(Identifier) {
		return params, nil
	}
	for {
		if !p.advanceIfMatch(Identifier) {
			return nil, p.emitParsingError("missing func parameter after \",\"")
		}
		params = append(params, p.previous())
		if !p.advanceIfMatch(Comma) {
			return params, nil
		}
	}
}

func (p *RDParser) statement() (Stmt, error) {
	if p.match(Print) {
		return p.printStmt()
	}
	if p.match(LeftBrace) {
		return p.blockStmt()
	}
	if p.match(If) {
		return p.ifStmt()
	}
	if p.match(While) {
		return p.whileStmt()
	}
	if p.match(For) {
		return p.forStmt()
	}
	if p.match(Return) {
		return p.returnStmt()
	}
	if p.match(Break) {
		return p.breakStmt()
	}
	return p.expressionStmt()
}

func (p *RDParser) printStmt() (Stmt, error) {
	var expr Expr
	var err error

	if !p.advanceIfMatch(Print) {
		return nil, p.emitParsingError("missing \"print\" keyword")
	}
	if expr, err = p.expression(); err != nil {
		return nil, err
	}
	if !p.advanceIfMatch(SemiColon) {
		return nil, p.emitParsingError("missing \";\"")
	}
	return &PrintStmt{Child: expr}, nil
}

func (p *RDParser) expressionStmt() (Stmt, error) {
	var expr Expr
	var err error

	if expr, err = p.expression(); err != nil {
		return nil, err
	}
	if !p.advanceIfMatch(SemiColon) {
		return nil, p.emitParsingError("missing \";\"")
	}
	return &InlineExprStmt{Child: expr}, nil
}

func (p *RDParser) blockStmt() (Stmt, error) {
	var stmts []Stmt

	if !p.advanceIfMatch(LeftBrace) {
		return nil, p.emitParsingError("missing \"{\"")
	}
	for p.hasNext() && !p.match(RightBrace) {
		var stmt Stmt
		var err error

		if stmt, err = p.declaration(); err != nil {
			return nil, err
		}
		stmts = append(stmts, stmt)
	}
	if !p.advanceIfMatch(RightBrace) {
		return nil, p.emitParsingError("missing \"}\"")
	}
	return &BlockStmt{Stmts: stmts}, nil
}

func (p *RDParser) ifStmt() (Stmt, error) {
	var condition Expr
	var thenBranch Stmt
	var elseBranch Stmt
	var err error

	if !p.advanceIfMatch(If) {
		return nil, p.emitParsingError("missing \"if\" keyword")
	}
	if !p.advanceIfMatch(LeftParen) {
		return nil, p.emitParsingError("if condition missing \"{\"")
	}
	if condition, err = p.expression(); err != nil {
		return nil, err
	}
	if !p.advanceIfMatch(RightParen) {
		return nil, p.emitParsingError("if condition missing \"}\"")
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

func (p *RDParser) whileStmt() (Stmt, error) {
	var condition Expr
	var body Stmt
	var err error
	if !p.advanceIfMatch(While) {
		return nil, p.emitParsingError("missing \"while\" keyword")
	}
	if !p.advanceIfMatch(LeftParen) {
		return nil, p.emitParsingError("while loop missing \"{\"")
	}
	if condition, err = p.expression(); err != nil {
		return nil, err
	}
	if !p.advanceIfMatch(RightParen) {
		return nil, p.emitParsingError("while loop missing \"}\"")
	}
	if body, err = p.statement(); err != nil {
		return nil, err
	}
	return &WhileStmt{Condition: condition, Body: body}, nil
}

func (p *RDParser) forStmt() (Stmt, error) {
	var initializer Stmt
	var condition Expr
	var increment Expr
	var body Stmt
	var err error

	if !p.advanceIfMatch(For) {
		return nil, p.emitParsingError("missing \"for\" keyword")
	}
	if !p.advanceIfMatch(LeftParen) {
		return nil, p.emitParsingError("for loop missing \"(\"")
	}
	// Initializer clause
	if p.match(Var) {
		if initializer, err = p.varDecl(); err != nil {
			return nil, err
		}
	} else if !p.advanceIfMatch(SemiColon) {
		if initializer, err = p.expressionStmt(); err != nil {
			return nil, err
		}
	}
	// Condition clause
	if !p.advanceIfMatch(SemiColon) {
		if condition, err = p.expression(); err != nil {
			return nil, err
		}
		if !p.advanceIfMatch(SemiColon) {
			return nil, p.emitParsingError("for loop condition clause missing \"'\"")
		}
	}
	// Increment clause
	if !p.advanceIfMatch(RightParen) {
		if increment, err = p.expression(); err != nil {
			return nil, err
		}
		if !p.advanceIfMatch(RightParen) {
			return nil, p.emitParsingError("for loop condition clause missing \")\"")
		}
	}
	// Body
	if body, err = p.statement(); err != nil {
		return nil, err
	}

	// C-style for-loop is just a syntactic sugar of a while-loop
	// If condition is omitted, then it's default to true
	if condition == nil {
		condition = &LiteralExpr{Value: true}
	}
	if increment != nil {
		body = &BlockStmt{Stmts: []Stmt{body, &InlineExprStmt{Child: increment}}}
	}

	whileStmt := &WhileStmt{Condition: condition, Body: body}

	if initializer != nil {
		return &BlockStmt{Stmts: []Stmt{initializer, whileStmt}}, nil
	}
	return whileStmt, nil
}

func (p *RDParser) returnStmt() (Stmt, error) {
	var expr Expr
	var err error

	if !p.advanceIfMatch(Return) {
		return nil, p.emitParsingError("missing \"return\" keyword")
	}
	if !p.advanceIfMatch(SemiColon) {
		if expr, err = p.expression(); err != nil {
			return nil, err
		}
		if !p.advanceIfMatch(SemiColon) {
			return nil, p.emitParsingError("return statement missing \";\"")
		}
		return &ReturnStmt{Value: expr}, nil
	}
	return &ReturnStmt{}, nil
}

func (p *RDParser) breakStmt() (Stmt, error) {
	if !p.advanceIfMatch(Break) {
		return nil, p.emitParsingError("missing \"break\" keyword")
	}
	if !p.advanceIfMatch(SemiColon) {
		return nil, p.emitParsingError("break statement missing \";\"")
	}
	return &BreakStmt{}, nil
}

func (p *RDParser) expression() (Expr, error) {
	return p.assignment()
}

func (p *RDParser) assignment() (Expr, error) {
	var left Expr
	var err error

	// If an "assignment" rule is satisfied by an assignment expression,
	// The LHS of this expr is a l-value expr that evaluates to the storage location.
	// The RHS of this expr is a r-value expr that evaluates to a value.
	// A l-value expr happens to satisfy "logicOr" rule, so we can use "logicOr" rule to parse it
	// and filter out the well-defined variants.
	if left, err = p.logicOr(); err != nil {
		return nil, err
	}

	// Actually does assignment
	if p.advanceIfMatch(Equal) {
		var value Expr

		// This expr should start with an identifier followed by "logicOr" expr
		if value, err = p.assignment(); err != nil {
			return nil, err
		}

		switch left := left.(type) {
		case *VariableExpr:
			return &AssignExpr{Name: left.Name, Value: value}, nil
		case *GetPropertyExpr:
			return &SetPropertyExpr{
				Object:   left.Object,
				Property: left.Property,
				Value:    value,
			}, nil
		default:
			return nil, p.emitParsingError("invalid assignment target")
		}
	}
	return left, nil
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
	return p.call()
}

func (p *RDParser) call() (Expr, error) {
	var expr Expr
	var err error

	if expr, err = p.primary(); err != nil {
		return nil, err
	}

	for {
		if p.advanceIfMatch(Dot) {
			// Handle class property-get call
			if !p.advanceIfMatch(Identifier) {
				return nil, p.emitParsingError("missing identifier for property access")
			}
			expr = &GetPropertyExpr{Object: expr, Property: p.previous()}

		} else if p.advanceIfMatch(LeftParen) {
			// Handle regular function call
			if p.advanceIfMatch(RightParen) {
				expr = &CallExpr{Callee: expr}
			} else {
				var arguments []Expr
				if arguments, err = p.arguments(); err != nil {
					return nil, err
				}
				if len(arguments) >= MaxNumFunCallArguments {
					return nil, p.emitParsingError("func call argument list too long")
				}
				if !p.advanceIfMatch(RightParen) {
					return nil, p.emitParsingError("func call argument list missing \")\"")
				}
				expr = &CallExpr{Callee: expr, Arguments: arguments}
			}
		} else {
			break
		}
	}
	return expr, nil
}

func (p *RDParser) arguments() ([]Expr, error) {
	var exprs []Expr
	var expr Expr
	var err error

	for {
		if expr, err = p.expression(); err != nil {
			return nil, err
		}
		exprs = append(exprs, expr)
		if !p.advanceIfMatch(Comma) {
			break
		}
	}
	return exprs, nil
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
			return nil, p.emitParsingError("grouping expr missing \")\"")
		}
		p.advance()
		return &GroupingExpr{Child: expr}, nil
	}
	if p.advanceIfMatch(Identifier) {
		return &VariableExpr{Name: p.previous()}, nil
	}
	if p.advanceIfMatch(This) {
		return &ThisExpr{}, nil
	}
	if p.advanceIfMatch(Super) {
		if !p.advanceIfMatch(Dot) {
			return nil, p.emitParsingError("super missing a \".\"")
		}
		if !p.advanceIfMatch(Identifier) {
			return nil, p.emitParsingError("super missing an identifier")
		}
		return &SuperExpr{Property: p.previous()}, nil
	}
	return nil, p.emitParsingError("expect a valid primary expr")
}
