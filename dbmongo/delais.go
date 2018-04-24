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

	"github.com/cnf/structhash"
)

// Delais tuple fichier ursaff
type Delais struct {
	NumeroCompte      string    `json:"numero_compte" bson:"numero_compte"`
	NumeroContentieux string    `json:"numero_contentieux" bson:"numero_contentieux"`
	DateCreation      time.Time `json:"date_creation" bson:"date_creation"`
	DateEcheanche     time.Time `json:"date_echeance" bson:"date_echeance"`
	DureeDelai        int       `json:"duree_delai" bson:"duree_delai"`
	Denomination      string    `json:"denomination" bson:"denomination"`
	Indic6m           string    `json:"indic_6m" bson:"indic_6m"`
	AnneeCreation     int       `json:"annee_creation" bson:"annee_creation"`
	MontantEcheancier float64   `json:"montant_echeancier" bson:"montant_echeancier"`
	NumeroStructure   string    `json:"numero_structure" bson:"numero_structure"`
	Stade             string    `json:"stade" bson:"stade"`
	Action            string    `json:"action" bson:"action"`
}

func parseDelais(paths []string, batch string) chan Etablissement {
	outputChannel := make(chan Etablissement)

	field := map[string]int{
		"NumeroCompte":      0,
		"NumeroContentieux": 1,
		"DateCreation":      2,
		"DateEcheanche":     3,
		"DureeDelai":        4,
		"Denomination":      5,
		"Indic6m":           6,
		"AnneeCreation":     7,
		"MontantEcheancier": 8,
		"NumeroStructure":   9,
		"Stade":             10,
		"Action":            11,
	}

	go func() {
		for _, path := range paths {

			file, err := os.Open(path)
			if err != nil {
				fmt.Println("Error", err)
			}

			reader := csv.NewReader(bufio.NewReader(file))
			reader.Comma = ';'
			reader.Read()

			for {
				row, error := reader.Read()
				if error == io.EOF {
					break
				} else if error != nil {
					log.Fatal(error)
				}

				delais := Delais{}
				delais.NumeroCompte = row[field["NumeroCompte"]]
				delais.NumeroContentieux = row[field["NumeroContentieux"]]
				delais.DateCreation, err = time.Parse("2006-01-02", row[field["DateCreation"]])
				delais.DateEcheanche, err = time.Parse("2006-01-02", row[field["DateEcheanche"]])
				delais.DureeDelai, err = strconv.Atoi(row[field["DureeDelai"]])
				delais.Denomination = row[field["Denomination"]]
				delais.Indic6m = row[field["Indic6m"]]
				delais.AnneeCreation, err = strconv.Atoi(row[field["AnneeCreation"]])
				delais.MontantEcheancier, err = strconv.ParseFloat(row[field["MontantEcheancier"]], 64)
				delais.NumeroStructure = row[field["NumeroStructure"]]
				delais.Stade = row[field["Stade"]]
				delais.Action = row[field["Action"]]

				hash := fmt.Sprintf("%x", structhash.Md5(delais, 1))

				outputChannel <- Etablissement{
					Key: row[field["NumeroCompte"]],
					Batch: map[string]Batch{
						batch: Batch{
							Delais: map[string]Delais{
								hash: delais,
							},
						},
					},
				}

			}
		}
		close(outputChannel)
	}()

	return outputChannel
}
