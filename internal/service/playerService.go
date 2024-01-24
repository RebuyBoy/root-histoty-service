package service

import (
	"context"
	"errors"
	"fmt"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
	"root-histoty-service/internal"
	"root-histoty-service/internal/model"
	"time"
)

type PlayerService struct {
	userRepo   internal.UserRepo
	logger     *logrus.Logger
	secretWord string
}

type CustomClaims struct {
	UserID string `json:"user_id"`
	//Username string `json:"username"`
	jwt.RegisteredClaims
}

func NewPlayerService(userRepo internal.UserRepo, logger *logrus.Logger, secretWord string) *PlayerService {
	return &PlayerService{userRepo: userRepo, logger: logger, secretWord: secretWord}
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

func (p *PlayerService) Authorize(ctx context.Context, login string, pinCode int) (string, error) {
	user, err := p.userRepo.GetUserByName(ctx, login)
	if err != nil {
		return "", err
	}
	if user.PinCode != pinCode {
		return "", errors.New("invalid pin code")
	}
	now := time.Now()

	claims := CustomClaims{
		UserID: user.ID,
		RegisteredClaims: jwt.RegisteredClaims{
			Issuer:    user.ID,
			ExpiresAt: jwt.NewNumericDate(now.Add(10 * time.Second)),
			IssuedAt:  jwt.NewNumericDate(now),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString([]byte(p.secretWord))
	if err != nil {
		return "", fmt.Errorf("failed to sigh token: %w", err)
	}
	return tokenString, nil
}
