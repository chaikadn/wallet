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

type walletService interface {
	GetBalance(walletID string) (int, error)
	Deposit(walletID string, amount int) error
	Withdraw(walletID string, amount int) error
}

type Handler struct {
	walletService walletService
	log           *zap.Logger
}

func NewHandler(walletService walletService, logger *zap.Logger) *Handler {
	return &Handler{walletService: walletService, log: logger}
}

func (h *Handler) RegisterRoutes(r chi.Router) {
	r.Use(middleware.Recoverer)
	// r.Use(middleware.RequestID)
	r.Use(m.Logger(h.log))

	r.Route("/api/v1", func(r chi.Router) {
		r.Get("/wallets/{"+urlParamWalletUUID+"}", h.handleBalance)
		r.Post("/wallets", h.handleTransaction)
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
		// TODO: выбирать статус код исходя из ошибки
		h.respondError(w, "failed to get balance", http.StatusBadRequest, err)
		return
	}

	resp := balanseResponse{
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

	// TODO: валидировать структуру transactionRequest

	switch req.OperationType {
	case operationTypeDeposit:
		err := h.walletService.Deposit(req.WalletUUID, req.Amount)
		if errors.Is(err, service.ErrWalletNotFound) {
			h.respondError(w, "wallet not found", http.StatusNotFound, err)
			return
		}
		if errors.Is(err, service.ErrInvalidAmount) {
			h.respondError(w, "amount must be greater than zero", http.StatusBadRequest, err)
			return
		}
		if err != nil {
			h.respondError(w, "failed to deposit", http.StatusInternalServerError, err)
			return
		}

	case operationTypeWithdraw:
		err := h.walletService.Withdraw(req.WalletUUID, req.Amount)
		if errors.Is(err, service.ErrWalletNotFound) {
			h.respondError(w, "wallet not found", http.StatusNotFound, err)
			return
		}
		if errors.Is(err, service.ErrInvalidAmount) {
			h.respondError(w, "amount must be greater than zero", http.StatusBadRequest, err)
			return
		}
		if errors.Is(err, service.ErrInsufficientFunds) {
			h.respondError(w, "insufficient funds", http.StatusUnprocessableEntity, err)
			return
		}
		if err != nil {
			h.respondError(w, "failed to withdraw", http.StatusInternalServerError, err)
			return
		}

	default:
		h.respondError(
			w, "unknown operation type", http.StatusBadRequest,
			fmt.Errorf("%w: %s", errUnknownOperationType, req.OperationType),
		)
		return
	}
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

// валидировать uuid
// в http слое проверить не пустой ли он и соответствует ли формату
// в service слое проверить бизнес-правила - кошелёк существует, и т.д.
// в storage слое - обычно валидации нет (каждый слой доверяет вызвавшему его слою), но можно сделать "defensive programming", защита от ошибок самого разработчика (например, если кто-то вызвал storage напрямую, минуя handler).
