package parser

import (
	"brianhang.me/interpreter/tokenize"
)

// statement       ::= while
//                   | for
//                   | return
//                   | expression
//
// expression      ::= assignment
//                   | block
//                   | class
//                   | conditional
//                   | func
//
// conditional     ::= 'if' '(' expression ')' statement ('else' statement)?
// while           ::= 'while' '(' expression ')' statement
// for             ::= 'for' '(' expression? ';' expression? ';' expression? ')' statement
//
// block           ::= '{' statement* '}'
//
// assignment      ::= IDENTIFIER '=' assignment | disjunction
//
// params          ::= identifier (',' identifier)* ','?
// func            ::= 'func' '(' params? ')' block
// return          ::= 'return' expression?
//
// class           ::= 'class' ('<' identifier)? '{' classAssignment* '}'
// classAssignment ::= IDENTIFIER '=' (func | disjunction)
//
// disjunction     ::= conjunction ('or' conjunction)*
// conjunction     ::= equality ('and' equality)*
//
// equality        ::= comparison (('==' | '!=') comparison)*
// comparison      ::= term (('>=' | '>' | '<=' | '<') term)*
//
// term            ::= factor (('+' | '-') factor)*
// factor          ::= unary (('*' | '/') unary)*
// unary           ::= ('!' | '-') unary | call
// args            ::= expression (',' expression)* ','?
// call            ::= atom ('(' args? ')' | '.' IDENTIFIER)*
// atom            ::= '(' expression ')'
//                   | IDENTIFIER
//                   | NUMBER
//                   | STRING
//                   | 'true' | 'false'
//                   | 'nil'

type Parser struct {
	tokens      *[]tokenize.TokenHolder
	curTokenIdx int
}

func NewParser(tokens *[]tokenize.TokenHolder) *Parser {
	parser := &Parser{tokens: tokens}
	return parser
}

func (p *Parser) Parse() ([]Node, error) {
	statements := make([]Node, 0)
	for {
		statement, err := p.maybeStatement()
		if err != nil {
			return statements, err
		}
		if statement == nil {
			break
		}
		statements = append(statements, statement)
	}
	return statements, nil
}

func (p *Parser) statement() (Node, error) {
	statement, err := p.maybeStatement()
	if err != nil {
		return statement, err
	}
	if statement == nil {
		return nil, &ExpectedStatementError{last: p.last()}
	}
	return statement, nil
}

func (p *Parser) maybeStatement() (Node, error) {
	token := p.peek()
	if token == nil {
		return nil, nil
	}
	switch token.GetID() {
	case tokenize.TokenWhile:
		return p.while()
	case tokenize.TokenFor:
		return p.forStatement()
	case tokenize.TokenReturn:
		return p.returnStatement()
	default:
		return p.maybeExpression()
	}
}

func (p *Parser) expression() (Node, error) {
	expression, err := p.maybeExpression()
	if err != nil {
		return expression, err
	}
	if expression == nil {
		return nil, &ExpectedExpressionError{last: p.last()}
	}
	return expression, nil
}

func (p *Parser) maybeExpression() (Node, error) {
	token := p.peek()
	if token == nil {
		return nil, nil
	}
	switch token.GetID() {
	case tokenize.TokenLeftCurly:
		return p.block()
	case tokenize.TokenClass:
		return p.class()
	case tokenize.TokenIf:
		return p.conditional()
	case tokenize.TokenFunc:
		return p.funcExpr()
	default:
		return p.assignment()
	}
}

var atomicTokenIDs = []tokenize.TokenID{
	tokenize.TokenIdentifier,
	tokenize.TokenNumber,
	tokenize.TokenString,
	tokenize.TokenTrue,
	tokenize.TokenFalse,
	tokenize.TokenNil,
}

func (p *Parser) atom() (ExpressionNode, error) {
	leftParen := p.maybeMatch(tokenize.TokenLeftParen)
	if leftParen != nil {
		expression, err := p.expression()
		if err != nil {
			return expression, err
		}
		_, err = p.match(tokenize.TokenRightParen)
		return expression, err
	}
	node := LiteralNode{}
	for _, tokenID := range atomicTokenIDs {
		value := p.maybeMatch(tokenID)
		if value != nil {
			node.Value = value
			return node, nil
		}
	}
	return node, &NoValueError{last: p.last()}
}

