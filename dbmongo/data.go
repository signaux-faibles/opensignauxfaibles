package main

import (
	"github.com/signaux-faibles/opensignauxfaibles/dbmongo/lib/engine"
	"github.com/signaux-faibles/opensignauxfaibles/dbmongo/lib/naf"

	"github.com/gin-gonic/gin"
)

func reduceHandler(c *gin.Context) {
	var params struct {
		BatchKey string   `json:"batch"`
		Algo     string   `json:"algo"`
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
		c.JSON(404, "le batch "+params.BatchKey+" n'existe pas")
	}

	if params.Key == "" && params.From == "" && params.To == "" {
		err = engine.Reduce(batch, params.Algo, params.Types)
	} else {
		err = engine.ReduceOne(batch, params.Algo, params.Key, params.From, params.To, params.Types)
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

	batch := engine.AdminBatch{}
	err = batch.Load(params.BatchKey)
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
		BatchKey string `json:"batch"`
	}
	err := c.ShouldBind(&params)
	if err != nil {
		c.JSON(400, err.Error())
	}

	err = engine.Compact(params.BatchKey)
	if err != nil {
		c.JSON(500, err.Error())
		return
	}
	c.JSON(200, "ok")
}

func nafHandler(c *gin.Context) {
	c.JSON(200, naf.Naf)
}

func purgeHandler(c *gin.Context) {
	var params struct {
		AreYouSure string `json:"areyousure"`
	}
	c.Bind(&params)
	if params.AreYouSure == "yes" {
		info := engine.Purge()
		c.JSON(200, info)
		return
	}
	c.JSON(300, "Provide areyousure=yes")
}
