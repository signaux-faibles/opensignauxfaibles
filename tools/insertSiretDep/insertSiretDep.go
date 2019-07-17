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

func main() {
	viper.SetConfigName("config")
	viper.SetConfigType("toml")
	viper.AddConfigPath(".")
	viper.ReadInConfig()

	url := viper.GetString("url")
	password := viper.GetString("password")
	email := viper.GetString("email")

	client := daclient.DatapiServer{
		URL: url,
	}
	client.Connect(email, password)
	file, _ := os.Open(viper.GetString("file"))
	reader := csv.NewReader(bufio.NewReader(file))
	reader.Comma = ';'

	var objects []daclient.Object

	for {
		line, error := reader.Read()
		if error == io.EOF {
			break
		} else if error != nil {
			log.Fatal(error)
		}
		o := daclient.Object{
			Key: map[string]string{
				"siret": line[0],
				"type":  "detection",
				"batch": "1906_7",
			},
			Scope: []string{"detection", line[2]},
			Value: map[string]interface{}{
				"raison_sociale": line[1],
			},
		}
		objects = append(objects, o)
		if len(objects) > 1000 {
			err := client.Put("public", objects)
			objects = nil
			if err != nil {
				fmt.Println("failed: " + err.Error())
			}
		}
	}
	err := client.Put("public", objects)
	if err != nil {
		fmt.Println("failed: " + err.Error())
		os.Exit(255)
	}
	fmt.Println("ok")
}