var unaryOperatorTokenIDs = []tokenize.TokenID{
	tokenize.TokenBang,
	tokenize.TokenMinus,
}

func (p *Parser) call() (ExpressionNode, error) {
	var err error
	node, err := p.atom()
	if err != nil {
		return node, err
	}
	isFindingCalls := true
	for isFindingCalls {
		nextToken := p.peek()
		if nextToken == nil {
			break
		}
		switch nextToken.GetID() {
		case tokenize.TokenLeftParen:
			call := CallNode{Function: node, LeftParen: p.consume()}
			call.Args, err = p.args()
			if err != nil {
				return call, err
			}
			call.RightParen, err = p.match(tokenize.TokenRightParen)
			if err != nil {
				return call, err
			}
			node = call
		case tokenize.TokenDot:
			lookup := LookupNode{Value: node}
			lookup.Key, err = p.matchIdentifier(tokenize.TokenIdentifier)
			if err != nil {
				return lookup, err
			}
			node = lookup
		default:
			isFindingCalls = false
		}
	}
	return node, nil
}

func (p *Parser) args() ([]ExpressionNode, error) {
	isFirstArg := true
	args := make([]ExpressionNode, 0)
	for {
		if isFirstArg {
			isFirstArg = false
		} else {
			p.match(tokenize.TokenComma)
		}
		arg, err := p.maybeExpression()
		if err != nil {
			return args, err
		}
		if arg == nil {
			break
		}
		args = append(args, arg)
	}
	p.maybeMatch(tokenize.TokenComma)
	return args, nil
}

func (p *Parser) unary() (ExpressionNode, error) {
	var err error
	var operator tokenize.TokenHolder
	for _, tokenID := range unaryOperatorTokenIDs {
		operator = p.maybeMatch(tokenID)
		if operator != nil {
			break
		}
	}
	if operator == nil {
		return p.call()
	}
	unaryExpr := UnaryExprNode{Operator: operator.GetToken()}
	unaryExpr.Operand, err = p.unary()
	return unaryExpr, err
}

var factorOperatorTokenIDs = []tokenize.TokenID{
	tokenize.TokenStar,
	tokenize.TokenSlash,
}

func (p *Parser) factor() (ExpressionNode, error) {
	return p.binaryExpression(factorOperatorTokenIDs, p.unary)
}

var termOperatorTokenIDs = []tokenize.TokenID{
	tokenize.TokenPlus,
	tokenize.TokenMinus,
}

func (p *Parser) term() (ExpressionNode, error) {
	return p.binaryExpression(termOperatorTokenIDs, p.factor)
}

var comparisonOperatorTokenIDs = []tokenize.TokenID{
	tokenize.TokenGreater,
	tokenize.TokenGreaterEqual,
	tokenize.TokenLess,
	tokenize.TokenLessEqual,
}

func (p *Parser) comparison() (ExpressionNode, error) {
	return p.binaryExpression(comparisonOperatorTokenIDs, p.term)
}

var equalityOperatorTokenIDs = []tokenize.TokenID{
	tokenize.TokenEqualEqual,
	tokenize.TokenBangEqual,
}

func (p *Parser) equality() (ExpressionNode, error) {
	return p.binaryExpression(equalityOperatorTokenIDs, p.comparison)
}

var conjunctionOperatorTokenIDs = []tokenize.TokenID{
	tokenize.TokenAnd,
}

func (p *Parser) conjunction() (ExpressionNode, error) {
	return p.binaryExpression(conjunctionOperatorTokenIDs, p.equality)
}

var disjunctionOperatorTokenIDs = []tokenize.TokenID{
	tokenize.TokenOr,
}

func (p *Parser) disjunction() (ExpressionNode, error) {
	return p.binaryExpression(disjunctionOperatorTokenIDs, p.conjunction)
}

