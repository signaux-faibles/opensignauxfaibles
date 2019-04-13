package main

import (
	"dbmongo/lib/engine"

	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
)

func datapiExportHandler(c *gin.Context) {
	err := engine.ExportToDatapi(
		viper.GetString("datapiUrl"),
		viper.GetString("datapiUser"),
		viper.GetString("datapiPassword"),
	)
	if err != nil {
		c.JSON(500, err.Error())
		return
	}
	c.JSON(200, "ok")
}
