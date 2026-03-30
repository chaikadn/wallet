package handler

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"wallet/internal/service"

	"github.com/go-chi/chi/v5"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"
)

func TestHandler_handleBalance(t *testing.T) {
	tests := []struct {
		name                string
		walletID            string
		mockReturnedBalance int
		mockReturnedErr     error
		wantCode            int
		wantBody            string
	}{
		{
			name:                "success",
			walletID:            "w1",
			mockReturnedBalance: 500,
			mockReturnedErr:     nil,
			wantCode:            http.StatusOK,
			wantBody:            `{"walletId":"w1","balance":500}`,
		},
		{
			name:                "wallet not found",
			walletID:            "wrong-id",
			mockReturnedBalance: 0,
			mockReturnedErr:     service.ErrWalletNotFound,
			wantCode:            http.StatusNotFound,
			wantBody:            `{"error":"wallet not found"}`,
		},
		{
			name:                "service unavailable",
			walletID:            "w1",
			mockReturnedBalance: 0,
			mockReturnedErr:     service.ErrServiceUnavailable,
			wantCode:            http.StatusInternalServerError,
			wantBody:            `{"error":"failed to get balance"}`,
		},
		// 	TODO: empty wallet id, wrong id format
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockService := NewMockWalletService(t)

			if tt.walletID != "" {
				mockService.EXPECT().
					GetBalance(tt.walletID).
					Return(tt.mockReturnedBalance, tt.mockReturnedErr).
					Once()
			}

			r := httptest.NewRequest(http.MethodGet, "/"+tt.walletID, nil)
			w := httptest.NewRecorder()

			router := chi.NewRouter()
			router.Get("/{"+urlParamWalletUUID+"}", NewHandler(mockService, zap.NewNop()).handleBalance)

			router.ServeHTTP(w, r)

			assert.Equal(t, tt.wantCode, w.Code)
			assert.JSONEq(t, tt.wantBody, w.Body.String())
		})
	}
}

func TestHandler_handleTransaction(t *testing.T) {
	tests := []struct {
		name            string
		request         transactionRequest
		mockReturnedErr error
		wantCode        int
		wantBody        string
	}{
		{
			name: "success deposit",
			request: transactionRequest{
				WalletUUID:    "w1",
				OperationType: operationTypeDeposit,
				Amount:        100,
			},
			wantCode: http.StatusNoContent,
		},
		{
			name: "success withdraw",
			request: transactionRequest{
				WalletUUID:    "w1",
				OperationType: operationTypeWithdraw,
				Amount:        100,
			},
			wantCode: http.StatusNoContent,
		},
		{
			name: "wallet not found: deposit",
			request: transactionRequest{
				WalletUUID:    "wrong-id",
				OperationType: operationTypeDeposit,
				Amount:        100,
			},
			mockReturnedErr: service.ErrWalletNotFound,
			wantCode:        http.StatusNotFound,
			wantBody:        `{"error": "wallet not found"}`,
		},
		{
			name: "wallet not found: withdraw",
			request: transactionRequest{
				WalletUUID:    "wrong-id",
				OperationType: operationTypeWithdraw,
				Amount:        100,
			},
			mockReturnedErr: service.ErrWalletNotFound,
			wantCode:        http.StatusNotFound,
			wantBody:        `{"error": "wallet not found"}`,
		},
		{
			name: "insufficient funds",
			request: transactionRequest{
				WalletUUID:    "w1",
				OperationType: operationTypeWithdraw,
				Amount:        100,
			},
			mockReturnedErr: service.ErrInsufficientFunds,
			wantCode:        http.StatusUnprocessableEntity,
			wantBody:        `{"error": "insufficient funds"}`,
		},
		{
			name: "invalid amount: deposit",
			request: transactionRequest{
				WalletUUID:    "w1",
				OperationType: operationTypeDeposit,
				Amount:        0,
			},
			mockReturnedErr: service.ErrInvalidAmount,
			wantCode:        http.StatusBadRequest,
			wantBody:        `{"error": "amount must be greater than zero"}`,
		},
		{
			name: "invalid amount: withdraw",
			request: transactionRequest{
				WalletUUID:    "w1",
				OperationType: operationTypeWithdraw,
				Amount:        -100,
			},
			mockReturnedErr: service.ErrInvalidAmount,
			wantCode:        http.StatusBadRequest,
			wantBody:        `{"error": "amount must be greater than zero"}`,
		},
		{
			name: "unknown operation type",
			request: transactionRequest{
				WalletUUID:    "w1",
				OperationType: "wrong-operation-type",
				Amount:        100,
			},
			wantCode: http.StatusBadRequest,
			wantBody: `{"error": "unknown operation type"}`,
		},
		{
			name: "service unavailable",
			request: transactionRequest{
				WalletUUID:    "w1",
				OperationType: operationTypeWithdraw,
				Amount:        100,
			},
			mockReturnedErr: service.ErrServiceUnavailable,
			wantCode:        http.StatusInternalServerError,
			wantBody:        `{"error": "failed to process transaction"}`,
		},
		// TODO: wrong request format
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockService := NewMockWalletService(t)

			switch tt.request.OperationType {
			case operationTypeDeposit:
				mockService.EXPECT().
					Deposit(tt.request.WalletUUID, tt.request.Amount).
					Return(tt.mockReturnedErr).
					Once()
			case operationTypeWithdraw:
				mockService.EXPECT().
					Withdraw(tt.request.WalletUUID, tt.request.Amount).
					Return(tt.mockReturnedErr).
					Once()
			}

			var body bytes.Buffer
			err := json.NewEncoder(&body).Encode(&tt.request)
			require.NoError(t, err)

			r := httptest.NewRequest(http.MethodPost, "/", &body)
			w := httptest.NewRecorder()

			router := chi.NewRouter()
			router.Post("/", NewHandler(mockService, zap.NewNop()).handleTransaction)

			router.ServeHTTP(w, r)

			assert.Equal(t, tt.wantCode, w.Code)
			if w.Code != http.StatusNoContent {
				assert.JSONEq(t, tt.wantBody, w.Body.String())
			}
		})
	}
}
