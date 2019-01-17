package main

import (
	"errors"
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
		if match, _ := regexp.MatchString("^map.*js$", f.Name()); match {
			fp, err := ioutil.ReadFile("js/" + routine + "/" + scope + "/" + f.Name())
			if err != nil {
				return errors.New("Lecture impossible: js/" + routine + "/" + scope + "/" + f.Name())
			}
			mr.Map = mr.Map + string(fp)
		}
		if match, _ := regexp.MatchString("^reduce.*js$", f.Name()); match {
			fp, err := ioutil.ReadFile("js/" + routine + "/" + scope + "/" + f.Name())
			if err != nil {
				return errors.New("Lecture impossible: js/" + routine + "/" + scope + "/" + f.Name())
			}
			mr.Reduce = mr.Reduce + string(fp)
		}
		if match, _ := regexp.MatchString("^finalize.*js$", f.Name()); match {
			fp, err := ioutil.ReadFile("js/" + routine + "/" + scope + "/" + f.Name())
			if err != nil {
				return errors.New("Lecture impossible: js/" + routine + "/" + scope + "/" + f.Name())
			}
			mr.Finalize = mr.Finalize + string(fp)
		}
	}
	return nil
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

func reduceHandler(c *gin.Context) {
	batchKey := c.Params.ByName("batchKey")
	algo := c.Params.ByName("algo")
	err := reduce(algo, batchKey)
	if err != nil {
		c.JSON(500, err.Error())
	} else {
		c.JSON(200, "Traitement effectué")
	}
}

func reduce(batchKey string, algo string) error {
	// éviter les noms d'algo essayant de pervertir l'exploration des fonctions
	isAlphaNum := regexp.MustCompile(`^[A-Za-z0-9]+$`).MatchString
	if !isAlphaNum(algo) {
		return errors.New("nom d'algorithme invalide, alphanumérique sans espace exigé")
	}

	functions, err := loadJSFunctions("js/" + algo + "/")

	naf, err = loadNAF()
	if err != nil {
		return err
	}

	batch, err := getBatch(batchKey)
	if err != nil {
		return err
	}

	scope := bson.M{
		"date_debut":             batch.Params.DateDebut,
		"date_fin":               batch.Params.DateFin,
		"date_fin_effectif":      batch.Params.DateFinEffectif,
		"serie_periode":          genereSeriePeriode(batch.Params.DateDebut, batch.Params.DateFin),
		"serie_periode_annuelle": genereSeriePeriodeAnnuelle(batch.Params.DateDebut, batch.Params.DateFin),
		"offset_effectif":        (batch.Params.DateFinEffectif.Year()-batch.Params.DateFin.Year())*12 + int(batch.Params.DateFinEffectif.Month()-batch.Params.DateFin.Month()),
		"actual_batch":           batch.ID.Key,
		"naf":                    naf,
		"f":                      functions,
		"batches":                getBatchesID(),
		"types":                  getTypes(),
	}

	job := &mgo.MapReduce{
		Map:      functions["map"].Code,
		Reduce:   functions["reduce"].Code,
		Finalize: functions["finalize"].Code,
		Out:      bson.M{"replace": "Features"},
		Scope:    scope,
	}

	_, err = db.DB.C("RawData").Find(bson.M{"value.index.algo2": true}).MapReduce(job, nil)

	return err
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

	functions, err := loadJSFunctions("js/compact/")
	if err != nil {
		return err
	}
	// Traitement MR
	job := &mgo.MapReduce{
		Map:      functions["map"].Code,
		Reduce:   functions["reduce"].Code,
		Finalize: functions["finalize"].Code,
		Out:      bson.M{"replace": "RawData"},
		Scope: bson.M{
			"functions":     functions,
			"batches":       getBatchesID(),
			"types":         getTypes(),
			"completeTypes": completeTypes,
		},
	}

	_, err = db.DB.C("RawData").Find(nil).MapReduce(job, nil)

	return err
}

func getTypes() []string {
	return []string{
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
	}
}

func getNAF(c *gin.Context) {
	c.JSON(200, naf)
}
