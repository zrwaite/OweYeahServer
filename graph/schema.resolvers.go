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

func (r *mutationResolver) CreateUser(ctx context.Context, input model.UserInput) (*model.UserAuthResult, error) {
	return database.CreateUser(ctx, input), nil
}

func (r *mutationResolver) Login(ctx context.Context, input model.UserInput) (*model.UserAuthResult, error) {
	return database.Login(ctx, input), nil
}

func (r *mutationResolver) DeleteUser(ctx context.Context, username string) (*model.Result, error) {
	return database.DeleteUser(ctx, username), nil
}

func (r *mutationResolver) CreateInvoice(ctx context.Context, input model.InvoiceOrPaymentInput) (*model.InvoiceOrPaymentResult, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *mutationResolver) CreatePayment(ctx context.Context, input model.InvoiceOrPaymentInput) (*model.InvoiceOrPaymentResult, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *mutationResolver) AddContact(ctx context.Context, username string, contactUsername string) (*model.Result, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *mutationResolver) RemoveContact(ctx context.Context, username string, contactUsername string) (*model.Result, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *queryResolver) User(ctx context.Context, username string) (*model.UserResult, error) {
	return database.User(ctx, username), nil
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
