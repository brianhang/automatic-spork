package tokenize

import (
	"bufio"
	"io"
	"unicode"
)

const BufferLength = 4096

var singleRuneTokenType = map[rune]TokenID{
	'(': TokenLeftParen,
	')': TokenRightParen,
	'{': TokenLeftCurly,
	'}': TokenRightCurly,
	',': TokenComma,
	'.': TokenDot,
	'-': TokenMinus,
	'+': TokenPlus,
	'*': TokenStar,
	';': TokenSemicolon,
}

type Tokenizer struct {
	input  *bufio.Reader
	line   int
	column int
}

func NewTokenizer(input io.Reader) *Tokenizer {
	t := &Tokenizer{input: bufio.NewReader(input)}
	t.line = 1
	return t
}

func (t *Tokenizer) Tokenize() ([]Token, error) {
	tokens := make([]Token, 0)
	for {
		r, _, err := t.input.ReadRune()
		if err == io.EOF {
			break
		}
		if err != nil {
			return tokens, err
		}

		t.column += 1

		if tokenID, ok := singleRuneTokenType[r]; ok {
			tokens = append(tokens, t.token(tokenID))
			continue
		}
		if unicode.IsSpace(r) {
			if r == '\n' {
				t.column = 0
				t.line++
			}
			continue
		}
		var token Token
		switch r {
		case ',':
			token = t.token(TokenComma)
		case '.':
			token = t.token(TokenDot)
		case '-':
			token = t.token(TokenMinus)
		case '+':
			token = t.token(TokenPlus)
		case '*':
			token = t.token(TokenStar)
		case '/':
			if t.consumeIfNext('/') {
				t.consumeUntilEOL()
				continue
			}
			token = t.token(TokenSlash)
		case '!':
			if t.consumeIfNext('=') {
				token = t.token(TokenBangEqual)
			} else {
				token = t.token(TokenBang)
			}
		case '=':
			if t.consumeIfNext('=') {
				token = t.token(TokenEqualEqual)
			} else {
				token = t.token(TokenEqual)
			}
		case '>':
			if t.consumeIfNext('=') {
				token = t.token(TokenGreaterEqual)
			} else {
				token = t.token(TokenGreater)
			}
		case '<':
			if t.consumeIfNext('=') {
				token = t.token(TokenLessEqual)
			} else {
				token = t.token(TokenLess)
			}
		default:
			return tokens, &UnexpectedCharacterError{
				character: r,
				line:      t.line,
				column:    t.column,
			}
		}
		tokens = append(tokens, token)
	}
	return tokens, nil
}

func (t *Tokenizer) token(tokenID TokenID) Token {
	return Token{
		id:     tokenID,
		line:   t.line,
		column: t.column,
	}
}

func (t *Tokenizer) consumeIfNext(expected rune) bool {
	r, _, err := t.input.ReadRune()
	if err != nil || r != expected {
		t.input.UnreadRune()
		return false
	}
	return true
}

func (t *Tokenizer) consumeUntilEOL() error {
	for {
		r, _, err := t.input.ReadRune()
		if err != nil || r == '\n' {
			t.input.UnreadRune()
			return err
		}
	}
}
