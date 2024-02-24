package api

import (
	"github.com/gofiber/fiber/v2"
	"github.com/jackc/pgx/v5/pgxpool"
)

type App struct {
	Pool *pgxpool.Pool
}

func NewServer(pool *pgxpool.Pool) *fiber.App {
	app := fiber.New(fiber.Config{
		Prefork:       true,
		CaseSensitive: true,
	})

	a := &App{Pool: pool}

	app.Get("/clientes/:id/extrato", a.extrato)
	app.Post("/clientes/:id/transacoes", a.transacoes)

	return app
}