func (p *Parser) binaryExpression(
	operatorTokenIDs []tokenize.TokenID,
	getOperand func() (ExpressionNode, error),
) (ExpressionNode, error) {
	var err error
	node, err := getOperand()
	if err != nil {
		return node, err
	}
	for {
		var operator tokenize.TokenHolder
		for _, tokenID := range operatorTokenIDs {
			operator = p.maybeMatch(tokenID)
			if operator != nil {
				break
			}
		}
		if operator == nil {
			break
		}
		binExprNode := BinaryExprNode{Operator: operator.GetToken(), LHS: node}
		binExprNode.RHS, err = getOperand()
		if err != nil {
			return node, err
		}
		node = binExprNode
	}
	return node, nil
}

func (p *Parser) class() (ClassNode, error) {
	var err error
	node := ClassNode{}
	node.Class, err = p.matchIdentifier(tokenize.TokenClass)
	if err != nil {
		return node, err
	}

	node.Extends = p.maybeMatch(tokenize.TokenLess)
	if node.Extends != nil {
		node.ParentClass, err = p.matchIdentifier(tokenize.TokenIdentifier)
		if err != nil {
			return node, err
		}
	}

	if node.BodyStart, err = p.match(tokenize.TokenLeftCurly); err != nil {
		return node, err
	}
	for {
		identifier := p.maybeMatch(tokenize.TokenIdentifier)
		if identifier == nil {
			break
		}
		equal, err := p.match(tokenize.TokenEqual)
		if err != nil {
			return node, err
		}
		value, err := p.expression()
		if err != nil {
			return node, err
		}
		if value != nil {
			value, err = p.disjunction()
			if err != nil {
				return node, err
			}
		}
		assignment := AssignmentNode{LHS: identifier, Equal: equal, RHS: value}
		node.Body = append(node.Body, assignment)
	}
	if node.BodyEnd, err = p.match(tokenize.TokenRightCurly); err != nil {
		return node, err
	}
	return node, nil
}

func (p *Parser) returnStatement() (ReturnNode, error) {
	var err error
	node := ReturnNode{}
	node.Return, err = p.matchIdentifier(tokenize.TokenReturn)
	if err != nil {
		return node, err
	}
	node.Value, err = p.maybeExpression()
	if err != nil {
		return node, err
	}
	return node, nil
}

func (p *Parser) funcExpr() (FuncNode, error) {
	var err error
	node := FuncNode{}
	node.Func, err = p.match(tokenize.TokenFunc)
	if err != nil {
		return node, err
	}
	isFirstParam := true
	for {
		if isFirstParam {
			isFirstParam = false
		} else {
			p.match(tokenize.TokenComma)
		}
		param, err := p.match(tokenize.TokenIdentifier)
		if err != nil {
			return node, err
		}
		if param == nil {
			break
		}
		node.Params = append(node.Params, param.(tokenize.IdentifierToken))
	}
	p.maybeMatch(tokenize.TokenComma)
	if _, err = p.match(tokenize.TokenRightParen); err != nil {
		return node, err
	}
	return node, nil
}

func (p *Parser) assignment() (ExpressionNode, error) {
	expr, err := p.disjunction()
	if err != nil {
		return expr, err
	}
	equal := p.maybeMatch(tokenize.TokenEqual)
	if equal == nil {
		return expr, nil
	}
	identifier, ok := expr.(LiteralNode)
	if !ok || identifier.GetStartToken().GetID() != tokenize.TokenIdentifier {
		return identifier, &InvalidAssignmentTargetError{
			target: identifier.GetStartToken(),
		}
	}
	assignment := AssignmentNode{LHS: identifier.Value, Equal: equal}
	assignment.RHS, err = p.assignment()
	if err != nil {
		return assignment, err
	}
	return assignment, nil
}

func (p *Parser) block() (BlockNode, error) {
	var err error
	node := BlockNode{}
	if _, err = p.match(tokenize.TokenLeftCurly); err != nil {
		return node, err
	}
	for {
		statement, err := p.maybeStatement()
		if err != nil {
			return node, err
		}
		if statement == nil {
			break
		}
		node.Children = append(node.Children, statement)
	}
	if _, err = p.match(tokenize.TokenRightCurly); err != nil {
		return node, err
	}
	return node, nil
}

