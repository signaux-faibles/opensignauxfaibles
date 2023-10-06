package engine

import (
	"log"
	"sync"
	"time"

	"github.com/globalsign/mgo"
	"github.com/globalsign/mgo/bson"
	"github.com/spf13/viper"

	"opensignauxfaibles/lib/base"
)

// DB type centralisant les accès à une base de données
type DB struct {
	DB       *mgo.Database
	DBStatus *mgo.Database
}

func loadConfig() {
	// Note: viper.SetConfigType() and viper.AddConfigPath() are called by initConfig()
	viper.SetDefault("APP_DATA", "$HOME/data-raw/")
	viper.SetDefault("DB", "opensignauxfaibles")
	err := viper.ReadInConfig()
	if err != nil {
		log.Fatal(err) // /!\ en ci, seul le fichier config-sample.toml existe
	}
}

// InitDB Initialisation de la connexion MongoDB
func InitDB() DB {
	loadConfig()
	dbDial := viper.GetString("DB_DIAL")
	dbDatabase := viper.GetString("DB")

	// définition de 2 connexions pour isoler les requêtes (TODO: utile ?)
	mongostatus, err := mgo.Dial(dbDial)
	if err != nil {
		log.Fatal("Erreur de connexion (status) à MongoDB")
	}
	mongostatus.SetSocketTimeout(72000 * time.Second)

	mongodb, err := mgo.Dial(dbDial)
	if err != nil {
		log.Fatal("Erreur de connexion (data) à MongoDB")
	}
	mongodb.SetSocketTimeout(72000 * time.Second)
	dbstatus := mongostatus.DB(dbDatabase)
	db := mongodb.DB(dbDatabase)

	// Création d'index sur la collection Admin, pour selection et tri de GetBatches()
	_ = db.C("Admin").EnsureIndex(mgo.Index{
		Key: []string{"_id.type", "_id.key"},
	})

	firstBatchID := viper.GetString("FIRST_BATCH")
	if !base.IsBatchID(firstBatchID) {
		log.Fatal("Paramètre FIRST_BATCH incorrect, vérifiez la configuration.")
	}

	db.C("RawData").Create(&mgo.CollectionInfo{})

	// Création d'index sur la collection RawData, pour le filtrage du map-reduce de Public et Reduce
	db.C("RawData").EnsureIndex(mgo.Index{
		Name: "algo2",                        // trouvé sur la db de prod
		Key:  []string{"-value.index.algo2"}, // booléen
	})

	var firstBatch base.AdminBatch
	db.C("Admin").Find(bson.M{"_id.type": "batch", "_id.key": firstBatchID}).One(&firstBatch)
	// Si la table Admin n'existe pas, elle sera créée lors de l'insertion, ci-dessous

	if firstBatch.ID.Type == "" {
		firstBatch = base.AdminBatch{
			ID: base.AdminID{
				Key:  firstBatchID,
				Type: "batch",
			},
		}

		err := db.C("Admin").Insert(firstBatch)
		if err != nil {
			log.Fatal("Impossible de créer le premier batch: " + err.Error())
		}
	}

	return DB{
		DB:       db,
		DBStatus: dbstatus,
	}
}

var importing sync.WaitGroup
