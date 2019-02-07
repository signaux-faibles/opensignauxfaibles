package main

import (
	"dbmongo/lib/engine"
	"fmt"

	"github.com/globalsign/mgo/bson"

	"net/http"
	"time"

	jwt "github.com/appleboy/gin-jwt"
	"github.com/gin-contrib/cors"
	"github.com/gin-contrib/static"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"github.com/spf13/viper"

	ginSwagger "github.com/swaggo/gin-swagger"
	"github.com/swaggo/gin-swagger/swaggerFiles"

	_ "./docs"
)

var wsupgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin:     checkOrigin,
}

func checkOrigin(r *http.Request) bool {
	return true
}

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
// @title API openSignauxFaibles
// @version 1.1
// @description Cette API centralise toutes les fonctionnalités du module de traitement de données OpenSignauxFaibles
// @description Pour plus de renseignements: https://beta.gouv.fr/startups/signaux-faibles.html
// @license.name Licence MIT
// @license.url https://raw.githubusercontent.com/entrepreneur-interet-general/opensignauxfaibles/master/LICENSE
// @BasePath /
// @securityDefinitions.apikey ApiKeyAuth
// @in header
// @name Authorization
func main() {
	// Lancer Rserve en background

	// go r()
	engine.Db = engine.InitDB()
	go engine.MessageSocketAddClient()

	r := gin.New()

	if !viper.GetBool("DEV") {
		gin.SetMode(gin.ReleaseMode)
	} else {
		r.Use(gin.Recovery())
		r.Use(gin.Logger())
	}

	config := cors.DefaultConfig()
	if viper.GetBool("DEV") {
		config.AllowOrigins = []string{"http://localhost:8080"}
	} else {
		config.AllowOrigins = []string{"https://signaux.faibles.fr"}
	}
	config.AddAllowHeaders("Authorization")
	config.AddAllowMethods("GET", "POST")
	r.Use(cors.New(config))

	authMiddleware, err := jwt.New(&jwt.GinJWTMiddleware{
		Realm:           "Signaux-Faibles",
		Key:             []byte(viper.GetString("jwtSecret")),
		SendCookie:      false,
		Timeout:         time.Hour,
		MaxRefresh:      time.Hour,
		IdentityKey:     "id",
		PayloadFunc:     payload,
		IdentityHandler: identityHandler,
		Authenticator:   authenticator,
		Authorizator:    authorizator,
		Unauthorized:    unauthorizedHandler,
		TokenLookup:     "header: Authorization, query: token, cookie: jwt",
		TokenHeadName:   "Bearer",
		TimeFunc:        time.Now,
	})

	if err != nil {
		panic("Erreur lors de la mise en place de l'authentification:" + err.Error())
	}

	r.Use(static.Serve("/", static.LocalFile("static/", true)))
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	r.POST("/login", authMiddleware.LoginHandler)
	r.POST("/login/get", loginGetHandler)
	r.POST("/login/check", loginCheckHandler)
	r.POST("/login/recovery/get", getRecoveryEmailHandler)
	r.POST("/login/recovery/setPassword", checkRecoverySetPassword)
	r.GET("/ws/:jwt", func(c *gin.Context) {
		wshandler(c.Writer, c.Request, c.Params.ByName("jwt"))
	})

	api := r.Group("api")
	api.Use(authMiddleware.MiddlewareFunc())

	{
		api.GET("/refreshToken", authMiddleware.RefreshHandler)

		api.POST("/admin/batch", upsertBatchHandler)
		api.GET("/admin/batch", listBatchHandler)
		api.GET("/admin/files", adminFilesHandler)
		api.GET("/admin/types", listTypesHandler)

		api.GET("/admin/features", adminFeature)

		api.GET("/admin/getLogs", getLogsHandler)

		api.GET("/admin/batch/next", nextBatchHandler)
		api.GET("/admin/batch/process", processBatchHandler)
		api.POST("/admin/files", addFile)

		api.GET("/admin/batch/revert", revertBatchHandler)

		api.GET("/data/naf", nafHandler)
		api.GET("/data/batch/purge", purgeBatchHandler)
		api.POST("/data/import", importBatchHandler)
		api.GET("/data/compact", compactHandler)
		api.POST("/data/reduce", reduceHandler)
		api.POST("/data/search", searchRaisonSocialeHandler)
		api.GET("/data/purge", purgeHandler)
		api.GET("/data/purgeNotCompacted", deleteHandler)
		api.GET("/data/public/:batch", publicHandler)

		api.GET("/dashboard/tasks", getTasksHandler)
	}

	bind := viper.GetString("APP_BIND")
	r.Run(bind)
}

func deleteHandler(c *gin.Context) {
	var result []interface{}
	engine.Db.DB.C("RawData").RemoveAll(bson.M{"_id": bson.M{"$type": "objectId"}})
	c.JSON(200, result)
}

func getTasksHandler(c *gin.Context) {

}
