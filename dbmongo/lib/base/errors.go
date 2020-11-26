package base

// CriticityError is a generic error with criticity
type CriticityError interface {
	error
	Criticity() string
}

type criticError struct {
	err           error
	isFilterError bool
}

// Criticity returns the type of error
func (pe *criticError) Criticity() string {
	if pe == nil {
		return ""
	} else if pe.isFilterError {
		return "filter"
	}
	return "error"
}

func (pe *criticError) Error() string {
	return pe.err.Error()
}

func newCriticError(err error, isFilterError bool) CriticityError {
	if err == nil {
		return nil
	}
	return &criticError{err, isFilterError}
}

// NewFilterError returns a filter error (occurs when something goes wrong while filtering)
func NewFilterError(err error) CriticityError {
	return newCriticError(err, true)
}

// NewRegularError creates a regular error
func NewRegularError(err error) CriticityError {
	return newCriticError(err, false)
}
