package auth

import (
	"context"
	"fmt"
	"github.com/golang-jwt/jwt/v5"
	"microservices/task_5/server/internal/models"
	"time"
)

type Claims struct {
	UserID int64 `json:"user_id"`
	jwt.RegisteredClaims
}

type Service struct {
	salt      string        // соль для создания токенов (хранится в .env)
	accessExp time.Duration // время в часах, через которое протухнет access (также хранится в .env)
}

func NewService(salt string, accessExp time.Duration) *Service {
	return &Service{salt: salt, accessExp: accessExp}
}

/* Функция создает пару {Accesss, Refresh} для заданного User (временно только Accesss) */
func (s *Service) GetTokens(_ context.Context, user models.User) (models.Token, error) {
	claims := Claims{
		UserID: user.ID,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(s.accessExp)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}

	accessTokenString, err := s.makeAccessToken(claims)
	if err != nil {
		return models.Token{}, err
	}

	tokens := models.Token{
		AccessToken: accessTokenString,
	}

	return tokens, nil
}

/* Проверяет пару {Access, Refresh} на валидность */
func (s *Service) Authenticate(_ context.Context, tokens models.Token) (models.Claims, error) {
	claims := Claims{}

	_, err := jwt.ParseWithClaims(tokens.AccessToken, &claims, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(s.salt), nil
	})

	if err != nil {
		return models.Claims{}, err
	}

	out := models.Claims{
		UserID: claims.UserID,
	}

	return out, nil
}

/* Создает Access */
func (s *Service) makeAccessToken(claims Claims) (string, error) {
	accessToken := jwt.NewWithClaims(jwt.SigningMethodHS512, claims)
	accessTokenString, err := accessToken.SignedString([]byte(s.salt))
	if err != nil {
		return "", err
	}
	return accessTokenString, nil
}
