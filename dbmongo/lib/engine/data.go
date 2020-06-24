package engine

import (
	"errors"
	"fmt"
	"log"
	"sync"
	"time"

	"github.com/signaux-faibles/opensignauxfaibles/dbmongo/lib/misc"

	"github.com/globalsign/mgo"
	"github.com/globalsign/mgo/bson"
	"github.com/spf13/viper"
)

//go:generate go run js/loadJS.go

func loadJSFunctions(directoryNames ...string) (map[string]bson.JavaScript, error) {
	functions := make(map[string]bson.JavaScript)
	var err error

	// If encountering an error at following line, you probably forgot to
	// generate the file with "go generate" in ./lib/engine
	for k, v := range jsFunctions["common"] {
		functions[k] = bson.JavaScript{
			Code: string(v),
		}
	}

	for _, directoryName := range directoryNames {
		if _, ok := jsFunctions[directoryName]; !ok {
			err = errors.New("Map reduce javascript functions could not be found for " + directoryName)
		} else {
			err = nil
		}
		for k, v := range jsFunctions[directoryName] {
			functions[k] = bson.JavaScript{
				Code: string(v),
			}
		}
	}
	return functions, err
}

// PurgeNotCompacted permet de supprimer les objets non encore compactés
// c'est à dire, vider la collection ImportedData
func PurgeNotCompacted() error {
	_, err := Db.DB.C("ImportedData").RemoveAll(nil)
	return err
}

// MRWait centralise les variables nécessaires à l'isolation des traitements parallèlisés MR
type MRWait struct {
	waitGroup sync.WaitGroup
	running   sync.Map
	lock      sync.Mutex
	mergeLock sync.Mutex
}

func (w *MRWait) init() {
	w.waitGroup = sync.WaitGroup{}
	w.lock = sync.Mutex{}
	w.running = sync.Map{}
	w.running.Store("active", 0)
	w.running.Store("errors", 0)
	w.running.Store("total", 0)
}

// add incrémente le compteur désigné de la valeur choisie
// Retourne false si la valeur obtenue excède la valeur max
// Si max < 0 alors le test n'est pas effectué
func (w *MRWait) add(compteur string, val int, max int) bool {
	w.lock.Lock()
	defer w.lock.Unlock()
	total, _ := w.running.Load(compteur)
	if total.(int) < max || max < 0 {
		w.running.Store(compteur, total.(int)+val)
		return true
	}
	return false
}

// MRroutine travaille dans un pool pour exécuter des jobs de mapreduce. merge et nonAtomic recommandés.
func MRroutine(job *mgo.MapReduce, query bson.M, dbTemp string, collOrig string, w *MRWait, pipeChannel chan string) {
	w.add("total", 1, -1)

	for {
		ok := w.add("active", 1, viper.GetInt("MRthreads"))
		if ok {
			break
		}
		time.Sleep(time.Second)
	}
	fmt.Println(query)

	db, _ := mgo.Dial(viper.GetString("DB_DIAL"))
	db.SetSocketTimeout(720000 * time.Second)
	_, err := db.DB(viper.GetString("DB")).C(collOrig).Find(query).MapReduce(job, nil)

	if err == nil {
		pipeChannel <- dbTemp
	} else {
		fmt.Println(err)
		w.add("errors", 1, -1)
	}

	w.add("active", -1, -1)
	db.Close()
	w.waitGroup.Done()
}

// Compact traite le compactage de la base RawData
func Compact(batchKey string) error {
	// Détermination scope traitement
	batches, _ := GetBatches()

	var batchesID []string
	var completeTypes = make(map[string][]string)
	for _, b := range batches {
		completeTypes[b.ID.Key] = b.CompleteTypes
		batchesID = append(batchesID, b.ID.Key)
	}
	found := -1
	for ind, batchID := range batchesID {
		if batchID == batchKey {
			found = ind
			break
		}
	}
	// Si le numéro de batch n'est pas valide, erreur
	var batch AdminBatch
	if found == -1 {
		return errors.New("Le batch " + batchKey + "n'a pas été trouvé")
	}
	batch = batches[found]

	functions, err := loadJSFunctions("compact")
	if err != nil {
		return err
	}

	// Traitement MR
	job := &mgo.MapReduce{
		Map:      functions["map"].Code,
		Reduce:   functions["reduce"].Code,
		Finalize: functions["finalize"].Code,
		Out:      bson.M{"reduce": "RawData"},
		Scope: bson.M{
			"f":             functions,
			"batches":       batchesID,
			"completeTypes": completeTypes,
			"batchKey":      batchKey,
			"serie_periode": misc.GenereSeriePeriode(batch.Params.DateDebut, batch.Params.DateFin),
		},
	}

	chunks, err := ChunkCollection(viper.GetString("DB"), "RawData", viper.GetInt64("chunkByteSize"))

	for _, query := range chunks.ToQueries(nil, "value.key") {
		log.Println(query)
		_, err = Db.DB.C("ImportedData").Find(query).MapReduce(job, nil)
		if err != nil {
			return err
		}
	}

	err = PurgeNotCompacted()
	return err
}

