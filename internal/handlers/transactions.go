package handlers

import (
	"context"
	"encoding/json"
	"net/http"
	"time"

	"github.com/navyamani-ch/Account-Management/internal/services"
)

type transactionsHandler struct {
	service services.TransactionService
}

func NewTransactionHandler(service services.TransactionService) *transactionsHandler {
	return &transactionsHandler{service: service}
}

type TransactionCreatePayload struct {
	SourceAccountID      int    `json:"source_account_id"`
	DestinationAccountID int    `json:"destination_account_id"`
	Amount               string `json:"amount"`
}

func (h *transactionsHandler) CreateTransaction(w http.ResponseWriter, r *http.Request) {
	var transactionPay TransactionCreatePayload

	if err := json.NewDecoder(r.Body).Decode(&transactionPay); err != nil {
		w.WriteHeader(http.StatusBadRequest)

		_ = json.NewEncoder(w).Encode(map[string]string{
			"error": "Invalid Json",
		})

		return
	}

	ctx, cancelFunc := context.WithTimeout(r.Context(), 30*time.Second)
	defer cancelFunc()

	payload := &services.TransactionCreatePayload{
		SourceAccountID:      transactionPay.SourceAccountID,
		DestinationAccountID: transactionPay.DestinationAccountID,
		Amount:               transactionPay.Amount,
	}

	if err := h.service.Create(ctx, payload); err != nil {
		w.WriteHeader(http.StatusBadRequest)

		_ = json.NewEncoder(w).Encode(map[string]string{
			"error": err.Error(),
		})

		return
	}

	w.WriteHeader(http.StatusCreated)

	_ = json.NewEncoder(w).Encode(map[string]string{
		"response": "Transaction Successfully completed",
	})
}
