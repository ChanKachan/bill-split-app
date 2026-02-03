package handler

import (
	"bill-split/internal/domain/service/http"

	"github.com/gin-gonic/gin"
)

type HandlersInterface interface{}

type Handlers struct {
	authorization Authorization
	user          http.UserHttpService
}

func NewHandlers(
	userHttpService http.UserHttpService,
) *Handlers {
	return &Handlers{
		//authorization:,
		user: userHttpService,
	}
}

func (h *Handlers) InitRoutes() *gin.Engine {
	r := gin.Default()

	r.POST("/user/create", h.user.Create)

	r.Run("0.0.0.0:8080")

	return r
}
