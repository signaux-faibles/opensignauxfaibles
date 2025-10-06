package filter

import "opensignauxfaibles/lib/sfregexp"

const sirenLength = 9

// sirenFilter is a simple SirenFilter implementation
type sirenFilter map[string]bool

// ShouldSkip retourne `true` si le numéro SIREN/SIRET est hors périmètre.
// Tout siret ou siren incorrect est hors périmètre.
// Si aucun filtre n'est défini ("nil"), filtre uniquement les siret / siren
// incorrects.
func (f sirenFilter) ShouldSkip(siretOrSiren string) bool {
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
