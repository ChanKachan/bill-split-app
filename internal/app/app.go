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
	"github.com/gorilla/websocket"
	"github.com/redis/go-redis/v9"

	"github.com/gin-gonic/gin"
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
		websocket.Upgrader{
			ReadBufferSize:  1024,
			WriteBufferSize: 1024,
		},
		configWS,
		chatServ,
	)

	handlers := handler.NewHandlers(
		authHandler,
		groupHandler,
		optimizationHandler,
		userHandler,
		costHandler,
	)
	r := handlers.InitRoutes()

	// Добавляем CORS middleware (без использования дополнительных пакетов)
	r.Use(func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", os.Getenv("CORS_ALLOWED_ORIGINS"))
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

	r.POST("/register", handlers.AuthHandler.RegisterUser)
	r.POST("/auth", handlers.AuthHandler.Auth)
	r.GET("/ws", chatHandler.ConnectionWS)

	api := r.Group("/api", middleware.AuthMiddleware())
	{
		// chat
		//api.POST("/ws", chatHandler.ConnectionWS)
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
