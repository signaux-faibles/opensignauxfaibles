package engine

import (
	"errors"
	"strconv"
	"testing"

	"github.com/signaux-faibles/gournal"
	"github.com/signaux-faibles/opensignauxfaibles/dbmongo/lib/base"
)

func TestShouldBreak(t *testing.T) {
	tracker := gournal.NewTracker(
		map[string]string{},
		TrackerReports,
	)
	filterError := base.NewFilterError(errors.New("filterError"))
	fatalError := base.NewCriticError(errors.New("fatalError"), "fatal")
	errorError := base.NewCriticError(errors.New("errorError"), "error")
	noError := []error{}

	testCases := []struct {
		errors         map[int][]error
		expectedReport bool
	}{
		{map[int][]error{1: []error{filterError}}, false},
		{map[int][]error{1: []error{fatalError}}, true},
		{map[int][]error{1: []error{errorError}}, true},
		{map[int][]error{1: noError, 2: []error{filterError}}, false},
		{map[int][]error{1: noError, 2: []error{fatalError}}, true},
		{map[int][]error{1: noError, 2: []error{errorError}}, true},
		{map[int][]error{1: []error{filterError, fatalError}}, true},
		{map[int][]error{1: []error{errorError, filterError}}, true},
	}

	for ind, tc := range testCases {
		tracker.Errors = tc.errors
		actual := ShouldBreak(tracker, 0)
		if actual != tc.expectedReport {
			t.Error("Test case " + strconv.Itoa(ind) + " failed")
		}
	}
}
