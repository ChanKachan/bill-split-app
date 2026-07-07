package main

import (
	"github.com/ChanKachan/bill-split-app/internal/app"
	"log"
)

func main() {
	if err := app.Start(); err != nil {
		log.Println(err)
	}
}
