package parser

import (
	"fmt"

	"brianhang.me/interpreter/tokenize"
)

type UnexpectedTokenError struct {
	token tokenize.TokenHolder
}

func (e *UnexpectedTokenError) Error() string {
	token := e.token
	return fmt.Sprintf(
		"Unexpected token %s on line %d at column %d",
		token,
		token.GetLine(),
		token.GetColumn(),
	)
}

type ExpectedTokenError struct {
	expected tokenize.TokenID
	actual   tokenize.TokenHolder
	last     tokenize.TokenHolder
}

func (e *ExpectedTokenError) Error() string {
	actual := e.actual
	if actual == nil {
		last := e.last
		if last == nil {
			return fmt.Sprintf("Expected token %d", e.expected)
		}
		return fmt.Sprintf(
			"Expected token %d near line %d",
			e.expected,
			last.GetLine(),
		)
	}
	return fmt.Sprintf(
		"Expected token %d, but got %s instead on line %d at column %d",
		e.expected,
		actual,
		actual.GetLine(),
		actual.GetColumn(),
	)
}

type ExpectedStatementError struct {
	last tokenize.TokenHolder
}
type ExpectedExpressionError struct {
	last tokenize.TokenHolder
}

func (e *ExpectedStatementError) Error() string {
	last := e.last
	if last == nil {
		return "Expected a statement"
	}
	return fmt.Sprintf("Expected a statement near line %d", last.GetLine())
}
func (e *ExpectedExpressionError) Error() string {
	last := e.last
	if last == nil {
		return "Expected an expression"
	}
	return fmt.Sprintf("Expected an expression near line %d", last.GetLine())
}

type InvalidAssignmentTargetError struct {
	target tokenize.TokenHolder
}

func (e *InvalidAssignmentTargetError) Error() string {
	target := e.target
	return fmt.Sprintf(
		"Invalid left hand side for assignment on line %d at column %d",
		target.GetLine(),
		target.GetColumn(),
	)
}

type NoValueError struct {
	last tokenize.TokenHolder
}

func (e *NoValueError) Error() string {
	last := e.last
	if last == nil {
		return "Expected a value, but none was provided"
	}
	return fmt.Sprintf(
		"Expected a value near line %d at column %d, but none was provided",
		last.GetLine(),
		last.GetColumn(),
	)
}
