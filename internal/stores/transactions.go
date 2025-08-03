package stores

import (
	"context"
	"database/sql"
	"errors"
	"log"
)

type transactionStore struct {
	db           *sql.DB
	accountStore accountsStore
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
	tx, err := a.db.Begin()
	if err != nil {
		return errors.New("error while creating sql transaction")
	}

	defer func() {
		if err != nil {
			_ = tx.Rollback()
		}
	}()

	query := `INSERT INTO transactions (source_account_id, destination_account_id, amount) VALUES ($1, $2, $3)`

	_, err = tx.ExecContext(ctx, query, payload.SourceAccountID, payload.DestinationAccountID, payload.Amount)
	if err != nil {
		log.Printf("error while creating transaction err: %v", err.Error())

		return errors.New("DB error")
	}

	//Debit amount through source account
	err = a.accountStore.Update(ctx, &AccountPayload{
		AccountID: payload.SourceAccountID,
		Amount:    -1 * payload.Amount,
	}, tx)
	if err != nil {
		return errors.New("error while updating the sourceAccount Balance")
	}

	//Credit amount in destination account
	err = a.accountStore.Update(ctx, &AccountPayload{
		AccountID: payload.DestinationAccountID,
		Amount:    payload.Amount,
	}, tx)
	if err != nil {
		return errors.New("error while updating the destination Balance")
	}

	_ = tx.Commit()

	return nil
}
