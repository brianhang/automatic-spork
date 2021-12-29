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
		return tokens, &UnexpectedCharacterError{
			character: r,
			line:      t.line,
			column:    t.column,
		}
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
