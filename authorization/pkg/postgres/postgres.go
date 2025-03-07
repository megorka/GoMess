package postgres

import (
	"context"
	"fmt"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Config struct {
	Host     string `yaml:"POSTGRES_HOST" env:"POSTGRES_HOST" env-default:"localhost"`
	Port     int    `yaml:"POSTGRES_PORT" env:"POSTGRES_PORT" env-default:"5432"`
	Username string `yaml:"POSTGRES_USER" env:"POSTGRES_USER" env-default:"postgres"`
	Password string `yaml:"POSTGRES_PASS" env:"POSTGRES_PASS" env-default:"megorka"`
	Database string `yaml:"POSTGRES_DB" env:"POSTGRES_DB" env-default:"yandex"`
}

func New(ctx context.Context, config Config) (*pgxpool.Pool, error) {
	connString := fmt.Sprintf("postgres://%s:%s@%s:%d/%s?sslmode=disable",
		config.Username,
		config.Password,
		config.Host,
		config.Port,
		config.Database)

	conn, err := pgxpool.New(ctx, connString)
	if err != nil {
		return nil, fmt.Errorf("unable to connect to database: %w", err)
	}

	return conn, nil
}
