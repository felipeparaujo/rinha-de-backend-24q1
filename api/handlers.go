package main

import (
	"fmt"
	"net/http"
)

func (a *app) extrato(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "extrato")
}

func (a *app) transacoes(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, "transacoes")
}
