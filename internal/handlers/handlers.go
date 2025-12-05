package handlers

import (
	"bill-split/internal/config"

	"github.com/gin-gonic/gin"
)

type HandlersInterface interface{}

type Handlers struct {
	Authorization Authorization
}

func NewHandlers(dbpool *config.Interfaces) *Handlers {
	return &Handlers{
		Authorization: NewAuthorizationHandler(dbpool),
	}
}

func (h *Handlers) InitRoutes() *gin.Engine {
	r := gin.Default()

	r.Run("0.0.0.0:8080")

	return r
}
