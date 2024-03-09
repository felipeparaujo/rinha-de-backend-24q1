package api

import (
	"context"
	"errors"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/gofiber/fiber/v2"
)

var ErrInvalidTransacoesRequest = errors.New("invalid request")

type ExtratoResponse struct {
	Saldo             ExtratoSaldoResponse       `json:"saldo"`
	UltimasTransacoes []ExtratoTransacaoResponse `json:"ultimas_transacoes"`
}

type ExtratoSaldoResponse struct {
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

func (a *App) extrato(c *fiber.Ctx) error {
	id, err := strconv.Atoi(c.Params("id"))
	if err != nil || id < 1 || id > 5 {
		return c.SendStatus(http.StatusNotFound)
	}

	resp := ExtratoResponse{Saldo: ExtratoSaldoResponse{DataExtrato: time.Now()}, UltimasTransacoes: []ExtratoTransacaoResponse{}}

	row := a.DBs[id-1].QueryRow(`SELECT l, b FROM u`)
	err = row.Scan(&resp.Saldo.Limite, &resp.Saldo.Total)
	if err != nil {
		log.Print(err)
		return c.SendStatus(http.StatusInternalServerError)
	}

	rows, err := a.DBs[id-1].Query(`SELECT t, a, d FROM t ORDER BY id DESC LIMIT 10`)
	if err != nil {
		log.Print(err)
		return c.SendStatus(http.StatusInternalServerError)
	}
	defer rows.Close()

	for rows.Next() {
		transacao := ExtratoTransacaoResponse{}
		err := rows.Scan(&transacao.RealizadaEm, &transacao.Valor, &transacao.Descricao)
		if err != nil {
			return c.SendStatus(http.StatusInternalServerError)
		}

		// debit transactions are stored as negative values, but users expect them to be positive
		transacao.Tipo = "c"
		if transacao.Valor < 0 {
			transacao.Valor *= -1
			transacao.Tipo = "d"
		}

		resp.UltimasTransacoes = append(resp.UltimasTransacoes, transacao)
	}
	resp.Saldo.Limite *= -1

	return c.Status(http.StatusOK).JSON(resp)
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

func (t *TransacoesRequest) validate() error {
	if t.Valor < 1 {
		return ErrInvalidTransacoesRequest
	}

	descLen := len(t.Descricao)
	if descLen < 1 || descLen > 10 {
		return ErrInvalidTransacoesRequest
	}

	if t.Tipo != "d" && t.Tipo != "c" {
		return ErrInvalidTransacoesRequest
	}

	return nil
}

func (a *App) transacoes(c *fiber.Ctx) error {
	req := TransacoesRequest{}
	if err := c.BodyParser(&req); err != nil {
		return c.SendStatus(http.StatusBadRequest)
	}

	id, err := strconv.Atoi(c.Params("id"))
	if err != nil || id < 1 || id > 5 {
		return c.SendStatus(http.StatusNotFound)
	}

	if err := req.validate(); err != nil {
		return c.SendStatus(http.StatusUnprocessableEntity)
	}

	resp := TransacoesResponse{}
	balanceDiff := req.Valor
	if req.Tipo == "d" {
		balanceDiff *= -1
	}
	db := a.DBs[id-1]
	for i := 0; i < 10; i++ {
		tx, err := db.BeginTx(context.Background(), nil)
		if err != nil {
			log.Print(err)
			continue
		}

		// Transaction started, there will be no more retries.
		defer tx.Rollback()

		err = tx.QueryRow("SELECT l, b + ? FROM u LIMIT 1", balanceDiff).Scan(&resp.Limite, &resp.Saldo)
		if err != nil {
			log.Print(err)
			return c.SendStatus(http.StatusInternalServerError)
		}

		if resp.Saldo < resp.Limite {
			return c.SendStatus(http.StatusUnprocessableEntity)
		}

		_, err = tx.Exec("INSERT INTO t (a, d) VALUES (?, ?)", balanceDiff, req.Descricao)
		if err != nil {
			log.Print(err)
			return c.SendStatus(http.StatusInternalServerError)
		}
		_, err = tx.Exec("UPDATE u SET b = ? WHERE l = ?", resp.Saldo, resp.Limite)
		if err != nil {
			log.Print(err)
			return c.SendStatus(http.StatusInternalServerError)
		}
		tx.Commit()

		// we store limit as a negative number, but users expect it to be positive
		resp.Limite *= -1

		log.Print("attempts:", i)
		return c.Status(http.StatusOK).JSON(resp)
	}

	return c.SendStatus(http.StatusLocked)
}
