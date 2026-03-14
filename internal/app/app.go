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
	//userService := service.NewUserHttpService(userRepo)
	userGrpcService := grpcService.NewUserService(userRepo)

	// Group
	groupRepo := repository.NewGroupRepository(dbpool.GetPGXPool())
	groupService := service.NewGroupService(groupRepo)

	authService := service.NewAuthService(userRepo)

	handlers := handler.NewHandlers(authService, groupService)
	r := handlers.InitRoutes()

	r.POST("/register", handlers.RegisterUser)
	r.POST("/auth", handlers.Auth)

	// TODO: убрал  middleware.AuthMiddleware() на время тестов
	api := r.Group("/api", middleware.AuthMiddleware())
	{
		api.POST("/create/group", handlers.CreateGroup)
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
