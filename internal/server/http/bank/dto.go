package bank

import (
	"github.com/ALarutin/pointpay_test/api/grpc/account/pb"
)

type accountDTO struct {
	ID       string `json:"id"`
	WalletID uint64 `json:"walletID,omitempty"`
	Balance  string `json:"balance"`
}

type balanceDTO struct {
	Balance string `json:"balance,omitempty"`
}

func accountToDTO(account *pb.Account) accountDTO {
	return accountDTO{
		ID:       account.GetID(),
		WalletID: account.GetWalletId(),
		Balance:  account.GetBalance(),
	}
}
