// Code generated by github.com/99designs/gqlgen, DO NOT EDIT.

package model

import (
	"fmt"
	"io"
	"strconv"
	"time"
)

// A list of probabilities (one for each outcome) together with a timestamp and
// an explanation why you made this estimate. Every time you change your mind
// about a forecast you will create a new Estimate.
// All probabilities always add up to 100.
type Estimate struct {
	ID            string         `json:"id"`
	Created       time.Time      `json:"created"`
	Reason        string         `json:"reason"`
	Probabilities []*Probability `json:"probabilities"`
}

// A prediction about the future.
type Forecast struct {
	ID          string    `json:"id"`
	Title       string    `json:"title"`
	Description string    `json:"description"`
	Created     time.Time `json:"created"`
	// The point in time at which you predict you will be able to resolve whether
	// how the forecast resolved.
	Resolves time.Time `json:"resolves"`
	// The point in time at which you no longer want to update your probability
	// estimates for the forecast. In most cases you won't need this. One example
	// where you might is when you want to predict the outcome of an exam. You may
	// want to set 'closes' to the time right before the exam starts, even though
	// 'resolves' is several weeks later (when the exam results are published). This
	// way your prediction history will only reflect your estimations before you
	// took the exam, which is something you may want (or not, in which case you
	// could simply not set 'closes').
	Closes     *time.Time  `json:"closes"`
	Resolution Resolution  `json:"resolution"`
	Outcomes   []*Outcome  `json:"outcomes"`
	Estimates  []*Estimate `json:"estimates"`
}

type NewForecast struct {
	Title       string     `json:"title"`
	Description string     `json:"description"`
	Resolves    time.Time  `json:"resolves"`
	Closes      *time.Time `json:"closes"`
}

// The possible results of a forecast. In the simplest case you will only have
// two outcomes: Yes and No.
type Outcome struct {
	ID      string `json:"id"`
	Text    string `json:"text"`
	Correct bool   `json:"correct"`
}

// A number between 0 and 100 tied to a specific Outcome. It is always part of
// an Estimate.
type Probability struct {
	ID      string   `json:"id"`
	Value   int      `json:"value"`
	Outcome *Outcome `json:"outcome"`
}

type Resolution string

const (
	ResolutionResolved      Resolution = "RESOLVED"
	ResolutionNotApplicable Resolution = "NOT_APPLICABLE"
	ResolutionUnresolved    Resolution = "UNRESOLVED"
)

var AllResolution = []Resolution{
	ResolutionResolved,
	ResolutionNotApplicable,
	ResolutionUnresolved,
}

func (e Resolution) IsValid() bool {
	switch e {
	case ResolutionResolved, ResolutionNotApplicable, ResolutionUnresolved:
		return true
	}
	return false
}

func (e Resolution) String() string {
	return string(e)
}

func (e *Resolution) UnmarshalGQL(v interface{}) error {
	str, ok := v.(string)
	if !ok {
		return fmt.Errorf("enums must be strings")
	}

	*e = Resolution(str)
	if !e.IsValid() {
		return fmt.Errorf("%s is not a valid Resolution", str)
	}
	return nil
}

func (e Resolution) MarshalGQL(w io.Writer) {
	fmt.Fprint(w, strconv.Quote(e.String()))
}
