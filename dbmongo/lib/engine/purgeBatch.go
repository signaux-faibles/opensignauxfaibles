package engine

import (
	"errors"
	"strconv"
	"time"

	"github.com/globalsign/mgo"
	"github.com/globalsign/mgo/bson"
	"github.com/spf13/viper"
)

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

	chunks, err := ChunkCollection(viper.GetString("DB"), "RawData", viper.GetInt64("chunkByteSize"))
	if err != nil {
		return err
	}

	w := MRWait{}
	w.init()

	var tempDBs []string
	var tempDBChannel = make(chan string)

	i := 0
	for _, query := range chunks.ToQueries(nil, "_id") {
		w.waitGroup.Add(1)

		dbTemp := "purgeBatch" + strconv.Itoa(i)
		job := &mgo.MapReduce{
			Map:      functions["map"].Code,
			Reduce:   functions["reduce"].Code,
			Finalize: functions["finalize"].Code,
			Out:      bson.M{"replace": "TemporaryCollection", "db": dbTemp},
			Scope:    scope,
		}
		i++
		go MRroutine(job, query, dbTemp, "RawData", &w, tempDBChannel)

	}

	go func() {
		for tempDB := range tempDBChannel {
			tempDBs = append(tempDBs, tempDB)
		}
	}()

	w.waitGroup.Wait()
	close(tempDBChannel)

	db, _ := mgo.Dial(viper.GetString("DB_DIAL"))
	db.SetSocketTimeout(720000 * time.Second)

	for _, dbTemp := range tempDBs {
		pipeline := []bson.M{
			bson.M{
				"$merge": bson.M{
					"into": bson.M{
						"coll": "Public",
						"db":   viper.GetString("DB")},
				},
			},
		}

		pipe := db.DB(dbTemp).C("TemporaryCollection").Pipe(pipeline)
		var result []interface{}
		err = pipe.AllowDiskUse().All(&result)
		if err != nil {
			w.add("errors", 1, -1)
		} else {
			err = db.DB(dbTemp).DropDatabase()
			if err != nil {
				w.add("errors", 1, -1)
			}
		}
	}

	db.Close()

	errorcount, _ := w.running.Load("errors")
	if errorcount.(int) != 0 {
		return errors.New("erreurs constatées, consultez les journaux")
	}

	return nil
}
