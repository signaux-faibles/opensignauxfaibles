package engine

import (
	"dbmongo/lib/exportdatapi"
	"dbmongo/lib/naf"
	"fmt"

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
		detection, err := exportdatapi.Compute(data)

		if err != nil {
			continue
		}

		datas = append(datas, detection...)

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

// ExportPoliciesToDatapi exports standard policies to datapi
func ExportPoliciesToDatapi(url, user, password, batch string) error {
	var policies []daclient.Object
	policies = append(policies, daclient.Object{
		Key: map[string]string{
			"type": "policy",
			"name": "Accès Bourgogne Franche-Comté",
		},
		Value: map[string]interface{}{
			"match": "(public|reference)",
			"scope": []string{"Bourgogne Franche-Comté"},
			"promote": []string{
				"Côte d'or",
				"Doubs",
				"Haute-Saône",
				"Jura",
				"Nièvre",
				"Saône-et-Loire",
				"Territoire de Belfort",
				"Yonne",
			},
		},
	})

	policies = append(policies, daclient.Object{
		Key: map[string]string{
			"type": "policy",
			"name": "Accès Bourgogne",
		},
		Value: map[string]interface{}{
			"match": "(public|reference)",
			"scope": []string{"Bourgogne"},
			"promote": []string{
				"Côte d'or",
				"Nièvre",
				"Saône-et-Loire",
				"Yonne",
			},
		},
	})

	policies = append(policies, daclient.Object{
		Key: map[string]string{
			"type": "policy",
			"name": "Accès Franche-Comté",
		},
		Value: map[string]interface{}{
			"match": "(public|reference)",
			"scope": []string{"Franche-Comté"},
			"promote": []string{
				"Doubs",
				"Haute-Saône",
				"Jura",
				"Territoire de Belfort",
			},
		},
	})

	policies = append(policies, daclient.Object{
		Key: map[string]string{
			"type": "policy",
			"name": "Accès Pays de la Loire",
		},
		Value: map[string]interface{}{
			"match": "(public|reference)",
			"scope": []string{"Pays de la Loire"},
			"promote": []string{
				"Loire-Atlantique",
				"Maine-et-Loire",
				"Mayenne",
				"Sarthe",
				"Vendée",
			},
		},
	})

	policies = append(policies, daclient.Object{
		Key: map[string]string{
			"type": "policy",
			"name": "Limitation en écriture",
		},
		Value: map[string]interface{}{
			"match":  ".*",
			"key":    map[string]string{},
			"writer": []string{"datawriter"},
		},
	})

	policies = append(policies, daclient.Object{
		Key: map[string]string{
			"type": "policy",
			"name": "Limitation Accès",
		},
		Value: map[string]interface{}{
			"match":  "policy",
			"key":    map[string]string{},
			"writer": []string{"manager"},
			"reader": []string{"manager"},
		},
	})

	client := daclient.DatapiServer{
		URL: url,
	}
	err := client.Connect(user, password)

	if err != nil {
		return err
	}

	err = client.Put("system", policies)
	fmt.Println(err)
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

	batchData, err := GetBatch(batch)
	if err != nil {
		return err
	}

	batchObject := daclient.Object{
		Key: map[string]string{
			"key":   "batch",
			"batch": batchData.ID.Key,
		},
		Scope: []string{},
		Value: batchData.ToData(),
	}

	var data []daclient.Object
	data = append(data, nafCodes)
	data = append(data, types)
	data = append(data, batchObject)
	data = append(data, getRegions()...)
	err = client.Put("reference", data)

	return err
}

func getRegions() (regions []daclient.Object) {
	regions = append(regions, daclient.Object{
		Key: map[string]string{
			"type":   "region",
			"region": "Bourgogne",
		},
		Scope: []string{"Bourgogne"},
		Value: map[string]interface{}{
			"departements": []string{
				"Côte d'or",
				"Nièvre",
				"Saône-et-Loire",
				"Yonne",
			},
		},
	})

	regions = append(regions, daclient.Object{
		Key: map[string]string{
			"type":   "region",
			"region": "Franche-Comté",
		},
		Scope: []string{"Franche-Comté"},
		Value: map[string]interface{}{
			"departements": []string{
				"Doubs",
				"Haute-Saône",
				"Jura",
				"Territoire de Belfort",
			},
		},
	})

	regions = append(regions, daclient.Object{
		Key: map[string]string{
			"type":   "region",
			"region": "Pays de la Loire",
		},
		Scope: []string{"Pays de la Loire"},
		Value: map[string]interface{}{
			"departements": []string{
				"Loire-Atlantique",
				"Maine-et-Loire",
				"Mayenne",
				"Sarthe",
				"Vendée",
			},
		},
	})

	regions = append(regions, daclient.Object{
		Key: map[string]string{
			"type": "departements",
		},
		Scope: []string{"Doubs"},
		Value: map[string]interface{}{
			"Doubs": 25,
		},
	})

	regions = append(regions, daclient.Object{
		Key: map[string]string{
			"type": "departements",
		},
		Scope: []string{"Haute-Saône"},
		Value: map[string]interface{}{
			"Haute-Saône": 70,
		},
	})

	regions = append(regions, daclient.Object{
		Key: map[string]string{
			"type": "departements",
		},
		Scope: []string{"Jura"},
		Value: map[string]interface{}{
			"Jura": 39,
		},
	})

	regions = append(regions, daclient.Object{
		Key: map[string]string{
			"type": "departements",
		},
		Scope: []string{"Territoire de Belfort"},
		Value: map[string]interface{}{
			"Territoire de Belfort": 90,
		},
	})

	regions = append(regions, daclient.Object{
		Key: map[string]string{
			"type": "departements",
		},
		Scope: []string{"Côte d'or"},
		Value: map[string]interface{}{
			"Côte d'or": 21,
		},
	})

	regions = append(regions, daclient.Object{
		Key: map[string]string{
			"type": "departements",
		},
		Scope: []string{"Nièvre"},
		Value: map[string]interface{}{
			"Nièvre": 58,
		},
	})

	regions = append(regions, daclient.Object{
		Key: map[string]string{
			"type": "departements",
		},
		Scope: []string{"Saône-et-Loire"},
		Value: map[string]interface{}{
			"Saône-et-Loire": 71,
		},
	})

	regions = append(regions, daclient.Object{
		Key: map[string]string{
			"type": "departements",
		},
		Scope: []string{"Yonne"},
		Value: map[string]interface{}{
			"Yonne": 89,
		},
	})

	regions = append(regions, daclient.Object{
		Key: map[string]string{
			"type": "departements",
		},
		Scope: []string{"Loire-Atlantique"},
		Value: map[string]interface{}{
			"Loire-Atlantique": 44,
		},
	})

	regions = append(regions, daclient.Object{
		Key: map[string]string{
			"type": "departements",
		},
		Scope: []string{"Maine-et-Loire"},
		Value: map[string]interface{}{
			"Maine-et-Loire": 49,
		},
	})

	regions = append(regions, daclient.Object{
		Key: map[string]string{
			"type": "departements",
		},
		Scope: []string{"Mayenne"},
		Value: map[string]interface{}{
			"Mayenne": 53,
		},
	})

	regions = append(regions, daclient.Object{
		Key: map[string]string{
			"type": "departements",
		},
		Scope: []string{"Sarthe"},
		Value: map[string]interface{}{
			"Sarthe": 72,
		},
	})

	regions = append(regions, daclient.Object{
		Key: map[string]string{
			"type": "departements",
		},
		Scope: []string{"Vendée"},
		Value: map[string]interface{}{
			"Vendée": 85,
		},
	})

	return regions
}
