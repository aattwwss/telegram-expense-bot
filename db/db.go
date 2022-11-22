package db

import (
	"context"
	"fmt"
	"github.com/aattwwss/telegram-expense-bot/config"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Config struct {
	Username string
	Password string
	Host     string
	Port     string
	Database string
	Schema   string
}

func newConfig(cfg config.Config) Config {
	return Config{
		Username: cfg.DbUsername,
		Password: cfg.DbPassword,
		Host:     cfg.DbHost,
		Port:     cfg.DbPort,
		Database: cfg.DbDatabase,
		Schema:   cfg.DbSchema,
	}
}

func (c Config) connectionUrl() string {
	return fmt.Sprintf("postgresql://%s:%s@%s:%s/%s?search_path=%s", c.Username, c.Password, c.Host, c.Port, c.Database, c.Schema)
}

func LoadDB(ctx context.Context, cfg config.Config) (*pgxpool.Pool, error) {
	c := newConfig(cfg)
	return pgxpool.New(ctx, c.connectionUrl())
}
