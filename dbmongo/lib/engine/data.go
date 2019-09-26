package engine

import (
	"dbmongo/lib/misc"
	"dbmongo/lib/naf"
	"errors"
	"fmt"
	"regexp"

	"github.com/globalsign/mgo"
	"github.com/globalsign/mgo/bson"
)

//go:generate go run js/loadJS.go

func loadJSFunctions(directoryName string) (map[string]bson.JavaScript, error) {
	functions := make(map[string]bson.JavaScript)
	var err error
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
	return functions, err
}

// PurgeNotCompacted permet de supprimer les objets non encore compactés
// c'est à dire, vider la collection ImportedData
func PurgeNotCompacted() error {
	_, err := Db.DB.C("ImportedData").RemoveAll(nil)
	return err
}

// PurgeBatch permet de supprimer un batch dans les objets de RawData
func PurgeBatch(batchKey string) error {

	functions, err := loadJSFunctions("purgeBatch")
	if err != nil {
		return err
	}
	scope := bson.M{
		"currentBatch": batchKey,
		"f":            functions,
	}

	job := &mgo.MapReduce{
		Map:      functions["map"].Code,
		Reduce:   functions["reduce"].Code,
		Finalize: functions["finalize"].Code,
		Out:      bson.M{"merge": "RawData"},
		Scope:    scope,
	}

	var slices []string
	for i := 57; i <= 57; i++ {
		slices = append(slices, fmt.Sprintf("^%02d.*", i))
	}

	for _, s := range slices {
		fmt.Println("Purge des objets " + s)
		_, err = Db.DB.C("RawData").Find(bson.M{"_id": bson.RegEx{
			Pattern: s,
			Options: "",
		}}).MapReduce(job, nil)

		if err != nil {
			return err
		}
	}
	return nil
}

// Compact traite le compactage de la base RawData
//func Compact() error {
//	batches, _ := GetBatches()
//
//	// Détermination scope traitement
func Compact(batchKey string, types []string) error {
	// Détermination scope traitement
	batches, _ := GetBatches()

	var batchesID []string
	var completeTypes = make(map[string][]string)
	for _, b := range batches {
		completeTypes[b.ID.Key] = b.CompleteTypes
		batchesID = append(batchesID, b.ID.Key)
	}
	// Si le numéro de batch n'est pas valide, on prend le premier
	found := false
	for _, batchID := range batchesID {
		if batchID == batchKey {
			found = true
			break
		}
	}
	if !found {
		batchKey = batchesID[0]
	}

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
			"types":         types,
			"completeTypes": completeTypes,
			"batchKey":      batchKey,
		},
	}

	_, err = Db.DB.C("ImportedData").Find(nil).MapReduce(job, nil)

	if err != nil {
		return err
	}
	err = PurgeNotCompacted()
	return err
}

// Reduce alimente la base Features
func Reduce(batchKey string, algo string, query interface{}, collection string) error {

	// éviter les noms d'algo essayant de hacker l'exploration des fonctions ci-dessous
	isAlphaNum := regexp.MustCompile(`^[A-Za-z0-9]+$`).MatchString
	if !isAlphaNum(algo) {
		return errors.New("nom d'algorithme invalide, alphanumérique sans espace exigé")
	}

	functions, err := loadJSFunctions("reduce." + algo)

	naf, err := naf.LoadNAF()
	if err != nil {
		return err
	}

	batch, err := GetBatch(batchKey)
	if err != nil {
		return err
	}

	scope := bson.M{
		"date_debut":             batch.Params.DateDebut,
		"date_fin":               batch.Params.DateFin,
		"date_fin_effectif":      batch.Params.DateFinEffectif,
		"serie_periode":          misc.GenereSeriePeriode(batch.Params.DateDebut, batch.Params.DateFin),
		"serie_periode_annuelle": misc.GenereSeriePeriodeAnnuelle(batch.Params.DateDebut, batch.Params.DateFin),
		"offset_effectif":        (batch.Params.DateFinEffectif.Year()-batch.Params.DateFin.Year())*12 + int(batch.Params.DateFinEffectif.Month()-batch.Params.DateFin.Month()),
		"actual_batch":           batch.ID.Key,
		"naf":                    naf,
		"f":                      functions,
		"batches":                GetBatchesID(),
		"types":                  GetTypes(),
	}

	job := &mgo.MapReduce{
		Map:      functions["map"].Code,
		Reduce:   functions["reduce"].Code,
		Finalize: functions["finalize"].Code,
		//TODO merge into collection instead of replacing. Must be idempotent
		//transformation. Not the case now with agregation
		Out:   bson.M{"replace": collection}, // bson.M{"merge": collection},
		Scope: scope,
	}

	_, err = Db.DB.C("RawData").Find(query).MapReduce(job, nil)

	if err != nil {
		return err
	}

	// Separating different sirets in different objects
	query2 := []bson.M{{
		"$unwind": bson.M{"path": "$value", "preserveNullAndEmptyArrays": false},
	},
		{
			"$project": bson.M{"_id": 0.0, "info": "$_id", "value": 1.0},
		},
		{"$out": collection}}
	pipe := Db.DB.C(collection).Pipe(query2)
	resp := []bson.M{}
	err = pipe.All(&resp)

	return err
}

