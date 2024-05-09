package controllers

import (
	"access-backend/api/services"
	"access-backend/base"
	"github.com/gin-gonic/gin"
	"net/http"
	"strings"
)

type AuthController struct {
	AuthService services.BaseAuthorizationService
}

func (controller AuthController) Authorize(context *gin.Context) {
	tokenString := context.GetHeader("Authorization")
	if tokenString == "" {
		context.Error(base.ServiceError{
			Summary: "Authorization token required",
			Status:  http.StatusUnauthorized,
		})
		context.Abort()
		return
	}
	if strings.HasPrefix(tokenString, "Bearer ") {
		user, err := controller.AuthService.ParseToken(
			tokenString[7:],
		)
		if err != nil {
			context.Error(err)
			context.Abort()
			return
		}
		context.Set("auth", user)
		context.Next()
	} else {
		context.Error(base.ServiceError{
			Summary: "Invalid token type",
			Status:  http.StatusForbidden,
		})
		context.Abort()
		return
	}
}
