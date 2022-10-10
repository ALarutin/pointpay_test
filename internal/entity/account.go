package entity

import (
	"encoding/binary"
	"errors"

	stError "github.com/go-errors/errors"
	"github.com/google/uuid"
	"github.com/shopspring/decimal"
)

var (
	ErrorAccountCorrupted        = errors.New("account was corrupted")
	ErrorAccountNotFound         = errors.New("account not found")
	ErrorWalletIsNotZero         = errors.New("wallet is not zero")
	ErrorWalletIsZero            = errors.New("wallet is zero")
	ErrorNegativeBalance         = errors.New("can't do transaction, balance would became negative")
	ErrorOutdatedAccountToUpdate = errors.New("account update data is outdated")
)

type Account struct {
	ID       uuid.UUID
	WalletID uint64
	Balance  decimal.Decimal
	Version  uint32
}

func NewAccount() Account {
	return Account{ID: uuid.New(), Version: 1}
}

func (a *Account) GenerateWallet() error {
	if a.WalletID != 0 {
		return stError.New(ErrorWalletIsNotZero)
	}

	a.WalletID = binary.BigEndian.Uint64(a.ID[:])

	a.Version++

	return nil
}

func (a *Account) AddBalance(balance decimal.Decimal) error {
	if a.WalletID == 0 {
		return stError.New(ErrorWalletIsZero)
	}

	currentBalance := a.Balance.Add(balance)
	if currentBalance.IsNegative() {
		return stError.New(ErrorNegativeBalance)
	}

	a.Balance = currentBalance

	a.Version++

	return nil
}
