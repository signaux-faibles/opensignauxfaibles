package apartconso

// func Test_parseAPConso(t *testing.T) {
// 	batch := engine.AdminBatch{
// 		Files: engine.BatchFiles{
// 			"apconso": []string{"testData/apconso.excelsheet"},
// 		},
// 	}

// 	dataChannel, eventChannel := Parser(batch)
// 	engine.DiscardEvents(eventChannel)
// 	tuple := <-dataChannel
// 	data := tuple.(APConso)

// 	if *data.Effectif == 30 &&
// 		*data.HeureConsommee == 100 &&
// 		data.ID == "123456789" &&
// 		data.Siret == "12345678901234" &&
// 		*data.Montant == 200 {
// 		t.Log("Test parseAPConso: lecture ok")
// 	} else {
// 		t.Error("Erreur parseAPConso: données différentes du fichier de test")
// 	}

// 	tuple = <-dataChannel

// 	if tuple != nil {
// 		t.Error("Erreur parseAPConso: 2° ligne vide: le channel devrait être vide")
// 	} else {
// 		t.Log("Test parseAPConso: test deuxième ligne ok")
// 	}

// 	batch = engine.AdminBatch{
// 		Files: engine.BatchFiles{
// 			"apconso": []string{"testData/nonexistant.excelsheet"},
// 		},
// 	}

// 	dataChannel, eventChannel = Parser(batch)
// 	engine.DiscardEvents(eventChannel)
// 	tuple = <-dataChannel

// 	if tuple != nil {
// 		t.Error("Erreur parseAPConso: fichier absent: le channel devrait être vide")
// 	} else {
// 		t.Log("Test parseAPConso: fichier absent ok")
// 	}
// }
