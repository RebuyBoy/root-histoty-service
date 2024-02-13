package service

import (
	"context"
	"errors"
	"fmt"
	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
	"root-histoty-service/internal/model"
	"root-histoty-service/internal/repository"
	"root-histoty-service/pkg/auth"
	"time"
)

type PlayerService struct {
	userRepo     *repository.UserRepo
	logger       *logrus.Logger
	tokenManager auth.TokenManager
}

type Tokens struct {
	AccessToken          string
	RefreshToken         string
	AccessTokenExpireAt  int64
	RefreshTokenExpireAt int64
}

func NewPlayerService(userRepo *repository.UserRepo, logger *logrus.Logger, tokenManager auth.TokenManager) *PlayerService {
	return &PlayerService{
		userRepo:     userRepo,
		logger:       logger,
		tokenManager: tokenManager,
	}
}

func (p *PlayerService) Register(ctx context.Context, user *model.Player) error {
	exists, err := p.userRepo.IsExists(ctx, user.Name)
	if err != nil {
		return errors.New("failed check existed player")
	}

	if exists {
		return errors.New("user already exists")
	}

	if user.Name == "" {
		return errors.New("name cant be empty")
	}
	user.ID = uuid.New().String()
	user.RegisteredAt = time.Now()
	return p.userRepo.SaveUser(ctx, *user)
}

func (p *PlayerService) GetUserByName(ctx context.Context, name string) (*model.Player, error) {
	return p.userRepo.GetUserByName(ctx, name)
}

func (p *PlayerService) Authorize(ctx context.Context, login string, pinCode int) (*Tokens, error) {
	user, err := p.userRepo.GetUserByName(ctx, login)
	if err != nil {
		return nil, fmt.Errorf("user not found: %w", err)
	}
	if user.PinCode != pinCode {
		return nil, fmt.Errorf("invalid pin code: %w", err)
	}
	return p.createSession(user.ID)
}

func (p *PlayerService) createSession(userId string) (*Tokens, error) {
	var res Tokens
	accessToken, err := p.tokenManager.NewAccessToken(userId)

	res.AccessToken = accessToken.Token
	res.AccessTokenExpireAt = accessToken.ExpiredAt

	if err != nil {
		return nil, fmt.Errorf("failed to sigh token: %w", err)
	}

	refreshToken, err := p.tokenManager.NewRefreshToken(userId)
	if err != nil {
		return nil, fmt.Errorf("failed to sigh token: %w", err)
	}
	res.RefreshToken = refreshToken.Token
	res.RefreshTokenExpireAt = refreshToken.ExpiredAt
	if err != nil {
		return nil, fmt.Errorf("failed to sigh token: %w", err)
	}

	return &res, nil
}

func (p *PlayerService) RefreshTokens(ctx context.Context, refreshTokenId string) (*Tokens, error) {
	return nil, nil
}
