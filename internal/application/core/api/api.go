package api

import (
	"context"

	"github.com/diofanto33/cocosette-auth/internal/application/core/domain"
	"github.com/diofanto33/cocosette-auth/internal/ports"
	"github.com/pkg/errors"
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
	exists, err := a.db.Get(ctx, user.Email)
	if err != nil {
		switch errors.Cause(err) {
		case a.db.ErrRecordNotFound():
			user.Password = a.jwt.HashPassword(user.Password)
			err := a.db.Save(ctx, user)
			if err != nil {
				return domain.User{}, status.Error(codes.Internal, "could not save user")
			}
			return *user, nil
		default:
			return domain.User{}, status.Error(codes.Internal, "could not get user")
		}
	}

	return exists, status.Error(codes.AlreadyExists, "user already exists")
}

func (a Application) SignIn(ctx context.Context, user domain.User) (string, error) {
	exists, err := a.db.Get(ctx, user.Email)
	if err != nil {
		switch errors.Cause(err) {
		case a.db.ErrRecordNotFound():
			return "", status.Error(codes.NotFound, "user not found")
		default:
			return "", status.Error(codes.Internal, "could not get user")
		}
	}

	if !a.jwt.CheckPasswordHash(user.Password, exists.Password) {
		return "", status.Error(codes.InvalidArgument, "invalid password")
	}

	token, err := a.jwt.GenerateToken(ctx, &exists)
	if err != nil {
		return "", status.Error(codes.Internal, "could not generate token")
	}

	return token, nil
}

func (a Application) Authenticate(ctx context.Context, signedToken string) (int64, error) {
	email, err := a.jwt.ValidateToken(ctx, signedToken)
	if err != nil {
		return 0, status.Error(codes.Unauthenticated, "could not validate token")
	}

	user, err := a.db.Get(ctx, email)
	if err != nil {
		return 0, status.Error(codes.Internal, "could not get user")
	}

	return user.ID, nil
}
