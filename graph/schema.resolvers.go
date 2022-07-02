package graph

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"context"
	"fmt"

	"github.com/zrwaite/OweMate/database"
	"github.com/zrwaite/OweMate/graph/generated"
	"github.com/zrwaite/OweMate/graph/model"
)

func (r *mutationResolver) CreateUser(ctx context.Context, input model.UserInput) (*model.UserResult, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *mutationResolver) Login(ctx context.Context, input model.UserInput) (*model.UserResult, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *mutationResolver) UpdateUser(ctx context.Context, username string) (*model.UserResult, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *mutationResolver) DeleteUser(ctx context.Context, username string) (*model.Result, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *queryResolver) User(ctx context.Context, username string) (*model.UserResult, error) {
	user, status := database.GetUser(username)
	var userResult *model.UserResult
	if status == 200 {
		userResult = &model.UserResult{
			User:    user,
			Success: true,
		}
	} else if status == 404 {
		userResult = &model.UserResult{
			Success: false,
			Errors:  []string{"User not found"},
		}
	} else if status == 400 {
		userResult = &model.UserResult{
			Success: false,
			Errors:  []string{"Something went wrong"},
		}
	}

	return userResult, nil
}

// Mutation returns generated.MutationResolver implementation.
func (r *Resolver) Mutation() generated.MutationResolver { return &mutationResolver{r} }

// Query returns generated.QueryResolver implementation.
func (r *Resolver) Query() generated.QueryResolver { return &queryResolver{r} }

type mutationResolver struct{ *Resolver }
type queryResolver struct{ *Resolver }
