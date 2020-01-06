package main

import (
	"github.com/signaux-faibles/opensignauxfaibles/dbmongo/lib/engine"

	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
)

func datapiExportEtablissementHandler(c *gin.Context) {
	var params struct {
		Key   string `json:"key"`
		Batch string `json:"batch"`
	}
	err := c.Bind(&params)
	if err != nil {
		c.JSON(400, err.Error())
		return
	}

	b := engine.AdminBatch{}
	err = b.Load(params.Batch)
	if err != nil {
		c.JSON(404, "batch non trouvé, vérifiez le paramètre batch")
		return
	}

	if !(len(params.Key) == 14 || len(params.Key) == 0) {
		c.JSON(400, "siret de 14 caractères obligatoire si fourni")
		return
	}

	err = engine.ExportEtablissementToDatapi(
		viper.GetString("datapiUrl"),
		viper.GetString("datapiUser"),
		viper.GetString("datapiPassword"),
		params.Batch,
		params.Key,
	)
	if err != nil {
		c.JSON(500, err.Error())
		return
	}
	c.JSON(200, "ok")
}

func datapiExportDetectionHandler(c *gin.Context) {
	var params struct {
		Batch string `json:"batch"`
		Key   string `json:"key"`
		Algo  string `json:"algo"`
	}
	err := c.Bind(&params)
	if err != nil {
		c.JSON(400, err.Error())
		return
	}

	if params.Algo == "" || params.Batch == "" {
		c.JSON(400, "algo et batch obligatoires")
		return
	}

	err = engine.ExportDetectionToDatapi(
		viper.GetString("datapiUrl"),
		viper.GetString("datapiUser"),
		viper.GetString("datapiPassword"),
		params.Batch,
		params.Key,
		params.Algo,
	)
	if err != nil {
		c.JSON(500, err.Error())
		return
	}
	c.JSON(200, "ok")
}

func datapiExportPoliciesHandler(c *gin.Context) {
	var params struct {
		Filter string `json:"filter"`
	}

	err := c.Bind(&params)
	if err != nil {
		c.JSON(400, err.Error())
		return
	}

	err = engine.ExportPoliciesToDatapi(
		viper.GetString("datapiUrl"),
		viper.GetString("datapiUser"),
		viper.GetString("datapiPassword"),
		params.Filter,
	)

	if err != nil {
		c.JSON(500, err.Error())
		return
	}

	c.JSON(200, "ok")
}

func datapiExportReferenceHandler(c *gin.Context) {
	var params struct {
		Batch string `json:"batch"`
		Algo  string `json:"algo"`
	}
	err := c.Bind(&params)
	if err != nil {
		c.JSON(400, err.Error())
		return
	}
	if params.Algo == "" || params.Batch == "" {
		c.JSON(400, "batch et algo obligatoire")
	}

	err = engine.ExportReferencesToDatapi(
		viper.GetString("datapiUrl"),
		viper.GetString("datapiUser"),
		viper.GetString("datapiPassword"),
		params.Batch,
		params.Algo,
	)

	if err != nil {
		c.JSON(500, err.Error())
		return
	}
	c.JSON(200, "ok")
}
