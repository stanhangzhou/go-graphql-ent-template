package usecase

import (
	"context"
	"errors"

	"gitlab.com/trustify/core/pkg/entity/model"
	"gitlab.com/trustify/core/pkg/usercase/repository"
)

type user struct {
	userRepository repository.User
}

// User of usercase
type User interface {
	Get(ctx context.Context, id *model.ID) (*model.User, error)
	List(ctx context.Context, after *model.Cursor, first *int, before *model.Cursor, last *int, where *model.UserWhereInput) (*model.UserConnection, error)
	Create(ctx context.Context, input model.CreateUserInput) (*model.User, error)
	Update(ctx context.Context, input model.UpdateUserInput) (*model.User, error)
}

// NewUserUsecase returns user usecse
func NewUserUsecase(r repository.User) User {
	return &user{userRepository: r}
}

func (u *user) Get(ctx context.Context, id *model.ID) (*model.User, error) {
	return u.userRepository.Get(ctx, id)
}

func (u *user) List(ctx context.Context, after *model.Cursor, first *int, before *model.Cursor, last *int, where *model.UserWhereInput) (*model.UserConnection, error) {
	return u.userRepository.List(ctx, after, first, before, last, where)
}

func (u *user) Create(ctx context.Context, input model.CreateUserInput) (*model.User, error) {
	ex, err := u.userRepository.EmailExists(ctx, input.Email)
	if err != nil {
		return nil, err
	}

	if ex {
		return nil, errors.New("user with the given email already exists")
	}

	return u.userRepository.Create(ctx, input)
}

func (u *user) Update(ctx context.Context, input model.UpdateUserInput) (*model.User, error) {
	return u.userRepository.Update(ctx, input)
}
