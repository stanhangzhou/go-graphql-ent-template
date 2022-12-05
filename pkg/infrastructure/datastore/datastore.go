package datastore

import (
	"fmt"

	"entgo.io/ent/dialect"
	_ "github.com/lib/pq"
	"gitlab.com/trustify/core/config"
	"gitlab.com/trustify/core/ent"
)

// New returns data source name
func New() string {
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

// NewClient returns an orm client
func NewClient() (*ent.Client, error) {
	var entOpt []ent.Option
	entOpt = append(entOpt, ent.Debug())

	dsn := New()

	return ent.Open(dialect.Postgres, dsn, entOpt...)
}
