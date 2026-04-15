package handler

import (
	"github.com/gin-gonic/gin"
)

type Handlers struct {
	AuthHandler         AuthHandler
	GroupHandler        GroupHandler
	OptimizationHandler OptimizationHandler
	UserHandler         UserHandler
	CostHandler         CostHandler
}

func NewHandlers(
	authHandler AuthHandler, // Авторизация
	groupHandler GroupHandler, // Группа/Событие
	optimizationHandler OptimizationHandler,
	userHandler UserHandler,
	CostHandler CostHandler,
) *Handlers {
	return &Handlers{
		AuthHandler:         authHandler,
		GroupHandler:        groupHandler,
		OptimizationHandler: optimizationHandler,
		UserHandler:         userHandler,
		CostHandler:         CostHandler,
	}
}

func (h *Handlers) InitRoutes() *gin.Engine {
	r := gin.Default()
	return r
}
