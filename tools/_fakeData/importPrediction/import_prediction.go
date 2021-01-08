package main

import (
	"bufio"
	"encoding/csv"
	"fmt"
	"io"
	"log"
	"os"
	"strconv"
	"time"

	"github.com/globalsign/mgo/bson"

	"github.com/globalsign/mgo"
)

// Prediction prédiction
type Prediction struct {
	ID        bson.ObjectId `bson:"_id"`
	Batch     string        `bson:"batch"`
	Algo      string        `bson:"algo"`
	Siret     string        `bson:"siret"`
	Score     float64       `bson:"score"`
	Diff      float64       `bson:"diff"`
	Periode   time.Time     `bson:"periode"`
	Timestamp time.Time     `bson:"timestamp"`
	Alert     string        `bson:"alert"`
}

func main() {
	csvFile, _ := os.Open("/home/christophe/Project/data-fake/fake-prediction.csv")

	predictionDict := make(map[string]Prediction)

	reader := csv.NewReader(bufio.NewReader(csvFile))
	reader.LazyQuotes = true
	reader.Comma = ';'
	reader.Read()
	var prediction []interface{}

	for {
		line, error := reader.Read()
		if error == io.EOF {
			break
		} else if error != nil {
			log.Fatal(error)
		}
		if line[0] != "" {
			proba, _ := strconv.ParseFloat(line[1], 64)
			var alert = "Pas d'alerte"

			if proba > 0.13 {
				alert = "Alerte seuil F2"
			}
			if proba > 0.37 {
				alert = "Alerte seuil F1"
			}

			diff, _ := strconv.ParseFloat(line[2], 64)

			p := Prediction{
				ID:        bson.NewObjectId(),
				Siret:     line[0],
				Batch:     "1905",
				Algo:      "algo",
				Score:     proba,
				Diff:      diff,
				Timestamp: time.Now(),
				Periode: time.Date(
					2019, 05, 01, 0, 0, 0, 0, time.UTC),
				Alert: alert,
			}

			predictionDict[line[0]] = p
		}
	}

	for _, v := range predictionDict {
		prediction = append(prediction, v)
	}

	mongodb, err := mgo.Dial("")
	if err != nil {
		fmt.Println("Insertion interrompue: " + err.Error())
		return
	}

	db := mongodb.DB("fakesignauxfaibles")

	err = db.C("Scores").Insert(prediction...)

	if err != nil {
		fmt.Println("Insertion interrompue: " + err.Error())
		return
	}

	fmt.Println("Prédictions insérées: " + strconv.Itoa(len(prediction)))
}
