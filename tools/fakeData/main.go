package main

import (
	"flag"
	"fmt"
	"math/rand"
	"strings"
	"time"

	"github.com/spf13/viper"
)

// Pour déformer des fichiers réels pour créer un dataset consistent avec anonymisation des entreprises

type randomizer func(string, string, map[string]string) error

func contains(array []int, element int) bool {
	for _, i := range array {
		if i == element {
			return true
		}
	}
	return false
}

// run execute une fonction de randomisation
func run(name string, handler randomizer, mapping map[string]string) error {
	file := viper.GetString(name)
	outputFile := outputFileNamePrefixer(viper.GetString("prefixOutput"), file)
	fmt.Print("Fake " + name + ": ")
	err := handler(file, outputFile, mapping)
	if err != nil {
		fmt.Println("Fail : " + err.Error())
		fmt.Println("Interruption.")
		return fmt.Errorf("Interruption")
	}
	fmt.Println("OK -> " + outputFile)
	return nil
}

func main() {
	rand.Seed(time.Now().UTC().UnixNano())
	flag.Parse()

	viper.SetConfigName("config")
	viper.SetConfigType("toml")
	viper.AddConfigPath(".")
	err := viper.ReadInConfig()
	if err != nil {
		panic("Erreur à la lecture de la configuration:" + err.Error())
	}

	// Génération des numéros de comptes
	outputCompte := outputFileNamePrefixer(viper.GetString("prefixOutput"), viper.GetString("comptes"))
	fmt.Print("Fake comptes: ")
	mapping, err := readAndRandomComptes(viper.GetString("comptes"), outputCompte)
	if err != nil {
		fmt.Println("Fail : " + err.Error())
		fmt.Println("Interruption.")
	} else {
		fmt.Println("OK -> " + outputCompte)
	}

	randomizers := map[string]randomizer{
		"diane": readAndRandomDiane,
		// "apartdemande": readAndRandomApartDemande,
		// "apartconso":   readAndRandomApartConso,
		// "bdf":          readAndRandomBDF,
		// "emploi":       readAndRandomEmploi,
		// "delais":       readAndRandomDelais,
		// "sirene":       readAndRandomSirene,
		// "debits":       readAndRandomDebits,
		// "altares":      readAndRandomAltares,
		// "cotisations":  readAndRandomCotisations,
		// "prediction":   readAndRandomPrediction,
	}

	for k, r := range randomizers {
		err := run(k, r, mapping)
		if err != nil {
			panic(err)
		}
	}
}

func outputFileNamePrefixer(prefixOutput string, fileName string) string {
	path := strings.Split(fileName, "/")
	path[len(path)-1] = prefixOutput + path[len(path)-1]

	return strings.Join(path, "/")
}

const letterBytes = "0123456789"

func randStringBytesRmndr(n int) string {
	b := make([]byte, n)
	for i := range b {
		b[i] = letterBytes[rand.Int63()%int64(len(letterBytes))]
	}
	return string(b)
}
