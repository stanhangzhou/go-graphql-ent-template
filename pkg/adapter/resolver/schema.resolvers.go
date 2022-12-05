package resolver

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"context"

	"gitlab.com/trustify/core/ent"
	"gitlab.com/trustify/core/ent/schema/ulid"
	"gitlab.com/trustify/core/graph/generated"
)

func (r *queryResolver) Node(ctx context.Context, id ulid.ID) (ent.Noder, error) {
	return r.client.Noder(ctx, id, ent.WithNodeType(ent.IDToType))
}

// Mutation returns generated.MutationResolver implementation.
func (r *Resolver) Mutation() generated.MutationResolver { return &mutationResolver{r} }

// Query returns generated.QueryResolver implementation.
func (r *Resolver) Query() generated.QueryResolver { return &queryResolver{r} }

type mutationResolver struct{ *Resolver }
type queryResolver struct{ *Resolver }
