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
	params := struct {
		Algo     string `json:"algo"`
		Batch    string `json:"batch"`
		Naf1     string `json:"naf1"`
		Effectif int    `json:"effectif"`
		Suivi    bool   `json:"suivi"`
		Ccsf     bool   `json:"ccsf"`
		Procol   bool   `json:"procol"`
		Limit    int    `json:"limit"`
		Offset   int    `json:"offset"`
	}{}

	err := c.ShouldBind(&params)
	if err != nil {
		c.JSON(400, "Bad Request: "+err.Error())
	}

	result, err := engine.PredictionBrowse(params.Batch, params.Naf1, params.Effectif, params.Suivi, params.Ccsf, params.Procol, params.Limit, params.Offset)
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
