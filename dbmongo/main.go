package main

import (
	"github.com/signaux-faibles/opensignauxfaibles/dbmongo/lib/engine"

	"github.com/signaux-faibles/opensignauxfaibles/dbmongo/lib/naf"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"

	ginSwagger "github.com/swaggo/gin-swagger"
	"github.com/swaggo/gin-swagger/swaggerFiles"

	_ "github.com/signaux-faibles/opensignauxfaibles/dbmongo/docs"
)

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

	api := r.Group("api")

	{
		api.POST("/data/batch/purge", purgeBatchHandler)      // [ ] écrit dans Journal
		api.POST("/data/import", importBatchHandler)          // [x] écrit dans Journal
		api.POST("/data/check", checkBatchHandler)            // [x] écrit dans Journal
		api.POST("/data/compact", compactHandler)             // [x] écrit dans Journal
		api.POST("/data/reduce", reduceHandler)               // [x] écrit dans Journal
		api.POST("/data/public", publicHandler)               // [ ] écrit dans Journal
		api.POST("/data/pruneEntities", pruneEntitiesHandler) // [ ] écrit dans Journal
		api.POST("/data/validate", validateHandler)           // [ ] écrit dans Journal

		api.GET("/data/purgeNotCompacted", purgeNotCompactedHandler) // [ ] écrit dans Journal

		api.GET("/data/etablissements", exportEtablissementsHandler) // [ ] écrit dans Journal
		api.GET("/data/entreprises", exportEntreprisesHandler)       // [ ] écrit dans Journal
	}

	bind := viper.GetString("APP_BIND")
	r.Run(bind)
}
