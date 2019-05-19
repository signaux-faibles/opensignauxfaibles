package engine

import (
	"dbmongo/lib/exportdatapi"
	"dbmongo/lib/naf"

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
		detection, public, entreprise, err := exportdatapi.ComputeDetection(data)

		if err != nil {
			continue
		}

		datas = append(datas, detection, public, entreprise)

	}

	if datas != nil {
		err = client.Put("public", datas)
	}
	return err
}

// func getDepartement(b map[string]interface{}) (string, error) {
// 	sirene, ok := b["sirene"].(map[string]interface{})
// 	if !ok {
// 		return "", errors.New("no sirene")
// 	}

// 	dept, ok := sirene["departement"].(string)
// 	if !ok {
// 		return "", errors.New("no departement")
// 	}

// 	return dept, nil
// }

// // ExportPublicToDatapi sends public data to a datapi server
// func ExportPublicToDatapi(url string, user string, password string, batch string) error {
// 	client := daclient.DatapiServer{
// 		URL: url,
// 	}

// 	err := client.Connect(user, password)
// 	if err != nil {
// 		return err
// 	}

// 	cursor := Db.DB.C("Public").Find(bson.M{"_id.batch": batch})

// 	iter := cursor.Iter()

// 	var data struct {
// 		ID    map[string]string      `bson:"_id"`
// 		Value map[string]interface{} `bson:"value"`
// 	}

// 	var datas []daclient.Object

// 	var i int

// 	for iter.Next(&data) {
// 		i++

// 		departement, error := getDepartement
// 		key := map[string]string{
// 			"key":   data.ID["key"],
// 			"batch": data.ID["batch"],
// 			"type":  "detail",
// 			"scope": data.ID["scope"],
// 		}

// 		if data.Value != nil {
// 			o := daclient.Object{
// 				Key:   key,
// 				Value: data.Value,
// 			}

// 			datas = append(datas, o)
// 		}
// 	}

// 	if datas != nil {
// 		err = client.Put("public", datas)
// 	}

// 	return err
// }

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
				"key":   "batch",
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
