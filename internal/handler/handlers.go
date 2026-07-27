package handler

import (
	"github.com/ChanKachan/bill-split-app/internal/handler/ws/chat"
	"github.com/ChanKachan/bill-split-app/middleware"
	"github.com/gin-gonic/gin"
)

type Handlers struct {
	authHandler         AuthHandler
	groupHandler        GroupHandler
	optimizationHandler OptimizationHandler
	userHandler         UserHandler
	costHandler         CostHandler
	chatHandler         chat.ChatHandler
}

func NewHandlers(
	authHandler AuthHandler, // Авторизация
	groupHandler GroupHandler, // Группа/Событие
	optimizationHandler OptimizationHandler,
	userHandler UserHandler,
	costHandler CostHandler,
	chatHandler chat.ChatHandler,
) *Handlers {
	return &Handlers{
		authHandler:         authHandler,
		groupHandler:        groupHandler,
		optimizationHandler: optimizationHandler,
		userHandler:         userHandler,
		costHandler:         costHandler,
		chatHandler:         chatHandler,
	}
}

func (h *Handlers) InitRoutes(
	cors middleware.Cors,
	auth middleware.Middleware,
) *gin.Engine {
	r := gin.Default()

	r.Use(cors.Init())

	r.POST("/register", h.authHandler.RegisterUser)
	r.POST("/auth", h.authHandler.Auth)

	api := r.Group("/api", auth.AuthMiddleware())
	{
		// chat
		api.GET("/chats/:chatID/ws", h.chatHandler.ConnectionWS)
		// Оптимизация
		api.POST("/optimize", h.optimizationHandler.Optimize)
		// Пользователь
		apiUser := api.Group("/user")
		{
			apiUser.GET("/get/groups", h.groupHandler.GetUserGroups)
			apiUser.GET("/get", h.userHandler.GetUserData)
			apiUser.PATCH("/update", h.userHandler.UpdateUserData)
		}

		// Группа
		apiGroup := api.Group("/group")
		{
			apiGroup.GET("/get/members", h.groupHandler.GetUsersInGroup)
			apiGroup.GET("/:id", h.groupHandler.GetGroupWithMembers)
			apiGroup.POST("/create", h.groupHandler.CreateGroup)
			apiGroup.POST("/add/member", h.groupHandler.AddMember)
			apiGroup.POST("/leave/member", h.groupHandler.LeaveGroup)
			apiGroup.POST("/remove/member", h.groupHandler.RemoveUser)
			apiGroup.POST("/enter/:link", h.groupHandler.AddUserToGroup)

			// Расходы(Cost)
			apiCost := apiGroup.Group("/cost")
			{
				apiCost.POST("/create", h.costHandler.CreateCost)
				apiCost.GET("/my", h.costHandler.GetMyCosts)
				apiCost.GET("/group/:groupId", h.costHandler.GetGroupCosts)
				apiCost.GET("/:id", h.costHandler.GetCost)
				apiCost.PUT("/:id", h.costHandler.UpdateCost)
				apiCost.DELETE("/:id", h.costHandler.DeleteCost)
			}
		}
	}

	return r
}
