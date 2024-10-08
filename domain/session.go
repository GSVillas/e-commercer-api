package domain

import (
	"context"
	"errors"

	"github.com/google/uuid"
)

var (
	ErrTokenInvalid           = errors.New("invalid token")
	ErrSessionNotFound        = errors.New("token not found")
	ErrorUnexpectedMethod     = errors.New("unexpected signing method")
	ErrTokenNotFoundInContext = errors.New("token not found in context")
	ErrOTPNotFound            = errors.New("OTP not found")
	ErrOTPInvalid             = errors.New("OTP expires")
)

type Session struct {
	Token     string
	Name      string
	UserID    uuid.UUID
	Email     string
	AvatarURL string
}

type SessionResponse struct {
	Token string `json:"token"`
}

type SessionService interface {
	Create(ctx context.Context, user User) (string, error)
	GetUser(ctx context.Context, token string) (*Session, error)
	Update(ctx context.Context) error
	SaveOTP(ctx context.Context, email string, otp string) error
	GetOTP(ctx context.Context, email string) (string, error)
}

type SessionRepository interface {
	Create(ctx context.Context, user User, token string) error
	GetUser(ctx context.Context, userID string) (*Session, error)
	Update(ctx context.Context, user User, token string) error
	SaveOTP(ctx context.Context, email string, otp string) error
	GetOTP(ctx context.Context, email string) (string, error)
}

func (u *Session) ToResponse() *UserResponse {
	return &UserResponse{
		ID:        u.UserID.String(),
		Name:      u.Name,
		Email:     u.Email,
		AvatarURL: u.AvatarURL,
	}
}
