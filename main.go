package main

import (
	"github.com/signaux-faibles/opensignauxfaibles/lib/engine"

	"github.com/signaux-faibles/opensignauxfaibles/lib/naf"

	_ "github.com/signaux-faibles/opensignauxfaibles/docs"
)

// main Fonction Principale
func main() {
	engine.Db = engine.InitDB()
	go engine.MessageSocketAddClient()

	var err error
	naf.Naf, err = naf.LoadNAF()
	if err != nil {
		panic(err)
	}

	api.POST("/data/batch/purge", purgeBatchHandler)      // [x] écrit dans Journal
	api.POST("/data/check", checkBatchHandler)            // [x] écrit dans Journal
	api.POST("/data/import", importBatchHandler)          // [x] écrit dans Journal
	api.POST("/data/validate", validateHandler)           // [x] écrit dans Journal
	api.POST("/data/compact", compactHandler)             // [x] écrit dans Journal
	api.POST("/data/reduce", reduceHandler)               // [x] écrit dans Journal
	api.POST("/data/public", publicHandler)               // [x] écrit dans Journal
	api.POST("/data/pruneEntities", pruneEntitiesHandler) // [x] écrit dans Journal

	api.GET("/data/purgeNotCompacted", purgeNotCompactedHandler)

	api.GET("/data/etablissements", exportEtablissementsHandler)
	api.GET("/data/entreprises", exportEntreprisesHandler)
}
