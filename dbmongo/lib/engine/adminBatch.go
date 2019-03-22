package engine

import (
	"errors"
	"fmt"
	"time"

	"github.com/cnf/structhash"
	"github.com/globalsign/mgo/bson"
)

// AdminID Collection key
type AdminID struct {
	Key  string `json:"key" bson:"key"`
	Type string `json:"type" bson:"type"`
}

// AdminBatch metadata Batch
type AdminBatch struct {
	ID            AdminID    `json:"id" bson:"_id"`
	Files         BatchFiles `json:"files" bson:"files"`
	Readonly      bool       `json:"readonly" bson:"readonly"`
	CompleteTypes []string   `json:"complete_types" bson:"complete_types"`
	Params        struct {
		DateDebut       time.Time `json:"date_debut" bson:"date_debut"`
		DateFin         time.Time `json:"date_fin" bson:"date_fin"`
		DateFinEffectif time.Time `json:"date_fin_effectif" bson:"date_fin_effectif"`
	} `json:"params" bson:"param"`
}

// BatchFiles fichiers mappés par type
type BatchFiles map[string][]string

type newFile struct {
	FileName string `json:"filename"`
	Type     string `json:"type"`
	BatchKey string `json:"batch"`
}

// Load charge les données d'un batch depuis la base de données
func (batch *AdminBatch) Load(batchKey string) error {
	err := Db.DB.C("Admin").Find(bson.M{"_id.type": "batch", "_id.key": batchKey}).One(batch)
	return err
}

// Save écrit les données d'un batch vers la base de données
func (batch *AdminBatch) Save() error {
	_, err := Db.DB.C("Admin").Upsert(bson.M{"_id": batch.ID}, batch)
	return err
}

// New crée un nouveau batch
func (batch *AdminBatch) New(batchKey string) error {
	if !isBatchID(batchKey) {
		return errors.New("Valeur de batch non autorisée")
	}
	batch.ID.Key = batchKey
	batch.ID.Type = "batch"
	batch.Files = BatchFiles{}
	return nil
}

func isBatchID(batchID string) bool {
	_, err := time.Parse("0601", batchID[0:4])
	return err == nil
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

func addFileToBatch() chan newFile {
	channel := make(chan newFile)

	go func() {
		for file := range channel {
			batch, _ := GetBatch(file.BatchKey)
			batch.Files[file.Type] = append(batch.Files[file.Type], file.FileName)
			batch.Save()
			batches, _ := GetBatches()

			MainMessageChannel <- SocketMessage{
				JournalEvent: Event{
					Code:     "addFileToBatch",
					Priority: Info,
					Date:     time.Now(),
					Comment:  file.FileName + ": ajout de type" + file.Type + " au batch " + batch.ID.Key,
				},
				Batches: batches,
			}
		}
	}()

	return channel
}

// ImportBatch lance tous les parsers sur le batch fourni
func ImportBatch(batch AdminBatch, parsers []Parser) error {
	for _, parser := range parsers {
		outputChannel, eventChannel := parser(batch)
		go RelayEvents(eventChannel)
		for tuple := range outputChannel {
			hash := fmt.Sprintf("%x", structhash.Md5(tuple, 1))
			value := Value{
				Value: Data{
					Scope: tuple.Scope(),
					Key:   tuple.Key(),
					Batch: map[string]Batch{
						batch.ID.Key: Batch{
							tuple.Type(): map[string]Tuple{
								hash: tuple,
							}}}}}
			Db.ChanData <- &value
		}
	}

	Db.ChanData <- &Value{}
	return nil
}

// ProcessBatch traitement ad-hoc modifiable pour les besoins du développement
func ProcessBatch(batchList []string, parsers []Parser, types []string) error {

	for _, v := range batchList {
		batch, errBatch := GetBatch(v)
    if errBatch !=nil {
      return errors.New("Erreur de lecture du batch: " + errBatch.Error())
    }
    ImportBatch(batch, parsers)
    time.Sleep(5 * time.Second) // TODO: trouver une façon de synchroniser l'insert des paquets
    err := Compact(v, types)
    if err != nil {
      return errors.New("Erreur de compactage: " + err.Error())
    }
  }

	batch := LastBatch()
	return Reduce(batch.ID.Key, "algo2", nil, "Features")
}

// LastBatch retourne le dernier batch
func LastBatch() AdminBatch {
	batches, _ := GetBatches()
	l := len(batches)
	batch := batches[l-1]
	return batch
}

// NextBatch crée le batch suivant le dernier batch existant
func NextBatch() error {
	batch := LastBatch()
	// spew.Dump(batch)
	newBatchID, err := NextBatchID(batch.ID.Key)
	if err != nil {
		return fmt.Errorf("Mauvais numéro de batch: " + err.Error())
	}
	newBatch := AdminBatch{
		ID: AdminID{
			Key:  newBatchID,
			Type: "batch",
		},
		CompleteTypes: batch.CompleteTypes,
	}

	batch.Readonly = true

	err = batch.Save()
	if err != nil {
		return fmt.Errorf("Erreur readonly Batch: " + err.Error())
	}

	err = newBatch.Save()
	if err != nil {
		return fmt.Errorf("Erreur newBatch: " + err.Error())
	}
	return nil
}

// RevertBatch purge le batch et supprime sa référence dans la collection Admin
func RevertBatch() error {
	batch := LastBatch()
	err := PurgeBatch(batch.ID.Key)
	if err != nil {
		return fmt.Errorf("Erreur lors de la purge: " + err.Error())
	}
	err = DropBatch(batch.ID.Key)
	if err != nil {
		return fmt.Errorf("Erreur lors de la purge: " + err.Error())
	}

	return nil
}

// DropBatch supprime une référence de batch dans la collection Admin
func DropBatch(batchKey string) error {
	_, err := Db.DB.C("Admin").RemoveAll(bson.M{"_id.key": batchKey, "_id.type": "batch"})
	return err
}
