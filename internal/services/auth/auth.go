package auth

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"sso/internal/domain/models"
	"sso/internal/jwt"
	"time"

	"golang.org/x/crypto/bcrypt"
)

type AuthService struct {
	log          *slog.Logger
	authProvider AuthProvider
	tokenTTL     time.Duration
}

var (
	ErrInvalidCredentials = errors.New("invalid credentials")
)

type AuthProvider interface {
	CreateUser(ctx context.Context, email string, password string) (uid string, err error)
	GetUserByEmail(ctx context.Context, email string) (user models.User, err error)
	GetAppById(ctx context.Context, id string) (app models.App, err error)
}

func New(
	log *slog.Logger,
	authProvider AuthProvider,
	tokenTTL time.Duration,
) *AuthService {
	return &AuthService{
		log:          log,
		authProvider: authProvider,
		tokenTTL:     tokenTTL,
	}
}

func (a *AuthService) Login(
	ctx context.Context,
	email string,
	password string,
	appID string,
) (string, error) {
	const op = "AuthService.Login"

	log := a.log.With(
		slog.String("op", op),
		slog.String("username", email),
	)

	user, err := a.authProvider.GetUserByEmail(ctx, email)
	if err != nil {
		log.Error("failed to get user", slog.Any("error", err))
		return "", fmt.Errorf("%s: %w", op, err)
	}

	app, err := a.authProvider.GetAppById(ctx, appID)
	if err != nil {
		log.Info("failed to get app", slog.Any("error", err))
		return "", fmt.Errorf("%s: %w", op, err)
	}

	err = bcrypt.CompareHashAndPassword(user.Password.Data, []byte(password))
	if err != nil {
		log.Info("invalid credentials", slog.Any("err", err))
		return "", fmt.Errorf("%s: %w", op, ErrInvalidCredentials)
	}

	token, err := jwt.NewToken(user, app, a.tokenTTL)
	if err != nil {
		log.Error("failed to generate jwt token", slog.Any("error", err))
		return "", fmt.Errorf("%s: %w", op, err)
	}

	return token, nil
}

func (a *AuthService) RegisterNewUser(ctx context.Context, email string, password string) (string, error) {
	const op = "AuthService.RegisterNewUser"

	log := a.log.With(
		slog.String("op", op),
		slog.String("email", email),
	)

	id, err := a.authProvider.CreateUser(ctx, email, password)
	if err != nil {
		log.Error("failed to save user", slog.Any("error", err))
		return "", fmt.Errorf("%s: %w", op, err)
	}

	return id, nil
}

func (a *AuthService) IsAdmin(ctx context.Context, userID string) (bool, error) {
	//TODO IsAdmin template
	return false, fmt.Errorf("not implemented")
}
