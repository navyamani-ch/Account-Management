package services

import (
	"context"
	"errors"
	"strconv"

	"github.com/navyamani-ch/Account-Management/internal/stores"
)

type transactionService struct {
	accountStore stores.AccountStore
	store        stores.TransactionStore
}

func NewTransactionService(store stores.TransactionStore, accountStore stores.AccountStore) *transactionService {
	return &transactionService{store: store, accountStore: accountStore}
}

type TransactionCreatePayload struct {
	SourceAccountID      int
	DestinationAccountID int
	Amount               string
}

type TransactionService interface {
	Create(ctx context.Context, payload *TransactionCreatePayload) error
}

func (t *transactionService) Create(ctx context.Context, payload *TransactionCreatePayload) error {
	if payload.SourceAccountID < 1 {
		return errors.New("invalid sourceAccountID")
	}

	if payload.DestinationAccountID < 1 {
		return errors.New("invalid destinationAccountID")
	}

	amount, err := strconv.ParseFloat(payload.Amount, 64)
	if err != nil || amount < 1 {
		return errors.New("invalid amount")
	}

	// Fetch source and destination account details
	accountDetails, err := t.accountStore.GetAccounts(ctx, []int{payload.SourceAccountID,
		payload.DestinationAccountID})
	if err != nil {
		return errors.New("error while getting sourceAccountDetails and destinationDetails")
	}

	if len(accountDetails) < 2 {
		return errors.New("sourceAccount or destinationAccount details not exist")
	}

	for i := range accountDetails {
		if accountDetails[i].AccountID == payload.SourceAccountID {
			// Check if the source account has enough balance
			if accountDetails[i].Balance < amount {
				return errors.New("insufficient balance")
			}
		}
	}

	// Create a new transaction record
	err = t.store.Create(ctx, &stores.TransactionCreatePayload{
		SourceAccountID:      payload.SourceAccountID,
		DestinationAccountID: payload.DestinationAccountID,
		Amount:               amount,
	})
	if err != nil {
		return err
	}

	return nil
}
