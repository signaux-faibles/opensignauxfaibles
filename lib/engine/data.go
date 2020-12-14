package engine

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"sort"
	"sync"
	"time"

	"github.com/signaux-faibles/opensignauxfaibles/lib/base"
	"github.com/signaux-faibles/opensignauxfaibles/lib/marshal"
	"github.com/signaux-faibles/opensignauxfaibles/lib/misc"

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

// PruneEntities permet de compter puis supprimer les entités de RawData
// qui auraient du être exclues par le Filtre de périmètre SIREN.
func PruneEntities(batchKey string, delete bool) (int, error) {
	startDate := time.Now()

	// Récupérer le batch
	batch := base.AdminBatch{}
	err := Load(&batch, batchKey)
	if err != nil {
		return -1, errors.New("Batch inexistant: " + err.Error())
	}
	// Charger le filtre
	var cache = marshal.NewCache()
	filter, err := marshal.GetSirenFilter(cache, &batch)
	if err != nil {
		return -1, err
	}
	if filter == nil {
		return -1, errors.New("Ce batch ne spécifie pas de filtre")
	}
	// Créer une expression régulière pour reconnaitre les SIRENs du périmètre
	sirens := []string{}
	for siren := range filter {
		sirens = append(sirens, siren)
	}
	// Lister les entités de RawData qui ne figurent pas dans le filtre
	pipeline := []bson.M{
		{
			"$project": bson.M{
				"_id":   true,
				"siren": bson.M{"$substr": []interface{}{"$_id", 0, 9}},
			},
		},
		{
			"$match": bson.M{
				"siren": bson.M{"$nin": sirens},
			},
		},
		{
			"$project": bson.M{
				"_id":   true,
				"siren": false,
			},
		},
	}
	// Éventuellement, supprimer ces entités
	if delete == true {
		// Enregistrer la liste d'entités dans la collection temporaire "EntitiesToPrune"
		tmpCollection := "EntitiesToPrune"
		pipeline = append(pipeline, bson.M{"$out": tmpCollection})
		iterator := Db.DB.C("RawData").Pipe(pipeline).AllowDiskUse().Iter()
		var item struct {
			ID string `json:"id"   bson:"_id"`
		}
		// Supprimer les entités en itérant sur "EntitiesToPrune"
		var nbDeleted = 0
		iterator = Db.DB.C(tmpCollection).Find(bson.M{}).Iter()
		for iterator.Next(&item) {
			if err := iterator.Err(); err != nil {
				return -1, err
			}
			if err = Db.DB.C("RawData").Remove(bson.M{"_id": item.ID}); err != nil {
				return -1, err
			}
			nbDeleted++
		}
		// Supprimer la collection temporaire
		_ = Db.DB.C(tmpCollection).DropCollection()
		return nbDeleted, err
	}
	var result struct {
		IdsToDelete int `json:"ids_to_delete"   bson:"ids_to_delete"`
	}
	pipeline = append(pipeline, bson.M{"$count": "ids_to_delete"})
	err = Db.DB.C("RawData").Pipe(pipeline).AllowDiskUse().One(&result)
	LogOperationEvent("PruneEntities", startDate)
	return result.IdsToDelete, err
}

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
func Compact(fromBatchKey string) error {
	startDate := time.Now()

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
		if batchID == fromBatchKey {
			found = ind
			break
		}
	}
	// Si le numéro de batch n'est pas valide, erreur
	var batch base.AdminBatch
	if found == -1 {
		return errors.New("Le batch " + fromBatchKey + "n'a pas été trouvé")
	}
	batch = batches[found]

	functions, err := loadJSFunctions("compact")
	if err != nil {
		return err
	}

	// Traitement MR
	job := &mgo.MapReduce{
		Map:      functions["map"].Code,
		Reduce:   functions["reduce"].Code, // 1st pass: reduce() will be called on the documents of ImportedData
		Finalize: functions["finalize"].Code,
		Out:      bson.M{"reduce": "RawData"}, // 2nd pass: for each siret/siren, reduce() will be called on the current RawData document and the result of the 1st pass, to merge the the new data with the existing data
		Scope: bson.M{
			"f":             functions,
			"batches":       batchesID,
			"completeTypes": completeTypes,
			"fromBatchKey":  fromBatchKey,
			"serie_periode": misc.GenereSeriePeriode(batch.Params.DateDebut, batch.Params.DateFin),
		},
	}

	chunks, err := ChunkCollection(viper.GetString("DB"), "RawData", viper.GetInt64("chunkByteSize"))
	if err != nil {
		return err
	}

	for _, query := range chunks.ToQueries(nil, "value.key") {
		log.Println(query)
		_, err = Db.DB.C("ImportedData").Find(query).MapReduce(job, nil)
		if err != nil {
			return err
		}
	}

	err = PurgeNotCompacted()

	LogOperationEvent("Compact", startDate)

	return err
}

// GetBatches retourne tous les objets base.AdminBatch de la base triés par ID
func GetBatches() ([]base.AdminBatch, error) {
	var batches []base.AdminBatch
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

// ExportEtablissements exporte les établissements dans un fichier.
func ExportEtablissements(key string) error {
	pipeline := GetEtablissementWithScoresPipeline(key)
	iter := Db.DB.C("Public").Pipe(pipeline).AllowDiskUse().Iter()
	return storeMongoPipelineResults(iter)
}

// ExportEntreprises exporte les entreprises dans un fichier.
func ExportEntreprises(key string) error {
	pipeline := GetEntreprisePipeline(key)
	iter := Db.DB.C("Public").Pipe(pipeline).AllowDiskUse().Iter()
	return storeMongoPipelineResults(iter)
}

// ValidateDataEntries affiche les entrées de données invalides détectées dans la collection spécifiée.
func ValidateDataEntries(jsonSchema map[string]bson.M, collection string) error {
	startDate := time.Now()

	w := sync.WaitGroup{}
	gzipWriter := getItemChannelToStdout(&w)

	// lister les entrées de données non définies (type: undefined au lieu de object)
	pipeline, err := GetUndefinedDataValidationPipeline()
	if err != nil {
		return err
	}
	err = iterateToChannel(gzipWriter, Db.DB.C(collection).Pipe(pipeline).AllowDiskUse().Iter())
	if err != nil {
		return err
	}

	// lister les entrées de données non conformes aux modèles JSON Schema
	pipeline, err = GetDataValidationPipeline(jsonSchema)
	if err != nil {
		return err
	}
	err = iterateToChannel(gzipWriter, Db.DB.C(collection).Pipe(pipeline).AllowDiskUse().Iter())
	if err != nil {
		return err
	}

	close(gzipWriter)
	w.Wait()

	LogOperationEvent("ValidateDataEntries", startDate)

	return nil
}

func iterateToChannel(channel chan interface{}, iterator *mgo.Iter) error {
	var item interface{}
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
