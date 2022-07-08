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

func (r *connectionResolver) Contact(ctx context.Context, obj *model.Connection) (*model.User, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *invoiceOrPaymentResolver) Connection(ctx context.Context, obj *model.InvoiceOrPayment) (*model.Connection, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *mutationResolver) CreateUser(ctx context.Context, input model.UserInput) (*model.UserAuthResult, error) {
	return database.CreateUser(ctx, input), nil
}

func (r *mutationResolver) Login(ctx context.Context, input model.UserInput) (*model.UserAuthResult, error) {
	return database.Login(ctx, input), nil
}

func (r *mutationResolver) DeleteUser(ctx context.Context, username string) (*model.Result, error) {
	return database.DeleteUser(ctx, username), nil
}

func (r *mutationResolver) CreateConnection(ctx context.Context, username1 string, username2 string) (*model.ConnectionResult, error) {
	return database.CreateConnection(ctx, username1, username2), nil
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

func (r *userResolver) Invoices(ctx context.Context, obj *model.User) ([]*model.InvoiceOrPayment, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *userResolver) Payments(ctx context.Context, obj *model.User) ([]*model.InvoiceOrPayment, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *userResolver) Connections(ctx context.Context, obj *model.User) ([]*model.Connection, error) {
	panic(fmt.Errorf("not implemented"))
}

// Connection returns generated.ConnectionResolver implementation.
func (r *Resolver) Connection() generated.ConnectionResolver { return &connectionResolver{r} }

// InvoiceOrPayment returns generated.InvoiceOrPaymentResolver implementation.
func (r *Resolver) InvoiceOrPayment() generated.InvoiceOrPaymentResolver {
	return &invoiceOrPaymentResolver{r}
}

// Mutation returns generated.MutationResolver implementation.
func (r *Resolver) Mutation() generated.MutationResolver { return &mutationResolver{r} }

// Query returns generated.QueryResolver implementation.
func (r *Resolver) Query() generated.QueryResolver { return &queryResolver{r} }

// User returns generated.UserResolver implementation.
func (r *Resolver) User() generated.UserResolver { return &userResolver{r} }

type connectionResolver struct{ *Resolver }
type invoiceOrPaymentResolver struct{ *Resolver }
type mutationResolver struct{ *Resolver }
type queryResolver struct{ *Resolver }
type userResolver struct{ *Resolver }
