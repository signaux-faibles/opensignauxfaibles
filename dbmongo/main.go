package main

import (
	"fmt"
	"opensignauxfaibles/dbmongo/lib/engine"

	"net/http"
	"opensignauxfaibles/dbmongo/lib/naf"

	"github.com/gin-contrib/cors"
	"github.com/gin-contrib/static"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"github.com/spf13/viper"

	ginSwagger "github.com/swaggo/gin-swagger"
	"github.com/swaggo/gin-swagger/swaggerFiles"

	_ "opensignauxfaibles/dbmongo/docs"
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
		config.AllowOrigins = []string{"https://signaux.faibles.fr"}
	}
	config.AddAllowHeaders("Authorization")
	config.AddAllowMethods("GET", "POST")
	r.Use(cors.New(config))

	r.Use(static.Serve("/", static.LocalFile("static/", true)))
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

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
		api.POST("/data/exportReference", datapiExportReferenceHandler)
		api.POST("/data/exportDetection", datapiExportDetectionHandler)
		api.POST("/data/exportPolicies", datapiExportPoliciesHandler)
		// TODO: mapreduce pour traiter le scope, modification des objets utilisateurs
		// TODO: écrire l'aggrégation qui va bien

	}

	bind := viper.GetString("APP_BIND")
	r.Run(bind)
}
