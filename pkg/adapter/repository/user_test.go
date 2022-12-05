package repository_test

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"gitlab.com/trustify/core/ent"
	"gitlab.com/trustify/core/ent/schema/ulid"
	"gitlab.com/trustify/core/pkg/adapter/repository"
	"gitlab.com/trustify/core/pkg/entity/model"
	"gitlab.com/trustify/core/testutil"
)

func setup(t *testing.T) (client *ent.Client, teardown func()) {
	testutil.ReadConfig()
	c := testutil.NewDBClient(t)

	return c, func() {
		testutil.DropUser(t, c)
		defer c.Close()
	}
}

func TestUserRepository__Get(t *testing.T) {
	t.Helper()

	client, teardown := setup(t)
	defer teardown()

	repo := repository.NewUserRepository(client)

	type args struct {
		ctx context.Context
	}

	tests := []struct {
		name    string
		arrange func(t *testing.T) model.ID
		act     func(ctx context.Context, t *testing.T, id model.ID) (u *model.User, err error)
		assert  func(t *testing.T, u *model.User, err error)
		args    struct {
			ctx context.Context
		}
		teardown func(t *testing.T)
	}{
		{
			name: "it should get users by id",
			arrange: func(t *testing.T) model.ID {
				ctx := context.Background()
				_, err := client.User.Delete().Exec(ctx)
				if err != nil {
					t.Error(err)
					t.FailNow()
				}
				user, err := repo.Create(ctx, model.CreateUserInput{
					FirstName: "John",
					LastName:  "Doe",
					Email:     "john@yourname.xzy",
					Password:  "secret",
				})
				if err != nil {
					t.Error(err)
					t.FailNow()
				}
				return user.ID
			},
			act: func(ctx context.Context, t *testing.T, id model.ID) (u *model.User, err error) {
				return repo.Get(ctx, &id)
			},
			assert: func(t *testing.T, got *model.User, err error) {
				assert.Nil(t, err)
				assert.NotNil(t, got.ID)
				assert.Equal(t, "John", got.FirstName)
				assert.Equal(t, "Doe", got.LastName)
				assert.Equal(t, "john@yourname.xzy", got.Email)
				assert.Equal(t, "secret", got.Password)
				assert.NotNil(t, got.CreatedAt)
				assert.NotNil(t, got.UpdatedAt)
			},
			args: args{
				ctx: context.Background(),
			},
			teardown: func(t *testing.T) {
				testutil.DropUser(t, client)
			},
		},
		{
			name: "it should err if user dose not exist",
			arrange: func(t *testing.T) model.ID {
				ctx := context.Background()
				_, err := client.User.Delete().Exec(ctx)
				if err != nil {
					t.Error(err)
					t.FailNow()
				}
				return ulid.MustNew("")
			},
			act: func(ctx context.Context, t *testing.T, id model.ID) (u *model.User, err error) {
				return repo.Get(ctx, &id)
			},
			assert: func(t *testing.T, u *model.User, err error) {
				assert.Nil(t, u)
				assert.NotNil(t, err)
				assert.Equal(t, err.Error(), "user not found")
			},
			args: args{
				ctx: context.Background(),
			},
			teardown: func(t *testing.T) {
				testutil.DropUser(t, client)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			arrange := tt.arrange(t)
			got, err := tt.act(tt.args.ctx, t, arrange)
			tt.assert(t, got, err)
			tt.teardown(t)
		})
	}
}

func TestUserRepository__List(t *testing.T) {
	t.Helper()

	client, teardown := setup(t)
	defer teardown()

	repo := repository.NewUserRepository(client)

	type args struct {
		ctx context.Context
	}

	tests := []struct {
		name    string
		arrange func(t *testing.T)
		act     func(ctx context.Context, t *testing.T) (uc *model.UserConnection, err error)
		assert  func(t *testing.T, uc *model.UserConnection, err error)
		args    struct {
			ctx context.Context
		}
		teardown func(t *testing.T)
	}{
		{
			name: "it should get user's list",
			arrange: func(t *testing.T) {
				ctx := context.Background()
				_, err := client.User.Delete().Exec(ctx)
				if err != nil {
					t.Error(err)
					t.FailNow()
				}

				users := []struct {
					firstName string
					lastName  string
					email     string
					password  string
				}{
					{firstName: "John", lastName: "Doe", email: "john.doe@yourname.xyz", password: "secret"},
					{firstName: "Jack", lastName: "Sparrow", email: "jack@yourname.xyz", password: "secret"},
					{firstName: "Harry", lastName: "Potter", email: "pottter@yourname.xyz", password: "james"},
				}
				bulk := make([]*ent.UserCreate, len(users))
				for i, u := range users {
					bulk[i] = client.User.Create().
						SetFirstName(u.firstName).
						SetLastName(u.lastName).
						SetEmail(u.email).
						SetPassword(u.password)
				}

				_, err = client.User.CreateBulk(bulk...).Save(ctx)
				if err != nil {
					t.Error(err)
					t.FailNow()
				}
			},
			act: func(ctx context.Context, _ *testing.T) (us *model.UserConnection, err error) {
				first := 5
				return repo.List(ctx, nil, &first, nil, nil, nil)
			},
			assert: func(t *testing.T, got *model.UserConnection, err error) {
				assert.Nil(t, err)
				assert.Equal(t, 3, len(got.Edges))
			},
			args: args{
				ctx: context.Background(),
			},
			teardown: func(t *testing.T) {
				testutil.DropUser(t, client)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.arrange(t)
			got, err := tt.act(tt.args.ctx, t)
			tt.assert(t, got, err)
			tt.teardown(t)
		})
	}
}

