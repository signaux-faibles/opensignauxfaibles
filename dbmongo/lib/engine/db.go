package engine

import (
	"log"
	"time"

	"github.com/globalsign/mgo"
	"github.com/globalsign/mgo/bson"
	"github.com/signaux-faibles/opensignauxfaibles/dbmongo/lib/base"
	"github.com/spf13/viper"
)

// Status statut de la base de données
type Status struct {
	ID     base.AdminID  `json:"id" bson:"_id"`
	Status *string       `json:"status" bson:"status"`
	Epoch  int           `json:"epoch" bson:"epoch"`
	DB     *mgo.Database `json:"-" bson:"-"`
}

// DB type centralisant les accès à une base de données
type DB struct {
	DB       *mgo.Database
	DBStatus *mgo.Database
	Status   Status
	ChanData chan *Value
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
	err := viper.ReadInConfig()
	if err != nil {
		panic("Erreur à la lecture de la configuration")
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
		log.Println("Erreur de connexion (status) à MongoDB")
		log.Panic(err)
	}
	mongostatus.SetSocketTimeout(72000 * time.Second)

	mongodb, err := mgo.Dial(dbDial)
	if err != nil {
		log.Println("Erreur de connexion (data) à MongoDB")
		log.Panic(err)
	}
	mongodb.SetSocketTimeout(72000 * time.Second)
	dbstatus := mongostatus.DB(dbDatabase)
	db := mongodb.DB(dbDatabase)

	// Création d'index sur la collection ImportedData, pour le découpage du map-reduce de Compact
	db.C("ImportedData").EnsureIndex(mgo.Index{
		Key:      []string{"value.key"}, // numéro SIRET ou SIREN
	})

	firstBatchID := viper.GetString("FIRST_BATCH")
	if !base.IsBatchID(firstBatchID) {
		panic("Paramètre FIRST_BATCH incorrect, vérifiez la configuration.")
	}

	db.C("RawData").Create(&mgo.CollectionInfo{})

	// firstBatch, err := getBatch(db, firstBatchID)
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
			panic("Impossible de créer le premier batch: " + err.Error())
		}
	}

	chanData := insert(db)

	// envoie un struct vide pour purger les channels au cas où il reste les objets non insérés
	go func() {
		for range time.Tick(1 * time.Second) {
			chanData <- &Value{}
		}
	}()

	dbConnect := DB{
		DB:       db,
		DBStatus: dbstatus,
		ChanData: chanData,
	}

	return dbConnect
}

func insert(db *mgo.Database) chan *Value {
	source := make(chan *Value, 10)

	go func(chan *Value) {
		buffer := make(map[string]*Value)
		objects := make([]interface{}, 0)
		i := 0

		for value := range source {
			if (value.Value.Batch == nil) || i >= 100 {
				for _, v := range buffer {
					objects = append(objects, *v)
				}
				if len(objects) > 0 {
					db.C("ImportedData").Insert(objects...)
				}
				buffer = make(map[string]*Value)
				objects = make([]interface{}, 0)
				i = 0
			}
			if value.Value.Batch != nil {
				if knownValue, ok := buffer[value.Value.Key]; ok {
					newValue, _ := (*knownValue).Merge(*value)
					buffer[value.Value.Key] = &newValue
				} else {
					value.ID = bson.NewObjectId()
					buffer[value.Value.Key] = value
					i++
				}
			}
		}
	}(source)

	return source
}
