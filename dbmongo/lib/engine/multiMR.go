package engine

import (
	"fmt"

	"github.com/globalsign/mgo/bson"
)

func splitCollection(collection string, chunckCount int, query bson.M, project bson.M) ([]string, error) {

	count, err := Db.DB.C(collection).Find(query).Count()
	if err != nil {
		return nil, err
	}
	pos := 0
	var chunks []string
	for {
		pos = pos + count/chunckCount
		if pos >= count {
			break
		}

		var res struct {
			ID string `bson:"_id"`
		}

		Db.DB.C(collection).Find(query).Sort("_id").Skip(pos).Select(project).One(&res)
		chunks = append(chunks, res.ID)
	}

	return chunks, nil
}

// Public alimente la collection Public avec les objets destinés à la diffusion
func Public(batch AdminBatch, siret string) error {
	chunks, err := splitCollection("RawData", 10, bson.M{"value.index.algo2": true}, bson.M{"_id": 1})
	if err != nil {
		return err
	}
	
	var queries []interface{}
	queries = append(queries, bson.M{
		"_id": bson.M{"$lt": chunks[0][0:9]},
	})
  
	for n := range chunks {
		if n == len(chunks)-2 {
			queries = append(queries, bson.M{
				"_id": bson.M{"$gt": chunks[n+1][0:9]},
			})
			break
		}
		queries = append(queries,
			[]interface{}{
				bson.M{"_id": bson.M{"$gt": chunks[n][0:9]}},
				bson.M{"_id": bson.M{"$lt": chunks[n+1][0:9]}},
			},
		)
	}
	fmt.Println(chunks)
	fmt.Println(queries)
	// functions, err := loadJSFunctions("public")

	// scope := bson.M{
	// 	"date_debut":             batch.Params.DateDebut,
	// 	"date_fin":               batch.Params.DateFin,
	// 	"date_fin_effectif":      batch.Params.DateFinEffectif,
	// 	"serie_periode":          misc.GenereSeriePeriode(batch.Params.DateFin.AddDate(0, -24, 0), batch.Params.DateFin),
	// 	"serie_periode_annuelle": misc.GenereSeriePeriodeAnnuelle(batch.Params.DateFin.AddDate(0, -24, 0), batch.Params.DateFin),
	// 	"offset_effectif":        (batch.Params.DateFinEffectif.Year()-batch.Params.DateFin.Year())*12 + int(batch.Params.DateFinEffectif.Month()-batch.Params.DateFin.Month()),
	// 	"actual_batch":           batch.ID.Key,
	// 	"naf":                    naf.Naf,
	// 	"f":                      functions,
	// 	"batches":                GetBatchesID(),
	// 	"types":                  GetTypes(),
	// }

	// job := &mgo.MapReduce{
	// 	Map:      functions["map"].Code,
	// 	Reduce:   functions["reduce"].Code,
	// 	Finalize: functions["finalize"].Code,
	// 	Out:      bson.M{"replace": "Public"},
	// 	Scope:    scope,
	// }
	// // exécution

	// if siret != "" {
	// 	_, err = Db.DB.C("RawData").Find(bson.M{
	// 		"$or": []interface{}{
	// 			bson.M{"_id": siret},
	// 			bson.M{"_id": siret[0:9]},
	// 		},
	// 	}).MapReduce(job, nil)
	// } else {
	// 	_, err = Db.DB.C("RawData").Find(bson.M{"value.index.algo2": true}).MapReduce(job, nil)
	// }
	// if err != nil {
	// 	return errors.New("Erreur dans l'exécution des jobs MapReduce" + err.Error())
	// }
	return nil
}
