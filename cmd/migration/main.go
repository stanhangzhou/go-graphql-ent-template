package main

import (
	"context"
	"fmt"
	"log"

	"entgo.io/ent/dialect"
	_ "github.com/lib/pq"
	"gitlab.com/trustify/core/config"
	"gitlab.com/trustify/core/ent"
	"gitlab.com/trustify/core/ent/migrate"
)

func main() {
	config.ReadConfig(config.ReadConfigOption{})

	client, err := newClient()
	if err != nil {
		log.Fatalf("error opening postgres connection: %v", err)
	}
	defer client.Close()

	createDBSchema(client)
}

func newDSN() string {
	dsn := fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		config.C.Database.Host,
		config.C.Database.Port,
		config.C.Database.User,
		config.C.Database.Password,
		config.C.Database.Name,
		config.C.Database.SSL,
	)
	return dsn
}

func newClient() (*ent.Client, error) {
	var entOpt []ent.Option
	entOpt = append(entOpt, ent.Debug())

	dsn := newDSN()

	return ent.Open(dialect.Postgres, dsn, entOpt...)
}

func createDBSchema(client *ent.Client) {
	if err := client.Schema.Create(
		context.Background(),
		migrate.WithDropIndex(true),
		migrate.WithDropColumn(true),
	); err != nil {
		log.Fatalf("error creating database resources: %v", err)
	}
}
