package engine

import (
	"dbmongo/lib/exportdatapi"
	"dbmongo/lib/naf"
	"fmt"
	"io/ioutil"
	"sort"
	"strings"

	daclient "github.com/signaux-faibles/datapi/client"
)

func readConnu() ([]string, error) {
	bfc, err := ioutil.ReadFile("/home/christophe/Project/data-raw/1905/sirets_connus_bfc.csv")
	if err != nil {
		return nil, err
	}
	pdl, err := ioutil.ReadFile("/home/christophe/Project/data-raw/1905/sirets_connus_pdl.csv")
	if err != nil {
		return nil, err
	}
	sirets := strings.Split(string(bfc), "\n")
	sirets = append(sirets, strings.Split(string(pdl), "\n")...)
	sort.Strings(sirets)
	return sirets, nil
}

func findString(s string, a []string) bool {
	for _, v := range a {
		if v != "" && s[0:9] == v[0:9] {
			return true
		}
	}
	return false
}

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
	connus, err := readConnu()
	if err != nil {
		return err
	}

	for iter.Next(&data) {
		detection, err := exportdatapi.Compute(data)
		if err != nil {
			continue
		}
		i++

		c := daclient.Object{
			Key: map[string]string{
				"siret": data.ID["key"],
				"batch": data.ID["batch"],
				"type":  "detection",
			},
			Scope: []string{"detection", "score", data.Etablissement.Value.Sirene.Departement},
			Value: map[string]interface{}{
				"connu": findString(data.ID["key"], connus),
			},
		}

		datas = append(datas, detection...)
		datas = append(datas, c)

		if i > 10000 {
			if datas != nil {
				err = client.Put("public", datas)
				datas = nil
			}
			i = 0
		}
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
			"name": "Accès France",
		},
		Value: map[string]interface{}{
			"match": "(public|reference)",
			"scope": []string{"France entière"},
			"promote": []string{
				"Pays de la Loire",
				"Bourgogne Franche-Comté",
				"Bourgogne",
				"Franche-Comté",
				"21",
				"25",
				"70",
				"39",
				"58",
				"71",
				"90",
				"89",
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
			"name": "Accès Bourgogne Franche-Comté",
		},
		Value: map[string]interface{}{
			"match": "(public|reference)",
			"scope": []string{"Bourgogne Franche-Comté"},
			"promote": []string{
				"Bourgogne",
				"Franche-Comté",
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

	procol := daclient.Object{
		Key: map[string]string{
			"key":   "procol",
			"batch": batch,
		},
		Scope: []string{},
		Value: map[string]interface{}{
			"in_bonis":          "In bonis",
			"continuation":      "Plan de continuation",
			"plan_redressement": "Redressement judiciaire",
			"plan_sauvegarde":   "Plan de sauvegarde",
			"liquidation":       "Liquidation judiciaire",
			"sauvegarde":        "Plan de sauvegarde",
		},
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
	data = append(data, procol)
	data = append(data, batchObject)
	data = append(data, getRegions(batch)...)
	err = client.Put("reference", data)

	return err
}

func getRegions(batch string) (regions []daclient.Object) {
	regions = append(regions, daclient.Object{
		Key: map[string]string{
			"type":   "region",
			"region": "France entière",
			"batch":  batch,
		},
		Scope: []string{"France entière"},
		Value: map[string]interface{}{
			"departements": []string{
				"25",
				"70",
				"39",
				"90",
				"21",
				"58",
				"71",
				"89",
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
			"type":   "region",
			"region": "Bourgogne Franche-Comté",
			"batch":  batch,
		},
		Scope: []string{"Bourgogne Franche-Comté"},
		Value: map[string]interface{}{
			"departements": []string{
				"25",
				"70",
				"39",
				"90",
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
			"region": "Bourgogne",
			"batch":  batch,
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
			"batch":  batch,
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
			"batch":  batch,
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
			"type":  "departements",
			"batch": batch,
		},
		Scope: []string{"25"},
		Value: map[string]interface{}{
			"25": "Doubs",
		},
	})

	regions = append(regions, daclient.Object{
		Key: map[string]string{
			"type":  "departements",
			"batch": batch,
		},
		Scope: []string{"70"},
		Value: map[string]interface{}{
			"70": "Haute-Saône",
		},
	})

	regions = append(regions, daclient.Object{
		Key: map[string]string{
			"type":  "departements",
			"batch": batch,
		},
		Scope: []string{"39"},
		Value: map[string]interface{}{
			"39": "Jura",
		},
	})

	regions = append(regions, daclient.Object{
		Key: map[string]string{
			"type":  "departements",
			"batch": batch,
		},
		Scope: []string{"90"},
		Value: map[string]interface{}{
			"90": "Territoire de Belfort",
		},
	})

	regions = append(regions, daclient.Object{
		Key: map[string]string{
			"type":  "departements",
			"batch": batch,
		},
		Scope: []string{"21"},
		Value: map[string]interface{}{
			"21": "Côte d'or",
		},
	})

	regions = append(regions, daclient.Object{
		Key: map[string]string{
			"type":  "departements",
			"batch": batch,
		},
		Scope: []string{"58"},
		Value: map[string]interface{}{
			"58": "Nièvre",
		},
	})

	regions = append(regions, daclient.Object{
		Key: map[string]string{
			"type":  "departements",
			"batch": batch,
		},
		Scope: []string{"71"},
		Value: map[string]interface{}{
			"71": "Saône et Loire",
		},
	})

	regions = append(regions, daclient.Object{
		Key: map[string]string{
			"type":  "departements",
			"batch": batch,
		},
		Scope: []string{"89"},
		Value: map[string]interface{}{
			"89": "Yonne",
		},
	})

	regions = append(regions, daclient.Object{
		Key: map[string]string{
			"type":  "departements",
			"batch": batch,
		},
		Scope: []string{"44"},
		Value: map[string]interface{}{
			"44": "Loire-Atlantique",
		},
	})

	regions = append(regions, daclient.Object{
		Key: map[string]string{
			"type":  "departements",
			"batch": batch,
		},
		Scope: []string{"49"},
		Value: map[string]interface{}{
			"49": "Maine-et-Loire",
		},
	})

	regions = append(regions, daclient.Object{
		Key: map[string]string{
			"type":  "departements",
			"batch": batch,
		},
		Scope: []string{"53"},
		Value: map[string]interface{}{
			"53": "Mayenne",
		},
	})

	regions = append(regions, daclient.Object{
		Key: map[string]string{
			"type":  "departements",
			"batch": batch,
		},
		Scope: []string{"72"},
		Value: map[string]interface{}{
			"72": "Sarthe",
		},
	})

	regions = append(regions, daclient.Object{
		Key: map[string]string{
			"type":  "departements",
			"batch": batch,
		},
		Scope: []string{"85"},
		Value: map[string]interface{}{
			"85": "Vendée",
		},
	})

	return regions
}
