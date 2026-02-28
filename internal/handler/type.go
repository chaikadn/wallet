package handler

import "errors"

const urlParamWalletUUID = "WALLET_UUID"

var (
	errEmptyUUID            = errors.New("empty uuid")
	errUnknownOperationType = errors.New("unknown operation type")
)

type transactionRequest struct {
	WalletUUID    string        `json:"walletId"`
	OperationType operationType `json:"operationType"` // валидировать
	Amount        int           `json:"amount"`
}

type operationType string

const (
	operationTypeDeposit  operationType = "DEPOSIT"
	operationTypeWithdraw operationType = "WITHDRAW"
)

type balanseResponse struct {
	WalletUUID string `json:"walletId"`
	Balance    int    `json:"balance"`
}

type errorMessage struct {
	Err string `json:"error"`
}
