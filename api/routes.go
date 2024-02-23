package api

import "net/http"

func (a *App) routes() http.Handler {
	mux := http.NewServeMux()

	mux.HandleFunc("GET /clientes/{id}/extrato", a.extrato)
	mux.HandleFunc("POST /clientes/{id}/transacoes", a.transacoes)

	return mux
}
