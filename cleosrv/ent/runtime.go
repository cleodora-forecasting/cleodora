// Code generated by ent, DO NOT EDIT.

package ent

import (
	"time"

	"github.com/cleodora-forecasting/cleodora/cleosrv/ent/estimate"
	"github.com/cleodora-forecasting/cleodora/cleosrv/ent/forecast"
	"github.com/cleodora-forecasting/cleodora/cleosrv/ent/outcome"
	"github.com/cleodora-forecasting/cleodora/cleosrv/ent/probability"
	"github.com/cleodora-forecasting/cleodora/cleosrv/ent/schema"
)

// The init function reads all schema descriptors with runtime code
// (default values, validators, hooks and policies) and stitches it
// to their package variables.
func init() {
	estimateFields := schema.Estimate{}.Fields()
	_ = estimateFields
	// estimateDescReason is the schema descriptor for reason field.
	estimateDescReason := estimateFields[0].Descriptor()
	// estimate.DefaultReason holds the default value on creation for the reason field.
	estimate.DefaultReason = estimateDescReason.Default.(string)
	// estimateDescCreated is the schema descriptor for created field.
	estimateDescCreated := estimateFields[1].Descriptor()
	// estimate.DefaultCreated holds the default value on creation for the created field.
	estimate.DefaultCreated = estimateDescCreated.Default.(func() time.Time)
	forecastFields := schema.Forecast{}.Fields()
	_ = forecastFields
	// forecastDescTitle is the schema descriptor for title field.
	forecastDescTitle := forecastFields[0].Descriptor()
	// forecast.TitleValidator is a validator for the "title" field. It is called by the builders before save.
	forecast.TitleValidator = forecastDescTitle.Validators[0].(func(string) error)
	// forecastDescDescription is the schema descriptor for description field.
	forecastDescDescription := forecastFields[1].Descriptor()
	// forecast.DefaultDescription holds the default value on creation for the description field.
	forecast.DefaultDescription = forecastDescDescription.Default.(string)
	// forecastDescCreated is the schema descriptor for created field.
	forecastDescCreated := forecastFields[2].Descriptor()
	// forecast.DefaultCreated holds the default value on creation for the created field.
	forecast.DefaultCreated = forecastDescCreated.Default.(func() time.Time)
	outcomeFields := schema.Outcome{}.Fields()
	_ = outcomeFields
	// outcomeDescText is the schema descriptor for text field.
	outcomeDescText := outcomeFields[0].Descriptor()
	// outcome.TextValidator is a validator for the "text" field. It is called by the builders before save.
	outcome.TextValidator = outcomeDescText.Validators[0].(func(string) error)
	// outcomeDescCorrect is the schema descriptor for correct field.
	outcomeDescCorrect := outcomeFields[1].Descriptor()
	// outcome.DefaultCorrect holds the default value on creation for the correct field.
	outcome.DefaultCorrect = outcomeDescCorrect.Default.(bool)
	probabilityFields := schema.Probability{}.Fields()
	_ = probabilityFields
	// probabilityDescValue is the schema descriptor for value field.
	probabilityDescValue := probabilityFields[0].Descriptor()
	// probability.ValueValidator is a validator for the "value" field. It is called by the builders before save.
	probability.ValueValidator = probabilityDescValue.Validators[0].(func(int) error)
}