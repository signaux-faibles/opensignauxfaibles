package main

import (
	"errors"
	"fmt"
	"io/ioutil"
	"regexp"
	"sort"

	"github.com/gin-gonic/gin"
	"github.com/globalsign/mgo"
	"github.com/globalsign/mgo/bson"
)

// MapReduceJS Ensemble de fonctions JS pour mongodb
type MapReduceJS struct {
	Routine  string
	Scope    string
	Map      string
	Reduce   string
	Finalize string
}

func loadMR(typeJob string, target string) (*mgo.MapReduce, error) {
	mr := &mgo.MapReduce{}

	file, err := ioutil.ReadDir("js/" + typeJob + "/" + target)
	sort.Slice(file, func(i, j int) bool {
		return file[i].Name() < file[j].Name()
	})

	if err != nil {
		return nil, errors.New("Chemin introuvable")
	}

	mr.Map = ""
	mr.Reduce = ""
	mr.Finalize = ""

	for _, f := range file {
		if match, _ := regexp.MatchString("^map.*js", f.Name()); match {
			fp, err := ioutil.ReadFile("js/" + typeJob + "/" + target + "/" + f.Name())
			if err != nil {
				return nil, errors.New("Lecture impossible: js/" + typeJob + "/" + target + "/" + f.Name())
			}
			mr.Map = mr.Map + string(fp)
		}
		if match, _ := regexp.MatchString("^reduce.*js", f.Name()); match {
			fp, err := ioutil.ReadFile("js/" + typeJob + "/" + target + "/" + f.Name())
			if err != nil {
				return nil, errors.New("Lecture impossible: js/" + typeJob + "/" + target + "/" + f.Name())
			}
			mr.Reduce = mr.Reduce + string(fp)
		}
		if match, _ := regexp.MatchString("^finalize.*js", f.Name()); match {
			fp, err := ioutil.ReadFile("js/" + typeJob + "/" + target + "/" + f.Name())
			if err != nil {
				return nil, errors.New("Lecture impossible: js/" + typeJob + "/" + target + "/" + f.Name())
			}
			mr.Finalize = mr.Finalize + string(fp)
		}
	}
	return mr, nil

}

func (mr *MapReduceJS) load(routine string, scope string) error {
	file, err := ioutil.ReadDir("js/" + routine + "/" + scope)
	sort.Slice(file, func(i, j int) bool {
		return file[i].Name() < file[j].Name()
	})

	if err != nil {
		return errors.New("Chemin introuvable")
	}

	mr.Routine = routine
	mr.Scope = scope
	mr.Map = ""
	mr.Reduce = ""
	mr.Finalize = ""

	for _, f := range file {
		if match, _ := regexp.MatchString("^map.*js", f.Name()); match {
			fp, err := ioutil.ReadFile("js/" + routine + "/" + scope + "/" + f.Name())
			if err != nil {
				return errors.New("Lecture impossible: js/" + routine + "/" + scope + "/" + f.Name())
			}
			mr.Map = mr.Map + string(fp)
		}
		if match, _ := regexp.MatchString("^reduce.*js", f.Name()); match {
			fp, err := ioutil.ReadFile("js/" + routine + "/" + scope + "/" + f.Name())
			if err != nil {
				return errors.New("Lecture impossible: js/" + routine + "/" + scope + "/" + f.Name())
			}
			mr.Reduce = mr.Reduce + string(fp)
		}
		if match, _ := regexp.MatchString("^finalize.*js", f.Name()); match {
			fp, err := ioutil.ReadFile("js/" + routine + "/" + scope + "/" + f.Name())
			if err != nil {
				return errors.New("Lecture impossible: js/" + routine + "/" + scope + "/" + f.Name())
			}
			mr.Finalize = mr.Finalize + string(fp)
		}
	}
	return nil
}

func reduceHandler(c *gin.Context) {
	algo := c.Params.ByName("algo")
	batchKey := c.Params.ByName("batch")
	siret := c.Params.ByName("siret")

	batch, _ := getBatch(batchKey)
	result, err := reduce(batch, algo, siret)

	if err != nil {
		c.JSON(500, err.Error())
	} else {
		c.JSON(200, result)
	}
}

