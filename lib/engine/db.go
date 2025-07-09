package engine

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/globalsign/mgo"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/spf13/viper"
)

// DB type centralisant les accès à une base de données
type DB struct {
	DB         *mgo.Database
	DBStatus   *mgo.Database
	PostgresDB *pgxpool.Pool
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
func InitDB() (DB, error) {
	loadConfig()
	dbDial := viper.GetString("DB_DIAL")
	dbDatabase := viper.GetString("DB")

	// définition de 2 connexions pour isoler les requêtes (TODO: utile ?)
	mongostatus, err := mgo.Dial(dbDial)
	if err != nil {
		return DB{}, fmt.Errorf("erreur de connexion (status) à MongoDB : %w", err)
	}
	mongostatus.SetSocketTimeout(72000 * time.Second)

	mongodb, err := mgo.Dial(dbDial)
	if err != nil {
		return DB{}, fmt.Errorf("erreur de connexion (data) à MongoDB : %w", err)
	}
	mongodb.SetSocketTimeout(72000 * time.Second)
	dbstatus := mongostatus.DB(dbDatabase)
	db := mongodb.DB(dbDatabase)

	// Création d'index sur la collection Admin, pour selection et tri de GetBatches()
	_ = db.C("Admin").EnsureIndex(mgo.Index{
		Key: []string{"_id.type", "_id.key"},
	})

	conn, err := pgxpool.New(context.Background(), viper.GetString("POSTGRES_DB_URL"))

	if err != nil {
		// TODO currently we don't want the PostgreSQL database to be mandatory,
		// e.g. for the e2e tests
		log.Printf("Erreur de connexion à la base de données PostgreSQL: %v", err)
	}

	return DB{
		DB:         db,
		DBStatus:   dbstatus,
		PostgresDB: conn,
	}, nil
}
