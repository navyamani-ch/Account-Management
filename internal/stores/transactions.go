package stores

import (
	"context"
	"database/sql"
	"errors"
	"log"
)

type transactionStore struct {
	db *sql.DB
}

func NewTransactionStore(db *sql.DB) *transactionStore {
	return &transactionStore{db: db}
}

type TransactionCreatePayload struct {
	ID                   int
	SourceAccountID      int
	DestinationAccountID int
	Amount               float64
}

type TransactionStore interface {
	Create(ctx context.Context, payload *TransactionCreatePayload) error
}

func (a *transactionStore) Create(ctx context.Context, payload *TransactionCreatePayload) error {
	query := `INSERT INTO transactions (source_account_id, destination_account_id, amount) VALUES ($1, $2, $3)`

	_, err := a.db.ExecContext(ctx, query, payload.SourceAccountID, payload.DestinationAccountID, payload.Amount)
	if err != nil {
		log.Printf("error while creating transaction err: %v", err.Error())

		return errors.New("DB error")
	}

	return nil
}
