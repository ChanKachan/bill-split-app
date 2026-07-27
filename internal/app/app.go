package app

import (
	"fmt"
	"log"
	"net"
	"os"
	"sync"

	"github.com/ChanKachan/bill-split-app/internal/config"
	"github.com/ChanKachan/bill-split-app/internal/domain/repository/postgres"
	chatRepository "github.com/ChanKachan/bill-split-app/internal/domain/repository/postgres/chat"
	"github.com/ChanKachan/bill-split-app/internal/domain/repository/redis/cache"
	"github.com/ChanKachan/bill-split-app/internal/domain/service"
	chatService "github.com/ChanKachan/bill-split-app/internal/domain/service/chat"
	grpcService "github.com/ChanKachan/bill-split-app/internal/domain/service/grpc"
	"github.com/ChanKachan/bill-split-app/internal/handler"
	"github.com/ChanKachan/bill-split-app/internal/handler/ws/chat"
	"github.com/ChanKachan/bill-split-app/middleware"
	proto "github.com/ChanKachan/bill-split-app/proto/this"
	"github.com/ChanKachan/bill-split-app/repository"
	"github.com/redis/go-redis/v9"

	"google.golang.org/grpc"
)

func Start() error {
	var wg sync.WaitGroup

	// configs
	dbpool := config.NewPostgres(config.InitDb())
	configWS := config.GetConfigWebSocket()
	cfgRedis := config.InitRedis()

	defer dbpool.DbClose()

	chatRedisDB := repository.NewRedisDB(
		&redis.Options{
			Password: cfgRedis.Password,
			DB:       1,
			Addr:     fmt.Sprintf("%s:%d", cfgRedis.Host, cfgRedis.Port),
		},
	)

	// Chat
	chatCache := cache.NewChatCache(chatRedisDB)
	chatRepo := chatRepository.NewChatRepository(dbpool.GetPGXPool())
	chatServ := chatService.NewChatService(chatCache, chatRepo)

	// User
	userRepo := postgres.NewUserRepository(dbpool.GetPGXPool())
	userService := service.NewUserHttpService(userRepo)
	userGrpcService := grpcService.NewUserService(userRepo)

	// Group
	groupRepo := postgres.NewGroupRepository(dbpool.GetPGXPool())
	groupService := service.NewGroupService(groupRepo)

	// auth
	authService := service.NewAuthService(userRepo)
	authHandler := handler.NewAuthHandler(authService)

	// optimization
	optimizationService := service.NewOptimizationService()

	// cost
	costRepository := postgres.NewCostRepository(dbpool.GetPGXPool())
	costService := service.NewCostService(costRepository, groupService)

	//handlers
	userHandler := handler.NewUserHandler(userService)
	groupHandler := handler.NewGroupHandler(groupService)
	optimizationHandler := handler.NewOptimizationHandler(optimizationService)
	costHandler := handler.NewCostHandler(costService)
	chatHandler := chat.NewChatHandler(
		chat.NewUpgrader(
			1024,
			1024,
		),
		configWS,
		chatServ,
	)

	handlers := handler.NewHandlers(
		authHandler,
		groupHandler,
		optimizationHandler,
		userHandler,
		costHandler,
		chatHandler,
	)

	cors := middleware.NewCors()
	auth := middleware.NewMiddleware()
	r := handlers.InitRoutes(cors, auth)

	// Запускаем HTTP сервер в горутине
	wg.Add(1)
	go func() {
		log.Printf("HTTP Server started on port %s", os.Getenv("HTTP_PORT"))
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

	wg.Wait()

	return nil
}
