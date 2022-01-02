package parser

import (
	"fmt"
	"strings"
	"testing"

	"brianhang.me/interpreter/tokenize"
	"github.com/stretchr/testify/assert"
)

func TestParser(t *testing.T) {
	cases := []struct {
		source      string
		expectedAST string
	}{
		{
			"1 + 2 * 3 / 4 - 5",
			"[(- (+ (number 1) (/ (* (number 2) (number 3)) (number 4))) (number 5))]",
		},
		{
			"true and y or z and false",
			"[(or (and (true ) (identifier y)) (and (identifier z) (false )))]",
		},
		{
			"y = func(x){ x = x + 1 return x + 'hello' }",
			"[(= y (func x (block [(= x (+ (identifier x) (number 1))) (return (+ (identifier x) (string \"hello\")))])))]",
		},
	}
	for _, test := range cases {
		tokenizer := tokenize.NewTokenizer(strings.NewReader(test.source))
		tokens, err := tokenizer.Tokenize()
		if err != nil {
			t.Errorf("Failed to tokenize: %s\n", err)
			continue
		}
		parser := NewParser(&tokens)
		nodes, err := parser.Parse()
		if err != nil {
			fmt.Printf("Failed to parse: %s\n", err)
		}
		assert.Equal(t, test.expectedAST, fmt.Sprintf("%s", nodes))
	}
}
