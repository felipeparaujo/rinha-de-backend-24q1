package api

import (
	"errors"
	"net/http"
	"strconv"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/jackc/pgx/v5"
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
	if err := a.runInTransaction(c, func(tx pgx.Tx) error {
		rows, err := tx.Query(c.Context(), Select10LatestTransactionsForUser, id)
		defer func() { rows.Close() }()
		if err != nil {
			return c.SendStatus(http.StatusInternalServerError)
		}

		for rows.Next() {
			transacao := ExtratoTransacaoResponse{}
			err := rows.Scan(&transacao.Valor, &transacao.Tipo, &transacao.Descricao, &transacao.RealizadaEm)
			if err != nil {
				return c.SendStatus(http.StatusInternalServerError)
			}

			// debit transactions are stored as negative values, but users expect them to be positive
			if transacao.Tipo == "d" {
				transacao.Valor *= -1
			}

			resp.UltimasTransacoes = append(resp.UltimasTransacoes, transacao)
		}

		if err = tx.
			QueryRow(c.Context(), SelectBalanceAndLimitForUser, id).
			Scan(&resp.Saldo.Total, &resp.Saldo.Limite); err != nil {
			return c.SendStatus(http.StatusInternalServerError)
		}

		return nil
	}); err != nil {
		return err
	}

	// we store limit as a negative number, but users expect it to be positive
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
	rows, err := a.transact(c, id, req.Valor, req.Descricao, req.Tipo)
	defer func() { rows.Close() }()

	if err != nil {
		return c.SendStatus(http.StatusInternalServerError)
	}

	for rows.Next() {
		var rowsUpdated int32
		if err = rows.Scan(&resp.Limite, &resp.Saldo, &rowsUpdated); err != nil {
			return c.SendStatus(http.StatusInternalServerError)
		}

		if rowsUpdated != 1 {
			return c.SendStatus(http.StatusUnprocessableEntity)
		}
	}

	// we store limit as a negative number, but users expect it to be positive
	resp.Limite *= -1

	return c.Status(http.StatusOK).JSON(resp)
}

func (a *App) transact(c *fiber.Ctx, id int, valor int32, descricao string, tipo string) (pgx.Rows, error) {
	if tipo == "d" {
		valor *= -1
	}

	return a.pool.Query(c.Context(), CreateTransaction, id, valor, tipo, descricao)
}

func (a *App) runInTransaction(c *fiber.Ctx, h func(pgx.Tx) error) error {
	tx, err := a.pool.BeginTx(c.Context(), pgx.TxOptions{})
	if err != nil {
		return err
	}

	if err = h(tx); err != nil {
		tx.Rollback(c.Context())
		return err
	}

	return tx.Commit(c.Context())
}
