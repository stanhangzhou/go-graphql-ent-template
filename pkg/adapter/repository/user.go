package repository

import (
	"context"
	"errors"
	"time"

	"gitlab.com/trustify/core/ent"
	"gitlab.com/trustify/core/ent/user"
	"gitlab.com/trustify/core/pkg/entity/model"
	usecaseRepository "gitlab.com/trustify/core/pkg/usercase/repository"
)

type userRepository struct {
	client *ent.Client
}

func NewUserRepository(client *ent.Client) usecaseRepository.User {
	return &userRepository{client: client}
}

func (r *userRepository) Get(ctx context.Context, id *model.ID) (*model.User, error) {
	u, err := r.client.User.Query().Where(user.IDEQ(*id)).Only(ctx)
	if err != nil {
		return nil, errors.New("user not found")
	}
	return u, nil
}

func (r *userRepository) List(ctx context.Context, after *model.Cursor, first *int, before *model.Cursor, last *int, where *model.UserWhereInput) (*model.UserConnection, error) {
	return r.client.User.Query().Paginate(ctx, after, first, before, last, ent.WithUserFilter(where.Filter))
}

func (r *userRepository) Create(ctx context.Context, input model.CreateUserInput) (*model.User, error) {
	u, err := r.client.User.Create().SetInput(input).Save(ctx)
	if err != nil {
		return nil, errors.New("failed to create user")
	}
	return u, nil
}

func (r *userRepository) Update(ctx context.Context, input model.UpdateUserInput) (*model.User, error) {
	u, err := r.client.User.UpdateOneID(input.ID).SetInput(input).SetUpdatedAt(time.Now()).Save(ctx)
	if err != nil {
		return nil, errors.New("failed to update user")
	}
	return u, nil
}

func (r *userRepository) EmailExists(ctx context.Context, email string) (bool, error) {
	return r.client.User.Query().Where(user.Email(email)).Exist(ctx)
}
