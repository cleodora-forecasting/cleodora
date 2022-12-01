// Code generated by github.com/99designs/gqlgen, DO NOT EDIT.

package model

import (
	"fmt"
	"io"
	"strconv"
	"time"
)

type Forecast struct {
	ID          string     `json:"id"`
	Title       string     `json:"title"`
	Description string     `json:"description"`
	Created     time.Time  `json:"created"`
	Resolves    time.Time  `json:"resolves"`
	Closes      *time.Time `json:"closes"`
	Resolution  Resolution `json:"resolution"`
}

type NewForecast struct {
	Title       string     `json:"title"`
	Description string     `json:"description"`
	Resolves    time.Time  `json:"resolves"`
	Closes      *time.Time `json:"closes"`
}

type Resolution string

const (
	ResolutionTrue          Resolution = "TRUE"
	ResolutionFalse         Resolution = "FALSE"
	ResolutionNotApplicable Resolution = "NOT_APPLICABLE"
	ResolutionUnresolved    Resolution = "UNRESOLVED"
)

var AllResolution = []Resolution{
	ResolutionTrue,
	ResolutionFalse,
	ResolutionNotApplicable,
	ResolutionUnresolved,
}

func (e Resolution) IsValid() bool {
	switch e {
	case ResolutionTrue, ResolutionFalse, ResolutionNotApplicable, ResolutionUnresolved:
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
