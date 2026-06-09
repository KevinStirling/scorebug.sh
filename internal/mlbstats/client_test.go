package mlbstats

import (
	"testing"
	"time"
)

// dateOperand represents a date operand
type dateOperand int

// date operand constants
const (
	sameDay dateOperand = iota
	dayPrior
)

func Test_determineTime(t *testing.T) {
	type test struct {
		currentDate  time.Time
		expectedDate time.Time
		comparison   dateOperand
	}

	year, month, day := time.Now().Date()
	tt := map[string]test{
		"current date - equal": {
			currentDate:  time.Date(year, month, day, 12, 0, 0, 0, time.Now().Location()),
			expectedDate: time.Date(year, month, day, 12, 0, 0, 0, time.Now().Location()),
			comparison:   sameDay,
		},
		"current date - pre 5am": {
			currentDate:  time.Date(year, month, day, 5, 0, 0, 0, time.Now().Location()),
			expectedDate: time.Now().AddDate(0, 0, -1),
			comparison:   dayPrior,
		},
		"current date - post 5am": {
			currentDate:  time.Date(year, month, day, 6, 0, 0, 0, time.Now().Location()),
			expectedDate: time.Now(),
			comparison:   sameDay,
		},
	}

	for k, v := range tt {
		t.Run(k, func(t *testing.T) {
			currentDate := determineTime(v.currentDate)

			switch v.comparison {
			case dayPrior:
				if !currentDate.After(v.expectedDate) {
					t.Errorf("expected %v to be less than %v", currentDate, v.expectedDate)
				}
			case sameDay:
				if currentDate.Day() != v.expectedDate.Day() {
					t.Errorf("expected %v to equal %v", currentDate, v.expectedDate)
				}
			}
		})
	}
}
