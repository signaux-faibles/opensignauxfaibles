package main

import (
	"fmt"
	"time"

	"github.com/globalsign/mgo/bson"

	"github.com/signaux-faibles/opensignauxfaibles/dbmongo/lib/engine"

	"net/http"

	"github.com/signaux-faibles/opensignauxfaibles/dbmongo/lib/naf"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"github.com/spf13/viper"

	ginSwagger "github.com/swaggo/gin-swagger"
	"github.com/swaggo/gin-swagger/swaggerFiles"

	_ "github.com/signaux-faibles/opensignauxfaibles/dbmongo/docs"
)

var wsupgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin:     checkOrigin,
}

func checkOrigin(r *http.Request) bool {
	return true
}

// wshandler connecteur WebSocket
func wshandler(w http.ResponseWriter, r *http.Request, jwt string) {
	conn, err := wsupgrader.Upgrade(w, r, nil)
	if err != nil {
		fmt.Printf("Failed to set websocket upgrade: %+v", err)
		return
	}
	channel := make(chan engine.SocketMessage)
	engine.AddClientChannel <- channel

	for event := range channel {
		conn.WriteJSON(event)
	}
}

// main Fonction Principale
func main() {
	// Lancer Rserve en background

	// go r()
	engine.Db = engine.InitDB()
	go engine.MessageSocketAddClient()

	var err error
	naf.Naf, err = naf.LoadNAF()
	if err != nil {
		panic(err)
	}

	r := gin.New()

	if !viper.GetBool("DEV") {
		gin.SetMode(gin.ReleaseMode)
	} else {
		r.Use(gin.Recovery())
		r.Use(gin.Logger())
	}

	config := cors.DefaultConfig()
	if viper.GetBool("DEV") {
		config.AllowOrigins = []string{"*"}
	} else {
		config.AllowOrigins = []string{viper.GetString("corsDomain")}
	}
	config.AddAllowHeaders("Authorization")
	config.AddAllowMethods("GET", "POST")
	r.Use(cors.New(config))

	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler)) // serves interactive API documentation on /swagger/index.html

	r.POST("/datapi/exportReference", datapiExportReferenceHandler)
	r.POST("/datapi/exportDetection", datapiExportDetectionHandler)
	r.POST("/datapi/exportEtablissement", datapiExportEtablissementHandler)
	r.POST("/datapi/exportPolicies", datapiExportPoliciesHandler)

	api := r.Group("api")

	{
		api.POST("/admin/batch", upsertBatchHandler)
		api.GET("/admin/batch", listBatchHandler)
		api.GET("/admin/batch/next", nextBatchHandler)
		api.POST("/admin/batch/process", processBatchHandler)
		api.GET("/admin/batch/revert", revertBatchHandler)
		api.GET("/admin/regions", adminRegionHandler)
		api.GET("/admin/files", adminFilesHandler)
		api.POST("/admin/files", addFile)

		api.GET("/admin/types", listTypesHandler)
		api.GET("/admin/features", adminFeature)
		api.GET("/admin/events", eventsHandler)

		api.GET("/data/naf", nafHandler)
		api.POST("/data/batch/purge", purgeBatchHandler)
		api.POST("/data/import", importBatchHandler)
		api.POST("/data/check", checkBatchHandler)
		api.POST("/data/compact", compactHandler)
		api.POST("/data/reduce", reduceHandler)
		api.POST("/data/public", publicHandler)

		api.POST("/data/search", searchHandler)

		api.POST("/data/purge", purgeHandler)
		api.GET("/data/purgeNotCompacted", purgeNotCompactedHandler)

		api.POST("/data/copyScores", copyScores)

		api.GET("/data/etablissements", exportEtablissementsHandler)
		api.GET("/data/entreprises", exportEntreprisesHandler)

		// api.GET("/debug", debug)
	}

	bind := viper.GetString("APP_BIND")
	r.Run(bind)
}

func copyScores(c *gin.Context) {
	var params struct {
		From string `json:"from"`
		To   string `json:"to"`
		Algo string `json:"algo"`
	}
	c.Bind(&params)

	type Score struct {
		ID        bson.ObjectId `bson:"_id"`
		Siret     string        `bson:"siret"`
		Periode   string        `bson:"periode"`
		Score     float64       `bson:"score"`
		ScoreDiff float64       `bson:"score_diff"`
		Algo      string        `bson:"algo"`
		Batch     string        `bson:"batch"`
		Timestamp time.Time     `bson:"timestamp"`
		Alert     string        `bson:"alert"`
	}

	to, _ := engine.Db.DB.C("Scores").Find(bson.M{
		"batch": params.To,
		"algo":  params.Algo,
	}).Count()

	if to > 0 {
		c.JSON(300, "il existe déjà des objets dans la destination")
		return
	}

	from := engine.Db.DB.C("Scores").Find(bson.M{
		"batch": params.From,
		"algo":  params.Algo,
	}).Iter()

	var f Score
	var dest []interface{}
	for from.Next(&f) {
		f.ID = bson.NewObjectId()
		f.Batch = params.To
		dest = append(dest, f)
	}
	engine.Db.DB.C("Scores").Insert(dest...)
}
