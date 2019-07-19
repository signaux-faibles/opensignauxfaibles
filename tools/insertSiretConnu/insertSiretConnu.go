package main

import (
	"bufio"
	"encoding/csv"
	"fmt"
	"io"
	"log"
	"os"

	daclient "github.com/signaux-faibles/datapi/client"
	"github.com/spf13/viper"
)

func readSiretDep(siretDepFile string) map[string]string {
	file, err := os.Open(siretDepFile)
	if err != nil {
		panic(err)
	}
	reader := csv.NewReader(file)
	siretDep := make(map[string]string)
	reader.Comma = ';'
	for {
		row, err := reader.Read()
		if err == io.EOF {
			break
		} else if err != nil {
			panic(err)
		} else {
			siretDep[row[0]] = row[1]
		}
	}
	return siretDep
}

func main() {
	viper.SetConfigName("config")
	viper.SetConfigType("toml")
	viper.AddConfigPath(".")
	viper.ReadInConfig()

	// connexion datapi
	url := viper.GetString("url")
	password := viper.GetString("password")
	email := viper.GetString("email")
	client := daclient.DatapiServer{
		URL: url,
	}
	client.Connect(email, password)

	// fichier de référence siret/département
	siretDepFile := viper.GetString("siretDepFile")
	siretDep := readSiretDep(siretDepFile)

	// fichier à traiter
	file, _ := os.Open(viper.GetString("file"))
	reader := csv.NewReader(bufio.NewReader(file))
	reader.Comma = ';'

	var objects []daclient.Object
	var sirets = make(map[string]struct{})
	for {
		line, error := reader.Read()
		if error == io.EOF {
			break
		} else if error != nil {
			log.Fatal(error)
		}
		for k := range siretDep {
			if line[0][0:9] == k[0:9] {
				sirets[k] = struct{}{}
			}
		}
	}

	for siret := range sirets {
		if _, ok := siretDep[siret]; ok {
			o := daclient.Object{
				Key: map[string]string{
					"siret": siret,
					"type":  "detection",
					"batch": "1906_7",
				},
				Scope: []string{"detection", siretDep[siret]},
				Value: map[string]interface{}{
					"connu": true,
				},
			}

			objects = append(objects, o)
			fmt.Println(o)
		}
	}

	err := client.Put("public", objects)
	if err != nil {
		fmt.Println("failed: " + err.Error())
		os.Exit(255)
	}
	fmt.Println("ok")
}
