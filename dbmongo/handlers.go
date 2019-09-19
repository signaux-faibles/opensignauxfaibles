package main

import (
	"dbmongo/lib/altares"
	"dbmongo/lib/apartconso"
	"dbmongo/lib/apartdemande"
	"dbmongo/lib/bdf"
	"dbmongo/lib/diane"
	"dbmongo/lib/engine"
	"dbmongo/lib/files"
	"dbmongo/lib/interim"
	"dbmongo/lib/repeatableOrder"
	"dbmongo/lib/sirene"
	"dbmongo/lib/sirene_ul"
	"dbmongo/lib/urssaf"
	"errors"
	"fmt"
	"io"
	"os"
	"sort"

	"github.com/davecgh/go-spew/spew"
	"github.com/gin-gonic/gin"
	"github.com/globalsign/mgo/bson"
	"github.com/spf13/viper"
)

//
func adminFeature(c *gin.Context) {
	c.JSON(200, []string{"algo1", "algo2"})
}

//
func listTypesHandler(c *gin.Context) {
	c.JSON(200, engine.GetTypes())
}

//
func addFile(c *gin.Context) {
	file, err := c.FormFile("file")
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

func sp(s string) *string {
	return &s
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
	types := query.Parsers
	if err != nil {
		c.JSON(404, err.Error())
	}
	sort.Strings(query.Batches)
	err = engine.ProcessBatch(query.Batches, parsers, types)
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

// RegisteredParsers liste des parsers disponibles
var registeredParsers = map[string]engine.Parser{
	"urssaf":          urssaf.Parser,
	"apconso":         apartconso.Parser,
	"apdemande":       apartdemande.Parser,
	"bdf":             bdf.Parser,
	"altares":         altares.Parser,
	"repeatableOrder": repeatableOrder.Parser,
	"sirene":          sirene.Parser,
	"sirene_ul":       sirene_ul.Parser,
	"diane":           diane.Parser,
	"interim":         interim.Parser,
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

//
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
	engine.PurgeNotCompacted()
	c.JSON(200, result)
}

func getTasksHandler(c *gin.Context) {
	c.JSON(501, "Not implemented (for the moment)")
}

func browsePublicHandler(c *gin.Context) {
	data := engine.BrowsePublic(nil)
	c.JSON(200, data)
}

// Vérifie et charge les parsers
func resolveParsers(parserNames []string) ([]engine.Parser, error) {
	var parsers []engine.Parser
	if parserNames == nil {
		for _, f := range registeredParsers {
			parsers = append(parsers, f)
		}
	} else {
		for _, p := range parserNames {
			if f, ok := registeredParsers[p]; ok {
				parsers = append(parsers, f)
			} else {
				return parsers, errors.New(p + " n'est pas un parser reconnu.")
			}
		}
	}
	return parsers, nil
}

func getCommentsHandler(c *gin.Context) {
	var params struct {
		Siret string `json:"siret"`
	}
	err := c.Bind(&params)
	if err != nil {
		return
	}
	comments, err := engine.GetComments(params.Siret)
	if err != nil {
		c.JSON(500, err.Error())
		return
	}
	c.JSON(200, comments)
}

func getCommentHistoryHandler(c *gin.Context) {
	var comment engine.Comment
	err := c.Bind(&comment)
	if err != nil {
		return
	}
	comments, err := engine.GetCommentHistory(comment)
	if err != nil {
		c.JSON(500, err.Error())
		return
	}
	c.JSON(200, comments)
}

func setCommentHandler(c *gin.Context) {
	var comment engine.Comment
	c.Bind(&comment)
	err := engine.SetComment(comment)
	if err != nil {
		c.JSON(500, err.Error())
		return
	}
	c.JSON(200, "ok")
}

func publicHandler(c *gin.Context) {
	params := struct {
		Batch string `json:"batch"`
		Siret string `json:"siret"`
	}{}
	c.Bind(&params)
	batch := engine.AdminBatch{}
	err := batch.Load(params.Batch)
	if err != nil {
		c.JSON(404, "batch non trouvé")
		return
	}

	err = engine.Public(batch, params.Siret)
	if err != nil {
		c.JSON(500, err.Error())
		return
	}

	c.JSON(200, "ok")
}
