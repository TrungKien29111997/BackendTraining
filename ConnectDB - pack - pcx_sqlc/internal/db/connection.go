package db

import (
	"context"
	"fmt"
	"hoc-gin/internal/config"
	"log"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

var DB string

func InitDB() error {
	connStr := config.NewConfig().DNS()
	conf, err := pgxpool.ParseConfig(connStr)
	if err != nil {
		return fmt.Errorf("error parse config: %w", err)
	}

	conf.MaxConns = 50
	conf.MinConns = 5
	conf.MaxConnLifetime = 30 * time.Minute
	conf.MaxConnIdleTime = 5 * time.Minute
	conf.HealthCheckPeriod = 1 * time.Minute

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	DBPool, err := pgxpool.NewWithConfig(ctx, conf)
	if err != nil {
		return fmt.Errorf("unable to connect to database: %w", err)
	}
	if err := DBPool.Ping(ctx); err != nil {
		return fmt.Errorf("unable to ping database: %w", err)
	}
	log.Println("Connected")
	return nil
}
