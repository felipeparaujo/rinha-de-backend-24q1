package main

import (
	"context"
	"log"
	"time"

	"github.com/felipeparaujo/rinha-de-backend-24q1/api"
	"github.com/jackc/pgx/v5"
)

const maxRetries = 5

func main() {
	ctx := context.Background()

	config, err := pgx.ParseConfig("postgres://admin:password@db:5432/rinha")
	if err != nil {
		log.Fatal(err)
	}

	retryCount := 0
	conn, err := pgx.ConnectConfig(ctx, config)
	for err != nil && retryCount < maxRetries {
		conn, err = pgx.ConnectConfig(ctx, config)
		time.Sleep(1 * time.Second)
		retryCount++
	}

	if err != nil {
		log.Fatal(err)
	}

	a := &api.App{Conn: conn, Ctx: ctx}
	if err := a.ServeHTTP(); err != nil {
		log.Fatal(err)
	}
}
