package api

import (
	"context"

	"github.com/diofanto33/cocosette-auth/internal/application/core/domain"
	"github.com/diofanto33/cocosette-auth/internal/ports"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type Application struct {
	db  ports.DBPort
	jwt ports.JWTPort
}

func NewApplication(db ports.DBPort, jwt ports.JWTPort) *Application {
	return &Application{
		db:  db,
		jwt: jwt,
	}
}

func (a Application) Enroll(ctx context.Context, user *domain.User) (domain.User, error) {
	_, err := a.db.Get(ctx, user.Email)
	if err == nil {
		return domain.User{}, status.Error(codes.AlreadyExists, "user already exists")
	}

	user.Password = a.jwt.HashPassword(user.Password)
	err = a.db.Save(ctx, user)
	if err != nil {
		return domain.User{}, err
	}
	return *user, nil
}

func (a Application) SignIn(ctx context.Context, user domain.User) (string, error) {
	u, err := a.db.Get(ctx, user.Email)
	if err != nil {
		if status.Code(err) == codes.NotFound {
			return "", status.Error(codes.NotFound, "user not found")
		}
		return "", err
	}

	if !a.jwt.CheckPasswordHash(user.Password, u.Password) {
		return "", status.Error(codes.InvalidArgument, "invalid password")
	}

	token, err := a.jwt.GenerateToken(ctx, &u)
	if err != nil {
		return "", err
	}

	return token, nil
}

func (a Application) Authenticate(ctx context.Context, signedToken string) (int64, error) {
	email, err := a.jwt.ValidateToken(ctx, signedToken)
	if err != nil {
		return 0, err
	}

	user, err := a.db.Get(ctx, email)
	if err != nil {
		return 0, err
	}

	return user.ID, nil
}
