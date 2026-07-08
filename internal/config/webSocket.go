package config

import (
	"log"
	"os"
	"strconv"
)

type CfgWebSocket struct {
	PongWaitSeconds int
}

func GetConfigWebSocket() *CfgWebSocket {
	waitSecond, err := strconv.Atoi(os.Getenv("PONG_WAIT_SECONDS"))
	if err != nil {
		log.Printf("Error convert data in config web socket: %w", err)
		return nil
	}
	return &CfgWebSocket{
		PongWaitSeconds: waitSecond,
	}
}
