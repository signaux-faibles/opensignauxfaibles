package base

import (
	"fmt"
)

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

// newCriticError creates an error with the provided criticity
func newCriticError(err error, criticity string) CriticityError {
	if err == nil {
		return nil
	}
	return &CriticError{err, criticity}
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

// FilterError occurs when something goes wrong while filtering
type FilterError struct {
	*CriticError
}

// NewFilterError returns a filter error
func NewFilterError(err error) error {
	return newCriticError(err, "filter")
}

// NewRegularError creates a regular error
func NewRegularError(err error) error {
	return newCriticError(err, "error")
}

// NewFatalError creates a fatal error
func NewFatalError(err error) error {
	return newCriticError(err, "fatal")
}

// MappingError occurs when something goes wrong while looking for a mapping
type MappingError struct {
	*CriticError
}

// NewMappingError return a mapping Error
func NewMappingError(err error, criticity string) error {
	if err == nil {
		return nil
	}
	return &MappingError{newCriticError(err, criticity).(*CriticError)}
}

func (pe *MappingError) Error() string {
	return fmt.Sprintf("Error while loading or applying the key mapping: %v", pe.err)
}
