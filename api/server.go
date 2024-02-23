package api

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"log/slog"
	"net/http"
	"os"
	"sync"
	"time"

	"github.com/jackc/pgx/v5"
)

const (
	defaultIdleTimeout    = time.Minute
	defaultReadTimeout    = 5 * time.Second
	defaultWriteTimeout   = 10 * time.Second
	defaultShutdownPeriod = 30 * time.Second
)

type App struct {
	Conn *pgx.Conn
	Ctx  context.Context
	wg   sync.WaitGroup
}

func (a *App) ServeHTTP() error {
	srv := &http.Server{
		Addr:         fmt.Sprintf(":%s", os.Getenv("PORT")),
		Handler:      a.routes(),
		IdleTimeout:  defaultIdleTimeout,
		ReadTimeout:  defaultReadTimeout,
		WriteTimeout: defaultWriteTimeout,
	}

	log.Print("starting server", slog.Group("server", "addr", srv.Addr))
	return srv.ListenAndServe()
}

func JSON(w http.ResponseWriter, data any) error {
	js, err := json.MarshalIndent(data, "", "\t")
	if err != nil {
		return err
	}

	js = append(js, '\n')

	w.Header().Set("Content-Type", "application/json")
	w.Write(js)

	return nil
}
