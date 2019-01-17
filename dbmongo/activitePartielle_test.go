package main

import (
	"fmt"
	"testing"
)

func Test_parseAPConso(t *testing.T) {
	testChannel := parseAPConso("testData/test-apconso.xlsx")
	data := <-testChannel
	fmt.Println(db)
	if *data.Effectif == 30 &&
		*data.HeureConsommee == 100 &&
		data.ID == "123456789" &&
		data.Siret == "12345678901234" &&
		*data.Montant == 200 {
		t.Log("Test parseAPConso: lecture ok")
	} else {
		t.Error("Erreur parseAPConso: données différentes du fichier de test")
	}

	data = <-testChannel

	if data != nil {
		t.Error("Erreur parseAPConso: le channel devrait être vide")
	} else {
		t.Log("Test parseAPConso: test deuxième ligne ok")
	}

}
