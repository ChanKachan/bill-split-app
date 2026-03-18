package app

import (
	"bill-split/internal/config"
	"bill-split/internal/domain/service"
	grpcService "bill-split/internal/domain/service/grpc"
	"bill-split/internal/handler"
	"bill-split/internal/repository"
	"bill-split/middleware"
	proto "bill-split/proto/this"
	"google.golang.org/grpc"
	"log"
	"net"
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

	//handlers
	userHandler := handler.NewUserHandler(userService)                         // Пользователь
	groupHandler := handler.NewGroupHandler(groupService)                      // Группа
	optimizationHandler := handler.NewOptimizationHandler(optimizationService) // Оптимизация handler

	handlers := handler.NewHandlers(
		authService,
		groupHandler,
		optimizationHandler,
		userHandler,
	)
	r := handlers.InitRoutes()

	r.POST("/register", handlers.RegisterUser)
	r.POST("/auth", handlers.Auth)

	api := r.Group("/api", middleware.AuthMiddleware())
	{
		// Пользователь
		apiUser := api.Group("/user")
		{
			// GET
			apiUser.GET("/get/groups", handlers.GroupHandler.GetUserGroups)

			// PATCH
			apiUser.PATCH("/update", handlers.UserHandler.UpdateUserData)
		}

		// Группа
		apiGroup := api.Group("/group")
		{
			// GET
			apiGroup.GET("/get/members", handlers.GroupHandler.GetUsersInGroup)

			// POST
			apiGroup.POST("/create", handlers.GroupHandler.CreateGroup)
		}
	}

	r.Run("0.0.0.0:8080")

	grpcServer := grpc.NewServer()

	proto.RegisterUserServiceServer(grpcServer, userGrpcService)

	// Запускаем на порту 30000
	lis, err := net.Listen("tcp", ":30000")
	if err != nil {
		log.Fatal(err)
	}

	log.Println("gRPC Server started on port 30000")
	if err := grpcServer.Serve(lis); err != nil {
		log.Fatal(err)
	}

	return nil
}
