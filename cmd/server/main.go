package main

import (
	"log"

	_ "github.com/lib/pq"
	"gitlab.com/trustify/core/config"
	"gitlab.com/trustify/core/ent"
	"gitlab.com/trustify/core/pkg/adapter/controller"
	"gitlab.com/trustify/core/pkg/infrastructure/datastore"
	"gitlab.com/trustify/core/pkg/infrastructure/graphql"
	"gitlab.com/trustify/core/pkg/infrastructure/router"
	"gitlab.com/trustify/core/pkg/registry"
)

func main() {
	config.ReadConfig(config.ReadConfigOption{})

	client := newDBClient()
	ctrl := newController(client)

	srv := graphql.NewServer(client, ctrl)
	e := router.New(srv)

	e.Logger.Fatal(e.Start(":" + config.C.HttpServer.Port))
}

func newDBClient() *ent.Client {
	client, err := datastore.NewClient()
	if err != nil {
		log.Fatalf("failed opening postgres client: %v", err)
	}

	return client
}

func newController(client *ent.Client) controller.Controller {
	r := registry.New(client)
	return r.NewController()
}
