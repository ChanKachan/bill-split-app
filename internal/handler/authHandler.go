package handler

import (
	"bill-split/internal/config"
	"bill-split/internal/domain/service"
)

type Authorization interface {
}

type AuthorizationHandler struct {
	Dbpg        config.Postgres
	authService service.AuthService
}

func NewAuthorizationHandler(dbpg config.Postgres) Authorization {
	return &AuthorizationHandler{
		Dbpg:        dbpg,
		authService: service.NewAuthService(),
	}
}
