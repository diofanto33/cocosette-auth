package ports

import (
	"context"

	"github.com/diofanto33/cocosette-auth/internal/application/core/domain"
)

type DBPort interface {
	Get(ctx context.Context, email string) (domain.User, error)
	Save(context.Context, *domain.User) error
	ErrRecordNotFound() error
}
