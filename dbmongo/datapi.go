package main

import (
	"dbmongo/lib/engine"

	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
)

func datapiExportDetectionHandler(c *gin.Context) {
	var params struct {
		Batch string `json:"batch"`
	}
	err := c.Bind(&params)
	if err != nil {
		c.JSON(400, err.Error())
		return
	}

	err = engine.ExportDetectionToDatapi(
		viper.GetString("datapiUrl"),
		viper.GetString("datapiUser"),
		viper.GetString("datapiPassword"),
		params.Batch,
	)
	if err != nil {
		c.JSON(500, err.Error())
		return
	}
	c.JSON(200, "ok")
}

func datapiExportPublicHandler(c *gin.Context) {
	var params struct {
		Batch string `json:"batch"`
	}
	err := c.Bind(&params)
	if err != nil {
		c.JSON(400, err.Error())
		return
	}

	err = engine.ExportPublicToDatapi(
		viper.GetString("datapiUrl"),
		viper.GetString("datapiUser"),
		viper.GetString("datapiPassword"),
		params.Batch,
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
	}
	err := c.Bind(&params)
	if err != nil {
		return
	}

	err = engine.ExportReferencesToDatapi(
		viper.GetString("datapiUrl"),
		viper.GetString("datapiUser"),
		viper.GetString("datapiPassword"),
		params.Batch,
	)
	if err != nil {
		c.JSON(500, err.Error())
		return
	}
	c.JSON(200, "ok")
}
