package handler

import (
	"encoding/json"
	"net/http"
	"wallet/internal/service"

	"github.com/go-chi/chi/v5"
)

type Handler struct {
	service service.Service
}

func (h *Handler) RegisterRoutes(router *chi.Mux) {
	// общие middleware можно в main.go
	router.Route("/api/v1", func(r chi.Router) {
		r.Group(func(r chi.Router) {
			// сюда middleware для этой группы

			r.Get("/wallets/{WALLET_UUID}", h.handleBalance)
			r.Post("/wallet", h.handleTransaction)
		})
	})
}

func (h *Handler) handleBalance(w http.ResponseWriter, r *http.Request) {
	walletUUID := chi.URLParam(r, "WALLET_UUID")
	// валидировать uuid

	balance, err := h.service.GetBalance(walletUUID)
	if err != nil {
		// TODO: выбирать статус код исходя из ошибки
		// handle error
		h.respondError(w, "failed to get balance", http.StatusBadRequest)
		return
	}

	resp := balanseResponse{
		WalletUUID: walletUUID,
		Balance:    balance,
	}

	h.respondJSON(w, resp, http.StatusOK)
}

func (h *Handler) handleTransaction(w http.ResponseWriter, r *http.Request) {
	req := transactionRequest{}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		// handle error
		h.respondError(w, "failed to decode request", http.StatusBadRequest)
		return
	}

	switch req.OperationType {
	case operationTypeDeposit:
		if err := h.service.Deposit(req.WalletUUID, req.Amount); err != nil {
			h.respondError(w, "failed to deposit", http.StatusInternalServerError)
		}
	case operationTypeWithdraw:
		if err := h.service.Withdraw(req.WalletUUID, req.Amount); err != nil {
			h.respondError(w, "failed to withdraw", http.StatusInternalServerError)
		}
	default:
		h.respondError(w, "unknown operation type", http.StatusBadRequest)
		return
	}
}

// передавать ссылку на data для больших data
func (h *Handler) respondJSON(w http.ResponseWriter, data any, statusCode int) {
	// TODO
}

func (h *Handler) respondError(w http.ResponseWriter, errMessage string, statusCode int) {
	// TODO
}
