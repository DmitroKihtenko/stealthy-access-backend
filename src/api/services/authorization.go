package services

import (
	"access-backend/api"
	"access-backend/base"
	"net/http"
)

type BaseAuthorizationService interface {
	ParseToken(tokenString string) (*api.AdminUser, error)
}

type AuthService struct {
	BaseAuthorizationService
	AuthConfig *base.AuthorizationConfig
}

func (service AuthService) ParseToken(tokenString string) (*api.AdminUser, error) {
	if tokenString != service.AuthConfig.AccessToken {
		return nil, base.ServiceError{
			Summary: "Invalid token",
			Status:  http.StatusForbidden,
		}
	}

	return &api.AdminUser{
		Username: "admin",
	}, nil
}
