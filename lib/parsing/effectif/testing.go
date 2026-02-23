package effectif

import (
	"fmt"
	"sort"
	"strings"
	"time"
)

// MakeEffectifEntCSV generates a CSV string simulating the raw content of an effectif_ent file.
// `periods` defines the list of monthly periods to include as columns.
// `effectifs` maps each SIREN to a slice of effectif values, one per period.
//
// Use MakeEffectifEntCSVWithMissing if missing data needs to be included.
func MakeEffectifEntCSV(periods []time.Time, effectifs map[string][]int) string {
	withMissing := make(map[string][]*int, len(effectifs))
	for siren, values := range effectifs {
		ptrs := make([]*int, len(values))
		for i := range values {
			v := values[i]
			ptrs[i] = &v
		}
		withMissing[siren] = ptrs
	}
	return MakeEffectifEntCSVWithMissing(periods, withMissing)
}

// MakeEffectifEntCSVWithMissing generates a CSV string simulating the raw content of an effectif_ent file.
// `periods` defines the list of monthly periods to include as columns.
// `effectifs` maps each SIREN to a slice of nullable effectif values, one per period.
// A nil pointer in the slice represents an empty (absent) value in the CSV.
func MakeEffectifEntCSVWithMissing(periods []time.Time, effectifs map[string][]*int) string {
	var sb strings.Builder

	// Header
	header := []string{"siren"}
	for _, p := range periods {
		header = append(header, periodToEffColName(p))
	}
	header = append(header, "raison_soc")

	sb.WriteString(strings.Join(header, ";"))
	sb.WriteString("\n")

	// Rows — sorted by SIREN for deterministic output
	sirens := make([]string, 0, len(effectifs))
	for siren := range effectifs {
		sirens = append(sirens, siren)
	}
	sort.Strings(sirens)

	for _, siren := range sirens {
		values := effectifs[siren]
		fields := []string{siren}

		for i := range periods {
			val := ""
			if i < len(values) {
				if values[i] != nil {
					val = fmt.Sprintf("%d", *values[i])
				} else {
					val = ""
				}
			}
			fields = append(fields, val)
		}
		fields = append(fields, "MON ENTREPRISE") // raison_soc placeholder
		sb.WriteString(strings.Join(fields, ";"))
		sb.WriteString("\n")
	}

	return sb.String()
}

// periodToEffColName converts a date to an effectif column name in the eff{YYYYQM} format.
// Example: 2010-01-01 → "eff201011", 2010-04-01 → "eff201021".
func periodToEffColName(t time.Time) string {
	year := t.Year()
	month := int(t.Month())
	quarter := (month-1)/3 + 1
	monthOfQuarter := (month-1)%3 + 1
	return fmt.Sprintf("eff%04d%d%d", year, quarter, monthOfQuarter)
}
