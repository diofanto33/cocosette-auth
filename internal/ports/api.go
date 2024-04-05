package ports

import (
	"context"

	"github.com/diofanto33/cocosette-auth/internal/application/core/domain"
)

type APIPort interface {
	Enroll(ctx context.Context, user *domain.User) (domain.User, error)
	SignIn(ctx context.Context, user domain.User) (string, error)
	Authenticate(ctx context.Context, token string) (int64, error)
}
