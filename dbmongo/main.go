package main

import (
	"fmt"

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

var db = initDB()

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
	channel := make(chan socketMessage)
	addClientChannel <- channel

	for event := range channel {
		conn.WriteJSON(event)
	}
}

const identityKey = "id"

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
	go messageSocketAddClient()

	r := gin.New()

	if !viper.GetBool("DEV") {
		gin.SetMode(gin.ReleaseMode)
	} else {
		r.Use(gin.Recovery())
		r.Use(gin.Logger())
	}

	config := cors.DefaultConfig()
	config.AllowOrigins = []string{"http://localhost:3000", "http://localhost:8080", "https://signaux.faibles.fr"}
	config.AddAllowHeaders("Authorization")
	config.AddAllowMethods("GET", "POST")
	r.Use(cors.New(config))

	authMiddleware, err := jwt.New(&jwt.GinJWTMiddleware{
		Realm:           "Signaux-Faibles",
		Key:             []byte(viper.GetString("jwtSecret")),
		SendCookie:      false,
		Timeout:         time.Hour,
		MaxRefresh:      time.Hour,
		IdentityKey:     identityKey,
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

		api.POST("/admin/batch", upsertBatch)
		api.GET("/admin/batch", listBatch)
		api.GET("/admin/files", adminFiles)
		api.GET("/admin/types", listTypes)
		// api.GET("/admin/clone/:to", cloneDB)
		api.GET("/admin/features", adminFeature)
		api.GET("/admin/status", getDBStatus)
		api.GET("/admin/getLogs", getLogsHandler)
		api.GET("/admin/epoch", epoch)
		api.GET("/admin/batch/next", nextBatchHandler)
		api.GET("/admin/batch/process", processBatchHandler)
		api.POST("/admin/files", addFile)

		api.GET("/admin/batch/revert", revertBatchHandler)

		api.GET("/data/naf", getNAF)
		api.GET("/data/batch/purge", purgeBatchHandler)
		api.GET("/data/import/:batch", importBatchHandler)
		api.GET("/data/compact", compactHandler)
    api.GET("/data/reduce/:algo/:batchKey", reduceHandler)
    api.GET("/data/reduce/:algo/:batchKey/:key", reduceHandler)
		api.POST("/data/search", searchRaisonSociale)
		api.GET("/data/purge", purge)

		api.GET("/data/public/:batch", publicHandler)

		api.GET("/data/purgeNotCompacted", deleteHandler)
		api.GET("/dashboard/tasks", getTasks)
	}

	bind := viper.GetString("APP_BIND")
	r.Run(bind)
}

func loadConfig() {
	viper.SetConfigName("config")
	viper.SetConfigType("toml")
	viper.AddConfigPath("/etc/opensignauxfaibles")
	viper.AddConfigPath("$HOME/.opensignauxfaibles")
	viper.AddConfigPath(".")
	viper.SetDefault("APP_BIND", ":3000")
	viper.SetDefault("APP_DATA", "$HOME/data-raw/")
	viper.SetDefault("DB_HOST", "127.0.0.1")
	viper.SetDefault("DB_PORT", "27017")
	viper.SetDefault("DB", "opensignauxfaibles")
	viper.SetDefault("JWT_SECRET", "Secret à changer")
	// viper.SetDefault("KANBOARD_ENDPOINT", "http://localhost/kanboard/jsonrpc.php")
	// viper.SetDefault("KANBOARD_USERNAME", "admin")
	// viper.SetDefault("KANBOARD_PASSWORD", "admin")
	err := viper.ReadInConfig()
	if err != nil {
		panic("Erreur à la lecture de la configuration")
	}
}
