package app

import (
	"bill-split/internal/config"
	grpc2 "bill-split/internal/domain/service/grpc"
	"bill-split/internal/domain/service/http"
	"bill-split/internal/handler"
	"bill-split/internal/repository"
	proto "bill-split/proto/this"
	"log"
	"net"

	"google.golang.org/grpc"
)

func Start() error {
	conn, connStr := config.InitDb()
	dbpool := config.New(conn, connStr)

	defer dbpool.DbClose()

	// gRPC
	userRepo := repository.NewUserRepository(dbpool.GetSql())
	userService := grpc2.NewUserService(userRepo)

	// http
	userHttpService := http.NewUserHttpService(userRepo)

	handlers := handler.NewHandlers(userHttpService)
	handlers.InitRoutes()

	grpcServer := grpc.NewServer()

	proto.RegisterUserServiceServer(grpcServer, userService)

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
