package base

import (
	"errors"
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

// NewCriticError creates a critical error
func NewCriticError(err error, criticity string) error {
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

// ParseError occurs when something goes wrong while parsing
type ParseError struct {
	*CriticError
	Filename string
}

// NewCorruptedRowError creates an error for when a corrupted row is encountered.
func NewCorruptedRowError(Filename string) error {
	return NewParseError(errors.New("corrupted line"), Filename)
}

// NewParseError error parser
func NewParseError(err error, Filename string) error {
	if err == nil {
		return nil
	}
	c := NewCriticError(err, "error")
	return &ParseError{c.(*CriticError), Filename}
}

func (pe *ParseError) Error() string {
	return fmt.Sprintf("Error while parsing %s: %v", pe.Filename, pe.err)
}

// FilterError occurs when something goes wrong while filtering
type FilterError struct {
	*CriticError
}

// NewFilterError returns a filter error
func NewFilterError(err error) *FilterError {
	if err == nil {
		return nil
	}
	c := NewCriticError(err, "filter")
	return &FilterError{c.(*CriticError)}
}

func (pe *FilterError) Error() string {
	return fmt.Sprintf("Error while loading or applying filter: %v", pe.err)
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
	return &MappingError{NewCriticError(err, criticity).(*CriticError)}
}

func (pe *MappingError) Error() string {
	return fmt.Sprintf("Error while loading or applying the key mapping: %v", pe.err)
}
