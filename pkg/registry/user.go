package registry

import (
	"gitlab.com/trustify/core/pkg/adapter/controller"
	"gitlab.com/trustify/core/pkg/adapter/repository"
	"gitlab.com/trustify/core/pkg/usercase/usecase"
)

func (r *registry) NewUserController() controller.User {
	repo := repository.NewUserRepository(r.client)
	u := usecase.NewUserUsecase(repo)

	return controller.NewUserController(u)
}
