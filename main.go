package main

import (
	"fmt"
	"os"

	"brianhang.me/interpreter/parser"
	"brianhang.me/interpreter/tokenize"
)

func main() {
	tokenizer := tokenize.NewTokenizer(os.Stdin)
	tokens, err := tokenizer.Tokenize()
	if err != nil {
		fmt.Printf("Failed to tokenize: %s\n", err)
		return
	}
	parser := parser.NewParser(&tokens)
	nodes, err := parser.Parse()
	if err != nil {
		fmt.Printf("Failed to parse: %s\n", err)
	}
	fmt.Printf("%s\n", nodes)
}
