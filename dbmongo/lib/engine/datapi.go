package engine

import (
	"io/ioutil"
	"log"
	"sort"
	"strings"
	"time"

	"github.com/signaux-faibles/opensignauxfaibles/dbmongo/lib/exportdatapi"
	"github.com/signaux-faibles/opensignauxfaibles/dbmongo/lib/naf"

	"github.com/davecgh/go-spew/spew"
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

// ExportPoliciesToDatapi exports standard policies to datapi
func ExportPoliciesToDatapi(url, user, password, batch string) error {
	return exportdatapi.ExportPoliciesToDatapi(url, user, password, batch)
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
	data = append(data, exportdatapi.GetRegions(batch)...)
	err = client.Put("reference", data)

	return err
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

	iter := Db.DB.C("Scores").Pipe(pipeline).AllowDiskUse().Iter()

	var data exportdatapi.Detection

	var datas []daclient.Object

	// var i int
	connus, err := readConnu()
	if err != nil {
		return err
	}

	i := 0
	for iter.Next(&data) {
		spew.Dump(data)
		detection, err := exportdatapi.Compute(data)
		//spew.Dump(detection)
		if err != nil {
			log.Println(err)
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

		// envoi de tronçons de 2000 entreprises
		if i == 2000 {
			i = 0
			datapiSecureSend(user, password, &client, &datas)
			datas = nil
		}
	}

	if datas != nil {
		client.Connect(user, password)
		err = client.Put("public", datas)
		if err != nil {
			log.Println(err)
		}
	}

	return err
}

func datapiSecureSend(user string, password string, client *daclient.DatapiServer, datas *[]daclient.Object) error {
	if datas != nil {
		err := client.Connect(user, password)
		for err != nil {
			log.Println("erreur de connexion datapi: " + err.Error())
			log.Println("tentative de reconnexion")
			time.Sleep(5 * time.Second)
			err = client.Connect(user, password)
		}

		err = client.Put("public", *datas)

		if err != nil {
			log.Println(err.Error())
			return err
		}
	}
	return nil
}
