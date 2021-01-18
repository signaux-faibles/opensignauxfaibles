package bdf

import (
	"bufio"
	"encoding/csv"
	"io"
	"os"
	"strings"
	"time"

	"github.com/signaux-faibles/opensignauxfaibles/lib/base"
	"github.com/signaux-faibles/opensignauxfaibles/lib/marshal"
	"github.com/signaux-faibles/opensignauxfaibles/lib/misc"
)

// BDF Information Banque de France
type BDF struct {
	Siren               string    `col:"D1"              json:"siren"                 bson:"siren"`
	Annee               *int      `col:"ANNEE"           json:"annee_bdf"             bson:"annee_bdf"`
	ArreteBilan         time.Time `col:"ARRETE_BILAN"    json:"arrete_bilan_bdf"      bson:"arrete_bilan_bdf"`
	RaisonSociale       string    `col:"DENOM"           json:"raison_sociale"        bson:"raison_sociale"`
	Secteur             string    `col:"SECTEUR"         json:"secteur"               bson:"secteur"`
	PoidsFrng           *float64  `col:"POIDS_FRNG"      json:"poids_frng"            bson:"poids_frng"`
	TauxMarge           *float64  `col:"TX_MARGE"        json:"taux_marge"            bson:"taux_marge"`
	DelaiFournisseur    *float64  `col:"DELAI_FRS"       json:"delai_fournisseur"     bson:"delai_fournisseur"`
	DetteFiscale        *float64  `col:"POIDS_DFISC_SOC" json:"dette_fiscale"         bson:"dette_fiscale"`
	FinancierCourtTerme *float64  `col:"POIDS_FIN_CT"    json:"financier_court_terme" bson:"financier_court_terme"`
	FraisFinancier      *float64  `col:"POIDS_FRAIS_FIN" json:"frais_financier"       bson:"frais_financier"`
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

// Parser fournit une instance utilisable par ParseFilesFromBatch.
var Parser = &bdfParser{}

type bdfParser struct {
	file   *os.File
	reader *csv.Reader
	idx    marshal.ColMapping
}

func (parser *bdfParser) GetFileType() string {
	return "bdf"
}

func (parser *bdfParser) Init(cache *marshal.Cache, batch *base.AdminBatch) error {
	return nil
}

func (parser *bdfParser) Open(filePath string) (err error) {
	parser.file, parser.reader, err = openFile(filePath)
	if err == nil {
		parser.idx, err = marshal.IndexColumnsFromCsvHeader(parser.reader, BDF{})
	}
	return err
}

func (parser *bdfParser) Close() error {
	return parser.file.Close()
}

func openFile(filePath string) (*os.File, *csv.Reader, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return file, nil, err
	}
	reader := csv.NewReader(bufio.NewReader(file))
	reader.Comma = ';'
	reader.LazyQuotes = true
	return file, reader, nil
}

func (parser *bdfParser) ParseLines(parsedLineChan chan marshal.ParsedLineResult) {
	for {
		parsedLine := marshal.ParsedLineResult{}
		row, err := parser.reader.Read()
		if err == io.EOF {
			close(parsedLineChan)
			break
		} else if err != nil {
			parsedLine.AddRegularError(err)
		} else {
			parseBdfLine(row, parser.idx, &parsedLine)
			if len(parsedLine.Errors) > 0 {
				parsedLine.Tuples = []marshal.Tuple{}
			}
		}
		parsedLineChan <- parsedLine
	}
}

func parseBdfLine(row []string, field marshal.ColMapping, parsedLine *marshal.ParsedLineResult) {
	var err error
	bdf := BDF{}
	bdf.Siren = strings.Replace(row[field["D1"]], " ", "", -1)
	bdf.Annee, err = misc.ParsePInt(row[field["ANNEE"]])
	parsedLine.AddRegularError(err)
	var arrete = row[field["ARRETE_BILAN"]]
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
	parsedLine.AddRegularError(err)
	bdf.RaisonSociale = row[field["DENOM"]]
	bdf.Secteur = row[field["SECTEUR"]]
	if len(row) > field["POIDS_FRNG"] {
		bdf.PoidsFrng, err = misc.ParsePFloat(row[field["POIDS_FRNG"]])
		parsedLine.AddRegularError(err)
	} else {
		bdf.PoidsFrng = nil
	}
	if len(row) > field["TX_MARGE"] {
		bdf.TauxMarge, err = misc.ParsePFloat(row[field["TX_MARGE"]])
		parsedLine.AddRegularError(err)
	} else {
		bdf.TauxMarge = nil
	}
	if len(row) > field["DELAI_FRS"] {
		bdf.DelaiFournisseur, err = misc.ParsePFloat(row[field["DELAI_FRS"]])
		parsedLine.AddRegularError(err)
	} else {
		bdf.DelaiFournisseur = nil
	}
	if len(row) > field["POIDS_DFISC_SOC"] {
		bdf.DetteFiscale, err = misc.ParsePFloat(row[field["POIDS_DFISC_SOC"]])
		parsedLine.AddRegularError(err)
	} else {
		bdf.DetteFiscale = nil
	}
	if len(row) > field["POIDS_FIN_CT"] {
		bdf.FinancierCourtTerme, err = misc.ParsePFloat(row[field["POIDS_FIN_CT"]])
		parsedLine.AddRegularError(err)
	} else {
		bdf.FinancierCourtTerme = nil
	}
	if len(row) > field["POIDS_FRAIS_FIN"] {
		bdf.FraisFinancier, err = misc.ParsePFloat(row[field["POIDS_FRAIS_FIN"]])
		parsedLine.AddRegularError(err)
	} else {
		bdf.FraisFinancier = nil
	}
	parsedLine.AddTuple(bdf)
}
