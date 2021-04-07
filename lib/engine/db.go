package engine

import (
	"fmt"
	"log"
	"sync"
	"time"

	"github.com/globalsign/mgo"
	"github.com/globalsign/mgo/bson"
	"github.com/signaux-faibles/opensignauxfaibles/lib/base"
	"github.com/spf13/viper"
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
		log.Println(err) // /!\ en ci, seul le fichier config-sample.toml existe
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

	// Création d'index sur la collection Admin, pour selection et tri de GetBatches()
	db.C("Admin").EnsureIndex(mgo.Index{
		Key: []string{"_id.type", "_id.key"},
	})

	if err = CreateImportedDataCollection(db, "ImportedData"); err != nil {
		log.Fatal("échec d'initialisation de ImportedData: " + err.Error())
	}

	firstBatchID := viper.GetString("FIRST_BATCH")
	if !base.IsBatchID(firstBatchID) {
		panic("Paramètre FIRST_BATCH incorrect, vérifiez la configuration.")
	}

	db.C("RawData").Create(&mgo.CollectionInfo{})

	// Création d'index sur la collection RawData, pour le filtrage du map-reduce de Public et Reduce
	db.C("RawData").EnsureIndex(mgo.Index{
		Name: "algo2",                        // trouvé sur la db de prod
		Key:  []string{"-value.index.algo2"}, // booléen
	})

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

	return DB{
		DB:       db,
		DBStatus: dbstatus,
	}
}

// CreateImportedDataCollection créée une collection "ImportedData" avec index et
// validation de documents.
func CreateImportedDataCollection(db *mgo.Database, colName string) error {
	// Création d'index sur la collection ImportedData, pour le découpage du map-reduce de Compact
	db.C(colName).EnsureIndex(mgo.Index{
		Name: "value.key_1",         // trouvé sur la db de prod
		Key:  []string{"value.key"}, // numéro SIRET ou SIREN
	})

	// Injection du schéma de validation de données JSON dans ImportedData
	jsonSchemas, err := LoadJSONSchemaFiles()
	if err != nil {
		return err
	}
	schemaPerHashedDataType := MakeValidationSchemaPerHashedDataType(jsonSchemas)
	jsonSchema := MakeValidationSchemaForImportedData(schemaPerHashedDataType)
	return setupDocValidation(db, colName, jsonSchema)
}

// setupDocValidation configure la validation de documents pour une collection existante.
func setupDocValidation(db *mgo.Database, colName string, jsonSchema bson.M) error {
	var validRes struct {
		Ok            bool          `bson:"ok" json:"ok"`
		Errmsg        string        `bson:"errmsg" json:"errmsg"`
		Code          int           `bson:"code" json:"code"`
		CodeName      string        `bson:"codeName" json:"codeName"`
		OperationTime time.Duration `bson:"operationTime" json:"operationTime"`
	}
	db.Run(bson.D{
		{Name: "collMod", Value: colName},
		{Name: "validator", Value: bson.M{"$jsonSchema": jsonSchema}},
		{Name: "validationLevel", Value: "strict"},
		{Name: "validationAction", Value: "error"}}, &validRes)
	if !validRes.Ok {
		return fmt.Errorf("error %v while trying to setup doc validation on %v: %v", validRes.Code, colName, validRes.Errmsg)
	}
	return nil
}

var importing sync.WaitGroup

// InsertIntoImportedData retourne un canal dont les objets seront ajoutés à
// la collection ImportedData, par paquets de 100.
func InsertIntoImportedData(db *mgo.Database) chan *Value {
	importing.Add(1)
	source := make(chan *Value, 10)

	go func(chan *Value) {
		defer importing.Done()
		buffer := make(map[string]*Value)
		objects := make([]interface{}, 0)
		i := 0
		insertObjectsIntoImportedData := func() {
			for _, v := range buffer {
				objects = append(objects, *v)
			}
			if len(objects) > 0 {
				if err := db.C("ImportedData").Insert(objects...); err != nil {
					log.Println("Erreur lors de l'insertion de certains documents dans ImportedData: " + err.Error()) // ex: document invalide, cf CreateImportedDataCollection()
				}
			}
			buffer = make(map[string]*Value)
			objects = make([]interface{}, 0)
			i = 0
		}

		for value := range source {
			if i >= 100 {
				insertObjectsIntoImportedData()
			}
			if knownValue, ok := buffer[value.Value.Key]; ok {
				newValue, _ := (*knownValue).Merge(*value)
				buffer[value.Value.Key] = &newValue
			} else {
				value.ID = bson.NewObjectId()
				buffer[value.Value.Key] = value
				i++
			}
		}
		// le canal a été fermé => importer les données restantes avant de rendre la main
		insertObjectsIntoImportedData()
	}(source)

	return source
}

// FlushImportedData finalise l'insertion des données dans ImportedData.
func FlushImportedData(channel chan *Value) {
	close(channel)
	importing.Wait()
}
