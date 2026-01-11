package handler

// const urlParamWalletUUID = "WALLET_UUID"

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
