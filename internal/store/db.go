package store

import (
	"context"
	"fmt"
	//"os"	
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

var DB *pgxpool.Pool

func InitDB() error {
	dsn := "postgres://nomad:secret@localhost:5433/nomad"

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	config, err := pgxpool.ParseConfig(dsn)
	if err != nil {
		return fmt.Errorf("unable to parse config: %w", err)
	}

	DB, err = pgxpool.NewWithConfig(ctx, config)
	if err != nil {
		return fmt.Errorf("unable to connect to database: %w", err)
	}

	err = DB.Ping(ctx)
	if err != nil {
		return fmt.Errorf("database ping failed: %w", err)
	}

	fmt.Println("✅ Connected to Nomad PostgreSQL")
	return nil
}