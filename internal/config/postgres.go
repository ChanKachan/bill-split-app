package config

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/joho/godotenv"
)

type Postgres interface {
	GetPGXPool() *pgxpool.Pool
	DbClose()
}

type postgres struct {
	Dbpg *pgxpool.Pool
}

func NewPostgres(db *pgxpool.Pool) Postgres {
	return &postgres{
		Dbpg: db,
	}
}

func (i *postgres) DbClose() {
	i.Dbpg.Close()
}

func (i *postgres) GetPGXPool() *pgxpool.Pool {
	return i.Dbpg
}

func InitDb() *pgxpool.Pool {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	connStr := fmt.Sprintf("postgres://%s:%s@%s:%s/%s",
		os.Getenv("DB_USER"),
		os.Getenv("DB_PASSWORD"),
		os.Getenv("DB_HOST"),
		os.Getenv("DB_PORT"),
		os.Getenv("DB_NAME"),
	)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)

	defer cancel()

	dbpool, err := pgxpool.New(ctx, connStr)

	err = dbpool.Ping(ctx)

	if err != nil {
		log.Printf("Не удалось подключиться к базе данных: %v\n", err)
	}

	//Конвертируем в тип *sqlx.DB
	//sqlDB := stdlib.OpenDBFromPool(dbpool)
	//db := sqlx.NewDb(sqlDB, "postgres")

	log.Println("Подключение к базе данных прошла успешно.")

	return dbpool
}
