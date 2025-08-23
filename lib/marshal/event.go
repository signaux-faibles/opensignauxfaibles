package marshal

import (
	"time"
)

var gitCommit string

// Event est un objet de journal
type Event struct {
	StartDate  time.Time `json:"start_date"`
	CommitHash string    `json:"commit_hash,omitempty"`
	Report     *Report   `json:"report"`
	Parser     string    `json:"parser"`
	ReportType string    `json:"report_type"`
}

// SetGitCommit spécifie la valeur à stocker dans le CommitHash de chaque événement.
func SetGitCommit(hash string) {
	gitCommit = hash
}

// CreateEvent initialise un évènement avec les valeurs par défaut.
func CreateEvent() Event {
	return Event{
		Parser:     "unknown",
		CommitHash: gitCommit,
	}
}

// CreateReportEvent initialise un évènement contenant un rapport de parsing.
func CreateReportEvent(parserType string, report Report) Event {
	event := CreateEvent()
	event.Parser = parserType
	event.Report = &report
	return event
}
