package main

import (
	"context"
	"log"
	"time"

	"github.com/felipeparaujo/rinha-de-backend-24q1/api"
	"github.com/jackc/pgx/v5/pgxpool"
)

const maxRetries = 100

func main() {
	ctx := context.Background()

	config, err := pgxpool.ParseConfig("postgres://admin:password@db:5432/rinha")
	if err != nil {
		log.Fatal(err)
	}

	pool, err := pgxpool.NewWithConfig(ctx, config)
	for pool.Ping(ctx) != nil {
		time.Sleep(250 * time.Millisecond)
	}

	if err != nil {
		log.Fatal(err)
	}

	a := &api.App{Ctx: ctx, Pool: pool}
	if err := a.ServeHTTP(); err != nil {
		log.Fatal(err)
	}
}
