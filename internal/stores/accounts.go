package stores

import (
	"context"
	"database/sql"
	"errors"
	"log"
	"strconv"
)

type accountsStore struct {
	db *sql.DB
}

func NewAccountStore(db *sql.DB) *accountsStore {
	return &accountsStore{db: db}
}

type AccountDetails struct {
	ID        int
	AccountID int
	Balance   float64
}

type AccountPayload struct {
	AccountID int
	Amount    float64
}

type AccountStore interface {
	Read(ctx context.Context, id int) (*AccountDetails, error)
	Create(ctx context.Context, payload *AccountPayload) error
	Update(ctx context.Context, payload *AccountPayload) error
}

func (a *accountsStore) Read(ctx context.Context, id int) (*AccountDetails, error) {
	var accountDetails AccountDetails

	query := "SELECT id, account_id, balance FROM accounts where account_id = $1"
	err := a.db.QueryRowContext(ctx, query, id).Scan(&accountDetails.ID, &accountDetails.AccountID, &accountDetails.Balance)
	if err != nil {
		log.Printf("error while get accountDetails err: %v", err.Error())

		if err == sql.ErrNoRows {
			return nil, errors.New("account not found " + strconv.Itoa(id))
		}

		return nil, errors.New("DB error")
	}

	return &accountDetails, nil
}

func (a *accountsStore) Create(ctx context.Context, payload *AccountPayload) error {
	query := `INSERT INTO accounts(account_id, balance) VALUES ($1, $2)`

	_, err := a.db.ExecContext(ctx, query, payload.AccountID, payload.Amount)
	if err != nil {
		log.Printf("error while inserting accountDetails err: %v", err.Error())

		return errors.New("DB error")
	}

	return nil
}

func (a *accountsStore) Update(ctx context.Context, payload *AccountPayload) error {
	query := `UPDATE accounts SET balance = $1 WHERE account_id = $2`

	_, err := a.db.ExecContext(ctx, query, payload.Amount, payload.AccountID)
	if err != nil {
		log.Printf("error while updating accountDetails err: %v", err.Error())

		return errors.New("DB error")
	}

	return nil
}
