package main

import (
	"errors"

	"github.com/signaux-faibles/opensignauxfaibles/lib/apconso"
	"github.com/signaux-faibles/opensignauxfaibles/lib/apdemande"
	"github.com/signaux-faibles/opensignauxfaibles/lib/base"
	"github.com/signaux-faibles/opensignauxfaibles/lib/bdf"
	"github.com/signaux-faibles/opensignauxfaibles/lib/diane"
	"github.com/signaux-faibles/opensignauxfaibles/lib/ellisphere"
	"github.com/signaux-faibles/opensignauxfaibles/lib/engine"
	"github.com/signaux-faibles/opensignauxfaibles/lib/marshal"
	"github.com/signaux-faibles/opensignauxfaibles/lib/sirene"
	"github.com/signaux-faibles/opensignauxfaibles/lib/urssaf"

	sireneul "github.com/signaux-faibles/opensignauxfaibles/lib/sirene_ul"

	"github.com/globalsign/mgo/bson"
)

type purgeBatchParams struct {
	FromBatchKey           string `json:"fromBatch"`
	Key                    string `json:"debugForKey"`
	IUnderstandWhatImDoing bool   `json:"IUnderstandWhatImDoing"`
}

func purgeBatchHandler(params purgeBatchParams) error {

	if params.FromBatchKey == "" {
		return errors.New("paramètre `fromBatch` obligatoire")
	}

	var batch base.AdminBatch
	err := engine.Load(&batch, params.FromBatchKey)
	if err != nil {
		return errors.New("le batch " + params.FromBatchKey + " n'est pas accessible: " + err.Error())
	}

	if params.Key != "" {
		err = engine.PurgeBatchOne(batch, params.Key)
		if err != nil {
			return errors.New("erreur pendant le MapReduce: " + err.Error())
		}
	} else {
		if !params.IUnderstandWhatImDoing {
			return errors.New("pour une purge de la base complète, IUnderstandWhatImDoing doit être `true`")
		}
		err = engine.PurgeBatch(batch)
		if err != nil {
			return errors.New("(✖╭╮✖) le traitement n'a pas abouti: " + err.Error())
		}
	}
	return nil
}

type importBatchParams struct {
	BatchKey string   `json:"batch"`
	Parsers  []string `json:"parsers"`
	NoFilter bool     `json:"noFilter"`
}

// importBatchHandler traite les demandes d'import par l'API
// on peut demander l'exécution de tous les parsers sans fournir d'option
// ou demander l'exécution de parsers particuliers en fournissant une liste de leurs codes.
func importBatchHandler(params importBatchParams) error {
	batch := base.AdminBatch{}
	err := engine.Load(&batch, params.BatchKey)
	if err != nil {
		return errors.New("Batch inexistant: " + err.Error())
	}

	parsers, err := resolveParsers(params.Parsers)
	if err != nil {
		return err
	}
	err = engine.ImportBatch(batch, parsers, params.NoFilter)
	if err != nil {
		return err
	}
	return nil
}

type checkBatchParams struct {
	BatchKey string   `json:"batch"`
	Parsers  []string `json:"parsers"`
}

func checkBatchHandler(params checkBatchParams) (bson.M, error) {

	batch := base.AdminBatch{}
	err := engine.Load(&batch, params.BatchKey)
	if err != nil {
		return bson.M{}, errors.New("Batch inexistant: " + err.Error())
	}

	parsers, err := resolveParsers(params.Parsers)
	if err != nil {
		return bson.M{}, err
	}

	reports, err := engine.CheckBatch(batch, parsers)
	if err != nil {
		return bson.M{}, errors.New("Erreurs détectées: " + err.Error())
	}

	return bson.M{"reports": reports}, nil
}

func purgeNotCompactedHandler() error {
	return engine.PurgeNotCompacted()
}

type pruneEntitiesParams struct {
	BatchKey string `json:"batch"`
	Delete   bool   `json:"delete"`
}

// Count – then delete – companies from RawData that should have been
// excluded by the SIREN Filter.
func pruneEntitiesHandler(params pruneEntitiesParams) (bson.M, error) {
	count, err := engine.PruneEntities(params.BatchKey, params.Delete)
	return bson.M{"count": count}, err
}

// RegisteredParsers liste des parsers disponibles
var registeredParsers = map[string]marshal.Parser{
	"debit":        urssaf.ParserDebit,
	"ccsf":         urssaf.ParserCCSF,
	"cotisation":   urssaf.ParserCotisation,
	"admin_urssaf": urssaf.ParserCompte,
	"delai":        urssaf.ParserDelai,
	"effectif":     urssaf.ParserEffectif,
	"effectif_ent": urssaf.ParserEffectifEnt,
	"procol":       urssaf.ParserProcol,
	"apconso":      apconso.Parser,
	"apdemande":    apdemande.Parser,
	"bdf":          bdf.Parser,
	"sirene":       sirene.Parser,
	"sirene_ul":    sireneul.Parser,
	"diane":        diane.Parser,
	"ellisphere":   ellisphere.Parser,
}

// Vérifie et charge les parsers
func resolveParsers(parserNames []string) ([]marshal.Parser, error) {
	var parsers []marshal.Parser
	if parserNames == nil {
		for _, fileParser := range registeredParsers {
			parsers = append(parsers, fileParser)
		}
	} else {
		for _, fileType := range parserNames {
			if fileParser, ok := registeredParsers[fileType]; ok {
				parsers = append(parsers, fileParser)
			} else {
				return parsers, errors.New(fileType + " n'est pas un parser reconnu.")
			}
		}
	}
	return parsers, nil
}
