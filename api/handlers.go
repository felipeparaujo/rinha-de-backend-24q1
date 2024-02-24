package api

import (
	"errors"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/jackc/pgx/v5"
)

var ErrInvalidTransacoesRequest = errors.New("invalid request")

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
	Valor       *int32     `json:"valor"`
	Tipo        *string    `json:"tipo"`
	Descricao   *string    `json:"descricao"`
	RealizadaEm *time.Time `json:"realizada_em"`
}

func validateID(id string) int {
	idNum, err := strconv.Atoi(id)
	if err != nil {
		log.Print(err)
		return http.StatusBadRequest
	}

	if idNum < 1 || idNum > 5 {
		return http.StatusNotFound
	}

	return http.StatusOK
}

func (a *App) extrato(c *fiber.Ctx) error {
	id, err := strconv.Atoi(c.Params("id"))
	if err != nil || id < 1 || id > 5 {
		return c.SendStatus(http.StatusNotFound)
	}

	rows, err := a.pool.Query(
		c.Context(), `
		SELECT
			c.saldo, c.limite, t.valor, t.tipo, t.descricao, t.realizada_em
		FROM clientes c LEFT JOIN (
			SELECT * FROM transacoes
			WHERE cliente_id = $1
			ORDER BY realizada_em DESC
			LIMIT 10
		) t ON c.id = t.cliente_id
		WHERE c.id = $1
		`,
		id,
	)
	defer rows.Close()
	if err != nil {
		log.Print(err)
		return c.SendStatus(http.StatusInternalServerError)
	}

	resp := ExtratoResponse{Saldo: ExtratSaldoResponse{DataExtrato: time.Now()}, UltimasTransacoes: []ExtratoTransacaoResponse{}}
	for rows.Next() {
		transacao := ExtratoTransacaoResponse{}
		err := rows.Scan(&resp.Saldo.Total, &resp.Saldo.Limite, &transacao.Valor, &transacao.Tipo, &transacao.Descricao, &transacao.RealizadaEm)
		if err != nil {
			log.Print(err)
			return c.SendStatus(http.StatusInternalServerError)
		}

		// even if there are no transactions we'll get a row back because of the left join
		// so check if we've actually got a transaction before appending
		if transacao.RealizadaEm != nil {
			resp.UltimasTransacoes = append(resp.UltimasTransacoes, transacao)
		}
	}

	return c.Status(http.StatusOK).JSON(resp)
}

type TransacoesRequest struct {
	Valor     int32  `json:"valor"`
	Tipo      string `json:"tipo"`
	Descricao string `json:"descricao"`
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

type TransacoesResponse struct {
	Limite int32 `json:"limite"`
	Saldo  int32 `json:"saldo"`
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
	if err := a.runInTransaction(c, func(tx pgx.Tx) error {
		rows, err := tx.Query(c.Context(), "SELECT saldo, limite FROM clientes WHERE id = $1 FOR UPDATE", id)
		if err != nil {
			return c.SendStatus(http.StatusInternalServerError)
		}

		for rows.Next() {
			if err = rows.Scan(&resp.Saldo, &resp.Limite); err != nil {
				return c.SendStatus(http.StatusInternalServerError)
			}
		}

		saldoFinal := resp.Saldo + req.Valor
		if req.Tipo == "d" {
			saldoFinal = resp.Saldo - req.Valor
		}

		if saldoFinal < -resp.Limite {
			return c.SendStatus(http.StatusUnprocessableEntity)
		}

		update, err := tx.Query(c.Context(), "UPDATE clientes SET saldo = $1 WHERE id = $2", saldoFinal, id)
		update.Next()
		if err != nil {
			return c.SendStatus(http.StatusInternalServerError)
		}

		resp.Saldo = saldoFinal

		return nil
	}); err != nil {
		return err
	}

	insert, err := a.pool.Query(
		c.Context(),
		"INSERT INTO transacoes (cliente_id, valor, tipo, descricao) VALUES ($1, $2, $3, $4)", id, req.Valor, req.Tipo, req.Descricao,
	)
	insert.Next()
	if err != nil {
		c.Status(http.StatusInternalServerError)
	}

	return c.Status(http.StatusOK).JSON(resp)
}

func (a *App) runInTransaction(c *fiber.Ctx, h func(pgx.Tx) error) error {
	tx, err := a.pool.BeginTx(c.Context(), pgx.TxOptions{})
	if err = h(tx); err != nil {
		tx.Rollback(c.Context())
		return err
	}

	return tx.Commit(c.Context())
}
