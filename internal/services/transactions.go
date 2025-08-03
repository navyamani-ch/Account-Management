package services

import (
	"context"
	"errors"
	"strconv"

	"github.com/navyamani-ch/Account-Management/internal/stores"
)

type transactionService struct {
	accountService AccountService
	store          stores.TransactionStore
}

func NewTransactionService(store stores.TransactionStore, accountService AccountService) *transactionService {
	return &transactionService{store: store, accountService: accountService}
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
	if payload.SourceAccountID < 0 {
		return errors.New("invalid sourceAccountID")
	}

	if payload.DestinationAccountID < 0 {
		return errors.New("invalid destinationAccountID")
	}

	amount, err := strconv.ParseFloat(payload.Amount, 64)
	if err != nil || amount < 1 {
		return errors.New("invalid initialBalance")
	}

	// Fetch source account details
	sourceAccountDetails, err := t.accountService.Read(ctx, payload.SourceAccountID)
	if err != nil {
		return errors.New("error while getting sourceAccountDetails")
	}

	// Check if the source account has enough balance
	if sourceAccountDetails.Balance < amount {
		return errors.New("error while getting destinationAccountDetails")
	}

	// Check if destination account exists
	_, err = t.accountService.Read(ctx, payload.DestinationAccountID)
	if err != nil {
		return err
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

	// Debit amount from source account
	err = t.accountService.Update(ctx, &AccountUpdatePayload{
		AccountID:       payload.SourceAccountID,
		Amount:          amount,
		TransactionType: "DEBIT",
	})
	if err != nil {
		return err
	}

	// Credit amount to destination account
	err = t.accountService.Update(ctx, &AccountUpdatePayload{
		AccountID:       payload.DestinationAccountID,
		Amount:          amount,
		TransactionType: "CREDIT",
	})
	if err != nil {
		return err
	}

	return nil
}
