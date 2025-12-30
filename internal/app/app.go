package app

import (
	"bill-split/internal/config"
	"bill-split/internal/handler"
)

func Start() error {
	dbpool := config.NewInterfaces()

	defer dbpool.DbClose()

	handlers := handler.NewHandlers(dbpool)
	handlers.InitRoutes()

	return nil
}
