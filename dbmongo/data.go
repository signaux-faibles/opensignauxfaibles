package main

import (
	"dbmongo/lib/engine"
	"dbmongo/lib/naf"
	"fmt"
	"strconv"

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
			Pattern: "^" + params.Key[0:9],
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

func reduceSlicedHandler(c *gin.Context) {
	var params struct {
		BatchKey string `json:"batch"`
		Algo     string `json:"features"`
	}
	err := c.ShouldBind(&params)
	if err != nil {
		c.JSON(400, err.Error())
	}

	var queries [10]bson.M
	var collection string
	for i, _ := range queries {
		queries[i] = bson.M{
			"_id": bson.RegEx{
				Pattern: "^" + strconv.Itoa(i) + ".*",
				Options: "",
			},
			"value.index." + params.Algo: true,
		}

		fmt.Println(queries[i])
		collection = "Features_aux"
		err = engine.Reduce(params.BatchKey, params.Algo, queries[i], collection)
		if err != nil {
			c.JSON(500, err.Error())
			return
		}
		fmt.Println("Features_aux full of new stuff")
		err = engine.ReduceMergeAux()
		if err != nil {
			c.JSON(500, err.Error())
			return
		}
		fmt.Println("Merge completed")
	}

	if err != nil {
		c.JSON(500, err.Error())
	} else {
		c.JSON(200, "Traitement effectué")
	}
}

func publicSlicedHandler(c *gin.Context) {
	var params struct {
		BatchKey string `json:"batch"`
		Algo     string `json:"features"`
	}
	err := c.ShouldBind(&params)
	if err != nil {
		c.JSON(400, err.Error())
	}

	batch := engine.AdminBatch{}
	err = batch.Load(params.BatchKey)
	if err != nil {
		c.JSON(404, "batch non trouvé")
		return
	}

	var queries []bson.M
	var collection string
	slices := []string{
		"^0.*", "^1.*", "^2.*", "^3[0-4].*", "^3[5-9].*", "^4.*", "^5.*", "^6.*", "^7.*", "^8.*", "^9.*",
	}
	for _, s := range slices {
		query := bson.M{
			"_id": bson.RegEx{
				Pattern: s,
				Options: "",
			},
			"value.index." + params.Algo: true,
		}
		queries = append(queries, query)

		collection = "Public_aux"
		err = engine.Public(batch, params.Algo, query, collection)
		if err != nil {
			c.JSON(500, err.Error())
			return
		}
		fmt.Println("Public_aux full of new stuff")
		err = engine.PublicMergeAux()
		if err != nil {
			c.JSON(500, err.Error())
			return
		}
		fmt.Println("Public Merge completed")
	}

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
