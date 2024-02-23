package api

import (
	"encoding/json"
	"log"
	"net/http"
	"time"
)

type ExtratoResponse struct {
	Saldo             ExtratSaldoResponse        `json:"saldo"`
	UltimasTransacoes []ExtratoTransacaoResponse `json:"ultimas_transacoes"`
}

type ExtratSaldoResponse struct {
	Total       int32     `json:"total"`
	DataExtrato time.Time `json:"data_extrato"`
	Limite      int32     `json:"limite"`
}

type ExtratoTransacaoResponse struct {
	Valor       int32     `json:"valor"`
	Tipo        string    `json:"tipo"`
	Descricao   string    `json:"descricao"`
	RealizadaEm time.Time `json:"realizada_em"`
}

func (a *App) extrato(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")

	rows, err := a.Conn.Query(
		a.Ctx, `
		SELECT
			c.saldo, c.limite, t.valor, t.tipo, t.descricao, t.realizada_em
		FROM clientes c JOIN (
			SELECT * FROM transacoes
			WHERE cliente_id = $1
			ORDER BY realizada_em DESC
			LIMIT 10
		) t ON c.id = t.cliente_id`,
		id,
	)

	if err != nil {
		log.Print(err)
		return
	}
	defer rows.Close()

	resp := ExtratoResponse{Saldo: ExtratSaldoResponse{DataExtrato: time.Now()}}
	for rows.Next() {
		transacao := ExtratoTransacaoResponse{}
		err := rows.Scan(&resp.Saldo.Total, &resp.Saldo.Limite, &transacao.Valor, &transacao.Tipo, &transacao.Descricao, &transacao.RealizadaEm)
		if err != nil {
			log.Print(err)
		}

		resp.UltimasTransacoes = append(resp.UltimasTransacoes, transacao)
	}

	if err := JSON(w, 200, resp); err != nil {
		log.Print(err)
		a.errorMessage(w, r, 500, "internal server error")
		return
	}
}

type TransacoesRequest struct {
	Valor     int32  `json:"valor"`
	Tipo      string `json:"tipo"`
	Descricao string `json:"descricao"`
}

type TransacoesResponse struct {
	Limite int32 `json:"limite"`
	Saldo  int32 `json:"saldo"`
}

func (a *App) transacoes(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	req := TransacoesRequest{}
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		log.Print(err)
		a.errorMessage(w, r, http.StatusBadRequest, "invalid request")
		return
	}

	valor := req.Valor
	if req.Tipo == "d" {
		valor = -valor
	}

	id := r.PathValue("id")
	rows, err := a.Conn.Query(a.Ctx, "SELECT * FROM process_transaction($1, $2, $3, $4)", id, req.Tipo, req.Descricao, valor)
	if err != nil {
		log.Print(err)
		a.errorMessage(w, r, 500, "internal server error")
		return
	}
	defer rows.Close()

	for rows.Next() {
		var resp TransacoesResponse
		err := rows.Scan(&resp.Saldo, &resp.Limite)
		if err != nil {
			log.Print(err)
		}

		if err := JSON(w, 200, resp); err != nil {
			log.Print(err)
			a.errorMessage(w, r, 500, "internal server error")
			return
		}

		return
	}

	a.errorMessage(w, r, 500, "internal server error")
}
