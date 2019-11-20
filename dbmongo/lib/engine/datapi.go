package engine

import (
	"fmt"
	"io/ioutil"
	"log"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/signaux-faibles/opensignauxfaibles/dbmongo/lib/exportdatapi"
	"github.com/signaux-faibles/opensignauxfaibles/dbmongo/lib/naf"

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
		if len(v) > 9 && len(s) > 9 && s[0:9] == v[0:9] {
			return true
		}
	}
	return false
}

// ExportPoliciesToDatapi exports standard policies to datapi
func ExportPoliciesToDatapi(url, user, password string) error {
	// return exportdatapi.ExportPoliciesToDatapi(url, user, password, batch)
	var policies = exportdatapi.GetPolicies()
	client := daclient.DatapiServer{
		URL: url,
	}

	err := datapiSecureSend(user, password, "system", &client, &policies, nil)
	return err
}

// ExportReferencesToDatapi pushes references (batches, types, etc.) to a datapi server
func ExportReferencesToDatapi(url string, user string, password string, batch string, algo string) error {
	var adminAlgo AdminAlgo
	err := adminAlgo.Load(algo)

	if err != nil {
		return fmt.Errorf("algorithme %s inconnu: %s", algo, err.Error())
	}

	var adminBatch AdminBatch
	err = adminBatch.Load(batch)

	if err != nil {
		return fmt.Errorf("batch %s inconnu: %s", batch, err.Error())
	}

	client := daclient.DatapiServer{
		URL: url,
	}

	nafCodes := daclient.Object{
		Key: map[string]string{
			"key":   "naf",
			"batch": batch + "." + algo,
		},
		Value: naf.Naf.ToData(),
	}

	procol := daclient.Object{
		Key: map[string]string{
			"key":   "procol",
			"batch": batch + "." + algo,
		},
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
			"batch": batch + "." + algo,
		},
		Value: GetTypes().ToData(),
	}

	batchData, err := GetBatch(batch)
	if err != nil {
		return err
	}

	batchObject := daclient.Object{
		Key: map[string]string{
			"key":   "batch",
			"batch": batchData.ID.Key + "." + algo,
		},
		Value: batchData.ToData(adminAlgo.Label),
	}

	var data []daclient.Object
	data = append(data, nafCodes)
	data = append(data, types)
	data = append(data, procol)
	data = append(data, batchObject)
	data = append(data, exportdatapi.GetRegions(batch, algo)...)
	err = datapiSecureSend(user, password, "reference", &client, &data, adminAlgo.Scope)

	return err
}

// ExportDetectionToDatapi sends detections with some informations to a datapi server
func ExportDetectionToDatapi(url, user, password, batch, key, algo string) error {
	var adminAlgo AdminAlgo
	err := adminAlgo.Load(algo)
	if err != nil {
		return fmt.Errorf("algorithme %s inconnu: %s", algo, err.Error())
	}

	var adminBatch AdminBatch
	err = adminBatch.Load(batch)
	if err != nil {
		return fmt.Errorf("batch %s inconnu: %s", batch, err.Error())
	}

	var pipeline = exportdatapi.GetPipeline(batch, key, algo)

	iter := Db.DB.C("Scores").Pipe(pipeline).AllowDiskUse().Iter()

	connus, err := readConnu()
	if err != nil {
		return err
	}

	i := 0
	var datas []daclient.Object
	var data exportdatapi.Detection

	client := daclient.DatapiServer{
		URL: url,
	}

	for iter.Next(&data) {
		i++
		detection, err := exportdatapi.Compute(data)
		if err != nil {
			log.Println(err)
			continue
		}

		datas = append(datas, detection...)

		// fast & dirty: intègre la notion de connu dans l'objet
		urssaf := exportdatapi.UrssafScope(data.Etablissement.Value.Compte.Numero, data.Etablissement.Value.Sirene.Departement)

		c := daclient.Object{
			Key: map[string]string{
				"siret": data.ID["siret"],
				"siren": data.ID["siret"][0:9],
				"batch": data.ID["batch"] + "." + data.ID["algo"],
				"type":  "detection",
				urssaf:  "true",
			},
			Scope: []string{"detection", "score", data.Etablissement.Value.Sirene.Departement},
			Value: map[string]interface{}{
				"connu": findString(data.ID["key"], connus),
			},
		}
		datas = append(datas, c)

		if i == viper.GetInt("datapiChunk") {
			i = 0
			err := datapiSecureSend(user, password, "public", &client, &datas, adminAlgo.Scope)
			if err != nil {
				return err
			}
			datas = nil
		}
	}

	if datas != nil {
		err = datapiSecureSend(user, password, "public", &client, &datas, adminAlgo.Scope)
	}

	return err
}

func datapiSecureSend(user string, password string, bucket string, client *daclient.DatapiServer, datas *[]daclient.Object, additionnalScope []string) error {
	var sendPacket []daclient.Object
	for _, d := range *datas {
		d.Scope = append(d.Scope, additionnalScope...)
		sendPacket = append(sendPacket, d)
	}

	if sendPacket != nil {
		err := client.Connect(user, password)

		i := 0
		for err != nil && i < 5 {
			i++
			log.Println("erreur de connexion datapi: " + err.Error())
			time.Sleep(5 * time.Second)

			log.Println("tentative de reconnexion: " + strconv.Itoa(i))
			err = client.Connect(user, password)

			if i == 5 {
				return err
			}
		}

		i = 0

		err = client.Put(bucket, sendPacket)
		for err != nil && i < 5 {
			i++
			log.Println("erreur de transmission datapi: " + err.Error())
			time.Sleep(5 * time.Second)

			log.Println("tentative de réémission:" + err.Error())
			err = client.Put(bucket, sendPacket)
		}
	}

	return nil
}
