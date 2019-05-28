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

	iter := Db.DB.C("Scores").Pipe(pipeline).Iter()

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
				"21",
				"25",
				"70",
				"39",
				"58",
				"71",
				"90",
				"89",
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
				"21",
				"58",
				"71",
				"89",
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
				"25",
				"70",
				"39",
				"90",
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
				"44",
				"49",
				"53",
				"72",
				"85",
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
				"21",
				"58",
				"71",
				"89",
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
				"25",
				"70",
				"39",
				"90",
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
				"44",
				"49",
				"53",
				"72",
				"85",
			},
		},
	})

	regions = append(regions, daclient.Object{
		Key: map[string]string{
			"type": "departements",
		},
		Scope: []string{"25"},
		Value: map[string]interface{}{
			"25": 25,
		},
	})

	regions = append(regions, daclient.Object{
		Key: map[string]string{
			"type": "departements",
		},
		Scope: []string{"70"},
		Value: map[string]interface{}{
			"70": 70,
		},
	})

	regions = append(regions, daclient.Object{
		Key: map[string]string{
			"type": "departements",
		},
		Scope: []string{"39"},
		Value: map[string]interface{}{
			"39": 39,
		},
	})

	regions = append(regions, daclient.Object{
		Key: map[string]string{
			"type": "departements",
		},
		Scope: []string{"90"},
		Value: map[string]interface{}{
			"90": 90,
		},
	})

	regions = append(regions, daclient.Object{
		Key: map[string]string{
			"type": "departements",
		},
		Scope: []string{"21"},
		Value: map[string]interface{}{
			"21": 21,
		},
	})

	regions = append(regions, daclient.Object{
		Key: map[string]string{
			"type": "departements",
		},
		Scope: []string{"58"},
		Value: map[string]interface{}{
			"58": 58,
		},
	})

	regions = append(regions, daclient.Object{
		Key: map[string]string{
			"type": "departements",
		},
		Scope: []string{"71"},
		Value: map[string]interface{}{
			"71": 71,
		},
	})

	regions = append(regions, daclient.Object{
		Key: map[string]string{
			"type": "departements",
		},
		Scope: []string{"89"},
		Value: map[string]interface{}{
			"89": 89,
		},
	})

	regions = append(regions, daclient.Object{
		Key: map[string]string{
			"type": "departements",
		},
		Scope: []string{"44"},
		Value: map[string]interface{}{
			"44": 44,
		},
	})

	regions = append(regions, daclient.Object{
		Key: map[string]string{
			"type": "departements",
		},
		Scope: []string{"49"},
		Value: map[string]interface{}{
			"49": 49,
		},
	})

	regions = append(regions, daclient.Object{
		Key: map[string]string{
			"type": "departements",
		},
		Scope: []string{"53"},
		Value: map[string]interface{}{
			"53": 53,
		},
	})

	regions = append(regions, daclient.Object{
		Key: map[string]string{
			"type": "departements",
		},
		Scope: []string{"72"},
		Value: map[string]interface{}{
			"72": 72,
		},
	})

	regions = append(regions, daclient.Object{
		Key: map[string]string{
			"type": "departements",
		},
		Scope: []string{"85"},
		Value: map[string]interface{}{
			"85": 85,
		},
	})

	return regions
}
