package engine

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"regexp"
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

// ExportPoliciesToDatapi exports standard policies to datapi
func ExportPoliciesToDatapi(url, email, password, filter string) error {
	re, err := regexp.Compile(filter)
	if err != nil {
		return err
	}

	var policies = exportdatapi.GetPolicies()
	client := daclient.DatapiServer{
		URL:      url,
		Email:    email,
		Password: password,
		Bucket:   "system",
	}

	var packet []daclient.Object
	for _, p := range policies {
		if re.MatchString(p.Key["name"]) {
			packet = append(packet, p)
		}
	}

	err = datapiSecureSend(client, &packet, nil)
	return err
}

// ExportReferencesToDatapi pushes references (batches, types, etc.) to a datapi server
func ExportReferencesToDatapi(url string, email string, password string, batch string, algo string) error {
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
		URL:      url,
		Email:    email,
		Password: password,
		Bucket:   "reference",
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
	err = datapiSecureSend(client, &data, adminAlgo.Scope)

	return err
}

// ExportDetectionToDatapi sends detections with some informations to a datapi server
func ExportDetectionToDatapi(url, email, password, batch, key, algo string) error {
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

	var pipeline = exportdatapi.GetDetectionPipeline(batch, key, algo)

	iter := Db.DB.C("Scores").Pipe(pipeline).AllowDiskUse().Iter()

	connus, err := readConnu()
	if err != nil {
		return err
	}

	i := 0
	var datas []daclient.Object
	var data exportdatapi.Detection

	client := daclient.DatapiServer{
		URL:      url,
		Email:    email,
		Password: password,
		Bucket:   "detection",
	}

	for iter.Next(&data) {
		i++
		detection, err := exportdatapi.ComputeDetection(data, &connus)
		if err != nil {
			log.Println(err)
			continue
		}

		datas = append(datas, detection...)

		if i == viper.GetInt("datapiChunk") {
			i = 0
			err := datapiSecureSend(client, &datas, adminAlgo.Scope)
			if err != nil {
				return err
			}
			datas = nil
		}
	}

	if datas != nil {
		err = datapiSecureSend(client, &datas, adminAlgo.Scope)
	}

	return err
}

func datapiSecureSend(client daclient.DatapiServer, datas *[]daclient.Object, additionnalScope []string) error {
	var sendPacket []daclient.Object
	for _, d := range *datas {
		d.Scope = append(d.Scope, additionnalScope...)
		sendPacket = append(sendPacket, d)
	}

	if sendPacket != nil {
		err := client.Connect()

		i := 0
		for err != nil && i < 5 {
			i++
			log.Println("erreur de connexion datapi: " + err.Error())
			time.Sleep(5 * time.Second)

			log.Println("tentative de reconnexion: " + strconv.Itoa(i))
			err = client.Connect()

			if err == nil {
				break
			}

			if i == 5 {
				return err
			}
		}

		i = 0

		err = client.Put(sendPacket)
		for err != nil && i < 5 {
			i++
			log.Println("erreur de transmission datapi: " + err.Error())
			time.Sleep(5 * time.Second)

			log.Println("tentative de réémission: " + err.Error())
			err = client.Put(sendPacket)

			if err == nil {
				break
			}

			if i == 5 {
				return err
			}
		}
	}

	return nil
}

// ExportEtablissementToDatapi exporte les objets
func ExportEtablissementToDatapi(url, email, password, key string) error {
	connus, err := readConnu()
	if err != nil {
		return err
	}

	client := daclient.DatapiServer{
		URL:      url,
		Email:    email,
		Password: password,
		Bucket:   "entreprise",
		SendSize: 1000,
	}

	pipeline := exportdatapi.GetEtablissementPipeline(key)
	iter := Db.DB.C("Public").Pipe(pipeline).AllowDiskUse().Iter()

	var data exportdatapi.Etablissement
	datapi, waiter := client.Worker()

	for iter.Next(&data) {
		if data.Value.Sirene.Departement != "" {
			for _, d := range exportdatapi.ComputeEtablissement(data, &connus) {
				datapi <- d
			}
		} else {
			log.Println("Pas d'information Sirene, établissement ignoré:", data.Value.Key)
		}
	}
	close(datapi)
	waiter.Wait()

	if client.Errors > 0 {
		return errors.New("Erreurs détectées, envoi incomplet, plus d'informations dans le journal")
	}
	return nil
}

// ExportEntrepriseToFile exporte les entreprises et etablissements avec leurs
// scores, dans un fichier.
func ExportEntrepriseToFile(filepath string) error {
	file, err := os.Create(filepath)
	if err != nil {
		return err
	}
	defer file.Close()

	const indent = "  "

	_, err = file.Write([]byte("[\n"))
	if err != nil {
		return err
	}

	pipeline := exportdatapi.GetEntreprisePipeline()
	iter := Db.DB.C("Public").Pipe(pipeline).AllowDiskUse().Iter()

	/*
		type EntrepriseAvecEtablissements struct {
			ID    map[string]string `json:"_id" bson:"_id"`
			Value struct {
				Key            string                       `json:"key" bson:"key"`
				IDEntreprise   string                       `json:"idEntreprise" bson:"idEntreprise"`
				Etablissements []exportdatapi.Etablissement `json:"etablissements" bson:"etablissements"`
			} `bson:"value"`
		}
	*/

	var wroteOneElement = false
	var entreprise interface{} //EntrepriseAvecEtablissements
	for iter.Next(&entreprise) {
		if wroteOneElement {
			_, err = file.Write([]byte(",\n"))
			if err != nil {
				return err
			}
		}
		_, err = file.Write([]byte(indent))
		if err != nil {
			return err
		}
		bytesToWrite, err := json.MarshalIndent(entreprise, indent, indent)
		if err != nil {
			return err
		}
		nbBytesWritten := 0
		for nbBytesWritten < len(bytesToWrite) {
			bytesToWrite = bytesToWrite[nbBytesWritten:]
			nbBytesWritten, err = file.Write(bytesToWrite)
			_, err = fmt.Println("Printed", nbBytesWritten, "bytes /", len(bytesToWrite))
			if err != nil {
				return err
			}
		}
		wroteOneElement = true
	}
	file.Write([]byte("\n]\n"))
	if err != nil {
		return err
	}

	return nil
}
