package parser

import (
	"brianhang.me/interpreter/tokenize"
)

// expression ::= assignment
// assignment ::= IDENTIFIER '=' assignment | disjunction
// disjunction ::= conjunction ('or' conjunction)*
// conjunction ::= equality ('and' equality)*
// equality ::= comparison (('==' | '!=') comparison)*
// comparison ::= term (('>=' | '>' | '<=' | '<') term)*
// term ::= factor (('+' | '-') factor)*
// factor ::= unary (('*' | '/') unary)*
// unary ::= ('!' | '-') unary | atom
// atom ::= '(' expression ')'
//      | IDENTIFIER
//      | NUMBER
//      | STRING
//      | 'true' | 'false'
//      | 'nil'
type Parser struct {
	tokens      *[]tokenize.TokenHolder
	curTokenIdx int
}

func NewParser(tokens *[]tokenize.TokenHolder) *Parser {
	parser := &Parser{tokens: tokens}
	return parser
}
