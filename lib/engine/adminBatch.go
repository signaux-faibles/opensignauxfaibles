package engine

import (
	"errors"
	"fmt"
	"os"
	"sync"
	"time"

	"github.com/globalsign/mgo/bson"
	"github.com/signaux-faibles/opensignauxfaibles/lib/base"
	"github.com/signaux-faibles/opensignauxfaibles/lib/marshal"
	"github.com/spf13/viper"
)

// Load charge les données d'un batch depuis la base de données
func Load(batch *base.AdminBatch, batchKey string) error {
	err := Db.DB.C("Admin").Find(bson.M{"_id.type": "batch", "_id.key": batchKey}).One(batch)
	return err
}

// Save écrit les données d'un batch vers la base de données
func Save(batch *base.AdminBatch) error {
	_, err := Db.DB.C("Admin").Upsert(bson.M{"_id": batch.ID}, batch)
	return err
}

// ImportBatch lance tous les parsers sur le batch fourni
func ImportBatch(batch base.AdminBatch, parsers []marshal.Parser, skipFilter bool, data chan *Value) error {
	var cache = marshal.NewCache()
	filter, err := marshal.GetSirenFilter(cache, &batch)
	if err != nil {
		return err
	}
	if !skipFilter && filter == nil {
		return errors.New("Veuillez inclure un filtre")
	}
	startDate := time.Now()
	var wg sync.WaitGroup
	for _, parser := range parsers {
		wg.Add(1)
		outputChannel, eventChannel := marshal.ParseFilesFromBatch(cache, &batch, parser) // appelle la fonction ParseFile() pour chaque type de fichier
		// Insert events (parsing logs) into the "Journal" collection
		go func() {
			defer wg.Done()
			RelayEvents(eventChannel, "ImportBatch", startDate)
		}()
		// Insert tuples (data) into the "ImportedData" collection
		for tuple := range outputChannel {
			hash := fmt.Sprintf("%x", GetMD5(tuple))
			value := Value{
				Value: Data{
					Scope: tuple.Scope(),
					Key:   tuple.Key(),
					Batch: map[string]Batch{
						batch.ID.Key: {
							tuple.Type(): map[string]marshal.Tuple{
								hash: tuple,
							}}}}}
			data <- &value
		}
	}
	wg.Wait() // wait for all events and tuples to be inserted
	FlushImportedData(data)
	return nil
}

// CheckBatchPaths checks if the filepaths of batch.Files exist
func CheckBatchPaths(batch *base.AdminBatch) error {
	var ErrorString string
	for _, filepaths := range batch.Files {
		for _, filepath := range filepaths {
			filepath = viper.GetString("APP_DATA") + filepath
			if _, err := os.Stat(filepath); err != nil {
				ErrorString += filepath + " is missing (" + err.Error() + ").\n"
			}
		}
	}
	if ErrorString != "" {
		return errors.New(ErrorString)
	}
	return nil

}

// CheckBatch checks batch
func CheckBatch(batch base.AdminBatch, parsers []marshal.Parser) (reports []string, err error) {
	if err := CheckBatchPaths(&batch); err != nil {
		return nil, err
	}
	var cache = marshal.NewCache()
	startDate := time.Now()
	for _, parser := range parsers {
		outputChannel, eventChannel := marshal.ParseFilesFromBatch(cache, &batch, parser) // appelle la fonction ParseFile() pour chaque type de fichier
		DiscardTuple(outputChannel)
		lastReport := RelayEvents(eventChannel, "CheckBatch", startDate)
		reports = append(reports, lastReport)
	}

	return reports, nil
}
