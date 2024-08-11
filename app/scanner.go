package main

import (
	"errors"
	"fmt"
	"strconv"
	"unicode"
)

const (
	// Non-alpha symbols
	LeftParen = iota
	RightParen
	LeftBrace
	RightBrace
	Comma
	Dot
	Minus
	Plus
	SemiColon
	Slash
	Star
	Bang
	BangEqual
	Equal
	EqualEqual
	Greater
	GreaterEqual
	Less
	LessEqual

	// Literals
	Identifier
	String
	Number

	// Reserved words
	And
	Class
	Else
	False
	Fun
	For
	If
	Nil
	Or
	Print
	Return
	Super
	This
	True
	Var
	While

	EOF
)

var reservedWords = map[string]int{
	"and":    And,
	"class":  Class,
	"else":   Else,
	"false":  False,
	"fun":    Fun,
	"for":    For,
	"if":     If,
	"nil":    Nil,
	"or":     Or,
	"print":  Print,
	"return": Return,
	"super":  Super,
	"this":   This,
	"true":   True,
	"var":    Var,
	"while":  While,
}

type Token struct {
	Type    int
	Lexeme  string
	Literal interface{}
	LineNo  int
}

func (t Token) ToString() string {
	return fmt.Sprintf("%d %s", t.Type, t.Lexeme)
}

type Scanner interface {
	Scan(src string) error
}

type ScannerImpl struct {
	tokens []*Token

	source   string
	startIdx int
	currIdx  int
	lineNo   int
}

func (s *ScannerImpl) emit(tokenType int, literal interface{}) {
	s.tokens = append(s.tokens, &Token{
		Type:    tokenType,
		Lexeme:  s.source[s.startIdx:s.currIdx],
		Literal: literal,
		LineNo:  s.lineNo,
	})
}

func (s *ScannerImpl) emitEOF() {
	s.tokens = append(s.tokens, &Token{
		Type: EOF,
	})
}

func (s *ScannerImpl) hasNext() bool {
	return s.currIdx < len(s.source)
}

func (s *ScannerImpl) advance() byte {
	c := s.source[s.currIdx]
	s.currIdx += 1
	return c
}

func (s *ScannerImpl) advanceIfMatch(expected byte) bool {
	if !s.hasNext() {
		return false
	}
	if s.source[s.currIdx] != expected {
		return false
	}
	s.currIdx++
	return true
}

func (s *ScannerImpl) peek() byte {
	if !s.hasNext() {
		return 0
	}
	return s.source[s.currIdx]
}

func (s *ScannerImpl) peekAhead(offset int) byte {
	if s.currIdx+offset >= len(s.source) {
		return 0
	}
	return s.source[s.currIdx+offset]
}

func (s *ScannerImpl) emitString() error {
	for s.hasNext() && s.peek() != '"' {
		s.advance()
	}
	if !s.hasNext() {
		return errors.New("string is unterminated")
	}
	// Handle closing "
	s.advance()
	s.emit(String, string(s.source[s.startIdx+1:s.currIdx-1]))
	return nil
}

func (s *ScannerImpl) emitNumber() {
	/*
		All numbers in Lox are floating point at runtime
		Supported formats: 1234, 1234.56
	*/
	for s.hasNext() && unicode.IsDigit(rune(s.peek())) {
		s.advance()
	}

	if s.peek() == '.' {
		// Make sure there's a digit after the decimal points
		if unicode.IsDigit(rune(s.peekAhead(1))) {
			// Consume decimal point
			s.advance()
			for s.hasNext() && unicode.IsDigit(rune(s.peek())) {
				s.advance()
			}
		}
	}

	num, _ := strconv.ParseFloat(s.source[s.startIdx:s.currIdx], 64)
	s.emit(Number, num)
}

func (s *ScannerImpl) emitIdentifier() {
	rangeTable := []*unicode.RangeTable{
		unicode.Letter,
		unicode.Number,
	}
	for unicode.IsOneOf(rangeTable, rune(s.peek())) || s.peek() == '_' {
		s.advance()
	}

	name := s.source[s.startIdx:s.currIdx]
	keywordType, found := reservedWords[name]
	if found {
		s.emit(keywordType, nil)
	} else {
		s.emit(Identifier, nil)
	}
}

func (s *ScannerImpl) scanToken() error {
	var err error
	c := s.advance()

	switch c {
	case '(':
		s.emit(LeftParen, nil)
	case ')':
		s.emit(RightParen, nil)
	case '{':
		s.emit(LeftBrace, nil)
	case '}':
		s.emit(RightBrace, nil)
	case ',':
		s.emit(Comma, nil)
	case '.':
		s.emit(Dot, nil)
	case '-':
		s.emit(Minus, nil)
	case '+':
		s.emit(Plus, nil)
	case ';':
		s.emit(SemiColon, nil)
	case '*':
		s.emit(Star, nil)
	case '!':
		if s.advanceIfMatch('=') {
			s.emit(BangEqual, nil)
		} else {
			s.emit(Bang, nil)
		}
	case '=':
		if s.advanceIfMatch('=') {
			s.emit(EqualEqual, nil)
		} else {
			s.emit(Equal, nil)
		}
	case '<':
		if s.advanceIfMatch('=') {
			s.emit(LessEqual, nil)
		} else {
			s.emit(Less, nil)
		}
	case '>':
		if s.advanceIfMatch('=') {
			s.emit(GreaterEqual, nil)
		} else {
			s.emit(Greater, nil)
		}
	case '/':
		if s.advanceIfMatch('/') {
			// Handle comment
			for s.hasNext() && s.peek() != '\n' {
				s.advance()
			}
		} else {
			s.emit(Slash, nil)
		}
	case '"':
		err = s.emitString()
	case ' ':
		fallthrough
	case '\t':
		fallthrough
	case '\r':
		// Do nothing
	case '\n':
		s.lineNo += 1
	default:
		if unicode.IsDigit(rune(c)) {
			s.emitNumber()
		} else if unicode.IsLetter(rune(c)) || c == '_' {
			s.emitIdentifier()
		} else {
			err = fmt.Errorf("unrecognized token: %c", c)
		}
	}
	return err
}

func (s *ScannerImpl) reset(source string) {
	s.tokens = []*Token{}
	s.source = source
	s.lineNo = 1
	s.currIdx = 0
}

func (s *ScannerImpl) Scan(source string) error {
	s.reset(source)
	for s.hasNext() {
		s.startIdx = s.currIdx
		if err := s.scanToken(); err != nil {
			return err
		}
	}
	s.emitEOF()
	return nil
}
