package marshal

import (
	"time"
)

var gitCommit string

// Priority test
type Priority string

// Code test
type Code string

// Event est un objet de journal
type Event struct {
	Date       time.Time `json:"date"`
	StartDate  time.Time `json:"startDate"`
	CommitHash string    `json:"commitHash,omitempty"`
	Report     *Report   `json:"report"`
	Priority   Priority  `json:"priority"`
	Code       Code      `json:"parserCode"`
	ReportType string    `json:"report_type"`
}

// SetGitCommit spécifie la valeur à stocker dans le CommitHash de chaque événement.
func SetGitCommit(hash string) {
	gitCommit = hash
}

// CreateEvent initialise un évènement avec les valeurs par défaut.
func CreateEvent() Event {
	return Event{
		Date:       time.Now(),
		Priority:   Priority("info"),
		Code:       Code("unknown"),
		CommitHash: gitCommit,
	}
}

// CreateReportEvent initialise un évènement contenant un rapport de parsing.
func CreateReportEvent(fileType string, report Report) Event {
	event := CreateEvent()
	event.Code = Code(fileType)
	event.Report = &report
	return event
}
