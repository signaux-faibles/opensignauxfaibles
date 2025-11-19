// Package sirenehisto exporte un parser pour extraire les données d'un fichier
// Sirene Historique, qui contient les informations des changements de statuts
// passés des établissements
package sirenehisto

import (
	"fmt"
	"io"
	"time"

	"opensignauxfaibles/lib/engine"
	"opensignauxfaibles/lib/parsing"
)

type SireneHistoParser struct{}

func NewSireneHistoParser() engine.Parser {
	return &SireneHistoParser{}
}

func (p *SireneHistoParser) Type() engine.ParserType { return engine.SireneHisto }
func (p *SireneHistoParser) New(r io.Reader) engine.ParserInst {
	return &parsing.CsvParserInst{
		Reader:     r,
		RowParser:  &sireneHistoRowParser{},
		Comma:      ',',
		LazyQuotes: false,
		DestTuple:  SireneHisto{},
	}
}

type sireneHistoRowParser struct{}

func (rp *sireneHistoRowParser) ParseRow(row []string, res *engine.ParsedLineResult, idx parsing.ColIndex) {
	var err error

	idxRow := idx.IndexRow(row)
	sirenehisto := SireneHisto{}
	sirenehisto.Siret = idxRow.GetVal("siret")

	// Malformed dates are errors
	// Missing both dates is an error
	// Missing only one is OK
	dateDebutStr := idxRow.GetVal("dateDebut")
	dateDebut, err := time.Parse("2006-01-02", dateDebutStr)
	if (dateDebutStr != "") && err != nil {
		res.AddRegularError(err)
		return
	} else if dateDebutStr == "" {
		sirenehisto.DateDebut = nil
	} else {
		sirenehisto.DateDebut = &dateDebut
	}

	dateFinStr := idxRow.GetVal("dateFin")
	dateFin, err := time.Parse("2006-01-02", dateFinStr)
	if (dateFinStr != "") && err != nil {
		res.AddRegularError(err)
		return
	} else if dateFinStr == "" {
		sirenehisto.DateFin = nil
	} else {
		sirenehisto.DateFin = &dateFin
	}

	if dateDebutStr == "" && dateFinStr == "" {
		res.AddRegularError(fmt.Errorf("dates de début et de fin toutes les deux manquantes"))
		return
	}

	etatAdmin := idxRow.GetVal("etatAdministratifEtablissement")
	switch etatAdmin {
	case "A":
		sirenehisto.EstActif = true
	case "F":
		sirenehisto.EstActif = false
	default:
		res.AddRegularError(fmt.Errorf("état administratif malformé : \"%s\" (attendu \"A\" ou \"F\"", etatAdmin))
		return
	}

	sirenehisto.ChangementStatutActif, err = idxRow.GetBool("changementEtatAdministratifEtablissement")
	if err != nil {
		res.AddRegularError(err)
		return
	}

	res.AddTuple(sirenehisto)
}
