package bdf

import (
	"bufio"
	"encoding/csv"
	"errors"
	"io"
	"os"
	"strings"
	"time"

	"github.com/signaux-faibles/opensignauxfaibles/dbmongo/lib/base"
	"github.com/signaux-faibles/opensignauxfaibles/dbmongo/lib/marshal"
	"github.com/signaux-faibles/opensignauxfaibles/dbmongo/lib/misc"
	"github.com/signaux-faibles/opensignauxfaibles/dbmongo/lib/sfregexp"
)

// BDF Information Banque de France
type BDF struct {
	Siren               string    `json:"siren" bson:"siren"`
	Annee               *int      `json:"annee_bdf" bson:"annee_bdf"`
	ArreteBilan         time.Time `json:"arrete_bilan_bdf" bson:"arrete_bilan_bdf"`
	RaisonSociale       string    `json:"raison_sociale" bson:"raison_sociale"`
	Secteur             string    `json:"secteur" bson:"secteur"`
	PoidsFrng           *float64  `json:"poids_frng" bson:"poids_frng"`
	TauxMarge           *float64  `json:"taux_marge" bson:"taux_marge"`
	DelaiFournisseur    *float64  `json:"delai_fournisseur" bson:"delai_fournisseur"`
	DetteFiscale        *float64  `json:"dette_fiscale" bson:"dette_fiscale"`
	FinancierCourtTerme *float64  `json:"financier_court_terme" bson:"financier_court_terme"`
	FraisFinancier      *float64  `json:"frais_financier" bson:"frais_financier"`
}

// Key id de l'objet
func (bdf BDF) Key() string {
	return bdf.Siren
}

// Type de données
func (bdf BDF) Type() string {
	return "bdf"
}

// Scope de l'objet
func (bdf BDF) Scope() string {
	return "entreprise"
}

// Parser expose le parseur et le type de fichier qu'il supporte.
var Parser = marshal.Parser{FileType: "bdf", FileParser: ParseFile}

// ParseFile permet de lancer le parsing du fichier demandé.
func ParseFile(filePath string, cache *marshal.Cache, batch *base.AdminBatch) marshal.OpenFileResult {
	filter := marshal.GetSirenFilterFromCache(*cache) // TODO: retirer filtre
	file, reader, err := openFile(filePath)
	return marshal.OpenFileResult{
		Error: err,
		ParseLines: func(parsedLineChan chan base.ParsedLineResult) {
			parseLines(reader, &filter, parsedLineChan)
		},
		Close: file.Close,
	}
}

func openFile(filePath string) (*os.File, *csv.Reader, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return nil, nil, err
	}
	reader := csv.NewReader(bufio.NewReader(file))
	reader.Comma = ';'
	reader.LazyQuotes = true
	_, err = reader.Read() // Sauter l'en-tête
	return file, reader, err
}

func parseLines(reader *csv.Reader, filter *marshal.SirenFilter, parsedLineChan chan base.ParsedLineResult) {
	for {
		parsedLine := base.ParsedLineResult{}
		row, err := reader.Read()
		if err == io.EOF {
			close(parsedLineChan)
			break
		} else if err != nil {
			parsedLine.AddError(err)
		} else {
			parseBdfLine(row, *filter, &parsedLine)
			if len(parsedLine.Errors) > 0 {
				parsedLine.Tuples = []base.Tuple{}
			}
		}
		parsedLineChan <- parsedLine
	}
}

var field = map[string]int{
	"siren":               0,
	"année":               1,
	"arrêtéBilan":         2,
	"raisonSociale":       3,
	"secteur":             6,
	"poidsFrng":           7,
	"tauxMarge":           8,
	"delaiFournisseur":    9,
	"detteFiscale":        10,
	"financierCourtTerme": 11,
	"fraisFinancier":      12,
}