func reduce(batch AdminBatch, algo string, siret string) (interface{}, error) {
	var queryEtablissement interface{}
	var queryEntreprise interface{}
	var output interface{}
	var result interface{}

	dateDebut := batch.Params.DateDebut
	dateFin := batch.Params.DateFin
	dateFinEffectif := batch.Params.DateFinEffectif

	if siret == "" {
		db.DB.C("Features").RemoveAll(bson.M{"_id.batch": batch.ID.Key, "_id.algo": algo})
		queryEtablissement = bson.M{"value.index." + algo: true}
		queryEntreprise = nil
		output = bson.M{"merge": "Features"}
	} else {
		queryEtablissement = bson.M{"value.siret": siret}
		queryEntreprise = bson.M{"value.siren": siret[0:9]}
		output = nil
	}

	MREtablissement := MapReduceJS{}
	MREntreprise := MapReduceJS{}
	MRUnion := MapReduceJS{}
	errEt := MREtablissement.load(algo, "etablissement")
	errEn := MREntreprise.load(algo, "entreprise")
	errUn := MRUnion.load(algo, "union")

	if errEt != nil || errEn != nil || errUn != nil {
		return nil, fmt.Errorf("Problème d'accès aux fichiers MapReduce")
	}

	naf, errNAF := loadNAF()
	if errNAF != nil {
		return nil, fmt.Errorf("Problème d'accès aux fichiers naf")
	}

	scope := bson.M{
		"date_debut":             dateDebut,
		"date_fin":               dateFin,
		"date_fin_effectif":      dateFinEffectif,
		"serie_periode":          genereSeriePeriode(dateDebut, dateFin),
		"serie_periode_annuelle": genereSeriePeriodeAnnuelle(dateDebut, dateFin),
		"actual_batch":           batch.ID.Key,
		"naf":                    naf,
	}

	jobEtablissement := &mgo.MapReduce{
		Map:      string(MREtablissement.Map),
		Reduce:   string(MREtablissement.Reduce),
		Finalize: string(MREtablissement.Finalize),
		Out:      bson.M{"replace": "MRWorkspace"},
		Scope:    scope,
	}

	_, err := db.DB.C("Etablissement").Find(queryEtablissement).MapReduce(jobEtablissement, nil)
	if err != nil {
		return nil, fmt.Errorf("Erreur du traitement MapReduce Etablissement: " + err.Error())
	}

	jobEntreprise := &mgo.MapReduce{
		Map:      string(MREntreprise.Map),
		Reduce:   string(MREntreprise.Reduce),
		Finalize: string(MREntreprise.Finalize),
		Out:      bson.M{"merge": "MRWorkspace"},
		Scope:    scope,
	}

	_, err = db.DB.C("Entreprise").Find(queryEntreprise).MapReduce(jobEntreprise, nil)
	if err != nil {
		return nil, fmt.Errorf("Erreur du traitement Entreprise: " + err.Error())
	}

	jobUnion := &mgo.MapReduce{
		Map:      string(MRUnion.Map),
		Reduce:   string(MRUnion.Reduce),
		Finalize: string(MRUnion.Finalize),
		Out:      output,
		Scope:    scope,
	}

	if output == nil {
		_, err = db.DB.C("MRWorkspace").Find(queryEntreprise).MapReduce(jobUnion, &result)
	} else {
		_, err = db.DB.C("MRWorkspace").Find(queryEntreprise).MapReduce(jobUnion, nil)
	}

	if err != nil {
		return result, fmt.Errorf("Erreur du traitement MapReduce Union")
	}

	return result, nil

}

func unionReduceHandler(c *gin.Context) {
	algo := c.Params.ByName("algo")
	batchKey := c.Params.ByName("batch")
	siret := c.Params.ByName("siret")

	batch, _ := getBatch(batchKey)
	result, err := unionReduce(batch, algo, siret)

	if err != nil {
		c.JSON(500, err.Error())
	} else {
		c.JSON(200, result)
	}
}

