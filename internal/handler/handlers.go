package handler

import (
	"bill-split/internal/domain/service"
	"github.com/gin-gonic/gin"
)

type Handlers struct {
	authorization service.AuthService
}

func NewHandlers(authService service.AuthService) *Handlers {
	return &Handlers{
		authorization: authService,
	}
}

func (h *Handlers) InitRoutes() *gin.Engine {
	r := gin.Default()
	return r
}
