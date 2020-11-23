package engine

import (
	"errors"
	"fmt"
	"os"

	"github.com/globalsign/mgo/bson"
	"github.com/signaux-faibles/opensignauxfaibles/dbmongo/lib/base"
	"github.com/signaux-faibles/opensignauxfaibles/dbmongo/lib/marshal"
	"github.com/spf13/viper"
)

// AdminAlgo décrit les qualités d'un algorithme
type AdminAlgo struct {
	ID          base.AdminID `json:"id" bson:"_id"`
	Label       string       `json:"label" bson:"label"`
	Description string       `json:"description" bson:"description"`
	Scope       []string     `json:"scope,omitempty" bson:"scope,omitempty"`
}

// Load charge un objet algo de la base
func (algo *AdminAlgo) Load(algoKey string) error {
	err := Db.DBStatus.C("Admin").Find(bson.M{"_id.type": "algo", "_id.key": algoKey}).One(algo)
	return err
}

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
func ImportBatch(batch base.AdminBatch, parsers []marshal.Parser, skipFilter bool) error {
	var cache = marshal.NewCache()
	filter, err := marshal.GetSirenFilter(cache, &batch)
	if err != nil {
		return err
	}
	if !skipFilter && filter == nil {
		return errors.New("Veuillez inclure un filtre")
	}
	for _, parser := range parsers {
		outputChannel, eventChannel := marshal.ParseFilesFromBatch(cache, &batch, parser) // appelle la fonction ParseFile() pour chaque type de fichier
		go RelayEvents(eventChannel, "ImportBatch")
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
			Db.ChanData <- &value
		}
	}

	Db.ChanData <- &Value{}
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
	for _, parser := range parsers {
		outputChannel, eventChannel := marshal.ParseFilesFromBatch(cache, &batch, parser) // appelle la fonction ParseFile() pour chaque type de fichier
		DiscardTuple(outputChannel)
		lastReport := RelayEvents(eventChannel, "CheckBatch")
		reports = append(reports, lastReport)
	}

	Db.ChanData <- &Value{}
	return reports, nil
}
