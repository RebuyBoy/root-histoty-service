package internal

import (
	"context"
	"root-histoty-service/internal/model"
	"root-histoty-service/internal/service"
)

type UserRepo interface {
	SaveUser(ctx context.Context, user model.Player) error
	IsExists(ctx context.Context, name string) (bool, error)
	GetUserByName(ctx context.Context, name string) (*model.Player, error)
	GetUserById(ctx context.Context, id string) (*model.Player, error)
}

type PlayerService interface {
	Register(ctx context.Context, user *model.Player) error
	Authorize(ctx context.Context, login string, pinCode int) (*service.Tokens, error)
	GetUserByName(ctx context.Context, name string) (*model.Player, error)
	RefreshTokens(ctx context.Context, refreshToken string) (*service.Tokens, error)
	//GetUserById(id string) (*model.Player, error)
}