func unionReduce(batch AdminBatch, algo string, siret string) (interface{}, error) {
	var queryEntreprise interface{}
	var output interface{}
	var result interface{}

	dateDebut := batch.Params.DateDebut
	dateFin := batch.Params.DateFin
	dateFinEffectif := batch.Params.DateFinEffectif

	if siret == "" {
		db.DB.C("Features").RemoveAll(bson.M{"_id.batch": batch.ID.Key, "_id.algo": algo})
		queryEntreprise = nil
		output = bson.M{"merge": "Features"}
	} else {
		queryEntreprise = bson.M{"value.siren": siret[0:9]}
		output = nil
	}

	MRUnion := MapReduceJS{}
	errUn := MRUnion.load(algo, "union")

	if errUn != nil {
		return nil, fmt.Errorf("Problème d'accès aux fichiers MapReduce")
	}

	naf, errNAF := loadNAF()
	if errNAF != nil {
		return nil, fmt.Errorf("Problème d'accès aux fichiers naf")
	}

	scope := bson.M{
		"date_debut":             dateDebut,
		"date_fin":               dateFin,
		"date_fin_effectif":      dateFinEffectif,
		"serie_periode":          genereSeriePeriode(dateDebut, dateFin),
		"serie_periode_annuelle": genereSeriePeriodeAnnuelle(dateDebut, dateFin),
		"actual_batch":           batch.ID.Key,
		"naf":                    naf,
	}

	jobUnion := &mgo.MapReduce{
		Map:      string(MRUnion.Map),
		Reduce:   string(MRUnion.Reduce),
		Finalize: string(MRUnion.Finalize),
		Out:      output,
		Scope:    scope,
	}

	var err error
	if output == nil {
		_, err = db.DB.C("MRWorkspace").Find(queryEntreprise).MapReduce(jobUnion, &result)
	} else {
		_, err = db.DB.C("MRWorkspace").Find(queryEntreprise).MapReduce(jobUnion, nil)
	}

	if err != nil {
		return result, fmt.Errorf("Erreur du traitement MapReduce Union")
	}

	return result, nil
}

func loadJSFunctions(path string) (map[string]bson.JavaScript, error) {
	files, err := ioutil.ReadDir(path)
	r := regexp.MustCompile(`(.*)\.js$`)

	functions := make(map[string]bson.JavaScript)

	for _, f := range files {
		if name := r.FindStringSubmatch(f.Name()); len(name) > 0 {
			b, err := ioutil.ReadFile(path + f.Name())
			if err == nil {
				functions[name[1]] = bson.JavaScript{
					Code: string(b),
				}
			}
		}
	}
	return functions, err
}

func compactHandler(c *gin.Context) {
	err := compact()
	if err != nil {
		c.JSON(500, err.Error())
		return
	}

	c.JSON(200, "ok")
}

func compact() error {
	batches, _ := getBatches()

	// Détermination scope traitement
	var completeTypes = make(map[string][]string)
	var batchesID []string

	for _, b := range batches {
		completeTypes[b.ID.Key] = b.CompleteTypes
		batchesID = append(batchesID, b.ID.Key)
	}

	query := bson.M{"value.key": "74655966775348"}
	output := bson.M{"replace": "testCollection"}
	functions, err := loadJSFunctions("js/compact/")
	if err != nil {
		return err
	}
	// Traitement MR
	job := &mgo.MapReduce{
		Map:      functions["map"].Code,
		Reduce:   functions["reduce"].Code,
		Finalize: functions["finalize"].Code,
		Out:      output,
		Scope: bson.M{
			"functions": functions,
			"batches":   batchesID,
			"types": []string{
				"altares",
				"apconso",
				"apdemande",
				"ccsf",
				"cotisation",
				"debit",
				"delai",
				"effectif",
				"sirene",
				"dpae",
				"bdf",
				"diane",
			},
			"completeTypes": completeTypes,
		},
	}

	_, err = db.DB.C("RawData").Find(query).MapReduce(job, nil)

	return err
}

func getFeatures(c *gin.Context) {
	var data []interface{}
	db.DB.C("Features").Find(nil).All(&data)
	c.JSON(200, data)
}

func getNAF(c *gin.Context) {
	c.JSON(200, naf)
}
