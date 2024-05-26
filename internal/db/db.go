// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.26.0

package db

import (
	"context"
	"database/sql"
	"fmt"
)

type DBTX interface {
	ExecContext(context.Context, string, ...interface{}) (sql.Result, error)
	PrepareContext(context.Context, string) (*sql.Stmt, error)
	QueryContext(context.Context, string, ...interface{}) (*sql.Rows, error)
	QueryRowContext(context.Context, string, ...interface{}) *sql.Row
}

func New(db DBTX) *Queries {
	return &Queries{db: db}
}

func Prepare(ctx context.Context, db DBTX) (*Queries, error) {
	q := Queries{db: db}
	var err error
	if q.createTransactionStmt, err = db.PrepareContext(ctx, createTransaction); err != nil {
		return nil, fmt.Errorf("error preparing query CreateTransaction: %w", err)
	}
	if q.deleteTransactionStmt, err = db.PrepareContext(ctx, deleteTransaction); err != nil {
		return nil, fmt.Errorf("error preparing query DeleteTransaction: %w", err)
	}
	if q.getTransactionStmt, err = db.PrepareContext(ctx, getTransaction); err != nil {
		return nil, fmt.Errorf("error preparing query GetTransaction: %w", err)
	}
	if q.listTransactionsStmt, err = db.PrepareContext(ctx, listTransactions); err != nil {
		return nil, fmt.Errorf("error preparing query ListTransactions: %w", err)
	}
	if q.listTransactionsByConfirmStmt, err = db.PrepareContext(ctx, listTransactionsByConfirm); err != nil {
		return nil, fmt.Errorf("error preparing query ListTransactionsByConfirm: %w", err)
	}
	if q.updateTransactionStmt, err = db.PrepareContext(ctx, updateTransaction); err != nil {
		return nil, fmt.Errorf("error preparing query UpdateTransaction: %w", err)
	}
	return &q, nil
}

func (q *Queries) Close() error {
	var err error
	if q.createTransactionStmt != nil {
		if cerr := q.createTransactionStmt.Close(); cerr != nil {
			err = fmt.Errorf("error closing createTransactionStmt: %w", cerr)
		}
	}
	if q.deleteTransactionStmt != nil {
		if cerr := q.deleteTransactionStmt.Close(); cerr != nil {
			err = fmt.Errorf("error closing deleteTransactionStmt: %w", cerr)
		}
	}
	if q.getTransactionStmt != nil {
		if cerr := q.getTransactionStmt.Close(); cerr != nil {
			err = fmt.Errorf("error closing getTransactionStmt: %w", cerr)
		}
	}
	if q.listTransactionsStmt != nil {
		if cerr := q.listTransactionsStmt.Close(); cerr != nil {
			err = fmt.Errorf("error closing listTransactionsStmt: %w", cerr)
		}
	}
	if q.listTransactionsByConfirmStmt != nil {
		if cerr := q.listTransactionsByConfirmStmt.Close(); cerr != nil {
			err = fmt.Errorf("error closing listTransactionsByConfirmStmt: %w", cerr)
		}
	}
	if q.updateTransactionStmt != nil {
		if cerr := q.updateTransactionStmt.Close(); cerr != nil {
			err = fmt.Errorf("error closing updateTransactionStmt: %w", cerr)
		}
	}
	return err
}

func (q *Queries) exec(ctx context.Context, stmt *sql.Stmt, query string, args ...interface{}) (sql.Result, error) {
	switch {
	case stmt != nil && q.tx != nil:
		return q.tx.StmtContext(ctx, stmt).ExecContext(ctx, args...)
	case stmt != nil:
		return stmt.ExecContext(ctx, args...)
	default:
		return q.db.ExecContext(ctx, query, args...)
	}
}

func (q *Queries) query(ctx context.Context, stmt *sql.Stmt, query string, args ...interface{}) (*sql.Rows, error) {
	switch {
	case stmt != nil && q.tx != nil:
		return q.tx.StmtContext(ctx, stmt).QueryContext(ctx, args...)
	case stmt != nil:
		return stmt.QueryContext(ctx, args...)
	default:
		return q.db.QueryContext(ctx, query, args...)
	}
}

func (q *Queries) queryRow(ctx context.Context, stmt *sql.Stmt, query string, args ...interface{}) *sql.Row {
	switch {
	case stmt != nil && q.tx != nil:
		return q.tx.StmtContext(ctx, stmt).QueryRowContext(ctx, args...)
	case stmt != nil:
		return stmt.QueryRowContext(ctx, args...)
	default:
		return q.db.QueryRowContext(ctx, query, args...)
	}
}

type Queries struct {
	db                            DBTX
	tx                            *sql.Tx
	createTransactionStmt         *sql.Stmt
	deleteTransactionStmt         *sql.Stmt
	getTransactionStmt            *sql.Stmt
	listTransactionsStmt          *sql.Stmt
	listTransactionsByConfirmStmt *sql.Stmt
	updateTransactionStmt         *sql.Stmt
}

func (q *Queries) WithTx(tx *sql.Tx) *Queries {
	return &Queries{
		db:                            tx,
		tx:                            tx,
		createTransactionStmt:         q.createTransactionStmt,
		deleteTransactionStmt:         q.deleteTransactionStmt,
		getTransactionStmt:            q.getTransactionStmt,
		listTransactionsStmt:          q.listTransactionsStmt,
		listTransactionsByConfirmStmt: q.listTransactionsByConfirmStmt,
		updateTransactionStmt:         q.updateTransactionStmt,
	}
}
