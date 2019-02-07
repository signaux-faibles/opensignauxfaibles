package apartdemande

import (
	"dbmongo/lib/engine"
	"testing"
)

func Test_parseAPDemande(t *testing.T) {
	batch := engine.AdminBatch{
		Files: engine.BatchFiles{
			"apdemande": []string{"testData/apdemande.excelsheet"},
		},
	}

	dataChannel, eventChannel := Parser(batch)
	engine.DiscardEvents(eventChannel)
	tuple := <-dataChannel
	data := tuple.(APDemande)
	if data.ID == "0210377020" &&
		data.Siret == "37215938570251" &&
		*data.EffectifEntreprise == 0 &&
		*data.Effectif == 0 {

		t.Log("Test parseAPDemande: lecture ok")
	} else {
		t.Error("Erreur parseAPDemande: données différentes du fichier de test")
	}

	tuple = <-dataChannel

	if tuple != nil {
		t.Error("Erreur parseAPDemande: le channel devrait être vide")
	} else {
		t.Log("Test parseAPDemande: test deuxième ligne ok")
	}

	batch = engine.AdminBatch{
		Files: engine.BatchFiles{
			"apdemande": []string{"testData/nonexisting.excelsheet"},
		},
	}

	dataChannel, eventChannel = Parser(batch)
	engine.DiscardEvents(eventChannel)
	tuple = <-dataChannel

	if tuple != nil {
		t.Error("Erreur parseAPDemande: fichier inexistant le channel devrait être vide")
	} else {
		t.Log("Test parseAPDemande: test fichier inexistant ok")
	}

	batch = engine.AdminBatch{
		Files: engine.BatchFiles{},
	}

	dataChannel, eventChannel = Parser(batch)
	engine.DiscardEvents(eventChannel)
	tuple = <-dataChannel

	if tuple != nil {
		t.Error("Erreur parseAPConso: aucun fichier, le channel devrait être vide")
	} else {
		t.Log("Test parseAPConso: test aucun fichier ok")
	}
}
