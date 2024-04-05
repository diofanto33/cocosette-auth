package ports

import (
	"context"

	"github.com/diofanto33/cocosette-auth/internal/application/core/domain"
)

type JWTPort interface {
	GenerateToken(ctx context.Context, user *domain.User) (string, error)
	ValidateToken(ctx context.Context, token string) (string, error)
	HashPassword(password string) string
	CheckPasswordHash(password, hash string) bool
}
