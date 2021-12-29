package tokenize

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

type Token struct {
	id     TokenID
	line   int
	column int
}
