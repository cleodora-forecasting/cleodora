package dbmodel

import (
	"time"

	"gorm.io/gorm"
)

type Forecast struct {
	gorm.Model
	Title       string
	Description string
	Created     time.Time
	Resolves    time.Time
	Closes      *time.Time
	Resolution  Resolution
	Estimates   []Estimate
}

type Estimate struct {
	gorm.Model
	ForecastID    uint
	Created       time.Time
	Reason        string
	Probabilities []Probability
}

type Probability struct {
	gorm.Model
	EstimateID uint
	Value      int
	Outcome    Outcome
	OutcomeID  uint
}

type Outcome struct {
	gorm.Model
	Probabilities []Probability
	Text          string
	Correct       bool
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
