package engine

import (
	"encoding/json"
	"fmt"
	"log"
	"sort"
	"sync"
	"time"

	"opensignauxfaibles/lib/base"

	"github.com/globalsign/mgo"
	"github.com/globalsign/mgo/bson"
	"github.com/spf13/viper"
)

// MRWait centralise les variables nécessaires à l'isolation des traitements parallèlisés MR
type MRWait struct {
	waitGroup sync.WaitGroup
	running   sync.Map
	lock      sync.Mutex
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
func MRroutine(job mgo.MapReduce, query bson.M, dbTemp string, collOrig string, w *MRWait, pipeChannel chan string) {
	w.add("total", 1, -1)

	for {
		ok := w.add("active", 1, viper.GetInt("MRthreads"))
		if ok {
			break
		}
		time.Sleep(time.Second)
	}
	log.Println(query) // TODO: supprimer cet affichage ?

	db, err := mgo.Dial(viper.GetString("DB_DIAL"))
	if err != nil {
		log.Println("erreur de connection pendant le MRroutine: " + err.Error())
	}
	db.SetSocketTimeout(720000 * time.Second)

	_, err = db.DB(viper.GetString("DB")).C(collOrig).Find(query).MapReduce(&job, nil)

	if err == nil {
		pipeChannel <- dbTemp
	} else {
		fmt.Println(err) // TODO: supprimer cet affichage ?
		w.add("errors", 1, -1)
	}

	w.add("active", -1, -1)
	db.Close()
	w.waitGroup.Done()
}

// GetBatches retourne tous les objets base.AdminBatch de la base triés par ID
func GetBatches() ([]base.AdminBatch, error) {
	var batches []base.AdminBatch
	err := Db.DB.C("Admin").Find(bson.M{"_id.type": "batch"}).Sort("_id.key").All(&batches)
	return batches, err
}

// GetBatch retourne le batch correspondant à la clé batchKey
func GetBatch(batchKey string) (base.AdminBatch, error) {
	var batch base.AdminBatch
	err := Load(&batch, batchKey)
	return batch, err
}

type splitKey struct {
	ID string `bson:"_id"`
}

// Chunks est le retour de la fonction mongodb SplitKeys
type Chunks struct {
	OK        int        `bson:"ok"`
	SplitKeys []splitKey `bson:"splitKeys"`
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
	// la base n'a pas besoin de split
	if len(chunks.SplitKeys) == 0 {
		return []bson.M{query}
	}

	// creation des clés de partage sans doublons
	var splitKeysMap = make(map[string]struct{})
	for i := 0; i < len(chunks.SplitKeys); i++ {
		splitKey := chunks.SplitKeys[i].ID[0:9]
		splitKeysMap[splitKey] = struct{}{}
	}
	var splitKeys []string
	for k := range splitKeysMap {
		splitKeys = append(splitKeys, k)
	}
	sort.Strings(splitKeys)

	// creation des requêtes à partir des clés de split
	var ret []bson.M
	ret = append(ret, bson.M{
		field: bson.M{
			"$lt": splitKeys[0],
		},
	})
	for i := 1; i < len(splitKeys); i++ {
		ret = append(ret, bson.M{
			"$and": []bson.M{
				{field: bson.M{"$gte": splitKeys[i-1]}},
				{field: bson.M{"$lt": splitKeys[i]}},
				query,
			},
		})
	}
	ret = append(ret, bson.M{
		field: bson.M{
			"$gte": splitKeys[len(splitKeys)-1],
		},
	})
	return ret
}

func getItemChannelToStdout(wait *sync.WaitGroup) chan interface{} {
	c := make(chan interface{})
	wait.Add(1)
	go func() {
		defer wait.Done()
		for item := range c {
			printJSON(item)
		}
	}()
	return c
}

func printJSON(object interface{}) {
	res, _ := json.Marshal(object)
	fmt.Println(string(res))
}

// ValidateDataEntries affiche les entrées de données invalides détectées dans la collection spécifiée.
func ValidateDataEntries(jsonSchema map[string]bson.M, collection string) error {
	startDate := time.Now()

	w := sync.WaitGroup{}
	writer := getItemChannelToStdout(&w)

	// lister les entrées de données non définies (type: undefined au lieu de object)
	pipeline, err := GetUndefinedDataValidationPipeline()
	if err != nil {
		return err
	}
	err = iterateToChannel(writer, Db.DB.C(collection).Pipe(pipeline).AllowDiskUse().Iter())
	if err != nil {
		return err
	}

	// lister les entrées de données non conformes aux modèles JSON Schema
	pipeline, err = GetDataValidationPipeline(jsonSchema)
	if err != nil {
		return err
	}
	err = iterateToChannel(writer, Db.DB.C(collection).Pipe(pipeline).AllowDiskUse().Iter())
	if err != nil {
		return err
	}

	close(writer)
	w.Wait()

	LogOperationEvent("ValidateDataEntries", startDate)

	return nil
}

func iterateToChannel(channel chan interface{}, iterator *mgo.Iter) error {
	var item interface{}
	if err := iterator.Err(); err != nil {
		return err // e.g. "Erreur: Unknown $jsonSchema keyword: _TODO"
	}
	for iterator.Next(&item) {
		if err := iterator.Err(); err != nil {
			return err
		}
		channel <- item
	}
	return nil
}

func storeMongoPipelineResults(iterator *mgo.Iter) error {
	wait := sync.WaitGroup{}
	gzipWriter := getItemChannelToStdout(&wait)
	err := iterateToChannel(gzipWriter, iterator)
	if err != nil {
		return err
	}
	close(gzipWriter)
	wait.Wait()
	return nil
}
