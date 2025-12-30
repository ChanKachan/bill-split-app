package main

import (
	"bill-split/internal/app"
	"log"
)

func main() {
	if err := app.Start(); err != nil {
		log.Println(err)
	}
}
