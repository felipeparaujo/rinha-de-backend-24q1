package main

import (
	"context"
	"log"
	"sync"
	"time"

	"github.com/jackc/pgx/v5"
)

type Cliente struct {
	ID           int    `json:"id"`
	Nome         string `json:"nome"`
	Limite       int32  `json:"limite"`
	SaldoInicial int32  `json:"saldo_inicial"`
}

type app struct {
	conn *pgx.Conn
	wg   sync.WaitGroup
	ctx  context.Context
}

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

	a := &app{conn: conn, ctx: ctx}
	if err := a.serveHTTP(); err != nil {
		log.Fatal(err)
	}
}
