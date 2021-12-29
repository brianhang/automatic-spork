package tokenize

import "fmt"

type UnexpectedCharacterError struct {
	character rune
	line      int
	column    int
}

func (e *UnexpectedCharacterError) Error() string {
	return fmt.Sprintf(
		"Unexpected character '%c' on line %d at column %d",
		e.character,
		e.line,
		e.column,
	)
}

type UnterminatedStringError struct {
	delimiter rune
	line      int
	column    int
}

func (e *UnterminatedStringError) Error() string {
	return fmt.Sprintf(
		"Expected a closing %c for string starting on line %d at column %d",
		e.delimiter,
		e.line,
		e.column,
	)
}