func parseBdfLine(row []string, filter marshal.SirenFilter, parsedLine *base.ParsedLineResult) {
	bdf := BDF{}
	bdf.Siren = strings.Replace(row[field["siren"]], " ", "", -1)

	if !sfregexp.ValidSiren(bdf.Siren) {
		parsedLine.AddError(errors.New("siren invalide : " + bdf.Siren))
		return
	}

	var err error
	bdf.Annee, err = misc.ParsePInt(row[field["année"]])
	parsedLine.AddError(err)
	var arrete = row[field["arrêtéBilan"]]
	arrete = strings.Replace(arrete, "janv", "-01-", -1)
	arrete = strings.Replace(arrete, "JAN", "-01-", -1)
	arrete = strings.Replace(arrete, "févr", "-02-", -1)
	arrete = strings.Replace(arrete, "FEB", "-02-", -1)
	arrete = strings.Replace(arrete, "mars", "-03-", -1)
	arrete = strings.Replace(arrete, "MAR", "-03-", -1)
	arrete = strings.Replace(arrete, "avr", "-04-", -1)
	arrete = strings.Replace(arrete, "APR", "-04-", -1)
	arrete = strings.Replace(arrete, "mai", "-05-", -1)
	arrete = strings.Replace(arrete, "MAY", "-05-", -1)
	arrete = strings.Replace(arrete, "juin", "-06-", -1)
	arrete = strings.Replace(arrete, "JUN", "-06-", -1)
	arrete = strings.Replace(arrete, "juil", "-07-", -1)
	arrete = strings.Replace(arrete, "JUL", "-07-", -1)
	arrete = strings.Replace(arrete, "août", "-08-", -1)
	arrete = strings.Replace(arrete, "AUG", "-08-", -1)
	arrete = strings.Replace(arrete, "sept", "-09-", -1)
	arrete = strings.Replace(arrete, "SEP", "-09-", -1)
	arrete = strings.Replace(arrete, "oct", "-10-", -1)
	arrete = strings.Replace(arrete, "OCT", "-10-", -1)
	arrete = strings.Replace(arrete, "nov", "-11-", -1)
	arrete = strings.Replace(arrete, "NOV", "-11-", -1)
	arrete = strings.Replace(arrete, "déc", "-12-", -1)
	arrete = strings.Replace(arrete, "DEC", "-12-", -1)
	bdf.ArreteBilan, err = time.Parse("02-01-2006", arrete)
	parsedLine.AddError(err)
	bdf.RaisonSociale = row[field["raisonSociale"]]
	bdf.Secteur = row[field["secteur"]]
	if len(row) > field["poidsFrng"] {
		bdf.PoidsFrng, err = misc.ParsePFloat(row[field["poidsFrng"]])
		parsedLine.AddError(err)
	} else {
		bdf.PoidsFrng = nil
	}
	if len(row) > field["tauxMarge"] {
		bdf.TauxMarge, err = misc.ParsePFloat(row[field["tauxMarge"]])
		parsedLine.AddError(err)
	} else {
		bdf.TauxMarge = nil
	}
	if len(row) > field["delaiFournisseur"] {
		bdf.DelaiFournisseur, err = misc.ParsePFloat(row[field["delaiFournisseur"]])
		parsedLine.AddError(err)
	} else {
		bdf.DelaiFournisseur = nil
	}
	if len(row) > field["detteFiscale"] {
		bdf.DetteFiscale, err = misc.ParsePFloat(row[field["detteFiscale"]])
		parsedLine.AddError(err)
	} else {
		bdf.DetteFiscale = nil
	}
	if len(row) > field["financierCourtTerme"] {
		bdf.FinancierCourtTerme, err = misc.ParsePFloat(row[field["financierCourtTerme"]])
		parsedLine.AddError(err)
	} else {
		bdf.FinancierCourtTerme = nil
	}
	if len(row) > field["fraisFinancier"] {
		bdf.FraisFinancier, err = misc.ParsePFloat(row[field["fraisFinancier"]])
		parsedLine.AddError(err)
	} else {
		bdf.FraisFinancier = nil
	}
	parsedLine.AddTuple(bdf)
}
