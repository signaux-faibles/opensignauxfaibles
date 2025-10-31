package filter

import (
	"bufio"
	"compress/gzip"
	"encoding/csv"
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"regexp"
	"strconv"
	"strings"
)

// Usage: $ ./create_filter --path testData/test_data.csv

// DefaultNbMois is the default number of the most recent months during which the effectif of the company must reach the threshold.
const DefaultNbMois = 100

// DefaultMinEffectif is the default effectif threshold, expressed in number of employees.
const DefaultMinEffectif = 10

// DefaultNbIgnoredCols is the default number of rightmost columns that don't contain effectif data.
const DefaultNbIgnoredCols = 2

// NbLeadingColsToSkip is the number of leftmost columns that don't contain effectif data.
const NbLeadingColsToSkip = 5 // column names: "compte", "siret", "rais_soc", "ape_ins" and "dep"

type filter func(string) bool

// Implementation of the create_filter command.
func main() {

	var path = flag.String("path", "", "Chemin d'accès au fichier effectif")
	var nbMois = flag.Int(
		"nbMois",
		DefaultNbMois,
		"Nombre de mois observés (avec effectif connu) pour déterminer si l'entreprise dépasse 10 salariés",
	)
	var minEffectif = flag.Int(
		"minEffectif",
		DefaultMinEffectif,
		"Si une entreprise atteint ou dépasse 'minEffectif' dans les 'nbMois' derniers mois, elle est inclue dans le périmètre du filtre.",
	)
	var nIgnoredCols = flag.Int(
		"nIgnoredCols",
		DefaultNbIgnoredCols,
		"Nombre de colonnes à ignorer à la fin du fichier effectif",
	)
	flag.Parse()

	// create filter
	err := Create(os.Stdout, *path, *nbMois, *minEffectif, *nIgnoredCols)
	if err != nil {
		log.Panic(err)
	}
}

// Create generates a "filter" from an "effectif" file.
func Create(writer io.Writer, effectifFileName string, nbMois, minEffectif int, nIgnoredCols int, filters ...filter) error {
	last := guessLastNMissing(effectifFileName, nIgnoredCols)
	r, f, err := makeEffectifReaderFromFile(effectifFileName)
	if err != nil {
		return err
	}

	perimeter := getInitialPerimeter(r, nbMois, minEffectif, nIgnoredCols+last)

	for _, filter := range filters {
		perimeter = applyFilter(perimeter, filter)
	}

	fmt.Fprintln(writer, "siren")
	for siren := range perimeter {
		fmt.Fprintln(writer, siren)
	}
	return f.Close()
}

func applyFilter(perimeter map[string]struct{}, f filter) map[string]struct{} {
	newPerimeter := make(map[string]struct{})
	for siren := range perimeter {
		if f(siren) {
			newPerimeter[siren] = struct{}{}
		}
	}
	return newPerimeter
}

// If the effectif file has a ".gz" suffix, it will be decompressed on the fly.
func makeEffectifReaderFromFile(effectifFileName string) (*csv.Reader, *os.File, error) {
	var fileReader io.Reader
	compressed := strings.HasSuffix(effectifFileName, "csv.gz")

	file, err := os.Open(effectifFileName)
	if err != nil {
		return nil, nil, err
	}
	if compressed {
		fileReader, err = gzip.NewReader(file)
		if err != nil {
			return nil, nil, err
		}
	} else {
		fileReader = bufio.NewReader(file)
	}
	return initializeEffectifReader(fileReader), file, err
}

func initializeEffectifReader(reader io.Reader) *csv.Reader {
	r := csv.NewReader(reader)
	r.LazyQuotes = true
	r.Comma = ';'
	return r
}

// getInitialPerimeter makes a perimeter on effectif criterias alone
func getInitialPerimeter(r *csv.Reader, nbMois, minEffectif, nIgnoredCols int) map[string]struct{} {
	detectedSirens := map[string]struct{}{} // smaller memory footprint than map[string]bool
	_, err := r.Read()                      // en tête
	if err != nil {
		log.Panic(err)
	}
	lineNumber, skippedLines := 0, 0
	for {
		lineNumber++
		record, err := r.Read()

		// Stop at EOF.
		if err == io.EOF {
			break
		} else if err != nil {
			log.Panic(err)
		}
		siret := record[1]
		shouldKeep := len(siret) == 14 &&
			isInsidePerimeter(record[NbLeadingColsToSkip:len(record)-nIgnoredCols], nbMois, minEffectif)

		var siren string
		if len(siret) >= 9 {
			siren = siret[0:9] // trim siret into a siren
			_, alreadyDetected := detectedSirens[siren]
			if shouldKeep && !alreadyDetected {
				detectedSirens[siren] = struct{}{}
			}
		} else {
			skippedLines++
			fmt.Printf("%d digit siret encountered, skipping line %d \n", len(siret), lineNumber)
		}
	}
	if skippedLines > 0 {
		fmt.Printf("%d lines with bad siret/siren skipped :( \n", skippedLines)
	}
	return detectedSirens
}

func isInsidePerimeter(record []string, nbMois, minEffectif int) bool {
	for i := len(record) - 1; i >= len(record)-nbMois && i >= 0; i-- {
		if record[i] == "" {
			continue
		}
		reg, err := regexp.Compile("[^0-9]")
		if err != nil {
			log.Fatal(err)
		}
		effectif, err := strconv.Atoi(reg.ReplaceAllString(record[i], ""))
		if err != nil {
			fmt.Println(record)
			log.Panic(err)
		}
		if effectif >= minEffectif {
			return true
		}
	}
	return false
}

func guessLastNMissing(path string, nIgnoredCols int) int {
	r, f, err := makeEffectifReaderFromFile(path)
	if err != nil {
		log.Panic(err)
	}
	defer f.Close()
	if _, err = r.Read(); err != nil { // en tête
		log.Panic(err)
	}
	return guessLastNMissingFromReader(r, nIgnoredCols)
}

// guessLastNMissingFromReader returns the number of rightmost columns
// (on top of nIgnoredCols columns) that never have a value.
func guessLastNMissingFromReader(r *csv.Reader, nIgnoredCols int) int {
	var lastConsideredCol int // index of the rightmost column of the last read row
	lastColWithValue := -1    // index of the rightmost column which had a value at least once
	for {
		record, err := r.Read()
		if err == io.EOF {
			break
		} else if err != nil {
			log.Panic(err)
		}
		lastConsideredCol = len(record) - 1 - nIgnoredCols
		for i := lastConsideredCol; i > lastColWithValue; i-- {
			if record[i] != "" {
				lastColWithValue = i
			}
		}
	}
	return lastConsideredCol - lastColWithValue
}
