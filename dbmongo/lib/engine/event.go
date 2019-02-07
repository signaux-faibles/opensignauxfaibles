package engine

import (
	"dbmongo/lib/files"
	"encoding/json"
	"errors"
	"log"
)

// PlugEvents connecte deux chan Event
func PlugEvents(source chan Event, dest chan Event) {
	for s := range source {
		dest <- s
	}
}

// SocketMessage permet la diffusion d'information vers tous les clients
type SocketMessage struct {
	JournalEvent Event               `json:"journalEvent" bson:"journalEvent"`
	Batches      []AdminBatch        `json:"batches,omitempty" bson:"batches,omitempty"`
	Types        []Type              `json:"types,omitempty" bson:"types,omitempty"`
	Features     []string            `json:"features,omitempty" bson:"features,omitempty"`
	Files        []files.FileSummary `json:"files,omitempty" bson:"files,omitempty"`
	Channel      chan SocketMessage
}

// MarshalJSON fournit un objet serialisable
func (message SocketMessage) MarshalJSON() ([]byte, error) {
	var tmp SocketMessage
	tmp.JournalEvent = message.JournalEvent
	tmp.Batches = message.Batches
	tmp.Types = message.Types
	tmp.Features = message.Features
	tmp.Files = message.Files

	return json.Marshal(tmp)
}

// Send transmet le message
func (message SocketMessage) Send() error {
	if message.Channel == nil {
		return errors.New("Aucun channel défini")
	}
	message.Channel <- message
	return nil
}

// PlugTuples connecte deux chan Tuple
func PlugTuples(source chan Tuple, dest chan Tuple) {
	for s := range source {
		dest <- s
	}
}

// GetEventsFromDB retourne les n derniers enregistrements correspondant à la requête
func GetEventsFromDB(query interface{}, n int) ([]Event, error) {
	var logs []Event
	err := Db.DB.C("Journal").Find(query).Sort("-date").Limit(n).All(&logs)
	return logs, err
}

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
	channel := make(messageChannel)
	go func() {
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
func RelayEvents(eventChannel chan Event) {
	if eventChannel == nil {
		return
	}
	for e := range eventChannel {
		MainMessageChannel <- SocketMessage{
			JournalEvent: e,
		}
	}
}
