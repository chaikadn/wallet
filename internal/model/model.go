package model

// type Wallet struct {
// 	ID      string `json:"walletId"` // TODO uuid
// 	Balance int    `json:"balance"`
// }

type ChangeBalanceRequest struct {
	WalletID      string        `json:"walletId"`
	OperationType OperationType `json:"operationType"` // валидировать через кастомный маршалер или через валидатор
	Amount        int           `json:"amount"`
}

type OperationType string

const (
	OperationTypeDeposit  = "DEPOSIT"
	OperationTypeWithdraw = "WITHDRAW"
)

type GetBalanseResponse struct {
	WalletID string `json:"walletId"`
	Balance  int    `json:"balance"`
}
