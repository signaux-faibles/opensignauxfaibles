package main

import (
	"errors"
	"fmt"
	"strconv"
	"time"

	"github.com/davecgh/go-spew/spew"
	"github.com/gin-gonic/gin"
	"github.com/globalsign/mgo"
	"github.com/globalsign/mgo/bson"
)

// AdminBatch metadata Batch
type AdminBatch struct {
	ID            AdminID    `json:"id" bson:"_id"`
	Files         BatchFiles `json:"files" bson:"files"`
	Readonly      bool       `json:"readonly" bson:"readonly"`
	CompleteTypes []string   `json:"complete_types" bson:"complete_types"`
	Params        struct {
		DateDebut       time.Time `json:"date_debut" bson:"date_debut"`
		DateFin         time.Time `json:"date_fin" bson:"date_fin"`
		DateFinEffectif time.Time `json:"date_fin_effectif" bson:"date_fin_effectif"`
	} `json:"params" bson:"param"`
}

// BatchFiles fichiers mappés par type
type BatchFiles map[string][]string

func (batchFiles BatchFiles) attachFile(fileType string, file string) {
	batchFiles[fileType] = append(batchFiles[fileType], file)
}

func isBatchID(batchID string) bool {
	_, err := time.Parse("0601", batchID)
	return err == nil
}

func (batch *AdminBatch) load(batchKey string) error {
	err := db.DB.C("Admin").Find(bson.M{"_id.type": "batch", "_id.key": batchKey}).One(batch)
	return err
}

func (batch *AdminBatch) save() error {
	_, err := db.DB.C("Admin").Upsert(bson.M{"_id": batch.ID}, batch)
	return err
}

func (batch *AdminBatch) new(batchID string) error {
	if batchID == "" {
		return errors.New("Valeur de batch non autorisée")
	}
	batch.ID.Key = batchID
	batch.ID.Type = "batch"
	batch.Files = BatchFiles{}
	return nil
}

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
	err := nextBatch()
	if err != nil {
		c.JSON(500, fmt.Errorf("Erreur nextBatch: "+err.Error()))
	}
	batches, _ := getBatches()
	mainMessageChannel <- socketMessage{
		Batches: batches,
	}
	c.JSON(200, "nextBatch ok")
}

func nextBatch() error {
	batch := lastBatch()
	// spew.Dump(batch)
	newBatchID, err := nextBatchID(batch.ID.Key)
	if err != nil {
		return fmt.Errorf("Mauvais numéro de batch: " + err.Error())
	}
	newBatch := AdminBatch{
		ID: AdminID{
			Key:  newBatchID,
			Type: "batch",
		},
		CompleteTypes: batch.CompleteTypes,
	}

	batch.Readonly = true

	err = batch.save()
	if err != nil {
		return fmt.Errorf("Erreur readonly Batch: " + err.Error())
	}

	err = newBatch.save()
	if err != nil {
		return fmt.Errorf("Erreur newBatch: " + err.Error())
	}
	return nil
}

