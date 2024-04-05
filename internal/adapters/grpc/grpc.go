package grpc

import (
	"context"

	"github.com/diofanto33/cocosette-auth/internal/application/core/domain"
	"github.com/diofanto33/cocosette-proto/golang/auth"

	log "github.com/sirupsen/logrus"
)

func (a Adapter) Register(ctx context.Context, req *auth.RegisterRequest) (*auth.RegisterResponse, error) {
	log.WithContext(ctx).Infof("Registering user with email: %s", req.Email)

	new_user := domain.NewUser(req.Email, req.Password)

	registered_user, err := a.api.Enroll(ctx, &new_user)
	if err != nil {
		log.WithContext(ctx).Errorf("Failed to register user with email: %s", req.Email)
		return &auth.RegisterResponse{}, err
	}

	log.WithContext(ctx).Infof("User registered successfully with Email: %d", registered_user.Email)

	return &auth.RegisterResponse{}, err
}

func (a Adapter) Login(ctx context.Context, req *auth.LoginRequest) (*auth.LoginResponse, error) {
	log.WithContext(ctx).Infof("Logging in user with email: %s", req.Email)

	user_req := domain.NewUser(req.Email, req.Password)

	token, err := a.api.SignIn(ctx, user_req)
	if err != nil {
		log.WithContext(ctx).Errorf("Failed to login user with email: %s", req.Email)
		return &auth.LoginResponse{}, err
	}

	log.WithContext(ctx).Infof("User logged in successfully with email: %s", req.Email)
	return &auth.LoginResponse{
		Token: token,
	}, err
}

func (a Adapter) Validate(ctx context.Context, req *auth.ValidateRequest) (*auth.ValidateResponse, error) {
	log.WithContext(ctx).Info("Validating token: %s", req.Token)
	userID, err := a.api.Authenticate(ctx, req.Token)
	if err != nil {
		log.WithContext(ctx).Error("Failed to validate token: %s", req.Token)
		return &auth.ValidateResponse{}, err
	}

	log.WithContext(ctx).Info("Token validated successfully")
	return &auth.ValidateResponse{
		UserId: userID,
	}, nil
}
