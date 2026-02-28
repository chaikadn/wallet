package handler

import (
	"encoding/json"
	"fmt"
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

			r.Get(fmt.Sprintf("/wallets/{%s}", urlParamWalletUUID), h.handleBalance)
			r.Post("/wallet", h.handleTransaction)
		})
	})
}

func (h *Handler) handleBalance(w http.ResponseWriter, r *http.Request) {
	walletUUID := chi.URLParam(r, urlParamWalletUUID)

	// валидировать uuid
	// в http слое проверить не пустой ли он и соответствует ли формату
	// в service слое проверить бизнес-правила - кошелёк существует, и т.д.
	// в storage слое - обычно валидации нет (каждый слой доверяет вызвавшему его слою), но можно сделать "defensive programming", защита от ошибок самого разработчика (например, если кто-то вызвал storage напрямую, минуя handler).

	if walletUUID == "" {
		h.respondError(w, "empty wallet uuid", http.StatusBadRequest, errEmptyUUID)
		return
	}
	// TODO: проверить соответствие формату uuid

	balance, err := h.service.GetBalance(walletUUID)
	if err != nil {
		// TODO: выбирать статус код исходя из ошибки
		// handle error
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
		// handle error
		h.respondError(w, "failed to decode request", http.StatusBadRequest, err)
		return
	}

	switch req.OperationType {
	case operationTypeDeposit:
		if err := h.service.Deposit(req.WalletUUID, req.Amount); err != nil {
			h.respondError(w, "failed to deposit", http.StatusInternalServerError, err)
			return
		}
	case operationTypeWithdraw:
		if err := h.service.Withdraw(req.WalletUUID, req.Amount); err != nil {
			h.respondError(w, "failed to withdraw", http.StatusInternalServerError, err)
			return
		}
	default:
		h.respondError(
			w, "unknown operation type", http.StatusBadRequest,
			fmt.Errorf("%v: %s", errUnknownOperationType, req.OperationType),
		)
		return
	}
}

func (h *Handler) respondJSON(w http.ResponseWriter, data any, statusCode int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)

	err := json.NewEncoder(w).Encode(data)
	if err != nil {
		// TODO: log error
	}
}

func (h *Handler) respondError(w http.ResponseWriter, errMessage string, statusCode int, err error) {
	// TODO: log error

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)

	err = json.NewEncoder(w).Encode(errorMessage{Err: errMessage})
	if err != nil {
		// TODO: log encoding error
	}
}
