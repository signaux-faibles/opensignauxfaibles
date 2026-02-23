package sireneul

import (
	"fmt"
	"sort"
	"strings"
)

// SireneULEntry represents one legal unit (unité légale) to include in a generated SireneUL CSV.
type SireneULEntry struct {
	Siren              string // 9-char SIREN
	APE                string // activité principale (e.g. "47.11A"); empty if unknown
	CategorieJuridique string // structure juridique (e.g. "5710")
}

// sireneULColumns is the ordered list of CSV columns matching the real SireneUL file format.
var sireneULColumns = []string{
	"siren",
	"statutDiffusionUniteLegale",
	"unitePurgeeUniteLegale",
	"dateCreationUniteLegale",
	"sigleUniteLegale",
	"sexeUniteLegale",
	"prenom1UniteLegale",
	"prenom2UniteLegale",
	"prenom3UniteLegale",
	"prenom4UniteLegale",
	"prenomUsuelUniteLegale",
	"pseudonymeUniteLegale",
	"identifiantAssociationUniteLegale",
	"trancheEffectifsUniteLegale",
	"anneeEffectifsUniteLegale",
	"dateDernierTraitementUniteLegale",
	"nombrePeriodesUniteLegale",
	"categorieEntreprise",
	"anneeCategorieEntreprise",
	"dateDebut",
	"etatAdministratifUniteLegale",
	"nomUniteLegale",
	"nomUsageUniteLegale",
	"denominationUniteLegale",
	"denominationUsuelle1UniteLegale",
	"denominationUsuelle2UniteLegale",
	"denominationUsuelle3UniteLegale",
	"categorieJuridiqueUniteLegale",
	"activitePrincipaleUniteLegale",
	"nomenclatureActivitePrincipaleUniteLegale",
	"nicSiegeUniteLegale",
	"economieSocialeSolidaireUniteLegale",
	"caractereEmployeurUniteLegale",
}

// MakeSireneULCSV generates a CSV string simulating the raw content of a SireneUL
// (unités légales) file. Entries are sorted by SIREN for deterministic output.
//
// For each entry:
//   - Siren identifies the legal unit.
//   - APE sets the activitePrincipaleUniteLegale column (nomenclature NAFRev2 is
//     set automatically when APE is non-empty).
//   - CategorieJuridique sets the categorieJuridiqueUniteLegale column.
func MakeSireneULCSV(entries []SireneULEntry) string {
	sorted := make([]SireneULEntry, len(entries))
	copy(sorted, entries)
	sort.Slice(sorted, func(i, j int) bool {
		return sorted[i].Siren < sorted[j].Siren
	})

	var sb strings.Builder
	sb.WriteString(strings.Join(sireneULColumns, ","))
	sb.WriteString("\n")
	for _, e := range sorted {
		sb.WriteString(sireneULEntryToRow(e))
		sb.WriteString("\n")
	}
	return sb.String()
}

// makeRow builds a CSV row in sireneULColumns order from the given field map.
// Returns an error if a key in fields does not correspond to any known column.
func makeRow(fields map[string]string) (string, error) {
	colSet := make(map[string]struct{}, len(sireneULColumns))
	for _, col := range sireneULColumns {
		colSet[col] = struct{}{}
	}
	for k := range fields {
		if _, ok := colSet[k]; !ok {
			return "", fmt.Errorf("unknown column %q", k)
		}
	}
	row := make([]string, len(sireneULColumns))
	for i, col := range sireneULColumns {
		row[i] = fields[col]
	}
	return strings.Join(row, ","), nil
}

func sireneULEntryToRow(e SireneULEntry) string {
	row := map[string]string{}

	row["siren"] = e.Siren
	row["etatAdministratifUniteLegale"] = "A"
	row["categorieJuridiqueUniteLegale"] = e.CategorieJuridique
	row["activitePrincipaleUniteLegale"] = e.APE
	if e.APE != "" {
		row["nomenclatureActivitePrincipaleUniteLegale"] = "NAFRev2"
	}

	fields := make([]string, len(sireneULColumns))
	for i, col := range sireneULColumns {
		fields[i] = row[col]
	}
	return strings.Join(fields, ",")
}
