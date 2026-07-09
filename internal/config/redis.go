package config

import (
	"log"
	"os"
	"strconv"
)

type CfgRedis struct {
	Port     int
	Host     string
	Password string
	DB       int
}

func InitRedis() *CfgRedis {
	port, err := strconv.Atoi(os.Getenv("REDIS_PORT"))
	if err != nil {
		log.Printf("InitRedis | error convert string to int: %v", err)
		return nil
	}

	return &CfgRedis{
		Port:     port,
		Password: os.Getenv("REDIS_PASSWORD"),
		Host:     os.Getenv("REDIS_HOST"),
	}
}
