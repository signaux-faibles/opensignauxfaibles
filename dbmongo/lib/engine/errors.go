package engine

import "fmt"

type CriticityError interface {
	error
	Criticity() string
}

// CriticError is a generic error with criticity
type CriticError struct {
	err       error
	criticity string
}

func NewCriticError(err error, criticity string) error {
	if err == nil {
		return nil
	}
	return &CriticError{err, criticity}
}

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
	ParsedVariable string
}

func NewParseError(err *CriticError, ParsedVariable string) error {
	if err == nil {
		return nil
	}
	return &ParseError{err, ParsedVariable}
}

func (pe *ParseError) Error() string {
	return fmt.Sprintf("Error while parsing %s: %v", pe.ParsedVariable, pe.err)
}

// FilterError occurs when something goes wrong while filtering
type FilterError struct {
	*CriticError
}

func NewFilterError(err error, criticity string) *FilterError {
	if err == nil {
		return nil
	}
	c := NewCriticError(err, criticity)
	return &FilterError{c.(*CriticError)}
}

func (pe *FilterError) Error() string {
	return fmt.Sprintf("Error while loading or applying filter: %v", pe.err)
}

// MappingError occurs when something goes wrong while filtering
type MappingError struct {
	*CriticError
}

func NewMappingError(err error, criticity string) error {
	if err == nil {
		return nil
	}
	return &MappingError{NewCriticError(err, criticity).(*CriticError)}
}

func (pe *MappingError) Error() string {
	return fmt.Sprintf("Error while loading or applying the key mapping: %v", pe.err)
}
