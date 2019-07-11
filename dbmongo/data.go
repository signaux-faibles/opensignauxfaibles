package main

import (
  "dbmongo/lib/engine"
  "dbmongo/lib/naf"

  "github.com/gin-gonic/gin"
  "github.com/globalsign/mgo/bson"
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

  var query bson.M
  var collection string
  if params.Key != "" {
    query = bson.M{"_id": bson.M{"$regex": bson.RegEx{
      Pattern:"^"+params.Key[0:9],
      Options: "",
    }}}
    collection = "Features_debug"
  } else {
    query = bson.M{"value.index." + params.Algo: true}
    collection = "Features"
  }

  err = engine.Reduce(params.BatchKey, params.Algo, query, collection)
  if err != nil {
    c.JSON(500, err.Error())
  } else {
    c.JSON(200, "Traitement effectué")
  }
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
  info := engine.Purge()
  c.JSON(200, info)
}
