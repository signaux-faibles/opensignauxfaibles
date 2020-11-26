package base

import "errors"

// CriticityError object
type CriticityError interface {
	error
	Criticity() string
}

// CriticError is a generic error with criticity
type CriticError struct {
	err           error
	isFilterError bool
}

// Criticity returns criticity
func (pe *CriticError) Criticity() string {
	if pe == nil {
		return ""
	} else if pe.isFilterError {
		return "filter"
	}
	return "error"
}

func (pe *CriticError) Error() string {
	return pe.err.Error()
}

// newCriticError creates an error with the provided criticity
func newCriticError(err error, isFilterError bool) CriticityError {
	if err == nil {
		return nil
	}
	return &CriticError{err, isFilterError}
}

// NewFilterError returns a filter error (occurs when something goes wrong while filtering)
func NewFilterError() CriticityError {
	return newCriticError(errors.New("(filtered)"), true)
}

// NewRegularError creates a regular error
func NewRegularError(err error) CriticityError {
	return newCriticError(err, false)
}
