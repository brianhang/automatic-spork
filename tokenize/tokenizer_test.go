package tokenize

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestTokenizer(t *testing.T) {
	cases := []struct {
		source           string
		expectedTokenIDs []TokenID
	}{
		{
			"1 + 2 - 3 / 4 * 5",
			[]TokenID{TokenNumber, TokenPlus, TokenNumber, TokenMinus, TokenNumber, TokenSlash, TokenNumber, TokenStar, TokenNumber},
		},
		{
			"1 < 2 <= 3   >=2>1",
			[]TokenID{TokenNumber, TokenLess, TokenNumber, TokenLessEqual, TokenNumber, TokenGreaterEqual, TokenNumber, TokenGreater, TokenNumber},
		},
		{
			"true == true and true or !false != false",
			[]TokenID{TokenTrue, TokenEqualEqual, TokenTrue, TokenAnd, TokenTrue, TokenOr, TokenBang, TokenFalse, TokenBangEqual, TokenFalse},
		},
		{
			"-3.14 // ratio",
			[]TokenID{TokenMinus, TokenNumber},
		},
		{
			"func foo() {while (true) { return }}",
			[]TokenID{TokenFunc, TokenIdentifier, TokenLeftParen, TokenRightParen, TokenLeftCurly, TokenWhile, TokenLeftParen, TokenTrue, TokenRightParen, TokenLeftCurly, TokenReturn, TokenRightCurly, TokenRightCurly},
		},
		{
			"ε = .0000001",
			[]TokenID{TokenIdentifier, TokenEqual, TokenNumber},
		},
	}
	for _, test := range cases {
		tokens := tokenizeString(t, test.source)
		tokenIDs := make([]TokenID, len(tokens))
		for idx, token := range tokens {
			tokenIDs[idx] = token.GetToken().id
		}
		assert.Equal(t, test.expectedTokenIDs, tokenIDs, "Tokens for \"%s\" did not match", test.source)
	}
}

func TestTokenizeNumber(t *testing.T) {
	cases := []struct {
		source   string
		expected float64
	}{
		{"3.1415926535897", 3.1415926535897},
		{".05", 0.05},
		{"1337", 1337},
		{"1337.24", 1337.24},
	}
	for _, test := range cases {
		tokens := tokenizeString(t, test.source)
		assert.Equal(t, 1, len(tokens), "Expected input \"%s\" to have one token output", test.source)
		token, ok := tokens[0].(NumberToken)
		assert.True(t, ok, "Expected input \"%s\" to result in a number token, got %s", test.source, tokens[0])
		assert.Equal(t, test.expected, token.value, "Wrong value for input \"%s\"", test.source)
	}
}

func TestTokenizeString(t *testing.T) {
	cases := []struct {
		source   string
		expected string
	}{
		{"\"How're you?\"", "How're you?"},
		{"'\\tfoo\\nbar'", "\tfoo\nbar"},
		{"'single \"quote\"'", "single \"quote\""},
	}
	for _, test := range cases {
		tokens := tokenizeString(t, test.source)
		assert.Equal(t, 1, len(tokens), "Expected input \"%s\" to have one token output", test.source)
		token, ok := tokens[0].(StringToken)
		assert.True(t, ok, "Expected input \"%s\" to result in a string token, got %s", test.source, tokens[0])
		assert.Equal(t, test.expected, token.value, "Wrong value for input \"%s\"", test.source)
	}

	tokenizer := NewTokenizer(strings.NewReader("x = 'unclosed string"))
	_, err := tokenizer.Tokenize()
	assert.Contains(t, err.Error(), "Expected a closing '")
}

func tokenizeString(t *testing.T, source string) []TokenHolder {
	tokenizer := NewTokenizer(strings.NewReader(source))
	tokens, err := tokenizer.Tokenize()
	if err != nil {
		t.Errorf("Unexpected error for input \"%s\": %s", source, err)
	}
	return tokens
}
