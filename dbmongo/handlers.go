package main

import (
	"errors"
	"fmt"
	"io"
	"os"
	"sort"

	"github.com/signaux-faibles/opensignauxfaibles/dbmongo/lib/apconso"
	"github.com/signaux-faibles/opensignauxfaibles/dbmongo/lib/apdemande"
	"github.com/signaux-faibles/opensignauxfaibles/dbmongo/lib/bdf"
	"github.com/signaux-faibles/opensignauxfaibles/dbmongo/lib/diane"
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
	var batch engine.AdminBatch
	err := c.ShouldBind(&batch)
	if err != nil {
		c.JSON(400, err.Error)
		return
	}

	err = batch.Save()
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
	var batch []engine.AdminBatch
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

	// TODO: valider que tous les batches demandés existent

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
		BatchKey string `json:"batch"`
	}
	err := c.ShouldBind(&params)
	if err != nil {
		c.JSON(400, "Requête malformée: "+err.Error())
		return
	}
	if params.BatchKey == "" {
		batch := engine.LastBatch()
		params.BatchKey = batch.ID.Key
	}
	err = engine.PurgeBatch(params.BatchKey)

	if err != nil {
		c.JSON(500, "Erreur dans la purge du batch: "+err.Error())
	} else {
		c.JSON(200, "ok")
	}
}

func revertBatchHandler(c *gin.Context) {
	err := engine.RevertBatch()
	if err != nil {
		c.JSON(500, err)
	}
	batches, _ := engine.GetBatches()
	engine.MainMessageChannel <- engine.SocketMessage{
		Batches: batches,
	}
	c.JSON(200, "ok")
}

//
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
	}
	err := c.ShouldBind(&params)
	if err != nil {
		c.JSON(400, "Requête malformée: "+err.Error())
		return
	}
	batch := engine.AdminBatch{}
	batch.Load(params.BatchKey)

	parsers, err := resolveParsers(params.Parsers)
	if err != nil {
		c.JSON(404, err.Error())
	}
	engine.ImportBatch(batch, parsers)
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
	batch := engine.AdminBatch{}
	batch.Load(params.BatchKey)

	parsers, err := resolveParsers(params.Parsers)
	if err != nil {
		c.JSON(404, err.Error())
	}
	err = engine.CheckBatch(batch, parsers)
	if err != nil {
		c.JSON(417, "Erreurs détectées: "+err.Error())
	} else {
		c.JSON(200, true)
	}
}

//
func eventsHandler(c *gin.Context) {
	logs, err := engine.GetEventsFromDB(nil, 250)
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
var registeredParsers = map[string]engine.Parser{
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
}

// encapsultateParser transforms parser options into a functional parser
func encapsulateParser(po *marshal.ParserOptions) engine.Parser {
	parser := func(c engine.Cache, ab *engine.AdminBatch) (chan engine.Tuple, chan engine.Event) {
		tuple, event, _ := marshal.GenericMarshal(po, c, ab)
		return tuple, event
	}
	return parser
}

// Vérifie et charge les parsers
func resolveParsers(parserNames []string) ([]engine.Parser, error) {
	var parsers []engine.Parser
	var registeredParserOptions map[string]*marshal.ParserOptions
	parserOptionsDir := viper.GetString("PARSEROPTIONS_DIR")
	if parserOptionsDir == "" {
		fmt.Println("No parser options could be read. Have you set PARSEROPTIONS_DIR config ?")
	} else {
		aux, err := marshal.RegisteredParserOptions(parserOptionsDir)
		if err != nil {
			return parsers, errors.New("Parser options could not be read at " + parserOptionsDir + ": " + err.Error())
		}
		registeredParserOptions = aux
	}

	if parserNames == nil {
		for _, f := range registeredParsers {
			parsers = append(parsers, f)
		}
		for _, po := range registeredParserOptions {
			parsers = append(parsers, encapsulateParser(po))
		}
	} else {
		for _, p := range parserNames {
			if f, ok := registeredParsers[p]; ok {
				parsers = append(parsers, f)
			} else if po, ok := registeredParserOptions[p]; ok {
				parsers = append(parsers, encapsulateParser(po))
			} else {
				return parsers, errors.New(p + " n'est pas un parser reconnu.")
			}
		}
	}
	return parsers, nil
}
