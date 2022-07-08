package graph

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"context"

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

func (r *mutationResolver) CreateConnection(ctx context.Context, username string, contactUsername string) (*model.ConnectionResult, error) {
	return database.CreateConnection(ctx, username, contactUsername), nil
}

func (r *mutationResolver) CreateInvoice(ctx context.Context, input model.InvoiceOrPaymentInput) (*model.InvoiceResult, error) {
	return database.CreateInvoice(ctx, input), nil
}

func (r *mutationResolver) CreatePayment(ctx context.Context, input model.InvoiceOrPaymentInput) (*model.PaymentResult, error) {
	return database.CreatePayment(ctx, input), nil
}

func (r *queryResolver) User(ctx context.Context, username string) (*model.UserResult, error) {
	return database.User(ctx, username), nil
}

func (r *queryResolver) GetFilteredUsers(ctx context.Context, partialUsername string) (*model.UsersResult, error) {
	return database.GetFilteredUsers(ctx, partialUsername), nil
}

// Mutation returns generated.MutationResolver implementation.
func (r *Resolver) Mutation() generated.MutationResolver { return &mutationResolver{r} }

// Query returns generated.QueryResolver implementation.
func (r *Resolver) Query() generated.QueryResolver { return &queryResolver{r} }

type mutationResolver struct{ *Resolver }
type queryResolver struct{ *Resolver }
