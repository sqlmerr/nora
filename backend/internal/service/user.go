package service

import (
	"context"
	e "nora/internal/error"
	"nora/internal/model"
	"nora/internal/util"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

func (s *Service) Register(ctx context.Context, data *model.UserCreate) (*model.User, error) {
	if data.Password == "" || data.Username == "" {
		return nil, e.ErrInvalidCredentials
	}

	user, err := s.FindOneUserByUsername(ctx, data.Username)
	if err == nil && user != nil {
		return nil, e.ErrUsernameOccupied
	}

	hashedPassword, err := util.HashPassword(data.Password)
	if err != nil {
		return nil, err
	}

	data.Password = hashedPassword

	id, err := s.u.Create(ctx, data)
	if err != nil {
		return nil, err
	}

	return s.FindOneUser(ctx, *id)
}

func (s *Service) Login(ctx context.Context, data *model.UserLogin) (string, error) {
	user, err := s.FindOneUserByUsername(ctx, data.Username)
	if err != nil {
		return "", err
	}
	if user == nil {
		return "", e.ErrInvalidCredentials
	}

	if status := util.CheckPasswordHash(data.Password, user.Password); !status {
		return "", e.ErrInvalidCredentials
	}

	jwtToken := jwt.New(jwt.SigningMethodHS512)
	claims := jwtToken.Claims.(jwt.MapClaims)
	claims["sub"] = data.Username
	claims["exp"] = time.Now().Add(time.Hour * 24).Unix()

	return jwtToken.SignedString([]byte(s.cfg.JwtSecret))
}

func (s *Service) FindOneUser(ctx context.Context, id uuid.UUID) (*model.User, error) {
	return s.u.FindOne(ctx, id)
}

func (s *Service) FindOneUserByUsername(ctx context.Context, username string) (*model.User, error) {
	return s.u.FindOneByUsername(ctx, username)
}
