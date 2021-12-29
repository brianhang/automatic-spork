package tokenize

import (
	"bufio"
	"io"
	"strings"
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

func (t *Tokenizer) Tokenize() ([]TokenHolder, error) {
	tokens := make([]TokenHolder, 0)
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
		var token TokenHolder
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
		case '"', '\'':
			token, err = t.string(r, t.line, t.column)
			if err != nil {
				return tokens, err
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

func (t *Tokenizer) string(delimiter rune, startLine int, startCol int) (StringToken, error) {
	token := StringToken{
		Token: Token{id: TokenString},
	}
	var sb strings.Builder
	isEscaping := false
	for {
		r, _, err := t.input.ReadRune()
		if err == io.EOF {
			return token, &UnterminatedStringError{
				delimiter: delimiter,
				line:      startLine,
				column:    startCol,
			}
		}
		if err != nil {
			return token, err
		}
		if isEscaping {
			switch r {
			case '\\', delimiter:
				sb.WriteRune(r)
			case 'n':
				sb.WriteRune('\n')
			case 'r':
				sb.WriteRune('\r')
			case 't':
				sb.WriteRune('\t')
			case 'b':
				sb.WriteRune('\b')
			case 'f':
				sb.WriteRune('\f')
			}
			isEscaping = false
			continue
		}
		if r == '\\' {
			isEscaping = true
			continue
		}
		if r == delimiter {
			break
		}
		sb.WriteRune(r)
	}
	token.value = sb.String()
	return token, nil
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
