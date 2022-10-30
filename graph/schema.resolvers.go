package graph

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"context"
	"fmt"
	"math/rand"
	"time"

	"github.com/cleodora-forecasting/cleodora/graph/generated"
	"github.com/cleodora-forecasting/cleodora/graph/model"
)

// CreateTodo is the resolver for the createTodo field.
func (r *mutationResolver) CreateTodo(ctx context.Context, input model.NewTodo) (*model.Todo, error) {
	todo := &model.Todo{
		Text:   input.Text,
		ID:     fmt.Sprintf("T%d", rand.Int()),
		User:   &model.User{ID: input.UserID, Name: "user " + input.UserID},
		UserID: input.UserID,
	}
	r.todos = append(r.todos, todo)
	return todo, nil
}

// CreateForecast is the resolver for the createForecast field.
func (r *mutationResolver) CreateForecast(ctx context.Context, input model.NewForecast) (*model.Forecast, error) {
	forecast := &model.Forecast{
		ID:          fmt.Sprintf("T%d", rand.Int()),
		Summary:     input.Summary,
		Description: input.Description,
		Closes:      input.Closes,
		Created:     time.Now(),
	}
	r.forecasts = append(r.forecasts, forecast)
	return forecast, nil
}

// Todos is the resolver for the todos field.
func (r *queryResolver) Todos(ctx context.Context) ([]*model.Todo, error) {
	return r.todos, nil
}

// Forecasts is the resolver for the forecasts field.
func (r *queryResolver) Forecasts(ctx context.Context) ([]*model.Forecast, error) {
	return r.forecasts, nil
}

// User is the resolver for the user field.
func (r *todoResolver) User(ctx context.Context, obj *model.Todo) (*model.User, error) {
	return &model.User{ID: obj.UserID, Name: "user " + obj.UserID}, nil
}

// Mutation returns generated.MutationResolver implementation.
func (r *Resolver) Mutation() generated.MutationResolver { return &mutationResolver{r} }

// Query returns generated.QueryResolver implementation.
func (r *Resolver) Query() generated.QueryResolver { return &queryResolver{r} }

// Todo returns generated.TodoResolver implementation.
func (r *Resolver) Todo() generated.TodoResolver { return &todoResolver{r} }

type mutationResolver struct{ *Resolver }
type queryResolver struct{ *Resolver }
type todoResolver struct{ *Resolver }
