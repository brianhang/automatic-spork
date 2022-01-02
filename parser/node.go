package parser

import (
	"fmt"
	"strings"

	"brianhang.me/interpreter/tokenize"
)

type Node interface {
	GetStartToken() tokenize.TokenHolder
	GetEndToken() tokenize.TokenHolder
	String() string
}

type StatementNode interface {
	Node
}

type ExpressionNode interface {
	Node
}

type AtomNode interface {
	ExpressionNode
	GetToken() tokenize.TokenHolder
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
	Body       BlockNode
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
func (n ConditionalNode) String() string {
	if n.FalseBody == nil {
		return fmt.Sprintf("(if %s %s)", n.Condition, n.TrueBody)
	}
	return fmt.Sprintf("(if %s %s else %s)", n.Condition, n.TrueBody, n.FalseBody)
}

func (n WhileNode) GetStartToken() tokenize.TokenHolder {
	return n.While
}
func (n WhileNode) GetEndToken() tokenize.TokenHolder {
	return n.Body.GetEndToken()
}
func (n WhileNode) String() string {
	return fmt.Sprintf("(while %s %s)", n.Condition, n.Body)
}

func (n ForNode) GetStartToken() tokenize.TokenHolder {
	return n.For
}
func (n ForNode) GetEndToken() tokenize.TokenHolder {
	return n.Body.GetEndToken()
}
func (n ForNode) String() string {
	return fmt.Sprintf("(for %s %s %s %s)", n.Init, n.Condition, n.Update, n.Body)
}

func (n BlockNode) GetStartToken() tokenize.TokenHolder {
	return n.BodyStart
}
func (n BlockNode) GetEndToken() tokenize.TokenHolder {
	return n.BodyEnd
}
func (n BlockNode) String() string {
	return fmt.Sprintf("(block %s)", n.Children)
}

func (n AssignmentNode) GetStartToken() tokenize.TokenHolder {
	return n.LHS
}
func (n AssignmentNode) GetEndToken() tokenize.TokenHolder {
	return n.RHS.GetEndToken()
}
func (n AssignmentNode) String() string {
	return fmt.Sprintf("(= %s %s)", n.LHS, n.RHS)
}

func (n CallNode) GetStartToken() tokenize.TokenHolder {
	return n.Function.GetStartToken()
}
func (n CallNode) GetEndToken() tokenize.TokenHolder {
	return n.RightParen
}
func (n CallNode) String() string {
	var args strings.Builder
	for _, child := range n.Args {
		args.WriteString(fmt.Sprintf(" %s", child))
	}
	return fmt.Sprintf("(call %s%s)", n.Function, args.String())
}

func (n FuncNode) GetStartToken() tokenize.TokenHolder {
	return n.Func
}
func (n FuncNode) GetEndToken() tokenize.TokenHolder {
	return n.RightParen
}
func (n FuncNode) String() string {
	var params strings.Builder
	for _, child := range n.Params {
		params.WriteString(fmt.Sprintf(" %s", child.GetValue()))
	}
	return fmt.Sprintf("(func%s %s)", params.String(), n.Body)
}

func (n ReturnNode) GetStartToken() tokenize.TokenHolder {
	return n.Return
}
func (n ReturnNode) GetEndToken() tokenize.TokenHolder {
	return n.Value.GetEndToken()
}
func (n ReturnNode) String() string {
	value := n.Value
	if value == nil {
		return "(return nil)"
	}
	return fmt.Sprintf("(return %s)", value)
}

func (n ClassNode) GetStartToken() tokenize.TokenHolder {
	return n.Class
}
func (n ClassNode) GetEndToken() tokenize.TokenHolder {
	return n.BodyEnd
}
func (n ClassNode) String() string {
	parentName := n.ParentClass.GetValue()
	if len(parentName) > 0 {
		return fmt.Sprintf("(class %s %s)", parentName, n.Body)
	}
	return fmt.Sprintf("(class %s)", n.Body)
}

func (n LogicalExprNode) GetStartToken() tokenize.TokenHolder {
	return n.LHS.GetStartToken()
}
func (n LogicalExprNode) GetEndToken() tokenize.TokenHolder {
	return n.RHS.GetEndToken()
}
func (n LogicalExprNode) String() string {
	return fmt.Sprintf("(%s %s %s)", n.Operator.GetValue(), n.LHS, n.RHS)
}

func (n BinaryExprNode) GetStartToken() tokenize.TokenHolder {
	return n.LHS.GetStartToken()
}
func (n BinaryExprNode) GetEndToken() tokenize.TokenHolder {
	return n.RHS.GetEndToken()
}
func (n BinaryExprNode) String() string {
	return fmt.Sprintf("(%s %s %s)", n.Operator.GetToken(), n.LHS, n.RHS)
}

func (n UnaryExprNode) GetStartToken() tokenize.TokenHolder {
	return n.Operator
}
func (n UnaryExprNode) GetEndToken() tokenize.TokenHolder {
	return n.Operand.GetEndToken()
}
func (n UnaryExprNode) String() string {
	return fmt.Sprintf("(%d %s)", n.Operator.GetToken(), n.Operand)
}

func (n LookupNode) GetStartToken() tokenize.TokenHolder {
	return n.Value.GetStartToken()
}
func (n LookupNode) GetEndToken() tokenize.TokenHolder {
	return n.Value.GetEndToken()
}
func (n LookupNode) String() string {
	return fmt.Sprintf("(lookup %s %s)", n.Value, n.Key.GetValue())
}

func (n LiteralNode) GetStartToken() tokenize.TokenHolder {
	return n.Value
}
func (n LiteralNode) GetEndToken() tokenize.TokenHolder {
	return n.Value
}
func (n LiteralNode) String() string {
	return fmt.Sprintf("(%s %s)", n.Value.GetID(), n.Value)
}
