package main

import (
	"errors"
	"io/ioutil"
	"regexp"

	"github.com/gin-gonic/gin"
	"github.com/globalsign/mgo"
	"github.com/globalsign/mgo/bson"
)

func loadJSFunctions(path string) (map[string]bson.JavaScript, error) {
	files, err := ioutil.ReadDir(path)
	r := regexp.MustCompile(`(.*)\.js$`)

	functions := make(map[string]bson.JavaScript)

	for _, f := range files {
		if name := r.FindStringSubmatch(f.Name()); len(name) > 0 {
			b, err := ioutil.ReadFile(path + f.Name())
			if err == nil {
				functions[name[1]] = bson.JavaScript{
					Code: string(b),
				}
			}
		}
	}
	return functions, err
}

//
// @summary Lance un traitement de réduction
// @description Alimente la collection Features
// @Tags Traitements
// @accept  json
// @produce  json
// @Param algo query string true "Identifiant du traitement"
// @Param batch query string true "Identifier du batch"
// @Security ApiKeyAuth
// @Success 200 {string} string ""
// @Router /api/data/reduce/{algo}/{batch} [get]
func reduceHandler(c *gin.Context) {
	batchKey := c.Params.ByName("batchKey")
	algo := c.Params.ByName("algo")
	err := reduce(algo, batchKey)
	if err != nil {
		c.JSON(500, err.Error())
	} else {
		c.JSON(200, "Traitement effectué")
	}
}

func reduce(batchKey string, algo string) error {
	// éviter les noms d'algo essayant de pervertir l'exploration des fonctions
	isAlphaNum := regexp.MustCompile(`^[A-Za-z0-9]+$`).MatchString
	if !isAlphaNum(algo) {
		return errors.New("nom d'algorithme invalide, alphanumérique sans espace exigé")
	}

	functions, err := loadJSFunctions("js/" + algo + "/")

	naf, err = loadNAF()
	if err != nil {
		return err
	}

	batch, err := getBatch(batchKey)
	if err != nil {
		return err
	}

	scope := bson.M{
		"date_debut":             batch.Params.DateDebut,
		"date_fin":               batch.Params.DateFin,
		"date_fin_effectif":      batch.Params.DateFinEffectif,
		"serie_periode":          genereSeriePeriode(batch.Params.DateDebut, batch.Params.DateFin),
		"serie_periode_annuelle": genereSeriePeriodeAnnuelle(batch.Params.DateDebut, batch.Params.DateFin),
		"offset_effectif":        (batch.Params.DateFinEffectif.Year()-batch.Params.DateFin.Year())*12 + int(batch.Params.DateFinEffectif.Month()-batch.Params.DateFin.Month()),
		"actual_batch":           batch.ID.Key,
		"naf":                    naf,
		"f":                      functions,
		"batches":                getBatchesID(),
		"types":                  getTypes(),
	}

	job := &mgo.MapReduce{
		Map:      functions["map"].Code,
		Reduce:   functions["reduce"].Code,
		Finalize: functions["finalize"].Code,
		Out:      bson.M{"replace": "Features"},
		Scope:    scope,
	}

	_, err = db.DB.C("RawData").Find(bson.M{"value.index.algo2": true}).MapReduce(job, nil)

	return err
}

//
// @summary Lance un traitement de compactage
// @description Alimente la collection Features
// @Tags Traitements
// @accept  json
// @produce  json
// @Param algo query string true "Identifiant du traitement"
// @Param batch query string true "Identifier du batch"
// @Success 200 {string} string ""
// @Router /api/data/compact [get]
// @Security ApiKeyAuth
func compactHandler(c *gin.Context) {
	err := compact()
	if err != nil {
		c.JSON(500, err.Error())
		return
	}

	c.JSON(200, "ok")
}

func compact() error {
	batches, _ := getBatches()

	// Détermination scope traitement
	var completeTypes = make(map[string][]string)
	var batchesID []string

	for _, b := range batches {
		completeTypes[b.ID.Key] = b.CompleteTypes
		batchesID = append(batchesID, b.ID.Key)
	}

	functions, err := loadJSFunctions("js/compact/")
	if err != nil {
		return err
	}
	// Traitement MR
	job := &mgo.MapReduce{
		Map:      functions["map"].Code,
		Reduce:   functions["reduce"].Code,
		Finalize: functions["finalize"].Code,
		Out:      bson.M{"replace": "RawData"},
		Scope: bson.M{
			"functions":     functions,
			"batches":       getBatchesID(),
			"types":         getTypes(),
			"completeTypes": completeTypes,
		},
	}

	_, err = db.DB.C("RawData").Find(nil).MapReduce(job, nil)

	return err
}

func getTypes() []string {
	return []string{
		"altares",
		"apconso",
		"apdemande",
		"ccsf",
		"cotisation",
		"debit",
		"delai",
		"effectif",
		"sirene",
		"dpae",
		"bdf",
		"diane",
	}
}

func getNAF(c *gin.Context) {
	c.JSON(200, naf)
}
