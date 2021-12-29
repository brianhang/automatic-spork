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

type TokenHolder interface {
	GetToken() Token
}

type Token struct {
	id     TokenID
	line   int
	column int
}

func (t Token) GetToken() Token {
	return t
}

type StringToken struct {
	Token
	value string
}

func (t StringToken) GetToken() Token {
	return t.Token
}

type NumberToken struct {
	Token
	value float64
}

func (t NumberToken) GetToken() Token {
	return t.Token
}

type IdentifierToken struct {
	Token
	value string
}

func (t IdentifierToken) GetToken() Token {
	return t.Token
}
