package handler

import (
	"bill-split/internal/domain/service"
	"github.com/gin-gonic/gin"
)

type Handlers struct {
	authorization service.AuthService
	groupService  service.GroupService
}

func NewHandlers(
	authService service.AuthService, // Авторизация
	groupService service.GroupService, // Группа/Событие
) *Handlers {
	return &Handlers{
		authorization: authService,
		groupService:  groupService,
	}
}

func (h *Handlers) InitRoutes() *gin.Engine {
	r := gin.Default()
	return r
}
