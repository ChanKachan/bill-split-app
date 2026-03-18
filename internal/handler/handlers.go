package handler

import (
	"bill-split/internal/domain/service"
	"github.com/gin-gonic/gin"
)

type Handlers struct {
	authorization       service.AuthService
	GroupHandler        GroupHandler
	OptimizationHandler OptimizationHandler
	UserHandler         UserHandler
}

func NewHandlers(
	authService service.AuthService, // Авторизация ToDO: Нужно убрать сервис отсюда и переписать его в отдельный handler, где будет и регистрация и авторизация
	groupHandler GroupHandler, // Группа/Событие
	optimizationHandler OptimizationHandler,
	userHandler UserHandler,
) *Handlers {
	return &Handlers{
		authorization:       authService,
		GroupHandler:        groupHandler,
		OptimizationHandler: optimizationHandler,
		UserHandler:         userHandler,
	}
}

func (h *Handlers) InitRoutes() *gin.Engine {
	r := gin.Default()
	return r
}
