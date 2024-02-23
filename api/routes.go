package api

import "net/http"

type handler func(w http.ResponseWriter, r *http.Request) int

func wrap(h handler) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		status := h(w, r)
		if status != http.StatusOK {
			w.WriteHeader(status)
		}
	}
}

func (a *App) routes() http.Handler {
	mux := http.NewServeMux()

	mux.HandleFunc("GET /clientes/{id}/extrato", wrap(a.extrato))
	mux.HandleFunc("POST /clientes/{id}/transacoes", wrap(a.transacoes))

	return mux
}
