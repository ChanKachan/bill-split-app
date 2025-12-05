package handlers

import (
	"bill-split/internal/config"
	"bill-split/internal/domain/service"
)

type Authorization interface {
}

type AuthorizationHandler struct {
	Dbpg        *config.Interfaces
	authService service.AuthServiceInterface
}

func NewAuthorizationHandler(dbpg *config.Interfaces) Authorization {
	return &AuthorizationHandler{
		Dbpg:        dbpg,
		authService: service.NewAuthService(),
	}
}
