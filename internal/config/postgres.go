package config

import (
	"context"
	"fmt"
	"log"
	"os"
	"runtime"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/joho/godotenv"
)

type Postgres interface {
	NewPoolConfig(maxConn int, connIdleTime, connLifeTime time.Duration) error
	GetSql() *pgxpool.Pool
	DbClose()
}

type postgres struct {
	conn         *pgxpool.Pool
	connStr      string
	poolConfig   *pgxpool.Config
	queryTimeout time.Duration
}

func New(db *pgxpool.Pool, connStr string) Postgres {
	return &postgres{
		conn:    db,
		connStr: connStr,
	}
}

func (p *postgres) NewPoolConfig(maxConn int, connIdleTime, connLifeTime time.Duration) error {
	poolConfig, err := pgxpool.ParseConfig(p.connStr)
	if err != nil {
		return err
	}

	// Смотрим, чтобы соединений к бд было меньше чем ядер
	cpu := runtime.NumCPU()
	if maxConn < cpu {
		maxConn = cpu
	}

	poolConfig.MaxConns = int32(maxConn)
	poolConfig.MaxConnIdleTime = connIdleTime
	poolConfig.MaxConnLifetime = connLifeTime
	p.poolConfig = poolConfig
	return nil
}

func (p *postgres) DbClose() {
	p.conn.Close()
}

func (p *postgres) GetSql() *pgxpool.Pool {
	return p.conn
}

func (p *postgres) Ping(ctx context.Context) error {
	return p.conn.Ping(ctx)
}

func (p *postgres) Exec(ctx context.Context, sql string, arguments ...any) (pgconn.CommandTag, error) {
	return p.conn.Exec(ctx, sql, arguments...)
}

func (p *postgres) Query(ctx context.Context, sql string, arguments ...any) (pgx.Rows, error) {
	return p.conn.Query(ctx, sql, arguments...)
}

func (p *postgres) QueryRow(ctxParent context.Context, sql string, arguments ...any) pgx.Row {
	ctx, cancel := context.WithTimeout(ctxParent, p.queryTimeout)
	defer cancel()
	return p.conn.QueryRow(ctx, sql, arguments...)
}

func InitDb() (*pgxpool.Pool, string) {
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

	log.Println("Подключение к базе данных прошла успешно.")

	return dbpool, connStr
}
