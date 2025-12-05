package main

import (
	"bill-split/internal/config"
	"bill-split/internal/handlers"
)

func main() {
	dbpool := config.NewInterfaces()

	defer dbpool.DbClose()

	handlers := handlers.NewHandlers(dbpool)
	handlers.InitRoutes()

}
