package api

import "database/sql"

type PreparedQueries struct {
	SelectLimitAndBalance    *sql.Stmt
	SelectLast10Transactions *sql.Stmt
	SelectLimitAndNewBalance *sql.Stmt
	InsertTransaction        *sql.Stmt
	UpdateBalance            *sql.Stmt
}

func PrepareQueries(db *sql.DB) (PreparedQueries, error) {
	p := PreparedQueries{}
	var err error
	p.SelectLimitAndBalance, err = db.Prepare(`SELECT l, b FROM u`)
	if err != nil {
		return p, err
	}
	p.SelectLast10Transactions, err = db.Prepare(`SELECT t, a, d FROM t ORDER BY id DESC LIMIT 10`)
	if err != nil {
		return p, err
	}
	p.SelectLimitAndNewBalance, err = db.Prepare("SELECT l, b + ? FROM u LIMIT 1")
	if err != nil {
		return p, err
	}
	p.InsertTransaction, err = db.Prepare("INSERT INTO t (a, d) VALUES (?, ?)")
	if err != nil {
		return p, err
	}
	p.UpdateBalance, err = db.Prepare("UPDATE u SET b = ? WHERE l = ?")
	if err != nil {
		return p, err
	}
	return p, nil
}