type object struct {
	Key struct {
		Siret string `json:"key" bson:"key"`
		Batch string `json:"batch" bson:"batch"`
	} `json:"key"`
	Value map[string]interface{} `json:"value" bson:"value"`
	Scope []string               `json:"scope" value:"scope"`
}

// ToDatapi exports data from database to datapi instance
func ToDatapi(batchKey string) error {

	prediction := Db.DB.C("Public").Find(bson.M{"_id.batch": batchKey})
	predictions := prediction.Iter()
	var p interface{}

	for predictions.Next(&p) {
		fmt.Println(p)
	}
	// public := Db.DB.C("Public").Find(bson.M{"_id.batch": batchKey})

	return nil
}

// GetBatches retourne tous les objets AdminBatch de la base triés par ID
func GetBatches() ([]AdminBatch, error) {
	var batches []AdminBatch
	err := Db.DB.C("Admin").Find(bson.M{"_id.type": "batch"}).Sort("_id.key").All(&batches)
	return batches, err
}

// GetBatchesID retourne les clés des batches contenus en base
func GetBatchesID() []string {
	batches, _ := GetBatches()
	var batchesID []string
	for _, b := range batches {
		batchesID = append(batchesID, b.ID.Key)
	}
	return batchesID
}

// GetBatch retourne le batch correspondant à la clé batchKey
func GetBatch(batchKey string) (AdminBatch, error) {
	var batch AdminBatch
	err := Db.DB.C("Admin").Find(bson.M{"_id.type": "batch", "_id.key": batchKey}).One(&batch)
	return batch, err
}

// Purge réinitialise la base, à utiliser avec modération
func Purge() interface{} {
	infoImportedData, errImportedData := Db.DB.C("ImportedData").RemoveAll(nil)
	infoRawData, errRawData := Db.DB.C("RawData").RemoveAll(nil)
	infoJournal, errJournal := Db.DB.C("Journal").RemoveAll(nil)
	infoFeatures, errFeatures := Db.DB.C("Features").RemoveAll(nil)
	infoPublic, errPublic := Db.DB.C("Public").RemoveAll(nil)

	returnData := map[string]map[string]interface{}{
		"ImportedData": map[string]interface{}{
			"info":  infoImportedData,
			"error": errImportedData,
		},
		"RawData": map[string]interface{}{
			"info":  infoRawData,
			"error": errRawData,
		},
		"Journal": map[string]interface{}{
			"info":  infoJournal,
			"error": errJournal,
		},
		"Features": map[string]interface{}{
			"info":  infoFeatures,
			"error": errFeatures,
		},
		"Public": map[string]interface{}{
			"info":  infoPublic,
			"error": errPublic,
		},
	}

	return returnData
}

// Chunks est le retour de la fonction mongodb SplitKeys
type Chunks struct {
	OK        int `bson:"ok"`
	SplitKeys []struct {
		ID string `bson:"_id"`
	} `bson:"splitKeys"`
}

// ChunkCollection exécute la fonction SplitKeys sur la collection passée en paramètres
func ChunkCollection(db string, collection string, chunkSize int64) (Chunks, error) {
	var result Chunks

	err := Db.DB.Run(
		bson.D{{Name: "splitVector", Value: db + "." + collection},
			{Name: "keyPattern", Value: bson.M{"_id": 1}},
			{Name: "maxChunkSizeBytes", Value: chunkSize}},
		&result)

	return result, err
}

// ToQueries translates chunks into bson queries to chunk collection by siren code
func (chunks Chunks) ToQueries(query bson.M, field string) []bson.M {
	if len(chunks.SplitKeys) != 0 {
		var ret []bson.M
		ret = append(ret, bson.M{
			field: bson.M{
				"$lt": chunks.SplitKeys[0].ID[0:9],
			},
		})
		for i := 1; i < len(chunks.SplitKeys); i++ {
			ret = append(ret, bson.M{
				"$and": []bson.M{
					bson.M{field: bson.M{"$gte": chunks.SplitKeys[i-1].ID[0:9]}},
					bson.M{field: bson.M{"$lt": chunks.SplitKeys[i].ID[0:9]}},
					query,
				},
			})
		}
		ret = append(ret, bson.M{
			field: bson.M{
				"$gte": chunks.SplitKeys[len(chunks.SplitKeys)-1].ID[0:9],
			},
		})
		return ret
	} else {
		return []bson.M{query}
	}
}
