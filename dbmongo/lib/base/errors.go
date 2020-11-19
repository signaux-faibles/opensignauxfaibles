package base

import "errors"

// CriticityError object
type CriticityError interface {
	error
	Criticity() string
}

// CriticError is a generic error with criticity
type CriticError struct {
	err       error
	criticity string
}

// Criticity returns criticity
func (pe *CriticError) Criticity() string {
	if pe == nil {
		return ""
	}
	return pe.criticity
}

func (pe *CriticError) Error() string {
	return pe.err.Error()
}

// newCriticError creates an error with the provided criticity
func newCriticError(err error, criticity string) CriticityError {
	if err == nil {
		return nil
	}
	return &CriticError{err, criticity}
}

// NewFilterError returns a filter error (occurs when something goes wrong while filtering)
func NewFilterError() CriticityError {
	return newCriticError(errors.New(""), "filter")
}

// NewRegularError creates a regular error
func NewRegularError(err error) CriticityError {
	return newCriticError(err, "error")
}

// NewFatalError creates a fatal error
func NewFatalError(err error) CriticityError {
	return newCriticError(err, "fatal")
}