func TestUserRepository__Create(t *testing.T) {
	t.Helper()

	client, teardown := setup(t)
	defer teardown()

	repo := repository.NewUserRepository(client)

	type args struct {
		ctx context.Context
	}

	tests := []struct {
		name    string
		arrange func(t *testing.T)
		act     func(ctx context.Context, t *testing.T) (uc *model.User, err error)
		assert  func(t *testing.T, u *model.User, err error)
		args    struct {
			ctx context.Context
		}
		teardown func(t *testing.T)
	}{
		{
			name: "it should create user",
			arrange: func(t *testing.T) {
				ctx := context.Background()
				_, err := client.User.Delete().Exec(ctx)
				if err != nil {
					t.Error(err)
					t.FailNow()
				}
			},
			act: func(ctx context.Context, _ *testing.T) (us *model.User, err error) {
				return repo.Create(ctx, model.CreateUserInput{
					FirstName: "John",
					LastName:  "Doe",
					Email:     "john@yourname.xyz",
					Password:  "secret",
				})
			},
			assert: func(t *testing.T, got *model.User, err error) {
				assert.Nil(t, err)
				assert.NotNil(t, got.ID)
				assert.Equal(t, "John", got.FirstName)
				assert.Equal(t, "Doe", got.LastName)
				assert.Equal(t, "john@yourname.xyz", got.Email)
				assert.Equal(t, "secret", got.Password)
				assert.NotNil(t, got.CreatedAt)
				assert.NotNil(t, got.UpdatedAt)
			},
			args: args{
				ctx: context.Background(),
			},
			teardown: func(t *testing.T) {
				testutil.DropUser(t, client)
			},
		},
		{
			name: "it should fail if email already exist",
			arrange: func(t *testing.T) {
				ctx := context.Background()
				_, err := repo.Create(ctx, model.CreateUserInput{
					FirstName: "John",
					LastName:  "Doe",
					Email:     "john@yourname.xyz",
					Password:  "secret",
				})
				if err != nil {
					t.Error(err)
					t.FailNow()
				}
			},
			act: func(ctx context.Context, _ *testing.T) (us *model.User, err error) {
				return repo.Create(ctx, model.CreateUserInput{
					FirstName: "John",
					LastName:  "Doe",
					Email:     "john@yourname.xyz",
					Password:  "secret",
				})
			},
			assert: func(t *testing.T, got *model.User, err error) {
				assert.Nil(t, got)
				assert.Equal(t, "failed to create user", err.Error())
			},
			args: args{
				ctx: context.Background(),
			},
			teardown: func(t *testing.T) {
				testutil.DropUser(t, client)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.arrange(t)
			got, err := tt.act(tt.args.ctx, t)
			tt.assert(t, got, err)
			tt.teardown(t)
		})
	}
}