// ReduceMergeAux merges collection reduced into its destination
func ReduceMergeAux() error {
	// job := &mgo.MapReduce{
	// 	Map:      "function map() { emit(this._id, {info: this.info, value: this.value}) }",
	// 	Reduce:   "function reduce(_, v) {return v}",
	// 	Finalize: "function finalize(_, v) { return v }",
	// 	Out:      bson.M{"merge": "Features"},
	// }
	// _, err := Db.DB.C("Features_aux").Find(bson.M{}).MapReduce(job, nil)

	query := []bson.M{{
		"$merge": bson.M{"into": "Features"},
	}}
	pipe := Db.DB.C("Features_aux").Pipe(query)
	resp := []bson.M{}
	err := pipe.All(&resp)

	if err != nil {
		return err
	}
	_, err = Db.DB.C("Features_aux").RemoveAll(nil)
	if err != nil {
		return err
	}
	return nil
}

// PublicMergeAux merges collection reduced into its destination
func PublicMergeAux() error {
	// job := &mgo.MapReduce{
	// 	Map:      "function map() { emit(this._id, {info: this.info, value: this.value}) }",
	// 	Reduce:   "function reduce(_, v) {return v}",
	// 	Finalize: "function finalize(_, v) { return v }",
	// 	Out:      bson.M{"merge": "Features"},
	// }
	// _, err := Db.DB.C("Features_aux").Find(bson.M{}).MapReduce(job, nil)

	query := []bson.M{{
		"$merge": bson.M{"into": "Public"},
	}}
	pipe := Db.DB.C("Public_aux").Pipe(query)
	resp := []bson.M{}
	err := pipe.All(&resp)

	if err != nil {
		return err
	}
	_, err = Db.DB.C("Public_aux").RemoveAll(nil)
	if err != nil {
		return err
	}
	return nil
}

// Public alimente la collection Public avec les objets destinés à la diffusion
func Public(batch AdminBatch, algo string, query bson.M, collection string) error {

	functions, err := loadJSFunctions("public")

	scope := bson.M{
		"date_debut":             batch.Params.DateDebut,
		"date_fin":               batch.Params.DateFin,
		"date_fin_effectif":      batch.Params.DateFinEffectif,
		"serie_periode":          misc.GenereSeriePeriode(batch.Params.DateFin.AddDate(0, -24, 0), batch.Params.DateFin),
		"serie_periode_annuelle": misc.GenereSeriePeriodeAnnuelle(batch.Params.DateFin.AddDate(0, -24, 0), batch.Params.DateFin),
		"offset_effectif":        (batch.Params.DateFinEffectif.Year()-batch.Params.DateFin.Year())*12 + int(batch.Params.DateFinEffectif.Month()-batch.Params.DateFin.Month()),
		"actual_batch":           batch.ID.Key,
		"naf":                    naf.Naf,
		"f":                      functions,
		"batches":                GetBatchesID(),
		"types":                  GetTypes(),
	}

	job := &mgo.MapReduce{
		Map:      functions["map"].Code,
		Reduce:   functions["reduce"].Code,
		Finalize: functions["finalize"].Code,
		Out:      bson.M{"replace": collection},
		Scope:    scope,
	}
	// exécution

	_, err = Db.DB.C("RawData").Find(query).MapReduce(job, nil)

	if err != nil {
		return errors.New("Erreur dans l'exécution des jobs MapReduce" + err.Error())
	}
	return nil
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
