package parser

import "brianhang.me/interpreter/tokenize"

type Node interface {
	GetStartToken() tokenize.TokenHolder
	GetEndToken() tokenize.TokenHolder
}

type StatementNode interface {
	Node
}

type ExpressionNode interface {
	Node
}

type ConditionalNode struct {
	If        tokenize.IdentifierToken
	Condition ExpressionNode
	TrueBody  StatementNode
	Else      tokenize.TokenHolder
	FalseBody StatementNode
}

type WhileNode struct {
	While     tokenize.IdentifierToken
	Condition ExpressionNode
	Body      StatementNode
}

type ForNode struct {
	For       tokenize.IdentifierToken
	Init      ExpressionNode
	Condition ExpressionNode
	Update    ExpressionNode
	Body      StatementNode
}

type BlockNode struct {
	BodyStart tokenize.Token
	Children  []StatementNode
	BodyEnd   tokenize.Token
}

type AssignmentNode struct {
	LHS   tokenize.TokenHolder
	Equal tokenize.TokenHolder
	RHS   ExpressionNode
}

type CallNode struct {
	Function   ExpressionNode
	LeftParen  tokenize.TokenHolder
	Args       []ExpressionNode
	RightParen tokenize.TokenHolder
}

type FuncNode struct {
	Func       tokenize.TokenHolder
	LeftParen  tokenize.TokenHolder
	Params     []tokenize.IdentifierToken
	RightParen tokenize.TokenHolder
}

type ReturnNode struct {
	Return tokenize.IdentifierToken
	Value  ExpressionNode
}

type ClassNode struct {
	Class       tokenize.IdentifierToken
	Extends     tokenize.TokenHolder
	ParentClass tokenize.IdentifierToken
	BodyStart   tokenize.TokenHolder
	Body        []AssignmentNode
	BodyEnd     tokenize.TokenHolder
}

type LogicalExprNode struct {
	LHS      ExpressionNode
	Operator tokenize.IdentifierToken
	RHS      ExpressionNode
}

type BinaryExprNode struct {
	LHS      ExpressionNode
	Operator tokenize.Token
	RHS      ExpressionNode
}

type UnaryExprNode struct {
	Operator tokenize.Token
	Operand  ExpressionNode
}

type LookupNode struct {
	Value ExpressionNode
	Key   tokenize.IdentifierToken
}

type LiteralNode struct {
	Value tokenize.TokenHolder
}

func (n ConditionalNode) GetStartToken() tokenize.TokenHolder {
	return n.If
}
func (n ConditionalNode) GetEndToken() tokenize.TokenHolder {
	if n.FalseBody != nil {
		return n.FalseBody.GetEndToken()
	}
	return n.TrueBody.GetEndToken()
}

func (n WhileNode) GetStartToken() tokenize.TokenHolder {
	return n.While
}
func (n WhileNode) GetEndToken() tokenize.TokenHolder {
	return n.Body.GetEndToken()
}

func (n ForNode) GetStartToken() tokenize.TokenHolder {
	return n.For
}
func (n ForNode) GetEndToken() tokenize.TokenHolder {
	return n.Body.GetEndToken()
}

func (n BlockNode) GetStartToken() tokenize.TokenHolder {
	return n.BodyStart
}
func (n BlockNode) GetEndToken() tokenize.TokenHolder {
	return n.BodyEnd
}

func (n AssignmentNode) GetStartToken() tokenize.TokenHolder {
	return n.LHS
}
func (n AssignmentNode) GetEndToken() tokenize.TokenHolder {
	return n.RHS.GetEndToken()
}

func (n CallNode) GetStartToken() tokenize.TokenHolder {
	return n.Function.GetStartToken()
}
func (n CallNode) GetEndToken() tokenize.TokenHolder {
	return n.RightParen
}

func (n FuncNode) GetStartToken() tokenize.TokenHolder {
	return n.Func
}
func (n FuncNode) GetEndToken() tokenize.TokenHolder {
	return n.RightParen
}

func (n ReturnNode) GetStartToken() tokenize.TokenHolder {
	return n.Return
}
func (n ReturnNode) GetEndToken() tokenize.TokenHolder {
	return n.Value.GetEndToken()
}

func (n ClassNode) GetStartToken() tokenize.TokenHolder {
	return n.Class
}
func (n ClassNode) GetEndToken() tokenize.TokenHolder {
	return n.BodyEnd
}

func (n LogicalExprNode) GetStartToken() tokenize.TokenHolder {
	return n.LHS.GetStartToken()
}
func (n LogicalExprNode) GetEndToken() tokenize.TokenHolder {
	return n.RHS.GetEndToken()
}

func (n BinaryExprNode) GetStartToken() tokenize.TokenHolder {
	return n.LHS.GetStartToken()
}
func (n BinaryExprNode) GetEndToken() tokenize.TokenHolder {
	return n.RHS.GetEndToken()
}

func (n UnaryExprNode) GetStartToken() tokenize.TokenHolder {
	return n.Operator
}
func (n UnaryExprNode) GetEndToken() tokenize.TokenHolder {
	return n.Operand.GetEndToken()
}

func (n LookupNode) GetStartToken() tokenize.TokenHolder {
	return n.Value.GetStartToken()
}
func (n LookupNode) GetEndToken() tokenize.TokenHolder {
	return n.Value.GetEndToken()
}

func (n LiteralNode) GetStartToken() tokenize.TokenHolder {
	return n.Value
}
func (n LiteralNode) GetEndToken() tokenize.TokenHolder {
	return n.Value
}
