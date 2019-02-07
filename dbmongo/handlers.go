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
	"dbmongo/lib/sirene"
	"dbmongo/lib/urssaf"
	"fmt"
	"io"
	"os"

	"github.com/davecgh/go-spew/spew"
	"github.com/gin-gonic/gin"
	"github.com/globalsign/mgo/bson"
	"github.com/spf13/viper"
)

// // AdminID Collection key
// type AdminID struct {
// 	Key  string `json:"key" bson:"key"`
// 	Type string `json:"type" bson:"type"`
// }

//
// @summary Identifiants des traitements de réduction
// @description Correspond aux mapReduces qui produisent les variables dans la collection Features
// @Tags Administration
// @accept  json
// @produce  json
// @Security ApiKeyAuth
// @Success 200 {string} string ""
// @Router /api/admin/features [get]
func adminFeature(c *gin.Context) {
	c.JSON(200, []string{"algo1", "algo2"})
}

//
// @summary Liste des types
// @description Correspond aux types disponibles dans les traitements d'importation
// @Tags Administration
// @accept  json
// @produce  json
// @Security ApiKeyAuth
// @Success 200 {string} string ""
// @Router /api/admin/types [get]
func listTypesHandler(c *gin.Context) {
	c.JSON(200, engine.GetTypes)
}

//
// @summary Upload d'un fichier
// @description Réalise l'upload d'un fichier dans un batch/type. Une fois l'upload effectué, le fichier est automatiquement attaché au batch.
// @Tags Administration
// @accept  json
// @produce  json
// @Security ApiKeyAuth
// @Success 200 {string} string ""
// @Router /api/admin/files [post]
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

	c.JSON(200, nil)
}

// // AdminBatch metadata Batch
// type AdminBatch struct {
// 	ID            engine.AdminID `json:"id" bson:"_id"`
// 	Files         BatchFiles     `json:"files" bson:"files"`
// 	Readonly      bool           `json:"readonly" bson:"readonly"`
// 	CompleteTypes []string       `json:"complete_types" bson:"complete_types"`
// 	Params        struct {
// 		DateDebut       time.Time `json:"date_debut" bson:"date_debut"`
// 		DateFin         time.Time `json:"date_fin" bson:"date_fin"`
// 		DateFinEffectif time.Time `json:"date_fin_effectif" bson:"date_fin_effectif"`
// 	} `json:"params" bson:"param"`
// }

//
// @summary Création du batch suivant
// @description Cloture le dernier batch et crée le batch suivant dans la collection admin
// @Tags Administration
// @accept  json
// @produce  json
// @Security ApiKeyAuth
// @Success 200 {string} string ""
// @Router /api/admin/batch/next [get]
func nextBatchHandler(c *gin.Context) {
	err := engine.NextBatch()
	if err != nil {
		c.JSON(500, fmt.Errorf("Erreur lors de la création du batch suivant: "+err.Error()))
	}
	batches, _ := engine.GetBatches()
	engine.MainMessageChannel <- engine.SocketMessage{
		Batches: batches,
	}
	c.JSON(200, "nextBatch ok")
}

func sp(s string) *string {
	return &s
}

//
// @summary Remplace un batch
// @description Alimente la collection Features
// @Tags Traitements
// @accept  json
// @produce  json
// @Param algo query string true "Identifiant du traitement"
// @Param batch query string true "Identifier du batch"
// @Security ApiKeyAuth
// @Success 200 {string} string ""
// @Router /api/reduce/{algo}/{batch} [get]
func upsertBatchHandler(c *gin.Context) {
	var batch engine.AdminBatch
	err := c.Bind(&batch)
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
// @summary Liste des batches
// @description Produit une extraction des objets batch de la collection Admin
// @Tags Administration
// @accept  json
// @produce  json
// @Security ApiKeyAuth
// @Success 200 {array} string ""
// @Router /api/admin/batch [get]
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
// @summary Traitement du dernier batch
// @description Exécute l'import, le compactage et la réduction du dernier batch
// @Tags Administration
// @accept  json
// @produce  json
// @Security ApiKeyAuth
// @Success 200 {string} string ""
// @Router /api/admin/batch/next [get]
func processBatchHandler(c *gin.Context) {
	var query struct {
		Batches []string `json:"batches"`
	}
	err := c.ShouldBind(query)
	if err != nil {
		c.JSON(400, err.Error)
	}

	// TODO: valider que tous les batches demandés existent
	err = engine.ProcessBatch(query.Batches)
	if err != nil {
		c.JSON(500, err.Error())
	}
	c.JSON(200, "ok !")
}

//
// @summary Traitement du dernier batch
// @description Exécute l'import, le compactage et la réduction du dernier batch
// @Tags Administration
// @accept  json
// @produce  json
// @Security ApiKeyAuth
// @Success 200 {string} string ""
// @Router /api/data/batch/purge [get]
func purgeBatchHandler(c *gin.Context) {
	batch := engine.LastBatch()
	err := engine.PurgeBatch(batch.ID.Key)

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
// TODO: composer automatiquement ce dictionnaire à l'import des parsers
var registeredParsers = map[string]engine.Parser{
	"urssaf":    urssaf.Parser,
	"apconso":   apartconso.Parser,
	"apdemande": apartdemande.Parser,
	"bdf":       bdf.Parser,
	"altares":   altares.Parser,
	"sirene":    sirene.Parser,
	"diane":     diane.Parser,
	"interim":   interim.Parser,
}

//
// @summary Liste les fichiers disponibles dans le dépot
// @description Tous ces fichiers sont contenu dans APP_DATA (voir config.toml)
// @Tags Administration
// @accept  json
// @produce  json
// @Security ApiKeyAuth
// @Success 200 {string} string ""
// @Router /api/admin/files [get]
func adminFilesHandler(c *gin.Context) {
	basePath := viper.GetString("APP_DATA")
	files, err := files.ListFiles(basePath)
	if err != nil {
		c.JSON(500, err)
	} else {
		c.JSON(200, files)
	}
}

//
// @summary Import de fichiers pour un batch
// @description Effectue l'import de tous les fichiers du batch donné en paramètre
// @Tags Traitements
// @accept  json
// @produce  json
// @Security ApiKeyAuth
// @Param batch query string true "Clé du batch"
// @Success 200 {string} string ""
// @Router /api/data/import/{batch} [get]
func importBatchHandler(c *gin.Context) {
	var params struct {
		BatchKey string   `json:"batch"`
		Parsers  []string `json:"parsers"`
	}
	err := c.ShouldBind(&params)
	if err != nil {
		c.JSON(400, "Requête malformée: "+err.Error())
	}
	batch := engine.AdminBatch{}
	batch.Load(params.BatchKey)

	var parsers []engine.Parser
	if params.Parsers == nil {
		for _, f := range registeredParsers {
			parsers = append(parsers, f)
		}
	} else {
		for _, p := range params.Parsers {
			if f, ok := registeredParsers[p]; ok {
				parsers = append(parsers, f)
			} else {
				c.JSON(404, p+" n'est pas un parser reconnu.")
				return
			}
		}
	}

	engine.ImportBatch(batch, parsers)
}

//
// @summary Journal d'évènements
// @description Liste les 100 derniers évènements du journal
// @Tags Administration
// @accept  json
// @produce  json
// @Security ApiKeyAuth
// @Success 200 {string} string ""
// @Router /api/admin/getLogs [get]
func getLogsHandler(c *gin.Context) {
	logs, err := engine.GetEventsFromDB(nil, 250)
	if err != nil {
		c.JSON(500, err.Error())
	} else {
		c.JSON(200, logs)
	}
}