func nextBatchID(batchID string) (string, error) {
	batchTime, err := time.Parse("0601", batchID)
	if err != nil {
		return "", err
	}
	nextBatchTime := time.Date(batchTime.Year(), time.Month(batchTime.Month()+1), 1, 0, 0, 0, 0, time.UTC)
	return nextBatchTime.Format("0601"), err
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
func upsertBatch(c *gin.Context) {
	status := db.Status

	batch := AdminBatch{}
	err := c.Bind(&batch)
	if err != nil {
		c.JSON(500, err)
		return
	}

	err = batch.save()
	if err != nil {
		c.JSON(500, "Erreur à l'enregistrement")
		return
	}

	batches, _ := getBatches()
	mainMessageChannel <- socketMessage{
		Batches: batches,
	}

	status.Epoch++
	status.write()

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
func listBatch(c *gin.Context) {
	var batch []AdminBatch
	err := db.DB.C("Admin").Find(bson.M{"_id.type": "batch"}).Sort("-_id.key").All(&batch)
	if err != nil {
		spew.Dump(err)
		c.JSON(500, err)
		return
	}
	c.JSON(200, batch)
}

func getBatchesID() []string {
	batches, _ := getBatches()
	var batchesID []string
	for _, b := range batches {
		batchesID = append(batchesID, b.ID.Key)
	}
	return batchesID
}

func getBatches() ([]AdminBatch, error) {
	var batches []AdminBatch
	err := db.DB.C("Admin").Find(bson.M{"_id.type": "batch"}).Sort("_id.key").All(&batches)
	return batches, err
}

// getBatch retourne le batch correspondant à la clé batchKey
func getBatch(batchKey string) (AdminBatch, error) {
	var batch AdminBatch
	err := db.DB.C("Admin").Find(bson.M{"_id.type": "batch", "_id.key": batchKey}).One(&batch)
	return batch, err
}

// batchToTime calcule la date de référence à partir de la référence de batch
func batchToTime(batch string) (time.Time, error) {
	year, err := strconv.Atoi(batch[0:2])
	if err != nil {
		return time.Time{}, err
	}

	month, err := strconv.Atoi(batch[2:4])
	if err != nil {
		return time.Time{}, err
	}

	date := time.Date(2000+year, time.Month(month), 1, 0, 0, 0, 0, time.UTC)
	return date, err
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
	go func() {
		processBatch()
	}()
	c.JSON(200, "ok !")
}

func processBatch() {
	journal("info", "processBatch", "Lancement de l'intégration du batch")
	status := db.Status
	batch := lastBatch()
	status.setDBStatus(sp("Import des fichiers"))
	importBatch(&batch)
	compact()
	// for _, algo := range []string{"algo1", "algo2"} {
	// 	_, err := reduce(batch, algo, "")
	// 	fmt.Println(err)
	// }
	status.setDBStatus(nil)
}

func lastBatch() AdminBatch {
	batches, _ := getBatches()
	l := len(batches)
	batch := batches[l-1]
	return batch
}

func createNextBatch() error {
	batchID, _ := nextBatchID(lastBatch().ID.Key)
	batch := AdminBatch{
		ID: AdminID{
			Key:  batchID,
			Type: "batch",
		},
	}
	err := batch.save()
	return err
}

type newFile struct {
	FileName string `json:"filename"`
	Type     string `json:"type"`
	BatchKey string `json:"batch"`
}

func addFileToBatch() chan newFile {
	channel := make(chan newFile)

	go func() {
		for file := range channel {
			batch, _ := getBatch(file.BatchKey)
			batch.Files[file.Type] = append(batch.Files[file.Type], file.FileName)
			batch.save()
			batches, _ := getBatches()
			db.Status.Epoch++
			db.Status.write()
			mainMessageChannel <- socketMessage{
				JournalEvent: journal(info, "addFileToBatch", "Fichier "+file.FileName+"du type "+file.Type+" ajouté au batch "+file.BatchKey),
				Batches:      batches,
			}
		}
	}()

	return channel
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
	batch := lastBatch()
	err := purgeBatch(batch.ID.Key)

	if err != nil {
		c.JSON(500, "Erreur dans la purge du batch: "+err.Error())
	} else {
		c.JSON(200, "ok")
	}
}

func purgeBatch(batchKey string) error {

	functions, err := loadJSFunctions("js/purgeBatch/")
	if err != nil {
		return err
	}
	scope := bson.M{
		"currentBatch": batchKey,
		"f":            functions,
	}

	job := &mgo.MapReduce{
		Map:      functions["map"].Code,
		Reduce:   functions["reduce"].Code,
		Finalize: functions["finalize"].Code,
		Out:      bson.M{"replace": "Toto"},
		Scope:    scope,
	}

	_, err = db.DB.C("RawData").Find(nil).MapReduce(job, nil)
	return err
}

func revertBatchHandler(c *gin.Context) {
	err := revertBatch()
	if err != nil {
		c.JSON(500, err)
	}
	batches, _ := getBatches()
	mainMessageChannel <- socketMessage{
		Batches: batches,
	}
	c.JSON(200, "ok")
}

func dropBatch(batchKey string) error {
	_, err := db.DB.C("Admin").RemoveAll(bson.M{"_id.key": batchKey, "_id.type": "batch"})
	return err
}

// revertBatch purge le batch et supprime sa référence dans la collection Admin
func revertBatch() error {
	batch := lastBatch()
	err := purgeBatch(batch.ID.Key)
	if err != nil {
		return fmt.Errorf("Erreur lors de la purge: " + err.Error())
	}
	err = dropBatch(batch.ID.Key)
	if err != nil {
		return fmt.Errorf("Erreur lors de la purge: " + err.Error())
	}

	return nil
}

var addFileChannel = addFileToBatch()
