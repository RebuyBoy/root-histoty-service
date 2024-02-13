package auth

import (
	"errors"
	"fmt"
	"time"

	"github.com/dgrijalva/jwt-go"
)

type Token struct {
	Token     string `json:"accessToken"`
	ExpiredAt int64  `json:"expiredAt"`
}
type TokenManager interface {
	NewAccessToken(userId string) (Token, error)
	NewRefreshToken(userId string) (Token, error)
	Parse(token string) (string, error)
}

type Manager struct {
	secretWord             string
	accessTokenTTLSeconds  int
	refreshTokenTTLSeconds int
}

func NewTokenManager(secretWord string, accessToken, refreshTokenTTL int) (*Manager, error) {
	if secretWord == "" {
		return nil, errors.New("empty signing key")
	}

	return &Manager{secretWord: secretWord, accessTokenTTLSeconds: accessToken, refreshTokenTTLSeconds: refreshTokenTTL}, nil
}

func (m *Manager) NewAccessToken(userId string) (Token, error) {
	//TODO set token type
	expiresAt := time.Now().Add(time.Duration(m.accessTokenTTLSeconds) * time.Second).Unix()
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.StandardClaims{
		ExpiresAt: expiresAt,
		Subject:   userId,
		IssuedAt:  time.Now().Unix(),
	})

	jwtToken, err := token.SignedString([]byte(m.secretWord))
	if err != nil {
		return Token{}, err
	}

	return Token{
		jwtToken,
		expiresAt,
	}, nil
}

func (m *Manager) NewRefreshToken(userId string) (Token, error) {
	expiresAt := time.Now().Add(time.Duration(m.refreshTokenTTLSeconds) * time.Second).Unix()
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.StandardClaims{
		ExpiresAt: expiresAt,
		Subject:   userId,
		IssuedAt:  time.Now().Unix(),
	})

	jwtToken, err := token.SignedString([]byte(m.secretWord))
	if err != nil {
		return Token{}, err
	}

	return Token{
		jwtToken,
		expiresAt,
	}, nil
}

func (m *Manager) Parse(accessToken string) (string, error) {
	token, err := jwt.Parse(accessToken, func(token *jwt.Token) (i interface{}, err error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}

		return []byte(m.secretWord), nil
	})
	if err != nil {
		return "", err
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return "", fmt.Errorf("error get user claims from token")
	}

	return claims["sub"].(string), nil
}
