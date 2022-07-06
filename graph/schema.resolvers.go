package graph

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"context"
	"fmt"

	"github.com/zrwaite/OweMate/auth"
	"github.com/zrwaite/OweMate/auth/tokens"
	"github.com/zrwaite/OweMate/database"
	"github.com/zrwaite/OweMate/graph/generated"
	"github.com/zrwaite/OweMate/graph/model"
)

func (r *mutationResolver) CreateUser(ctx context.Context, input model.UserInput) (*model.UserAuthResult, error) {
	errors, token := database.CreateUser(&input)
	if len(errors) > 0 {
		return &model.UserAuthResult{
			Success: false,
			Errors:  errors,
		}, nil
	} else {
		return &model.UserAuthResult{
			Success: true,
			Token:   token,
		}, nil
	}
}

func (r *mutationResolver) Login(ctx context.Context, input model.UserInput) (*model.UserAuthResult, error) {
	user, status := database.GetUser(input.Username)
	errors := []string{}
	if status == 404 {
		errors = append(errors, "User not found")
	} else if status == 400 {
		errors = append(errors, "Something went wrong")
	} else {
		if !auth.CheckPasswordHash(input.Password, user.Hash) {
			errors = append(errors, "Incorrect password")
		} else {
			token, tokenSuccess := tokens.EncodeToken(input.Username)
			if tokenSuccess {
				return &model.UserAuthResult{
					Success: true,
					Token:   token,
					User:    user,
				}, nil
			} else {
				errors = append(errors, "Failed to create token")
			}
		}
	}
	return &model.UserAuthResult{
		Success: false,
		Errors:  errors,
	}, nil
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
