package engine

import (
	"errors"
	"fmt"
	"os"
	"time"

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

// NextBatchID génère le batchKey suivant à partir d'un batchKey
func NextBatchID(batchID string) (string, error) {
	batchTime, err := time.Parse("0601", batchID)
	if err != nil {
		return "", err
	}
	nextBatchTime := time.Date(batchTime.Year(), time.Month(batchTime.Month()+1), 1, 0, 0, 0, 0, time.UTC)
	return nextBatchTime.Format("0601"), err
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
		outputChannel, eventChannel := parser(cache, &batch)
		go RelayEvents(eventChannel)
		for tuple := range outputChannel {
			hash := fmt.Sprintf("%x", GetMD5(tuple))
			value := Value{
				Value: Data{
					Scope: tuple.Scope(),
					Key:   tuple.Key(),
					Batch: map[string]Batch{
						batch.ID.Key: Batch{
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
		outputChannel, eventChannel := parser(cache, &batch)
		DiscardTuple(outputChannel)
		lastReport := RelayEvents(eventChannel)
		reports = append(reports, lastReport)
	}

	Db.ChanData <- &Value{}
	return reports, nil
}

// ProcessBatch traitement ad-hoc modifiable pour les besoins du développement
func ProcessBatch(batchList []string, parsers []marshal.Parser) error {

	for _, v := range batchList {
		batch, errBatch := GetBatch(v)
		if errBatch != nil {
			return errors.New("Erreur de lecture du batch: " + errBatch.Error())
		}
		importErr := ImportBatch(batch, parsers, false)
		if importErr != nil {
			return importErr
		}
		time.Sleep(5 * time.Second) // TODO: trouver une façon de synchroniser l'insert des paquets
		err := Compact(v)
		if err != nil {
			return errors.New("Erreur de compactage: " + err.Error())
		}
	}

	batch := LastBatch()
	return Reduce(batch, "algo2", []string{"all"})
}

// LastBatch retourne le dernier batch
func LastBatch() base.AdminBatch {
	batches, _ := GetBatches()
	l := len(batches)
	batch := batches[l-1]
	return batch
}

// NextBatch crée le batch suivant le dernier batch existant
func NextBatch() error {
	batch := LastBatch()
	newBatchID, err := NextBatchID(batch.ID.Key)
	if err != nil {
		return fmt.Errorf("Mauvais numéro de batch: " + err.Error())
	}
	newBatch := base.AdminBatch{
		ID: base.AdminID{
			Key:  newBatchID,
			Type: "batch",
		},
		CompleteTypes: batch.CompleteTypes,
	}

	batch.Readonly = true

	err = Save(&batch)
	if err != nil {
		return fmt.Errorf("Erreur readonly Batch: " + err.Error())
	}

	err = Save(&newBatch)
	if err != nil {
		return fmt.Errorf("Erreur newBatch: " + err.Error())
	}
	return nil
}
