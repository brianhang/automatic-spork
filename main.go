package main

import (
	"fmt"
	"os"

	"brianhang.me/interpreter/tokenize"
)

func main() {
	tokenizer := tokenize.NewTokenizer(os.Stdin)
	tokens, err := tokenizer.Tokenize()
	if err != nil {
		fmt.Printf("Failed to tokenize: %s\n", err)
		return
	}
	for _, token := range tokens {
		fmt.Println(token)
	}
}
