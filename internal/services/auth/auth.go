package auth

import (
	"context"
	"errors"
	"fmt"
	"github.com/cuongpiger/sso/internal/domain/models"
	"github.com/cuongpiger/sso/internal/lib/jwt"
	"github.com/cuongpiger/sso/internal/lib/logger/sl"
	"github.com/cuongpiger/sso/internal/storage"
	"log/slog"
	"time"

	"golang.org/x/crypto/bcrypt"
)

type Auth struct {
	log          *slog.Logger
	userSaver    UserSaver
	userProvider UserProvider
	appProvider  AppProvider
	tokenTTL     time.Duration
}

type UserSaver interface {
	SaveUser(ctx context.Context, email string, passHash []byte) (uid uint64, err error)
}

type UserProvider interface {
	GetUser(ctx context.Context, email string) (models.User, error)
}

type AppProvider interface {
	GetApp(ctx context.Context, appName string) (models.App, error)
}

var (
	ErrInvalidCredentials = errors.New("invalid credentials")
	ErrUserExists         = errors.New("user already exists")
	ErrUserNotFound       = errors.New("user not found")
)

// New returns a new instanse of the Auth service
func New(log *slog.Logger, userSaver UserSaver, userProvider UserProvider, appProvider AppProvider, tokenTTL time.Duration) *Auth {
	return &Auth{
		log:          log,
		userSaver:    userSaver,
		userProvider: userProvider,
		appProvider:  appProvider,
		tokenTTL:     tokenTTL,
	}
}

// Login checks if user with given credentials exists in the system
func (a *Auth) Login(ctx context.Context, email string, password string, appName string) (string, error) {
	const op = "auth.Login"
	log := a.log.With(slog.String("op", op))

	log.Info("attempting to login user")
	user, err := a.userProvider.GetUser(ctx, email)
	if err != nil {
		if errors.Is(err, storage.ErrUserNotFound) {
			log.Warn("user not found", sl.Err(err))
			return "", fmt.Errorf("%s:%w", op, ErrInvalidCredentials)
		}
		log.Error("failed to find user", sl.Err(err))
		return "", fmt.Errorf("%s:%w", op, err)
	}

	if err := bcrypt.CompareHashAndPassword(user.PassHash, []byte(password)); err != nil {
		log.Error("failed to compare passwords", sl.Err(err))
		return "", fmt.Errorf("%s:%w", op, ErrInvalidCredentials)
	}

	app, err := a.appProvider.GetApp(ctx, appName)
	if err != nil {
		if errors.Is(err, storage.ErrAppNotFound) {
			log.Warn("app not found", sl.Err(err))
			return "", fmt.Errorf("%s:%w", op, ErrInvalidCredentials)
		}
		log.Error("failed to find app", sl.Err(err))
		return "", fmt.Errorf("%s:%w", op, err)
	}

	log.Info("user logged in successfully")
	token, err := jwt.NewToken(user, app, a.tokenTTL)
	if err != nil {
		log.Error("failed to generate token", sl.Err(err))
		return "", fmt.Errorf("%s:%w", op, err)
	}
	return token, nil
}

// RegisterNewUser registers new user in the system and returns userID
func (a *Auth) RegisterNewUser(ctx context.Context, email string, password string) (uint64, error) {
	const op = "auth.RegisterNewUser"
	log := a.log.With(slog.String("op", op))

	log.Info("registering user")
	passHash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		log.Error("failed to generate password Hash", sl.Err(err))
		return 0, fmt.Errorf("%s:%w", op, err)
	}

	id, err := a.userSaver.SaveUser(ctx, email, passHash)
	if err != nil {
		if errors.Is(err, storage.ErrUserExists) {
			log.Warn("user already exists", sl.Err(err))
			return 0, fmt.Errorf("%s:%w", op, ErrUserExists)
		}
		log.Error("failed to save user", sl.Err(err))
		return 0, fmt.Errorf("%s:%w", op, err)
	}
	log.Info("user registered")

	return id, nil
}

func (a *Auth) GetUserID(ctx context.Context, email string) (uint64, error) {
	const op = "auth.GetUserID"
	log := a.log.With(slog.String("op", op))

	log.Info("attemting to get userID")
	user, err := a.userProvider.GetUser(ctx, email)
	if err != nil {
		if errors.Is(err, storage.ErrUserNotFound) {
			log.Warn("user not found", sl.Err(err))
			return 0, fmt.Errorf("%s:%w", op, ErrInvalidCredentials)
		}
		log.Error("failed to find user", sl.Err(err))
		return 0, fmt.Errorf("%s:%w", op, err)
	}
	log.Info("user founded")

	return user.ID, nil
}
