package graphql

import (
	"entgo.io/contrib/entgql"
	"github.com/99designs/gqlgen/graphql/handler"
	"gitlab.com/trustify/core/ent"
	"gitlab.com/trustify/core/pkg/adapter/controller"
	"gitlab.com/trustify/core/pkg/adapter/resolver"
)

// NewServer generates graphql server
func NewServer(client *ent.Client, controller controller.Controller) *handler.Server {
	srv := handler.NewDefaultServer(resolver.NewSchema(client, controller))
	srv.Use(entgql.Transactioner{TxOpener: client})

	return srv
}
