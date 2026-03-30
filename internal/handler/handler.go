package handler

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

	m "wallet/internal/middleware"
	"wallet/internal/service"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"go.uber.org/zap"
)

//go:generate go run github.com/vektra/mockery/v3@v3.7.0
type WalletService interface {
	GetBalance(walletID string) (int, error)
	Deposit(walletID string, amount int) error
	Withdraw(walletID string, amount int) error
}

type Handler struct {
	walletService WalletService
	log           *zap.Logger
}

func NewHandler(walletService WalletService, logger *zap.Logger) *Handler {
	return &Handler{walletService: walletService, log: logger}
}

func (h *Handler) RegisterRoutes(r chi.Router) {
	r.Use(middleware.Recoverer)
	// r.Use(middleware.RequestID)
	r.Use(m.Logger(h.log))

	r.Route("/api/v1", func(r chi.Router) {
		r.Post("/wallets", h.handleTransaction)
		r.Get("/wallets/{"+urlParamWalletUUID+"}", h.handleBalance)
		// r.Post("/wallets/create", h.handleCreate)
	})
}

func (h *Handler) handleBalance(w http.ResponseWriter, r *http.Request) {
	walletUUID := chi.URLParam(r, urlParamWalletUUID)

	if walletUUID == "" {
		h.respondError(w, "empty wallet uuid", http.StatusBadRequest, errEmptyUUID)
		return
	}
	// TODO: проверить соответствие формату uuid

	balance, err := h.walletService.GetBalance(walletUUID)

	if errors.Is(err, service.ErrWalletNotFound) {
		h.respondError(w, "wallet not found", http.StatusNotFound, err)
		return
	}
	if err != nil {
		h.respondError(w, "failed to get balance", http.StatusInternalServerError, err)
		return
	}

	resp := balanceResponse{
		WalletUUID: walletUUID,
		Balance:    balance,
	}

	h.respondJSON(w, &resp, http.StatusOK)
}

func (h *Handler) handleTransaction(w http.ResponseWriter, r *http.Request) {
	req := transactionRequest{}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.respondError(w, "failed to decode request", http.StatusBadRequest, err)
		return
	}

	var err error
	switch req.OperationType {
	case operationTypeDeposit:
		err = h.walletService.Deposit(req.WalletUUID, req.Amount)
	case operationTypeWithdraw:
		err = h.walletService.Withdraw(req.WalletUUID, req.Amount)
	default:
		h.respondError(
			w, "unknown operation type", http.StatusBadRequest,
			fmt.Errorf("%w: %s", errUnknownOperationType, req.OperationType),
		)
		return
	}

	switch {
	case errors.Is(err, service.ErrWalletNotFound):
		h.respondError(w, "wallet not found", http.StatusNotFound, err)
		return
	case errors.Is(err, service.ErrInvalidAmount):
		h.respondError(w, "amount must be greater than zero", http.StatusBadRequest, err)
		return
	case errors.Is(err, service.ErrInsufficientFunds):
		h.respondError(w, "insufficient funds", http.StatusUnprocessableEntity, err)
		return
	case err != nil:
		h.respondError(w, "failed to process transaction", http.StatusInternalServerError, err)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (h *Handler) respondJSON(w http.ResponseWriter, data any, statusCode int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)

	err := json.NewEncoder(w).Encode(data)
	if err != nil {
		h.log.Error("failed to encode response", zap.Error(err))
	}
}

func (h *Handler) respondError(w http.ResponseWriter, errMessage string, statusCode int, err error) {
	h.log.Error(errMessage, zap.Error(err))

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)

	err = json.NewEncoder(w).Encode(errorMessage{Err: errMessage})
	if err != nil {
		h.log.Error("failed to encode error response", zap.Error(err))
	}
}
