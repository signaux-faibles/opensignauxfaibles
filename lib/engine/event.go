package engine

import (
	"encoding/json"
	"log"
	"sync"
	"time"

	"github.com/globalsign/mgo/bson"
	"github.com/signaux-faibles/opensignauxfaibles/lib/base"
	"github.com/signaux-faibles/opensignauxfaibles/lib/marshal"
)

// SocketMessage permet la diffusion d'information vers tous les clients
type SocketMessage struct {
	JournalEvent marshal.Event     `json:"journalEvent" bson:"journalEvent"`
	Batches      []base.AdminBatch `json:"batches,omitempty" bson:"batches,omitempty"`
	Features     []string          `json:"features,omitempty" bson:"features,omitempty"`
}

var relaying sync.WaitGroup

type messageChannel chan SocketMessage

var messageClientChannels = []messageChannel{}

// MainMessageChannel permet d'envoyer un SocketMessage
var MainMessageChannel = messageDispatch()

// AddClientChannel enregistre un nouveau client
var AddClientChannel = make(chan messageChannel)

// MessageSocketAddClient surveille l'ajout de nouveaux clients pour les enregistrer dans la liste des clients
func MessageSocketAddClient() {
	for clientChannel := range AddClientChannel {
		messageClientChannels = append(messageClientChannels, clientChannel)
	}
}

// journal dispatch un event vers les clients et l'enregistre dans la bdd
func messageDispatch() chan SocketMessage {
	relaying.Add(1)
	channel := make(messageChannel)
	go func() {
		defer relaying.Done()
		for event := range channel {
			err := Db.DBStatus.C("Journal").Insert(event.JournalEvent)
			if err != nil {
				log.Print("Erreur critique d'insertion dans la base de données: " + err.Error())
				log.Print(json.Marshal(event.JournalEvent))
			}
			for _, clientChannel := range messageClientChannels {
				clientChannel <- event
			}
		}
	}()
	return channel
}

// RelayEvents transmet les messages
func RelayEvents(eventChannel chan marshal.Event, reportType string, startDate time.Time) (lastReport string) {
	if eventChannel == nil {
		return
	}
	for e := range eventChannel {
		if reportContainer, ok := e.Comment.(bson.M); ok {
			if strReport, ok := reportContainer["summary"].(string); ok {
				lastReport = strReport
			}
		}
		e.ReportType = reportType
		e.StartDate = startDate
		MainMessageChannel <- SocketMessage{
			JournalEvent: e,
		}
	}
	return lastReport
}

// LogOperationEvent rapporte la fin d'une opération effectuée par sfdata.
func LogOperationEvent(reportType string, startDate time.Time) {
	event := marshal.CreateEvent()
	event.StartDate = startDate
	event.ReportType = reportType
	MainMessageChannel <- SocketMessage{
		JournalEvent: event,
	}
}

// FlushEventQueue finalise l'insertion des événements dans Journal.
func FlushEventQueue() {
	close(MainMessageChannel)
	relaying.Wait()
}
