package engine

import (
	"log"
	"time"

	"github.com/globalsign/mgo"
	"github.com/spf13/viper"
)

// DB type centralisant les accès à une base de données
type DB struct {
	DB       *mgo.Database
	DBStatus *mgo.Database
}

// InitDB Initialisation de la connexion MongoDB
func InitDB() DB {
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

	return DB{
		DB:       db,
		DBStatus: dbstatus,
	}
}