func TestUserRepository__Update(t *testing.T) {
	t.Helper()

	client, teardown := setup(t)
	defer teardown()

	repo := repository.NewUserRepository(client)

	type args struct {
		ctx context.Context
	}

	tests := []struct {
		name    string
		arrange func(t *testing.T) model.ID
		act     func(ctx context.Context, t *testing.T, id model.ID) (uc *model.User, err error)
		assert  func(t *testing.T, u *model.User, err error)
		args    struct {
			ctx context.Context
		}
		teardown func(t *testing.T)
	}{
		{
			name: "it should update user",
			arrange: func(t *testing.T) model.ID {
				ctx := context.Background()
				_, err := client.User.Delete().Exec(ctx)
				if err != nil {
					t.Error(err)
					t.FailNow()
				}
				u, err := repo.Create(ctx, model.CreateUserInput{
					FirstName: "John",
					LastName:  "Doe",
					Email:     "john@yourname.xyz",
					Password:  "secret",
				})
				if err != nil {
					t.Error(err)
					t.FailNow()
				}
				return u.ID
			},
			act: func(ctx context.Context, _ *testing.T, id model.ID) (us *model.User, err error) {
				firstName, lastName, email, password := "Max", "Smith", "max@yourname.xyz", "supersecret"
				return repo.Update(ctx, model.UpdateUserInput{
					ID:        id,
					FirstName: &firstName,
					LastName:  &lastName,
					Email:     &email,
					Password:  &password,
				})
			},
			assert: func(t *testing.T, got *model.User, err error) {
				assert.Nil(t, err)
				assert.NotNil(t, got.ID)
				assert.Equal(t, "Max", got.FirstName)
				assert.Equal(t, "Smith", got.LastName)
				assert.Equal(t, "max@yourname.xyz", got.Email)
				assert.Equal(t, "supersecret", got.Password)
				assert.NotNil(t, got.CreatedAt)
				assert.NotNil(t, got.UpdatedAt)
			},
			args: args{
				ctx: context.Background(),
			},
			teardown: func(t *testing.T) {
				testutil.DropUser(t, client)
			},
		},
		{
			name: "it should return error on fail",
			arrange: func(t *testing.T) model.ID {
				ctx := context.Background()
				_, err := client.User.Delete().Exec(ctx)
				if err != nil {
					t.Error(err)
					t.FailNow()
				}
				u, err := repo.Create(ctx, model.CreateUserInput{
					FirstName: "John",
					LastName:  "Doe",
					Email:     "john@yourname.xyz",
					Password:  "secret",
				})
				repo.Create(ctx, model.CreateUserInput{
					FirstName: "Max",
					LastName:  "Doe",
					Email:     "max@yourname.xyz",
					Password:  "secret",
				})
				if err != nil {
					t.Error(err)
					t.FailNow()
				}
				return u.ID
			},
			act: func(ctx context.Context, _ *testing.T, id model.ID) (us *model.User, err error) {
				firstName, lastName, email, password := "Max", "Smith", "max@yourname.xyz", "supersecret"
				return repo.Update(ctx, model.UpdateUserInput{
					ID:        id,
					FirstName: &firstName,
					LastName:  &lastName,
					Email:     &email,
					Password:  &password,
				})
			},
			assert: func(t *testing.T, got *model.User, err error) {
				assert.Nil(t, got)
				assert.Equal(t, "failed to update user", err.Error())
			},
			args: args{
				ctx: context.Background(),
			},
			teardown: func(t *testing.T) {
				testutil.DropUser(t, client)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ar := tt.arrange(t)
			got, err := tt.act(tt.args.ctx, t, ar)
			tt.assert(t, got, err)
			tt.teardown(t)
		})
	}
}

func TestUserRepository__EmailExists(t *testing.T) {
	t.Helper()

	client, teardown := setup(t)
	defer teardown()

	repo := repository.NewUserRepository(client)

	type args struct {
		ctx context.Context
	}

	tests := []struct {
		name    string
		arrange func(t *testing.T)
		act     func(ctx context.Context, t *testing.T) (exists bool, err error)
		assert  func(t *testing.T, got bool, err error)
		args    struct {
			ctx context.Context
		}
		teardown func(t *testing.T)
	}{
		{
			name: "it return false if no user with the given email exists",
			arrange: func(t *testing.T) {
				ctx := context.Background()
				_, err := client.User.Delete().Exec(ctx)
				if err != nil {
					t.Error(err)
					t.FailNow()
				}
			},
			act: func(ctx context.Context, _ *testing.T) (exists bool, err error) {
				return repo.EmailExists(ctx, "john@yourname.xyz")
			},
			assert: func(t *testing.T, got bool, err error) {
				assert.Nil(t, err)
				assert.Equal(t, false, got)
			},
			args: args{
				ctx: context.Background(),
			},
			teardown: func(t *testing.T) {
				testutil.DropUser(t, client)
			},
		},
		{
			name: "it return true if user exists",
			arrange: func(t *testing.T) {
				ctx := context.Background()
				_, err := client.User.Delete().Exec(ctx)
				if err != nil {
					t.Error(err)
					t.FailNow()
				}
				_, err = repo.Create(ctx, model.CreateUserInput{
					FirstName: "John",
					LastName:  "Doe",
					Email:     "john@yourname.xyz",
					Password:  "secret",
				})
				if err != nil {
					t.Error(err)
					t.FailNow()
				}
			},
			act: func(ctx context.Context, _ *testing.T) (exists bool, err error) {
				return repo.EmailExists(ctx, "john@yourname.xyz")
			},
			assert: func(t *testing.T, got bool, err error) {
				assert.Nil(t, err)
				assert.Equal(t, true, got)
			},
			args: args{
				ctx: context.Background(),
			},
			teardown: func(t *testing.T) {
				testutil.DropUser(t, client)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.arrange(t)
			got, err := tt.act(tt.args.ctx, t)
			tt.assert(t, got, err)
			tt.teardown(t)
		})
	}
}
