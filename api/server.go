package api

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/jackc/pgx/v5/pgxpool"
)

type App struct {
	Port      int
	DBConnStr string
	pool      *pgxpool.Pool
}

func (a *App) Listen() error {
	ctx := context.Background()

	config, err := pgxpool.ParseConfig(a.DBConnStr)
	if err != nil {
		log.Fatal(err)
	}

	pool, err := pgxpool.NewWithConfig(ctx, config)
	if err != nil {
		return err
	}

	defer pool.Close()

	for pool.Ping(ctx) != nil {
		time.Sleep(250 * time.Millisecond)
	}

	a.pool = pool

	srv := fiber.New(fiber.Config{
		// Prefork:       true,
		CaseSensitive: true,
	})

	srv.Get("/clientes/:id/extrato", a.extrato)
	srv.Post("/clientes/:id/transacoes", a.transacoes)

	return srv.Listen(fmt.Sprintf(":%d", a.Port))
}
