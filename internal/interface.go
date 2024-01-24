package internal

import (
	"context"
	"root-histoty-service/internal/model"
)

type UserRepo interface {
	SaveUser(ctx context.Context, user model.Player) error
	IsExists(ctx context.Context, name string) (bool, error)
	GetUserByName(ctx context.Context, name string) (*model.Player, error)
	GetUserById(ctx context.Context, id string) (*model.Player, error)
}

type PlayerService interface {
	Register(ctx context.Context, user *model.Player) error
	Authorize(ctx context.Context, login string, pinCode int) (string, error)
	GetUserByName(ctx context.Context, name string) (*model.Player, error)
	//GetUserById(id string) (*model.Player, error)
}
