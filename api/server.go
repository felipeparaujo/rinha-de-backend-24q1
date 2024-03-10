package api

import (
	"database/sql"
	"fmt"
	"math/rand"
	"strconv"
	"time"

	"github.com/gofiber/fiber/v2"
	_ "github.com/mattn/go-sqlite3"
	"github.com/valyala/fasthttp/pprofhandler"
)

type App struct {
	Port            int
	DBs             []*sql.DB
	PreparedQueries []PreparedQueries // One constructed from each DB.
}

var maxRetries = 100

func (a *App) Listen() error {
	for i := 1; i <= 5; i++ {
		file_name := "/db/" + strconv.Itoa(i) + ".db"
		db, err := sql.Open("sqlite3", "file:"+file_name+"?_journal=wal&_synchronous=off&_txlock=exclusive")
		if err != nil {
			panic(err)
		}

		for db.Ping() != nil && maxRetries > 0 {
			maxRetries--
			time.Sleep(time.Duration(rand.Intn(250)) * time.Millisecond)
		}

		p, err := PrepareQueries(db)
		if err != nil {
			panic(err)
		}

		a.DBs = append(a.DBs, db)
		a.PreparedQueries = append(a.PreparedQueries, p)
	}

	srv := fiber.New(fiber.Config{
		CaseSensitive: true,
	})

	srv.Get("/clientes/:id/extrato", a.extrato)
	srv.Post("/clientes/:id/transacoes", a.transacoes)
	srv.Get("/debug/pprof/:profile?", func(c *fiber.Ctx) error { pprofhandler.PprofHandler(c.Context()); return nil })

	return srv.Listen(fmt.Sprintf(":%d", a.Port))
}
