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

func (r *mutationResolver) DeleteUser(ctx context.Context, username string) (*model.Result, error) {
	errors := database.DeleteUser(username)
	if len(errors) > 0 {
		return &model.Result{
			Success: false,
			Errors:  errors,
		}, nil
	} else {
		return &model.Result{
			Success: true,
		}, nil
	}
}

func (r *mutationResolver) CreateInvoice(ctx context.Context, input model.InvoiceInput) (*model.InvoiceResult, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *mutationResolver) CreatePayment(ctx context.Context, input model.PaymentInput) (*model.PaymentResult, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *mutationResolver) AddContact(ctx context.Context, username string, contactUsername string) (*model.Result, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *mutationResolver) RemoveContact(ctx context.Context, username string, contactUsername string) (*model.Result, error) {
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

func (r *queryResolver) GetFilteredUsers(ctx context.Context, partialUsername string) (*model.UsersResult, error) {
	panic(fmt.Errorf("not implemented"))
}

// Mutation returns generated.MutationResolver implementation.
func (r *Resolver) Mutation() generated.MutationResolver { return &mutationResolver{r} }

// Query returns generated.QueryResolver implementation.
func (r *Resolver) Query() generated.QueryResolver { return &queryResolver{r} }

type mutationResolver struct{ *Resolver }
type queryResolver struct{ *Resolver }

// !!! WARNING !!!
// The code below was going to be deleted when updating resolvers. It has been copied here so you have
// one last chance to move it out of harms way if you want. There are two reasons this happens:
//  - When renaming or deleting a resolver the old code will be put in here. You can safely delete
//    it when you're done.
//  - You have helper methods in this file. Move them out to keep these resolver files clean.
func (r *mutationResolver) UpdateUser(ctx context.Context, username string) (*model.UserResult, error) {
	panic(fmt.Errorf("not implemented"))
}
