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

// FilterError occurs when something goes wrong while filtering
type FilterError struct {
	*CriticError
}

// NewFilterNotice returns a filter error
func NewFilterNotice() error {
	return NewFilterError(errors.New("ligne filtrée"))
}

// NewFilterError returns a filter error
func NewFilterError(err error) error {
	if err == nil {
		return nil
	}
	return NewCriticError(err, "filter")
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
