package filter

import "opensignauxfaibles/lib/sfregexp"

const sirenLength = 9

// MapFilter is a simple MapFilter implementation
type MapFilter map[string]bool

// ShouldSkip retourne `true` si le numéro SIREN/SIRET est hors périmètre.
// Tout siret ou siren incorrect est hors périmètre.
// Si aucun filtre n'est défini ("nil"), filtre uniquement les siret / siren
// incorrects.
func (f MapFilter) ShouldSkip(siretOrSiren string) bool {
	sirenRe := sfregexp.RegexpDict["siren"]
	siretRe := sfregexp.RegexpDict["siret"]
	if !sirenRe.MatchString(siretOrSiren) && !siretRe.MatchString(siretOrSiren) {
		// siret / siren invalide
		return true
	}

	if f == nil {
		return false
	}

	siren := siretOrSiren

	if len(siretOrSiren) >= sirenLength {
		siren = siretOrSiren[:sirenLength]
	}

	return !f[siren]
}

// All returns all SIRENs in the filter as a set.
func (f MapFilter) All() map[string]struct{} {
	result := make(map[string]struct{}, len(f))
	for siren := range f {
		if f[siren] { // Only include true values
			result[siren] = struct{}{}
		}
	}
	return result
}
