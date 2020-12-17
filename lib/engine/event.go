// Ce fichier est responsable de collecter les messages et de les ajouter
// dans la collection Journal.

package engine

import (
	"encoding/json"
	"log"
	"sync"
	"time"

	"github.com/globalsign/mgo/bson"
	"github.com/signaux-faibles/opensignauxfaibles/lib/marshal"
)

type messageChannel chan marshal.Event

var mainMessageChannel = messageDispatch() // canal dans lequel on va émettre tous les messages

var relaying sync.WaitGroup // permet de savoir quand les messages ont fini d'être transmis

var messageClientChannels = []messageChannel{}

// AddClientChannel enregistre un nouveau client
var AddClientChannel = make(chan messageChannel)

// MessageSocketAddClient surveille l'ajout de nouveaux clients pour les enregistrer dans la liste des clients
func MessageSocketAddClient() {
	for clientChannel := range AddClientChannel {
		messageClientChannels = append(messageClientChannels, clientChannel)
	}
}

// Transmet les messages collectés vers les clients et l'enregistre dans la bdd
func messageDispatch() chan marshal.Event {
	relaying.Add(1)
	channel := make(messageChannel)
	go func() {
		defer relaying.Done()
		for event := range channel {
			err := Db.DBStatus.C("Journal").Insert(event)
			if err != nil {
				log.Print("Erreur critique d'insertion dans la base de données: " + err.Error())
				log.Print(json.Marshal(event))
			}
			for _, clientChannel := range messageClientChannels {
				clientChannel <- event
			}
		}
	}()
	return channel
}

// RelayEvents transmet les événements qui surviennent pendant le parsing d'un
// fichiers de données et retourne le rapport final du parsing de ce fichier.
func RelayEvents(eventChannel chan marshal.Event, reportType string, startDate time.Time) (lastReport string) {
	for e := range eventChannel {
		if reportContainer, ok := e.Comment.(bson.M); ok {
			if strReport, ok := reportContainer["summary"].(string); ok {
				lastReport = strReport
			}
		}
		e.ReportType = reportType
		e.StartDate = startDate
		mainMessageChannel <- e
	}
	return lastReport
}

// LogOperationEvent rapporte la fin d'une opération effectuée par sfdata.
func LogOperationEvent(reportType string, startDate time.Time) {
	event := marshal.CreateEvent()
	event.StartDate = startDate
	event.ReportType = reportType
	mainMessageChannel <- event
}

// FlushEventQueue finalise l'insertion des événements dans Journal.
func FlushEventQueue() {
	close(mainMessageChannel)
	relaying.Wait()
}
