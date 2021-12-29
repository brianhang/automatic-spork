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
