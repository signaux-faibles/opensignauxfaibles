package main

import (
	"errors"

	"github.com/signaux-faibles/opensignauxfaibles/dbmongo/lib/apconso"
	"github.com/signaux-faibles/opensignauxfaibles/dbmongo/lib/apdemande"
	"github.com/signaux-faibles/opensignauxfaibles/dbmongo/lib/base"
	"github.com/signaux-faibles/opensignauxfaibles/dbmongo/lib/bdf"
	"github.com/signaux-faibles/opensignauxfaibles/dbmongo/lib/diane"
	"github.com/signaux-faibles/opensignauxfaibles/dbmongo/lib/ellisphere"
	"github.com/signaux-faibles/opensignauxfaibles/dbmongo/lib/engine"
	"github.com/signaux-faibles/opensignauxfaibles/dbmongo/lib/marshal"
	"github.com/signaux-faibles/opensignauxfaibles/dbmongo/lib/sirene"
	"github.com/signaux-faibles/opensignauxfaibles/dbmongo/lib/urssaf"

	sireneul "github.com/signaux-faibles/opensignauxfaibles/dbmongo/lib/sirene_ul"

	"github.com/gin-gonic/gin"
	"github.com/globalsign/mgo/bson"
)

//
func purgeBatchHandler(c *gin.Context) {
	var params struct {
		FromBatchKey           string `json:"fromBatch"`
		Key                    string `json:"debugForKey"`
		IUnderstandWhatImDoing bool   `json:"IUnderstandWhatImDoing"`
	}

	err := c.ShouldBind(&params)
	if err != nil {
		c.JSON(400, "Requête malformée: "+err.Error())
		return
	}
	if params.FromBatchKey == "" {
		c.JSON(400, "paramètre `fromBatch` obligatoire")
		return
	}

	var batch base.AdminBatch
	err = engine.Load(&batch, params.FromBatchKey)
	if err != nil {
		c.JSON(400, "le batch "+params.FromBatchKey+" n'est pas accessible: "+err.Error())
		return
	}

	if params.Key != "" {
		err = engine.PurgeBatchOne(batch, params.Key)
		if err != nil {
			c.JSON(500, "erreur pendant le MapReduce: "+err.Error())
			return
		}
	} else {
		if !params.IUnderstandWhatImDoing {
			c.JSON(400, "pour une purge de la base complète, IUnderstandWhatImDoing doit être `true`")
			return
		}
		err = engine.PurgeBatch(batch)
		if err != nil {
			c.JSON(500, "(✖╭╮✖) le traitement n'a pas abouti: "+err.Error())
			return
		}
	}
}

// importBatchHandler traite les demandes d'import par l'API
// on peut demander l'exécution de tous les parsers sans fournir d'option
// ou demander l'exécution de parsers particuliers en fournissant une liste de leurs codes.
func importBatchHandler(c *gin.Context) {
	var params struct {
		BatchKey string   `json:"batch"`
		Parsers  []string `json:"parsers"`
		NoFilter bool     `json:"noFilter"`
	}
	err := c.ShouldBind(&params)
	if err != nil {
		c.JSON(400, "Requête malformée: "+err.Error())
		return
	}
	batch := base.AdminBatch{}
	err = engine.Load(&batch, params.BatchKey)
	if err != nil {
		c.JSON(404, "Batch inexistant: "+err.Error())
	}

	parsers, err := resolveParsers(params.Parsers)
	if err != nil {
		c.JSON(404, err.Error())
	}
	err = engine.ImportBatch(batch, parsers, params.NoFilter)
	if err != nil {
		c.JSON(500, err.Error())
	}
}

func checkBatchHandler(c *gin.Context) {
	var params struct {
		BatchKey string   `json:"batch"`
		Parsers  []string `json:"parsers"`
	}
	err := c.ShouldBind(&params)
	if err != nil {
		c.JSON(400, "Requête malformée: "+err.Error())
		return
	}
	batch := base.AdminBatch{}
	err = engine.Load(&batch, params.BatchKey)
	if err != nil {
		c.JSON(404, "Batch inexistant: "+err.Error())
	}

	parsers, err := resolveParsers(params.Parsers)
	if err != nil {
		c.JSON(404, err.Error())
	}

	reports, err := engine.CheckBatch(batch, parsers)
	if err != nil {
		c.JSON(417, "Erreurs détectées: "+err.Error())
	} else {
		c.JSON(200, bson.M{"reports": reports})
	}
}

func purgeNotCompactedHandler(c *gin.Context) {
	var result []interface{}
	err := engine.PurgeNotCompacted()
	if err != nil {
		c.JSON(500, err.Error())
		return
	}
	c.JSON(200, result)
}

// Count – then delete – companies from RawData that should have been
// excluded by the SIREN Filter.
func pruneEntitiesHandler(c *gin.Context) {
	var params struct {
		BatchKey string `json:"batch"`
		// TODO: DeleteEntities bool     `json:"delete"`
	}
	err := c.ShouldBind(&params)
	if err != nil {
		c.JSON(400, "Requête malformée: "+err.Error())
		return
	}
	count, err := engine.PruneEntities(params.BatchKey)
	if err != nil {
		c.JSON(500, err.Error())
		return
	}
	c.JSON(200, bson.M{"count": count})
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
