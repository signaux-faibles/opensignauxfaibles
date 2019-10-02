package main

import (
	"dbmongo/lib/engine"
	"dbmongo/lib/naf"

	"github.com/gin-gonic/gin"
)

func reduceHandler(c *gin.Context) {
	var params struct {
		BatchKey string `json:"batch"`
		Algo     string `json:"features"`
		Key      string `json:"key"`
	}
	err := c.ShouldBind(&params)
	if err != nil {
		c.JSON(400, err.Error())
	}

	batch, err := engine.GetBatch(params.BatchKey)
	if err != nil {
		c.JSON(404, "le batch "+params.BatchKey+" n'existe pas")
	}

	if params.Key == "" {
		err = engine.Reduce(batch, params.Algo)
	} else {
		err = engine.ReduceOne(batch, params.Algo, params.Key)
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
		Algo     string `json:"algo"`
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

	err = engine.Public(batch)

	if err != nil {
		c.JSON(500, err.Error())
		return
	}

	c.JSON(200, "ok")
}

func compactHandler(c *gin.Context) {
	var params struct {
		BatchKey string   `json:"batch"`
		Types    []string `json:"types"`
	}
	err := c.ShouldBind(&params)
	if err != nil {
		c.JSON(400, err.Error())
	}
	//TODO: verifier comportement si batch est vide
	//TODO: verifier comportement si types est vide
	err = engine.Compact(params.BatchKey, params.Types)
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
