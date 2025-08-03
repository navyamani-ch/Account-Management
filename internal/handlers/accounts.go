package handlers

import (
	"context"
	"encoding/json"
	"net/http"
	"strconv"
	"time"

	"github.com/gorilla/mux"

	"github.com/navyamani-ch/Account-Management/internal/services"
)

type AccountHandler struct {
	service services.AccountService
}

func NewAccountHandler(service services.AccountService) *AccountHandler {
	return &AccountHandler{service: service}
}

type AccountDetails struct {
	AccountID int    `json:"account_id"`
	Balance   string `json:"balance"`
}

type AccountCreatePayload struct {
	AccountID      int    `json:"account_id"`
	InitialBalance string `json:"initial_balance"`
}

func (h *AccountHandler) CreateAccount(w http.ResponseWriter, r *http.Request) {
	var acc AccountCreatePayload

	if err := json.NewDecoder(r.Body).Decode(&acc); err != nil {
		w.WriteHeader(http.StatusBadRequest)

		_ = json.NewEncoder(w).Encode(map[string]string{
			"error": "Invalid Json",
		})

		return
	}

	ctx, cancelFunc := context.WithTimeout(r.Context(), 30*time.Second)
	defer cancelFunc()

	payload := &services.AccountCreatePayload{
		AccountID:      acc.AccountID,
		InitialBalance: acc.InitialBalance,
	}

	if err := h.service.Create(ctx, payload); err != nil {
		w.WriteHeader(http.StatusBadRequest)

		_ = json.NewEncoder(w).Encode(map[string]string{
			"error": err.Error(),
		})

		return
	}

	w.WriteHeader(http.StatusCreated)
}

func (h *AccountHandler) GetAccount(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	idStr := vars["account_id"]

	accountID, err := strconv.Atoi(idStr)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)

		_ = json.NewEncoder(w).Encode(map[string]string{
			"error": "Invalid AccountId",
		})

		return
	}

	ctx, cancelFunc := context.WithTimeout(r.Context(), 30*time.Second)
	defer cancelFunc()

	acc, err := h.service.Read(ctx, accountID)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)

		_ = json.NewEncoder(w).Encode(map[string]string{
			"error": err.Error(),
		})

		return
	}

	if acc == nil {
		w.WriteHeader(http.StatusBadRequest)

		_ = json.NewEncoder(w).Encode(map[string]string{
			"error": "account not found",
		})

		return
	}

	resp := buildResponse(acc)

	w.Header().Set("Content-Type", "application/json")

	w.WriteHeader(http.StatusOK)

	_ = json.NewEncoder(w).Encode(resp)
}

func buildResponse(accountResp *services.AccountDetails) *AccountDetails {
	amount := strconv.FormatFloat(accountResp.Balance, 'f', 2, 64)

	var resp = &AccountDetails{
		AccountID: accountResp.AccountID,
		Balance:   amount,
	}

	return resp
}
