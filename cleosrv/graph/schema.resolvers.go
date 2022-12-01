package graph

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"context"
	"fmt"
	"math/rand"
	"time"

	"github.com/cleodora-forecasting/cleodora/cleosrv/graph/generated"
	"github.com/cleodora-forecasting/cleodora/cleosrv/graph/model"
)

// CreateForecast is the resolver for the createForecast field.
func (r *mutationResolver) CreateForecast(ctx context.Context, input model.NewForecast) (*model.Forecast, error) {
	forecast := &model.Forecast{
		ID:          fmt.Sprintf("T%d", rand.Int()),
		Title:       input.Title,
		Description: input.Description,
		Created:     time.Now(),
		Resolves:    input.Resolves,
		Closes:      input.Closes,
		Resolution:  model.ResolutionUnresolved,
	}

	r.forecasts = append(r.forecasts, forecast)
	return forecast, nil
}

// Forecasts is the resolver for the forecasts field.
func (r *queryResolver) Forecasts(ctx context.Context) ([]*model.Forecast, error) {
	return r.forecasts, nil
}

// Mutation returns generated.MutationResolver implementation.
func (r *Resolver) Mutation() generated.MutationResolver { return &mutationResolver{r} }

// Query returns generated.QueryResolver implementation.
func (r *Resolver) Query() generated.QueryResolver { return &queryResolver{r} }

type mutationResolver struct{ *Resolver }
type queryResolver struct{ *Resolver }
