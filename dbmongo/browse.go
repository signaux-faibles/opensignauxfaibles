package main

import (
	"github.com/signaux-faibles/opensignauxfaibles/dbmongo/lib/engine"

	"github.com/gin-gonic/gin"
)

func searchHandler(c *gin.Context) {
	var params engine.SearchParams
	err := c.ShouldBind(&params)
	if err != nil {
		c.JSON(400, err.Error())
		return
	}

	result, err := engine.Search(params)

	if err != nil {
		c.JSON(500, err.Error())
		return
	}
	c.JSON(200, result)
}
