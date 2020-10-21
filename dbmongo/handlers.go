package main

import (
	"errors"
	"fmt"
	"io"
	"os"
	"sort"

	"github.com/signaux-faibles/opensignauxfaibles/dbmongo/lib/apconso"
	"github.com/signaux-faibles/opensignauxfaibles/dbmongo/lib/apdemande"
	"github.com/signaux-faibles/opensignauxfaibles/dbmongo/lib/base"
	"github.com/signaux-faibles/opensignauxfaibles/dbmongo/lib/bdf"
	"github.com/signaux-faibles/opensignauxfaibles/dbmongo/lib/diane"
	"github.com/signaux-faibles/opensignauxfaibles/dbmongo/lib/ellisphere"
	"github.com/signaux-faibles/opensignauxfaibles/dbmongo/lib/engine"
	"github.com/signaux-faibles/opensignauxfaibles/dbmongo/lib/files"
	"github.com/signaux-faibles/opensignauxfaibles/dbmongo/lib/marshal"
	"github.com/signaux-faibles/opensignauxfaibles/dbmongo/lib/sirene"
	"github.com/signaux-faibles/opensignauxfaibles/dbmongo/lib/urssaf"

	sireneul "github.com/signaux-faibles/opensignauxfaibles/dbmongo/lib/sirene_ul"

	"github.com/davecgh/go-spew/spew"
	"github.com/gin-gonic/gin"
	"github.com/globalsign/mgo/bson"
	"github.com/spf13/viper"
)

//
func adminFeature(c *gin.Context) {
	c.JSON(200, []string{"algo_avec_urssaf", "algo_sans_urssaf"})
}

//
func listTypesHandler(c *gin.Context) {
	c.JSON(200, engine.GetTypes())
}

//
func addFile(c *gin.Context) {
	file, err := c.FormFile("file")
	if err != nil {
		c.JSON(500, err.Error())
		return
	}
	batch := c.Request.FormValue("batch")
	fileType := c.Request.FormValue("type")

	source, err := file.Open()
	if err != nil {
		c.JSON(500, err.Error())
		return
	}

	defer source.Close()

	os.MkdirAll(viper.GetString("APP_DATA")+"/"+batch+"/"+fileType+"/", os.ModePerm)
	destination, err := os.Create(viper.GetString("APP_DATA") + "/" + batch + "/" + fileType + "/" + file.Filename)
	if err != nil {
		c.JSON(500, err.Error())
		return
	}
	defer destination.Close()

	_, err = io.Copy(destination, source)

	if err != nil {
		c.JSON(500, err.Error())
		return
	}

	basePath := viper.GetString("APP_DATA")
	newFiles, err := files.ListFiles(basePath)

	if err != nil {
		c.JSON(500, err.Error())
		return
	}

	engine.MainMessageChannel <- engine.SocketMessage{
		Files: newFiles,
	}

	c.JSON(200, "ok")
}

//
func nextBatchHandler(c *gin.Context) {
	err := engine.NextBatch()
	if err != nil {
		c.JSON(500, fmt.Errorf("Erreur lors de la création du batch suivant: "+err.Error()))
	}
	batches, _ := engine.GetBatches()
	engine.MainMessageChannel <- engine.SocketMessage{
		Batches: batches,
	}
	c.JSON(200, batches)
}

//
func upsertBatchHandler(c *gin.Context) {
	var batch base.AdminBatch
	err := c.ShouldBind(&batch)
	if err != nil {
		c.JSON(400, err.Error)
		return
	}

	err = engine.Save(&batch)
	if err != nil {
		c.JSON(500, "Erreur à l'enregistrement: "+err.Error())
		return
	}

	batches, _ := engine.GetBatches()
	engine.MainMessageChannel <- engine.SocketMessage{
		Batches: batches,
	}

	c.JSON(200, batch)
}

//
func listBatchHandler(c *gin.Context) {
	var batch []base.AdminBatch
	err := engine.Db.DB.C("Admin").Find(bson.M{"_id.type": "batch"}).Sort("-_id.key").All(&batch)
	if err != nil {
		spew.Dump(err)
		c.JSON(500, err)
		return
	}
	c.JSON(200, batch)
}

//
func processBatchHandler(c *gin.Context) {
	var query struct {
		Batches []string `json:"batches"`
		Parsers []string `json:"parsers"`
	}

	err := c.ShouldBind(&query)

	if err != nil {
		c.JSON(400, err.Error())
		return
	}
	if query.Batches == nil {
		query.Batches = engine.GetBatchesID()
	}
	if query.Batches == nil {
		query.Batches = engine.GetBatchesID()
	}

	parsers, err := resolveParsers(query.Parsers)
	if err != nil {
		c.JSON(404, err.Error())
	}
	sort.Strings(query.Batches)
	err = engine.ProcessBatch(query.Batches, parsers)
	if err != nil {
		c.JSON(500, err.Error())
		return
	}
	c.JSON(200, "ok !")
}

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

func adminFilesHandler(c *gin.Context) {
	basePath := viper.GetString("APP_DATA")
	files, err := files.ListFiles(basePath)
	if err != nil {
		c.JSON(500, err)
	} else {
		c.JSON(200, files)
	}
}

func adminRegionHandler(c *gin.Context) {
	c.JSON(200, []string{
		"FR-BFC", "FR-PDL",
	})
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

//
func eventsHandler(c *gin.Context) {
	batchKey := c.Query("batchKey")
	var query interface{}
	if batchKey != "" {
		query = bson.M{"event.batchKey": batchKey}
	}
	logs, err := engine.GetEventsFromDB(query, 250)
	if err != nil {
		c.JSON(500, err.Error())
	} else {
		c.JSON(200, logs)
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
