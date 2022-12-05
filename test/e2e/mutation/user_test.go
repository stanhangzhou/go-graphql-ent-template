package mutation_test

import (
	"context"
	"net/http"
	"testing"

	"github.com/gavv/httpexpect/v2"
	"gitlab.com/trustify/core/ent"
	"gitlab.com/trustify/core/pkg/infrastructure/router"
	"gitlab.com/trustify/core/testutil"
	"gitlab.com/trustify/core/testutil/e2e"
)

func TestUser_CreateUser(t *testing.T) {
	expect, client, teardown := e2e.Setup(t, e2e.SetupOption{
		Teardown: func(t *testing.T, client *ent.Client) {
			testutil.DropUser(t, client)
		},
	})
	defer teardown()

	tests := []struct {
		name    string
		arrange func(t *testing.T)
		act     func(t *testing.T) *httpexpect.Response
		assert  func(t *testing.T, got *httpexpect.Response)
		args    struct {
			ctx context.Context
		}
		teardown func(t *testing.T)
	}{
		{
			name:    "it should create user",
			arrange: func(t *testing.T) {},
			act: func(t *testing.T) *httpexpect.Response {
				return expect.POST(router.QueryPath).WithJSON(map[string]string{
					"query": `
						mutation CreateUser {
							createUser(
								input: {firstName: "John", lastName: "Doe", email: "john@yourname.xyz", password: "secret12345"}
							) {
								id
								firstName
								lastName
								email
								createdAt
								updatedAt
							}
						}`,
				}).Expect()
			},
			assert: func(_ *testing.T, got *httpexpect.Response) {
				got.Status(http.StatusOK)
				res := e2e.GetData(got).Object()
				user := e2e.GetObject(res, "createUser")
				user.Value("id").String().NotEmpty()
				user.Value("firstName").String().Equal("John")
				user.Value("lastName").String().Equal("Doe")
				user.Value("email").String().Equal("john@yourname.xyz")
				user.Value("createdAt").String().NotEmpty()
				user.Value("updatedAt").String().NotEmpty()
			},
			teardown: func(t *testing.T) {
				testutil.DropUser(t, client)
			},
		},
		{
			name:    "it should fail if password is shorter than 8 characters",
			arrange: func(t *testing.T) {},
			act: func(t *testing.T) *httpexpect.Response {
				return expect.POST(router.QueryPath).WithJSON(map[string]string{
					"query": `
						mutation CreateUser {
							createUser(
								input: {firstName: "John", lastName: "Doe", email: "john@yourname.xyz", password: "secret"}
							) {
								id
								firstName
								lastName
								email
								createdAt
								updatedAt
							}
						}`,
				}).Expect()
			},
			assert: func(_ *testing.T, got *httpexpect.Response) {
				got.Status(http.StatusOK)
				res := e2e.GetData(got)
				res.Null()

				errors := e2e.GetErrors(got)
				errors.Array().Length().Equal(1)
				errors.Array().First().Object().Value("message").Equal("password must be at least 8 characters in length")
			},
			teardown: func(t *testing.T) {
				testutil.DropUser(t, client)
			},
		},
		{
			name:    "it should fail if email is invalid",
			arrange: func(t *testing.T) {},
			act: func(t *testing.T) *httpexpect.Response {
				return expect.POST(router.QueryPath).WithJSON(map[string]string{
					"query": `
						mutation CreateUser {
							createUser(
								input: {firstName: "John", lastName: "Doe", email: "john@yourname", password: "secret12345"}
							) {
								id
								firstName
								lastName
								email
								createdAt
								updatedAt
							}
						}`,
				}).Expect()
			},
			assert: func(_ *testing.T, got *httpexpect.Response) {
				got.Status(http.StatusOK)
				res := e2e.GetData(got)
				res.Null()

				errors := e2e.GetErrors(got)
				errors.Array().Length().Equal(1)
				errors.Array().First().Object().Value("message").Equal("email must be a valid email address")
			},
			teardown: func(t *testing.T) {
				testutil.DropUser(t, client)
			},
		},
		{
			name: "it should fail if user with the given email already exists",
			arrange: func(_ *testing.T) {
				client.User.Create().
					SetFirstName("John").
					SetLastName("Doe").
					SetEmail("john@yourname.xyz").
					SetPassword("secret1234").
					Save(context.Background())
			},
			act: func(t *testing.T) *httpexpect.Response {
				return expect.POST(router.QueryPath).WithJSON(map[string]string{
					"query": `
						mutation CreateUser {
							createUser(
								input: {firstName: "John", lastName: "Doe", email: "john@yourname.xyz", password: "secret1234"}
							) {
								id
								firstName
								lastName
								email
								createdAt
								updatedAt
							}
						}`,
				}).Expect()
			},
			assert: func(_ *testing.T, got *httpexpect.Response) {
				got.Status(http.StatusOK)
				res := e2e.GetData(got)
				res.Null()

				errors := e2e.GetErrors(got)
				errors.Array().Length().Equal(1)
				errors.Array().First().Object().Value("message").Equal("user with the given email already exists")
			},
			teardown: func(t *testing.T) {
				testutil.DropUser(t, client)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.arrange(t)
			got := tt.act(t)
			tt.assert(t, got)
			tt.teardown(t)
		})
	}
}