func (p *Parser) forStatement() (ForNode, error) {
	var err error
	node := ForNode{}
	node.For, err = p.matchIdentifier(tokenize.TokenFor)
	if err != nil {
		return node, err
	}
	if _, err = p.match(tokenize.TokenLeftParen); err != nil {
		return node, err
	}
	node.Init, err = p.maybeExpression()
	if err != nil {
		return node, err
	}
	if _, err = p.match(tokenize.TokenSemicolon); err != nil {
		return node, err
	}
	node.Condition, err = p.maybeExpression()
	if err != nil {
		return node, err
	}
	if _, err = p.match(tokenize.TokenSemicolon); err != nil {
		return node, err
	}
	node.Update, err = p.maybeExpression()
	if err != nil {
		return node, err
	}
	if _, err = p.match(tokenize.TokenRightParen); err != nil {
		return node, err
	}
	node.Body, err = p.statement()
	if err != nil {
		return node, err
	}
	return node, nil
}

func (p *Parser) while() (WhileNode, error) {
	var err error
	node := WhileNode{}
	node.While, err = p.matchIdentifier(tokenize.TokenWhile)
	if err != nil {
		return node, err
	}
	if _, err = p.match(tokenize.TokenLeftParen); err != nil {
		return node, err
	}
	node.Body, err = p.statement()
	if err != nil {
		return node, err
	}
	if _, err = p.match(tokenize.TokenRightParen); err != nil {
		return node, err
	}
	return node, nil
}

func (p *Parser) conditional() (ConditionalNode, error) {
	var err error
	node := ConditionalNode{}
	node.If, err = p.matchIdentifier(tokenize.TokenIf)
	if err != nil {
		return node, err
	}
	if _, err := p.match(tokenize.TokenLeftParen); err != nil {
		return node, err
	}
	node.Condition, err = p.expression()
	if err != nil {
		return node, err
	}
	if _, err = p.match(tokenize.TokenRightParen); err != nil {
		return node, err
	}
	node.TrueBody, err = p.statement()
	if err != nil {
		return node, err
	}
	node.Else = p.maybeMatch(tokenize.TokenElse)
	if node.Else == nil {
		return node, err
	}
	node.FalseBody, err = p.statement()
	if err != nil {
		return node, err
	}
	return node, nil
}

func (p *Parser) match(id tokenize.TokenID) (tokenize.TokenHolder, error) {
	token := p.maybeMatch(id)
	if token == nil {
		return nil, &ExpectedTokenError{
			expected: id,
			last:     p.last(),
		}
	}
	return token, nil
}

func (p *Parser) maybeMatch(id tokenize.TokenID) tokenize.TokenHolder {
	token := p.peek()
	if token == nil || token.GetID() != id {
		return nil
	}
	return p.consume()
}

func (p *Parser) matchIdentifier(id tokenize.TokenID) (tokenize.IdentifierToken, error) {
	token, err := p.match(id)
	identifier, ok := token.(tokenize.IdentifierToken)
	if err != nil {
		return identifier, err
	}
	if ok {
		return identifier, nil
	}
	return identifier, &ExpectedTokenError{
		expected: id,
		actual:   token,
		last:     p.last(),
	}
}

func (p *Parser) tokenAtOffset(offset int) tokenize.TokenHolder {
	idx := p.curTokenIdx + offset
	if p.tokens == nil || idx >= len(*p.tokens) {
		return nil
	}
	return (*p.tokens)[idx]
}

func (p *Parser) peek() tokenize.TokenHolder {
	return p.tokenAtOffset(0)
}

func (p *Parser) consume() tokenize.TokenHolder {
	token := p.peek()
	if token != nil {
		p.curTokenIdx++
	}
	return token
}

func (p *Parser) last() tokenize.TokenHolder {
	return p.tokenAtOffset(-1)
}
