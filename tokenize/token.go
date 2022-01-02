package tokenize

import (
	"fmt"
	"strconv"
)

type TokenID int

const (
	TokenEOF TokenID = iota

	TokenLeftParen
	TokenRightParen
	TokenLeftCurly
	TokenRightCurly

	TokenComma
	TokenDot
	TokenMinus
	TokenPlus
	TokenStar
	TokenSlash
	TokenSemicolon

	TokenBang
	TokenBangEqual
	TokenEqual
	TokenEqualEqual
	TokenGreater
	TokenGreaterEqual
	TokenLess
	TokenLessEqual

	TokenIdentifier
	TokenString
	TokenNumber
	TokenTrue
	TokenFalse
	TokenNil

	TokenAnd
	TokenOr
	TokenIf
	TokenElse
	TokenFor
	TokenWhile

	TokenFunc
	TokenReturn

	TokenClass
	TokenSuper
	TokenThis
)

var tokenToString = map[TokenID]string{
	TokenEOF:        "EOF",
	TokenLeftParen:  "(",
	TokenRightParen: ")",
	TokenLeftCurly:  "{",
	TokenRightCurly: "}",

	TokenComma:     ",",
	TokenDot:       ".",
	TokenMinus:     "-",
	TokenPlus:      "+",
	TokenStar:      "*",
	TokenSlash:     "/",
	TokenSemicolon: ";",

	TokenBang:         "!",
	TokenBangEqual:    "!=",
	TokenEqual:        "=",
	TokenEqualEqual:   "==",
	TokenGreater:      ">",
	TokenGreaterEqual: ">=",
	TokenLess:         "<",
	TokenLessEqual:    "<=",

	TokenIdentifier: "identifier",
	TokenString:     "string",
	TokenNumber:     "number",
	TokenTrue:       "true",
	TokenFalse:      "false",
	TokenNil:        "nil",

	TokenAnd:   "and",
	TokenOr:    "or",
	TokenIf:    "if",
	TokenElse:  "else",
	TokenFor:   "for",
	TokenWhile: "while",

	TokenFunc:   "func",
	TokenReturn: "return",

	TokenClass: "class",
	TokenSuper: "super",
	TokenThis:  "this",
}

func (id TokenID) String() string {
	return tokenToString[id]
}

type TokenHolder interface {
	GetToken() Token
	GetID() TokenID
	GetLine() int
	GetColumn() int
	String() string
}

type Token struct {
	TokenHolder
	id     TokenID
	line   int
	column int
}

func (t Token) GetToken() Token {
	return t
}

func (t Token) GetID() TokenID {
	return t.id
}
func (t Token) GetLine() int {
	return t.line
}
func (t Token) GetColumn() int {
	return t.column
}
func (t Token) String() string {
	return tokenToString[t.id]
}

type StringToken struct {
	Token
	value string
}

func (t StringToken) GetToken() Token {
	return t.Token
}
func (t StringToken) GetID() TokenID {
	return t.id
}
func (t StringToken) GetLine() int {
	return t.line
}
func (t StringToken) GetColumn() int {
	return t.column
}
func (t StringToken) String() string {
	return fmt.Sprintf("\"%s\"", t.value)
}

type NumberToken struct {
	Token
	value float64
}

func (t NumberToken) GetToken() Token {
	return t.Token
}
func (t NumberToken) GetID() TokenID {
	return t.id
}
func (t NumberToken) GetLine() int {
	return t.line
}
func (t NumberToken) GetColumn() int {
	return t.column
}
func (t NumberToken) String() string {
	return strconv.FormatFloat(t.value, 'f', -1, 64)
}

type IdentifierToken struct {
	Token
	value string
}

func (t IdentifierToken) GetToken() Token {
	return t.Token
}
func (t IdentifierToken) GetID() TokenID {
	return t.id
}
func (t IdentifierToken) GetLine() int {
	return t.line
}
func (t IdentifierToken) GetColumn() int {
	return t.column
}
func (t IdentifierToken) GetValue() string {
	return t.value
}
func (t IdentifierToken) String() string {
	return t.value
}
