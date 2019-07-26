package engine

import (
	"dbmongo/lib/exportdatapi"
	"dbmongo/lib/naf"
	"fmt"
	"io/ioutil"
	"sort"
	"strings"

	"github.com/spf13/viper"

	daclient "github.com/signaux-faibles/datapi/client"
)

func readConnu() ([]string, error) {
	connus, err := ioutil.ReadFile(viper.GetString("sirens"))
	if err != nil {
		return nil, err
	}

	sirets := strings.Split(string(connus), "\n")
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
			fmt.Println(err)
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

		if i > 3000 {
			if datas != nil {
				client.Connect(user, password)
				err = client.Put("public", datas)
				if err != nil {
					fmt.Println(err)
				}
				datas = nil
			}
			i = 0
		}
	}

	if datas != nil {
		client.Connect(user, password)
		err = client.Put("public", datas)
		if err != nil {
			fmt.Println(err)
		}
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
				"Auvergne-Rhône-Alpes",
				"Pays de la Loire",
				"Bourgogne-Franche-Comté",
				"Nouvelle-Aquitaine",
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
				"01",
				"03",
				"07",
				"15",
				"26",
				"38",
				"42",
				"43",
				"63",
				"69",
				"73",
				"74",
				"16",
				"17",
				"19",
				"23",
				"24",
				"33",
				"40",
				"47",
				"64",
				"79",
				"86",
				"87",
			},
		},
	})

	policies = append(policies, daclient.Object{
		Key: map[string]string{
			"type": "policy",
			"name": "Accès Bourgogne-Franche-Comté",
		},
		Value: map[string]interface{}{
			"match": "(public|reference)",
			"scope": []string{"Bourgogne-Franche-Comté"},
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
			"name": "Accès Auvergne-Rhône-Alpes",
		},
		Value: map[string]interface{}{
			"match": "(public|reference)",
			"scope": []string{"Auvergne-Rhône-Alpes"},
			"promote": []string{
				"01",
				"03",
				"07",
				"15",
				"26",
				"38",
				"42",
				"43",
				"63",
				"69",
				"73",
				"74",
			},
		},
	})

	policies = append(policies, daclient.Object{
		Key: map[string]string{
			"type": "policy",
			"name": "Accès Nouvelle-Aquitaine",
		},
		Value: map[string]interface{}{
			"match": "(public|reference)",
			"scope": []string{"Nouvelle-Aquitaine"},
			"promote": []string{
				"16",
				"17",
				"19",
				"23",
				"24",
				"33",
				"40",
				"47",
				"64",
				"79",
				"86",
				"87",
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
			"match":  "system",
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
				"01",
				"03",
				"07",
				"15",
				"26",
				"38",
				"42",
				"43",
				"63",
				"69",
				"73",
				"74",
				"16",
				"17",
				"19",
				"23",
				"24",
				"33",
				"40",
				"47",
				"64",
				"79",
				"86",
				"87",
			},
		},
	})

	regions = append(regions, daclient.Object{
		Key: map[string]string{
			"type":   "region",
			"region": "Bourgogne-Franche-Comté",
			"batch":  batch,
		},
		Scope: []string{"Bourgogne-Franche-Comté"},
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
			"type":   "region",
			"region": "Auvergne-Rhône-Alpes",
			"batch":  batch,
		},
		Scope: []string{"Auvergne-Rhône-Alpes"},
		Value: map[string]interface{}{
			"departements": []string{
				"01",
				"03",
				"07",
				"15",
				"26",
				"38",
				"42",
				"43",
				"63",
				"69",
				"73",
				"74",
			},
		},
	})

	regions = append(regions, daclient.Object{
		Key: map[string]string{
			"type":   "region",
			"region": "Nouvelle-Aquitaine",
			"batch":  batch,
		},
		Scope: []string{"Nouvelle-Aquitaine"},
		Value: map[string]interface{}{
			"departements": []string{
				"16",
				"17",
				"19",
				"23",
				"24",
				"33",
				"40",
				"47",
				"64",
				"79",
				"86",
				"87",
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

	regions = append(regions, daclient.Object{
		Key: map[string]string{
			"type":  "departements",
			"batch": batch,
		},
		Scope: []string{"01"},
		Value: map[string]interface{}{
			"01": "Ain",
		},
	})

	regions = append(regions, daclient.Object{
		Key: map[string]string{
			"type":  "departements",
			"batch": batch,
		},
		Scope: []string{"03"},
		Value: map[string]interface{}{
			"03": "Allier",
		},
	})

	regions = append(regions, daclient.Object{
		Key: map[string]string{
			"type":  "departements",
			"batch": batch,
		},
		Scope: []string{"07"},
		Value: map[string]interface{}{
			"07": "Ardèche",
		},
	})

	regions = append(regions, daclient.Object{
		Key: map[string]string{
			"type":  "departements",
			"batch": batch,
		},
		Scope: []string{"15"},
		Value: map[string]interface{}{
			"15": "Cantal",
		},
	})

	regions = append(regions, daclient.Object{
		Key: map[string]string{
			"type":  "departements",
			"batch": batch,
		},
		Scope: []string{"26"},
		Value: map[string]interface{}{
			"26": "Drôme",
		},
	})

	regions = append(regions, daclient.Object{
		Key: map[string]string{
			"type":  "departements",
			"batch": batch,
		},
		Scope: []string{"38"},
		Value: map[string]interface{}{
			"38": "Isère",
		},
	})

	regions = append(regions, daclient.Object{
		Key: map[string]string{
			"type":  "departements",
			"batch": batch,
		},
		Scope: []string{"42"},
		Value: map[string]interface{}{
			"42": "Loire",
		},
	})

	regions = append(regions, daclient.Object{
		Key: map[string]string{
			"type":  "departements",
			"batch": batch,
		},
		Scope: []string{"43"},
		Value: map[string]interface{}{
			"43": "Haute-Loire",
		},
	})

	regions = append(regions, daclient.Object{
		Key: map[string]string{
			"type":  "departements",
			"batch": batch,
		},
		Scope: []string{"63"},
		Value: map[string]interface{}{
			"63": "Puy-de-Dôme",
		},
	})

	regions = append(regions, daclient.Object{
		Key: map[string]string{
			"type":  "departements",
			"batch": batch,
		},
		Scope: []string{"69"},
		Value: map[string]interface{}{
			"69": "Rhône",
		},
	})

	regions = append(regions, daclient.Object{
		Key: map[string]string{
			"type":  "departements",
			"batch": batch,
		},
		Scope: []string{"73"},
		Value: map[string]interface{}{
			"73": "Savoie",
		},
	})

	regions = append(regions, daclient.Object{
		Key: map[string]string{
			"type":  "departements",
			"batch": batch,
		},
		Scope: []string{"74"},
		Value: map[string]interface{}{
			"74": "Haute-Savoie",
		},
	})

	regions = append(regions, daclient.Object{
		Key: map[string]string{
			"type":  "departements",
			"batch": batch,
		},
		Scope: []string{"16"},
		Value: map[string]interface{}{
			"16": "Charente",
		},
	})

	regions = append(regions, daclient.Object{
		Key: map[string]string{
			"type":  "departements",
			"batch": batch,
		},
		Scope: []string{"17"},
		Value: map[string]interface{}{
			"17": "Charente-Maritime",
		},
	})

	regions = append(regions, daclient.Object{
		Key: map[string]string{
			"type":  "departements",
			"batch": batch,
		},
		Scope: []string{"19"},
		Value: map[string]interface{}{
			"19": "Corrèze",
		},
	})

	regions = append(regions, daclient.Object{
		Key: map[string]string{
			"type":  "departements",
			"batch": batch,
		},
		Scope: []string{"23"},
		Value: map[string]interface{}{
			"23": "Creuse",
		},
	})

	regions = append(regions, daclient.Object{
		Key: map[string]string{
			"type":  "departements",
			"batch": batch,
		},
		Scope: []string{"24"},
		Value: map[string]interface{}{
			"24": "Dordogne",
		},
	})

	regions = append(regions, daclient.Object{
		Key: map[string]string{
			"type":  "departements",
			"batch": batch,
		},
		Scope: []string{"33"},
		Value: map[string]interface{}{
			"33": "Gironde",
		},
	})

	regions = append(regions, daclient.Object{
		Key: map[string]string{
			"type":  "departements",
			"batch": batch,
		},
		Scope: []string{"40"},
		Value: map[string]interface{}{
			"40": "Landes",
		},
	})

	regions = append(regions, daclient.Object{
		Key: map[string]string{
			"type":  "departements",
			"batch": batch,
		},
		Scope: []string{"47"},
		Value: map[string]interface{}{
			"47": "Lot-et-Garonne",
		},
	})

	regions = append(regions, daclient.Object{
		Key: map[string]string{
			"type":  "departements",
			"batch": batch,
		},
		Scope: []string{"64"},
		Value: map[string]interface{}{
			"64": "Pyrénées-Atlantique",
		},
	})

	regions = append(regions, daclient.Object{
		Key: map[string]string{
			"type":  "departements",
			"batch": batch,
		},
		Scope: []string{"79"},
		Value: map[string]interface{}{
			"79": "Deux-Sèvres",
		},
	})

	regions = append(regions, daclient.Object{
		Key: map[string]string{
			"type":  "departements",
			"batch": batch,
		},
		Scope: []string{"86"},
		Value: map[string]interface{}{
			"86": "Vienne",
		},
	})

	regions = append(regions, daclient.Object{
		Key: map[string]string{
			"type":  "departements",
			"batch": batch,
		},
		Scope: []string{"87"},
		Value: map[string]interface{}{
			"87": "Haute-Vienne",
		},
	})

	return regions
}
