package engine

import (
	"dbmongo/lib/exportdatapi"
	"dbmongo/lib/naf"

	"github.com/globalsign/mgo/bson"
	daclient "github.com/signaux-faibles/datapi/client"
)

// ExportDetectionToDatapi sends detections with some informations to a datapi server
func ExportDetectionToDatapi(url, user, password, batch string) error {
	client := daclient.DatapiServer{
		URL: url,
	}
	err := client.Connect(user, password)
	if err != nil {
		return err
	}

	var pipeline = exportdatapi.GetPipeline(batch)

	iter := Db.DB.C("Prediction").Pipe(pipeline).Iter()

	var data exportdatapi.Detection

	var datas []daclient.Object

	var i int

	for iter.Next(&data) {
		i++
		d, err := exportdatapi.ComputeDetection(data)
		if err != nil {
			continue
		}

		object := daclient.Object{
			Key:   d.Key,
			Scope: d.Scope,
			Value: d.Value,
		}

		datas = append(datas, object)
	}

	if datas != nil {
		err = client.Put("public", datas)
	}
	return err
}

// ExportPublicToDatapi sends public data to a datapi server
func ExportPublicToDatapi(url string, user string, password string, batch string) error {
	client := daclient.DatapiServer{
		URL: url,
	}
	err := client.Connect(user, password)
	if err != nil {
		return err
	}

	cursor := Db.DB.C("Public").Find(bson.M{"_id.batch": batch})

	iter := cursor.Iter()

	var data struct {
		ID    map[string]string      `bson:"_id"`
		Value map[string]interface{} `bson:"value"`
	}

	var datas []daclient.Object

	var i int

	for iter.Next(&data) {
		i++

		if data.Value != nil {
			o := daclient.Object{
				Key:   data.ID,
				Value: data.Value,
			}

			datas = append(datas, o)
		}
	}

	if datas != nil {
		err = client.Put("public", datas)
	}

	return err
}

// ExportReferencesToDatapi pushes references (batches, types, etc.) to a datapi server
func ExportReferencesToDatapi(url string, user string, password string, batch string) error {
	client := daclient.DatapiServer{
		URL: url,
	}

	err := client.Connect(user, password)

	nafCodes := daclient.Object{
		Key: map[string]string{
			"key":   "naf",
			"batch": batch,
		},
		Scope: []string{},
		Value: naf.Naf.ToData(),
	}

	types := daclient.Object{
		Key: map[string]string{
			"key":   "types",
			"batch": batch,
		},
		Scope: []string{},
		Value: GetTypes().ToData(),
	}

	var batchesData []daclient.Object
	batches, err := GetBatches()
	if err != nil {
		return err
	}
	for _, b := range batches {
		o := daclient.Object{
			Key: map[string]string{
				"key":   "batches",
				"batch": b.ID.Key,
			},
			Scope: []string{},
			Value: b.ToData(),
		}
		batchesData = append(batchesData, o)
	}

	var data []daclient.Object
	data = append(data, nafCodes)
	data = append(data, types)
	data = append(data, batchesData...)
	err = client.Put("reference", data)

	return err
}

// // ExportToDatapi exports Public to Datapi
// func ExportToDatapi(url string, user string, password string) error {
// 	client := daclient.DatapiServer{
// 		URL: url,
// 	}

// 	err := client.Connect(user, password)

// 	data := daclient.Object{
// 		Key: map[string]string{
// 			"type": "reference",
// 			"key":  "naf",
// 		},
// 		Scope: []string{},
// 		Value: naf.Naf.ToData(),
// 	}

// 	err = client.Put("reference", []daclient.Object{data})

// 	return err
// }
