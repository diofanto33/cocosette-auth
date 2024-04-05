package db

import (
	"context"
	"fmt"

	"github.com/diofanto33/cocosette-auth/internal/application/core/domain"
	"github.com/uptrace/opentelemetry-go-extra/otelgorm"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	Email    string `json:"email"`
	Password string `json:"password"`
}

type Adapter struct {
	db *gorm.DB
}

func (a *Adapter) ErrRecordNotFound() error {
	return gorm.ErrRecordNotFound
}

func (a Adapter) Get(ctx context.Context, email string) (domain.User, error) {
	var userEntity User
	res := a.db.WithContext(ctx).Where("email = ?", email).First(&userEntity)
	return domain.User{
		ID:       int64(userEntity.ID),
		Email:    userEntity.Email,
		Password: userEntity.Password,
	}, res.Error
}

func (a Adapter) Save(ctx context.Context, user *domain.User) error {
	userEntity := User{
		Email:    user.Email,
		Password: user.Password,
	}
	res := a.db.WithContext(ctx).Create(&userEntity)
	return res.Error
}

func NewAdapter(dataSourceUrl string) (*Adapter, error) {
	db, openErr := gorm.Open(postgres.Open(dataSourceUrl), &gorm.Config{})
	if openErr != nil {
		return nil, fmt.Errorf("db connection error: %v", openErr)
	}
	if err := db.Use(otelgorm.NewPlugin(otelgorm.WithDBName("auth"))); err != nil {
		return nil, fmt.Errorf("db otel plugin error: %v", err)
	}
	err := db.AutoMigrate(&User{})
	if err != nil {
		return nil, fmt.Errorf("db migration error: %v", err)
	}
	return &Adapter{db: db}, nil
}
