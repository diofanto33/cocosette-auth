package jwt

import (
	"context"
	"errors"
	"time"

	"github.com/diofanto33/cocosette-auth/internal/application/core/domain"
	"github.com/golang-jwt/jwt"
	"golang.org/x/crypto/bcrypt"
)

type Adapter struct {
	SecretKey       string
	Issuer          string
	ExpirationHours int64
	Claims          jwtClaims
}

type jwtClaims struct {
	jwt.StandardClaims
	Id    int64
	Email string
}

func NewAdapter(secretKey string, issuer string, expirationHours int64) (*Adapter, error) {
	return &Adapter{
		SecretKey:       secretKey,
		Issuer:          issuer,
		ExpirationHours: expirationHours,
		Claims:          jwtClaims{},
	}, nil
}

func (w *Adapter) GenerateToken(ctx context.Context, user *domain.User) (string, error) {
	claims := &jwtClaims{
		Id:    user.ID,
		Email: user.Email,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Local().Add(time.Hour * time.Duration(w.ExpirationHours)).Unix(),
			Issuer:    w.Issuer,
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	signedToken, err := token.SignedString([]byte(w.SecretKey))

	if err != nil {
		return "", err
	}

	return signedToken, nil
}

func (w *Adapter) ValidateToken(ctx context.Context, signedToken string) (string, error) {
	token, err := jwt.ParseWithClaims(
		signedToken,
		&jwtClaims{},
		func(token *jwt.Token) (interface{}, error) {
			return []byte(w.SecretKey), nil
		},
	)
	if err != nil {
		return "", err
	}

	claims, ok := token.Claims.(*jwtClaims)
	if !ok {
		return "", errors.New("Couldn't parse claims")
	}

	if claims.ExpiresAt < time.Now().Unix() {
		return "", errors.New("JWT is expired")
	}
	return claims.Email, nil
}

func (w *Adapter) HashPassword(password string) string {
	bytes, _ := bcrypt.GenerateFromPassword([]byte(password), 5)

	return string(bytes)
}

func (w *Adapter) CheckPasswordHash(password string, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))

	return err == nil
}
