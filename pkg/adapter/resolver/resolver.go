package resolver

import (
	"github.com/99designs/gqlgen/graphql"
	"gitlab.com/trustify/core/ent"
	"gitlab.com/trustify/core/graph/generated"
	"gitlab.com/trustify/core/pkg/adapter/controller"
	"gitlab.com/trustify/core/pkg/adapter/directives"
)

// This file will not be regenerated automatically.
//
// It serves as dependency injection for your app, add any dependencies you require here.

// Resolver is a context struct
type Resolver struct {
	client     *ent.Client
	controller controller.Controller
}

// NewSchema creates NewExecutableSchema
func NewSchema(client *ent.Client, controller controller.Controller) graphql.ExecutableSchema {
	c := generated.Config{
		Resolvers: &Resolver{
			client:     client,
			controller: controller,
		},
	}
	c.Directives.Binding = directives.Binding

	return generated.NewExecutableSchema(c)
}
