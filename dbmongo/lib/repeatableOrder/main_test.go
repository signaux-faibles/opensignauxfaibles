package repeatableOrder

import (
  "dbmongo/lib/engine"
  "testing"
  "time"
)

func Test_parseRepeatableOrder(t *testing.T) {
  batch := engine.AdminBatch{
    Files: engine.BatchFiles{
      "repeatable_order": []string{"testData/repeatable_test.csv"},
    },
  }

  dataChannel, eventChannel := Parser(batch)
  engine.DiscardEvents(eventChannel)
  tuple := <-dataChannel
  data := tuple.(RepeatableOrder)

  expectedPeriod, _ := time.Parse("2006-01-02", "2014-01-01")
  if data.Siret == "01234567891011" &&
    data.Periode == expectedPeriod &&
    *data.RandomOrder == 0.1234 {
    t.Log("Test RepeatableOrder: lecture ok")
  } else {
    t.Error("Erreur RepeatableOrder: données différentes du fichier de test")
  }

  tuple = <-dataChannel

  if tuple != nil {
    t.Error("Erreur RepeatableOrder: 2° ligne vide: le channel devrait être vide")
  } else {
    t.Log("Test RepeatableOrder: test deuxième ligne ok")
  }

  batch = engine.AdminBatch{
    Files: engine.BatchFiles{
      "apconso": []string{"testData/nonexistant.csv"},
    },
  }

  dataChannel, eventChannel = Parser(batch)
  engine.DiscardEvents(eventChannel)
  tuple = <-dataChannel

  if tuple != nil {
    t.Error("Erreur RepeatableOrder: fichier absent: le channel devrait être vide")
  } else {
    t.Log("Test RepeatableOrder: fichier absent ok")
  }
}
