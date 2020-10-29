package base

// ParsedLineResult est le résultat du parsing d'une ligne.
type ParsedLineResult struct {
	Tuples []Tuple
	Errors []error // TODO: utiliser CriticError ici
}

// AddTuple permet au parseur d'ajouter un tuple extrait depuis la ligne en cours.
func (res *ParsedLineResult) AddTuple(tuple Tuple) {
	if tuple != nil {
		res.Tuples = append(res.Tuples, tuple)
	}
}

// AddError permet au parseur d'ajouter un tuple extrait depuis la ligne en cours.
func (res *ParsedLineResult) AddError(err error) { // TODO: utiliser CriticError ici
	if err != nil {
		res.Errors = append(res.Errors, err)
	}
}

// Tuple unité de donnée à insérer dans un type
type Tuple interface {
	Key() string
	Scope() string
	Type() string
}
