package filter

const sirenLength = 9

// sirenFilter is a simple SirenFilter implementation
type sirenFilter map[string]bool

// ShouldSkip retourne `true` si le numéro SIREN/SIRET est hors périmètre.
// Si aucun filtre n'est défini, renvoi `false` par défaut (aucun filtrage).
func (f sirenFilter) ShouldSkip(siretOrSiren string) bool {
	if f == nil {
		return false
	}

	siren := siretOrSiren

	if len(siretOrSiren) >= sirenLength {
		siren = siretOrSiren[:sirenLength]
	}

	return !f[siren]
}
