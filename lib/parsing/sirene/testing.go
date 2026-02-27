package sirene

import (
	"fmt"
	"sort"
	"strings"
)

// SireneEntry represents one establishment to include in a generated Sirene CSV.
type SireneEntry struct {
	Siret    string
	Siege    bool // whether this establishment is the head office
	Etranger bool // if true, no valid French code commune is set (departement will be empty)
}

// sireneColumns is the ordered list of CSV columns matching the real Sirene file format.
var sireneColumns = []string{
	"siren",
	"nic",
	"siret",
	"statutDiffusionEtablissement",
	"dateCreationEtablissement",
	"trancheEffectifsEtablissement",
	"anneeEffectifsEtablissement",
	"activitePrincipaleRegistreMetiersEtablissement",
	"dateDernierTraitementEtablissement",
	"etablissementSiege",
	"nombrePeriodesEtablissement",
	"complementAdresseEtablissement",
	"numeroVoieEtablissement",
	"indiceRepetitionEtablissement",
	"typeVoieEtablissement",
	"libelleVoieEtablissement",
	"codePostalEtablissement",
	"libelleCommuneEtablissement",
	"libelleCommuneEtrangerEtablissement",
	"distributionSpecialeEtablissement",
	"codeCommuneEtablissement",
	"codeCedexEtablissement",
	"libelleCedexEtablissement",
	"codePaysEtrangerEtablissement",
	"libellePaysEtrangerEtablissement",
	"complementAdresse2Etablissement",
	"numeroVoie2Etablissement",
	"indiceRepetition2Etablissement",
	"typeVoie2Etablissement",
	"libelleVoie2Etablissement",
	"codePostal2Etablissement",
	"libelleCommune2Etablissement",
	"libelleCommuneEtranger2Etablissement",
	"distributionSpeciale2Etablissement",
	"codeCommune2Etablissement",
	"codeCedex2Etablissement",
	"libelleCedex2Etablissement",
	"codePaysEtranger2Etablissement",
	"libellePaysEtranger2Etablissement",
	"dateDebut",
	"etatAdministratifEtablissement",
	"enseigne1Etablissement",
	"enseigne2Etablissement",
	"enseigne3Etablissement",
	"denominationUsuelleEtablissement",
	"activitePrincipaleEtablissement",
	"nomenclatureActivitePrincipaleEtablissement",
	"caractereEmployeurEtablissement",
	"longitude",
	"latitude",
	"geo_score",
	"geo_type",
	"geo_adresse",
	"geo_id",
	"geo_ligne",
	"geo_l4",
	"geo_l5",
}

// MakeSireneCSV generates a CSV string simulating the raw content of a Sirene
// (établissements) file. Entries are sorted by SIRET for deterministic output.
//
//   - Siege controls the etablissementSiege column.
//   - When Etranger is false, a valid French code commune "75101" is set so that
//     departement "75" is extracted; when true, no valid code commune is set and
//     departement will be empty.
func MakeSireneCSV(entries []SireneEntry) string {
	sorted := make([]SireneEntry, len(entries))
	copy(sorted, entries)
	sort.Slice(sorted, func(i, j int) bool {
		return sorted[i].Siret < sorted[j].Siret
	})

	var sb strings.Builder
	sb.WriteString(strings.Join(sireneColumns, ","))
	sb.WriteString("\n")
	for _, e := range sorted {
		sb.WriteString(sireneEntryToRow(e))
		sb.WriteString("\n")
	}
	return sb.String()
}

// makeRow builds a CSV row in sireneColumns order from the given field map.
// Returns an error if a key in fields does not correspond to any known column.
func makeRow(fields map[string]string) (string, error) {
	colSet := make(map[string]struct{}, len(sireneColumns))
	for _, col := range sireneColumns {
		colSet[col] = struct{}{}
	}
	for k := range fields {
		if _, ok := colSet[k]; !ok {
			return "", fmt.Errorf("unknown column %q", k)
		}
	}
	row := make([]string, len(sireneColumns))
	for i, col := range sireneColumns {
		row[i] = fields[col]
	}
	return strings.Join(row, ","), nil
}

func sireneEntryToRow(e SireneEntry) string {
	row := map[string]string{}

	if len(e.Siret) == 14 {
		row["siren"] = e.Siret[:9]
		row["nic"] = e.Siret[9:]
		row["siret"] = e.Siret
	}

	if e.Siege {
		row["etablissementSiege"] = "true"
	} else {
		row["etablissementSiege"] = "false"
	}

	if !e.Etranger {
		row["codePostalEtablissement"] = "75001"
		row["codeCommuneEtablissement"] = "75101" // departement "75"
	}

	row["etatAdministratifEtablissement"] = "A"

	fields := make([]string, len(sireneColumns))
	for i, col := range sireneColumns {
		fields[i] = row[col]
	}
	return strings.Join(fields, ",")
}
