package app

import (
	"bill-split/internal/config"
	"bill-split/internal/domain/service"
	grpcService "bill-split/internal/domain/service/grpc"
	"bill-split/internal/handler"
	"bill-split/internal/repository"
	"bill-split/middleware"
	proto "bill-split/proto/this"
	"github.com/gin-gonic/gin"
	"google.golang.org/grpc"
	"log"
	"net"
	"os"
)

func Start() error {
	dbpool := config.NewPostgres(config.InitDb())

	defer dbpool.DbClose()

	// User
	userRepo := repository.NewUserRepository(dbpool.GetPGXPool())
	userService := service.NewUserHttpService(userRepo)
	userGrpcService := grpcService.NewUserService(userRepo)

	// Group
	groupRepo := repository.NewGroupRepository(dbpool.GetPGXPool())
	groupService := service.NewGroupService(groupRepo)

	// auth
	authService := service.NewAuthService(userRepo)

	// optimization
	optimizationService := service.NewOptimizationService()

	// cost
	costRepository := repository.NewCostRepository(dbpool.GetPGXPool())
	costService := service.NewCostService(costRepository, groupService)

	//handlers
	userHandler := handler.NewUserHandler(userService)
	groupHandler := handler.NewGroupHandler(groupService)
	optimizationHandler := handler.NewOptimizationHandler(optimizationService)
	costHandler := handler.NewCostHandler(costService)

	handlers := handler.NewHandlers(
		authService,
		groupHandler,
		optimizationHandler,
		userHandler,
		costHandler,
	)
	r := handlers.InitRoutes()

	// Добавляем CORS middleware (без использования дополнительных пакетов)
	r.Use(func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, PATCH, DELETE")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, Origin")
		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
		c.Writer.Header().Set("Access-Control-Expose-Headers", "Content-Length")

		// Handle preflight requests
		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	})

	r.POST("/register", handlers.RegisterUser)
	r.POST("/auth", handlers.Auth)

	api := r.Group("/api", middleware.AuthMiddleware())
	{
		// Оптимизация
		api.POST("/optimize", optimizationHandler.Optimize)
		// Пользователь
		apiUser := api.Group("/user")
		{
			apiUser.GET("/get/groups", handlers.GroupHandler.GetUserGroups)
			apiUser.GET("/get", handlers.UserHandler.GetUserData)
			apiUser.PATCH("/update", handlers.UserHandler.UpdateUserData)
		}

		// Группа
		apiGroup := api.Group("/group")
		{
			apiGroup.GET("/get/members", handlers.GroupHandler.GetUsersInGroup)
			apiGroup.GET("/:id", handlers.GroupHandler.GetGroupWithMembers)
			apiGroup.POST("/create", handlers.GroupHandler.CreateGroup)
			apiGroup.POST("/add/member", handlers.GroupHandler.AddMember)
			apiGroup.POST("/leave/member", handlers.GroupHandler.LeaveGroup)
			apiGroup.POST("/remove/member", handlers.GroupHandler.RemoveUser)
			apiGroup.POST("/enter/:link", handlers.GroupHandler.AddUserToGroup)

			// Расходы(Cost)
			apiCost := apiGroup.Group("/cost")
			{
				apiCost.POST("/create", handlers.CostHandler.CreateCost)
				apiCost.GET("/my", handlers.CostHandler.GetMyCosts)
				apiCost.GET("/group/:groupId", handlers.CostHandler.GetGroupCosts)
				apiCost.GET("/:id", handlers.CostHandler.GetCost)
				apiCost.PUT("/:id", handlers.CostHandler.UpdateCost)
				apiCost.DELETE("/:id", handlers.CostHandler.DeleteCost)
			}
		}
	}

	// Запускаем HTTP сервер в горутине
	go func() {
		log.Println("HTTP Server started on port 8080")
		if err := r.Run("0.0.0.0:" + os.Getenv("HTTP_PORT")); err != nil {
			log.Fatal(err)
		}
	}()

	// Запускаем gRPC сервер
	grpcServer := grpc.NewServer()
	proto.RegisterUserServiceServer(grpcServer, userGrpcService)

	lis, err := net.Listen("tcp", ":"+os.Getenv("GRPC_PORT"))
	if err != nil {
		log.Fatal(err)
	}

	log.Println("gRPC Server started on port 30000")
	if err := grpcServer.Serve(lis); err != nil {
		log.Fatal(err)
	}

	return nil
}
