package wallet

import "wallet/internal/storage"

type WalletService struct {
	st storage.Storage
}
