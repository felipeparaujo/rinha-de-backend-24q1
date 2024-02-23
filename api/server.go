package main

import (
	"fmt"
	"log"
	"log/slog"
	"net/http"
	"os"
	"strings"
	"time"
)

const (
	defaultIdleTimeout    = time.Minute
	defaultReadTimeout    = 5 * time.Second
	defaultWriteTimeout   = 10 * time.Second
	defaultShutdownPeriod = 30 * time.Second
)

func (a *app) serveHTTP() error {
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

func (a *app) routes() http.Handler {
	mux := http.NewServeMux()

	mux.HandleFunc("GET /clientes/{id}/extrato", a.extrato)
	mux.HandleFunc("POST /clientes/{id}/transacoes", a.transacoes)

	return mux
}

func (a *app) errorMessage(w http.ResponseWriter, r *http.Request, status int, message string) {
	message = strings.ToUpper(message[:1]) + message[1:]

	err := JSON(w, status, map[string]string{"error": message})
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
	}
}
