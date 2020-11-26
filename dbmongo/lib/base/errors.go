package base

// ParserError is an error that can be reported by any file parser
type ParserError interface {
	error
	IsFilterError() bool
}

type parserError struct {
	err           error
	isFilterError bool
}

// IsFilterError returns the type of error
func (pe *parserError) IsFilterError() bool {
	return pe.isFilterError
}

func (pe *parserError) Error() string {
	return pe.err.Error()
}

// NewFilterError returns a filter error (occurs when something goes wrong while filtering)
func NewFilterError(err error) ParserError {
	if err == nil {
		return nil
	}
	return &parserError{err, true}
}

// NewRegularError creates a regular error
func NewRegularError(err error) ParserError {
	if err == nil {
		return nil
	}
	return &parserError{err, false}
}
