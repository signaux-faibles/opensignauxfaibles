package main

import (
	"dbmongo/lib/engine"

	"github.com/gin-gonic/gin"
)

func publicHandler(c *gin.Context) {
	params := struct {
		Batch string `json:"batch"`
	}{}
	c.Bind(&params)

	batch := engine.AdminBatch{}
	err := batch.Load(params.Batch)
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

func predictionBrowseHandler(c *gin.Context) {
	var params engine.BrowseParams

	err := c.ShouldBind(&params)
	if err != nil {
		c.JSON(400, "Bad Request: "+err.Error())
	}

	result, err := engine.PredictionBrowse(params)
	if err != nil {
		c.JSON(500, err)
		return
	}
	c.JSON(200, result)
}

func searchRaisonSocialeHandler(c *gin.Context) {
	var params engine.SearchCriteria
	err := c.ShouldBind(&params)
	if err != nil {
		c.JSON(400, err.Error())
		return
	}

	result, err := engine.SearchRaisonSociale(params)

	if err != nil {
		c.JSON(500, err.Error())
		return
	}
	c.JSON(200, result)
}
