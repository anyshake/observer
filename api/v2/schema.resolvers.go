package v2

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.
// Code generated by github.com/99designs/gqlgen version v0.17.49

import (
	"context"
	"fmt"
)

// Test is the resolver for the test field.
func (r *queryResolver) Test(ctx context.Context) (*bool, error) {
	panic(fmt.Errorf("not implemented: Test - test"))
}

// Query returns QueryResolver implementation.
func (r *Resolver) Query() QueryResolver { return &queryResolver{r} }

type queryResolver struct{ *Resolver }
