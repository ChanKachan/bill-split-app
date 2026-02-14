package handler

import (
	"bill-split/internal/config"
	"bill-split/internal/domain/service"
)

type Authorization interface {
}

type authorizationHandler struct {
	Dbpg        config.Postgres
	authService service.AuthService
}

func NewAuthorizationHandler(dbpg config.Postgres, authService service.AuthService) Authorization {
	return &authorizationHandler{
		Dbpg:        dbpg,
		authService: authService,
	}
}
