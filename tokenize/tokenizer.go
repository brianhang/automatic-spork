package tokenize

import (
	"bufio"
	"io"
	"strconv"
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
	'-': TokenMinus,
	'+': TokenPlus,
	'*': TokenStar,
	';': TokenSemicolon,
}

var keywordTokenTypes = map[string]TokenID{
	"true":   TokenTrue,
	"false":  TokenFalse,
	"nil":    TokenNil,
	"and":    TokenAnd,
	"or":     TokenOr,
	"if":     TokenIf,
	"else":   TokenElse,
	"for":    TokenFor,
	"while":  TokenWhile,
	"func":   TokenFunc,
	"return": TokenReturn,
	"class":  TokenClass,
	"super":  TokenSuper,
	"this":   TokenThis,
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

		t.column++

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
			if numberToken, err := t.number(); err == nil {
				token = numberToken
			} else {
				token = t.token(TokenDot)
			}
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
			token, err = t.string(r)
			if err != nil {
				return tokens, err
			}
		default:
			if unicode.IsDigit(r) {
				numberToken, err := t.number()
				if err != nil {
					return tokens, err
				}
				token = numberToken
			} else if isRuneStartOfIdentifier(r) {
				identifierToken, err := t.identifier()
				if err != nil {
					return tokens, err
				}
				token = identifierToken
			} else {
				return tokens, &UnexpectedCharacterError{
					character: r,
					line:      t.line,
					column:    t.column,
				}
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

func (t *Tokenizer) string(delimiter rune) (StringToken, error) {
	token := StringToken{
		Token: Token{id: TokenString, line: t.line, column: t.column},
	}
	var sb strings.Builder
	isEscaping := false
	for {
		r, _, err := t.input.ReadRune()
		t.column++
		if err == io.EOF {
			return token, &UnterminatedStringError{
				delimiter: delimiter,
				line:      token.line,
				column:    token.column,
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

func (t *Tokenizer) number() (NumberToken, error) {
	token := NumberToken{
		Token: Token{id: TokenNumber, line: t.line, column: t.column},
	}
	err := t.input.UnreadRune()
	if err != nil {
		return token, err
	}
	isFractional := false
	var sb strings.Builder
	for {
		r, _, err := t.input.ReadRune()
		if r == '.' {
			if isFractional {
				t.input.UnreadRune()
				break
			}
			sb.WriteRune(r)
			isFractional = true
			continue
		}
		if err != nil || !unicode.IsDigit(r) {
			t.input.UnreadRune()
			break
		}
		t.column++
		sb.WriteRune(r)
	}
	value, err := strconv.ParseFloat(sb.String(), 64)
	if err != nil {
		return token, err
	}
	token.value = value
	return token, nil
}

func (t *Tokenizer) identifier() (IdentifierToken, error) {
	token := IdentifierToken{
		Token: Token{id: TokenIdentifier, line: t.line, column: t.column},
	}
	err := t.input.UnreadRune()
	t.column--
	if err != nil {
		return token, err
	}
	var sb strings.Builder
	for {
		r, _, err := t.input.ReadRune()
		if err == io.EOF {
			break
		}
		if err != nil {
			t.input.UnreadRune()
			return token, err
		}
		if !isIdentifierRune(r) {
			err := t.input.UnreadRune()
			if err != nil {
				return token, err
			}
			break
		}
		t.column++
		sb.WriteRune(r)
	}
	value := sb.String()
	if keywordTokenType, ok := keywordTokenTypes[value]; ok {
		token.id = keywordTokenType
	} else {
		token.value = value
	}
	return token, nil
}

func isRuneStartOfIdentifier(r rune) bool {
	return r == '_' || unicode.IsLetter(r)
}

func isIdentifierRune(r rune) bool {
	return isRuneStartOfIdentifier(r) || unicode.IsDigit(r)
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
