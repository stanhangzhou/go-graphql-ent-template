package testutil

import (
	"context"
	"testing"

	"entgo.io/ent/dialect"
	"gitlab.com/trustify/core/ent"
	"gitlab.com/trustify/core/ent/enttest"
	"gitlab.com/trustify/core/pkg/infrastructure/datastore"
)

// NewDBClient loads database for test
func NewDBClient(t *testing.T) *ent.Client {
	d := datastore.New()
	return enttest.Open(t, dialect.Postgres, d)
}

// DropAll drops all data from database
func DropAll(t *testing.T, client *ent.Client) {
	t.Log("drop data from database")
	DropUser(t, client)
}

// DropUser drops data from users
func DropUser(t *testing.T, client *ent.Client) {
	ctx := context.Background()
	_, err := client.User.Delete().Exec(ctx)
	if err != nil {
		t.Error(err)
		t.FailNow()
	}
}
