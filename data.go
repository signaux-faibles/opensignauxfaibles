package main

import (
	"errors"
	"strconv"
	"time"

	"github.com/signaux-faibles/opensignauxfaibles/lib/base"
	"github.com/signaux-faibles/opensignauxfaibles/lib/engine"
	"github.com/spf13/viper"

	"github.com/gin-gonic/gin"
)

func reduceHandler(c *gin.Context) {
	var params struct {
		BatchKey string   `json:"batch"`
		Key      string   `json:"key"`
		From     string   `json:"from"`
		To       string   `json:"to"`
		Types    []string `json:"types"`
		// Sélection des types de données qui vont être calculés ou recalculés.
		// Valeurs autorisées pour l'instant: "apart", "all"
	}

	err := c.ShouldBind(&params)
	if err != nil {
		c.JSON(400, err.Error())
	}

	batch, err := engine.GetBatch(params.BatchKey)
	if err != nil {
		c.JSON(404, "Batch inexistant: "+err.Error())
	}

	if params.Key == "" && params.From == "" && params.To == "" {
		err = engine.Reduce(batch, params.Types)
	} else {
		err = engine.ReduceOne(batch, params.Key, params.From, params.To, params.Types)
	}

	if err != nil {
		c.JSON(500, err.Error())
	} else {
		c.JSON(200, "Traitement effectué")
	}

}

func publicHandler(c *gin.Context) {
	var params struct {
		BatchKey string `json:"batch"`
		Key      string `json:"key"`
	}
	err := c.ShouldBind(&params)
	if err != nil || params.BatchKey == "" {
		c.JSON(400, err.Error())
	}

	batch := base.AdminBatch{}
	err = engine.Load(&batch, params.BatchKey)
	if err != nil {
		c.JSON(404, "batch non trouvé")
		return
	}

	if params.Key == "" {
		err = engine.Public(batch)
	} else if len(params.Key) >= 9 {
		err = engine.PublicOne(batch, params.Key[0:9])
	} else {
		c.JSON(400, "la clé fait moins de 9 caractères (siren)")
	}

	if err != nil {
		c.JSON(500, err.Error())
		return
	}

	c.JSON(200, "ok")
}

func compactHandler(c *gin.Context) {
	var params struct {
		FromBatchKey string `json:"fromBatchKey"`
	}
	err := c.ShouldBind(&params)
	if err != nil {
		c.JSON(400, err.Error())
	}

	err = engine.Compact(params.FromBatchKey)
	if err != nil {
		c.JSON(500, err.Error())
		return
	}
	c.JSON(200, "ok")
}

func getTimestamp() string {
	return strconv.FormatInt(time.Now().Unix(), 10)
}

func getKeyParam(c *gin.Context) (string, error) {
	key := c.Query("key")
	if !(len(key) == 9 || len(key) == 0) {
		return "", errors.New("si fourni, key doit être un numéro SIREN (9 chiffres)")
	}
	return key, nil
}

func exportEtablissementsHandler(c *gin.Context) {
	key, err := getKeyParam(c)
	if err != nil {
		c.JSON(400, err.Error())
		return
	}

	// On retourne le nom de fichier avant la fin du traitement, pour éviter erreur "Request timed out"
	var filepath = viper.GetString("exportPath") + "dbmongo-data-export-etablissements-" + getTimestamp() + ".json.gz"
	c.JSON(200, filepath)

	err = engine.ExportEtablissements(key, filepath)
	if err != nil {
		c.AbortWithError(500, err)
	}
}

func exportEntreprisesHandler(c *gin.Context) {
	key, err := getKeyParam(c)
	if err != nil {
		c.JSON(400, err.Error())
		return
	}

	// On retourne le nom de fichier avant la fin du traitement, pour éviter erreur "Request timed out"
	var filepath = viper.GetString("exportPath") + "dbmongo-data-export-entreprises-" + getTimestamp() + ".json.gz"
	c.JSON(200, filepath)

	err = engine.ExportEntreprises(key, filepath)
	if err != nil {
		c.AbortWithError(500, err)
	}
}

func validateHandler(c *gin.Context) {

	var params struct {
		Collection string `json:"collection"`
	}
	c.ShouldBind(&params)
	if params.Collection != "RawData" && params.Collection != "ImportedData" {
		c.JSON(400, "le paramètre collection doit valoir RawData ou ImportedData")
		return
	}

	// On retourne le nom de fichier avant la fin du traitement, pour éviter erreur "Request timed out"
	var filepath = viper.GetString("exportPath") + "dbmongo-" + params.Collection + "-validation-" + getTimestamp() + ".json.gz"

	jsonSchema, err := engine.LoadJSONSchemaFiles()
	if err != nil {
		c.JSON(500, err)
		return
	}
	c.JSON(200, filepath)

	err = engine.ValidateDataEntries(filepath, jsonSchema, params.Collection)
	if err != nil {
		c.AbortWithError(500, err)
	}
}
