package main

import (
	"dbmongo/lib/engine"

	"github.com/gin-gonic/gin"
)

//
// @summary Extraction des données publiques
// @description Lance le traitement mapReduce public pour alimenter la collection Public
// @Tags Traitements
// @accept  json
// @produce  json
// @Param batch query string true "Identifiant du batch"
// @Security ApiKeyAuth
// @Success 200 {string} string ""
// @Router /api/data/public/{batch} [get]
func publicHandler(c *gin.Context) {
	batchKey := c.Params.ByName("batch")
	batch := engine.AdminBatch{}
	batch.Load(batchKey)

	err := engine.Public(batch)
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

//
// @summary Rechercher une entreprise
// @description Effectue une recherche texte sur la collection Public et retourne les 15 premiers objets correspondants
// @Tags Traitements
// @accept  json
// @produce  json
// @Security ApiKeyAuth
// @Param guessRaisonSociale query string true "Chaine à chercher"
// @Success 200 {string} string ""
// @Router /api/data/search [post]
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
