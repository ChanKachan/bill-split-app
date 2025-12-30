package handler

import (
	"bill-split/internal/config"

	"github.com/gin-gonic/gin"
)

type HandlersInterface interface{}

type Handlers struct {
	authorization Authorization
}

func NewHandlers(dbpool config.Postgres) *Handlers {
	return &Handlers{
		authorization: NewAuthorizationHandler(dbpool),
	}
}

func (h *Handlers) InitRoutes() *gin.Engine {
	r := gin.Default()

	r.Run("0.0.0.0:8080")

	return r
}
