package services

import (
	"context"
	"errors"
	"strconv"

	"github.com/navyamani-ch/Account-Management/internal/stores"
)

type accountsService struct {
	store stores.AccountStore
}

func NewAccountService(store stores.AccountStore) *accountsService {
	return &accountsService{store: store}
}

type AccountDetails struct {
	AccountID int
	Balance   float64
}

type AccountCreatePayload struct {
	AccountID      int
	InitialBalance string
}

type AccountUpdatePayload struct {
	AccountID       int
	Amount          float64
	TransactionType string
}

type AccountService interface {
	Read(ctx context.Context, id int) (*AccountDetails, error)
	Create(ctx context.Context, payload *AccountCreatePayload) error
}

func (a *accountsService) Read(ctx context.Context, id int) (*AccountDetails, error) {
	resp, err := a.store.Read(ctx, id)
	if err != nil {
		return nil, err
	}

	response := buildResponse(resp)

	return response, nil
}

func (a *accountsService) Create(ctx context.Context, payload *AccountCreatePayload) error {
	if payload.AccountID == 0 {
		return errors.New("accountID is missing")
	}

	if payload.AccountID < 1 {
		return errors.New("invalid accountID")
	}

	initialBalance, err := strconv.ParseFloat(payload.InitialBalance, 64)
	if err != nil || initialBalance < 1 {
		return errors.New("invalid initialBalance")
	}

	accountDetails, _ := a.Read(ctx, payload.AccountID)
	if accountDetails != nil {
		return errors.New("account already exist")
	}

	err = a.store.Create(ctx, &stores.AccountPayload{
		AccountID: payload.AccountID,
		Amount:    initialBalance,
	})
	if err != nil {
		return err
	}

	return nil
}

func buildResponse(storeResp *stores.AccountDetails) *AccountDetails {
	var resp = &AccountDetails{
		AccountID: storeResp.AccountID,
		Balance:   storeResp.Balance,
	}

	return resp
}
