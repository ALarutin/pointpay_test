package account

import (
	stError "github.com/go-errors/errors"
	"github.com/google/uuid"
	"github.com/shopspring/decimal"

	"github.com/ALarutin/pointpay_test/internal/entity"
)

type accountDTO struct {
	ID       string `bson:"_id,omitempty"`
	WalletID int64  `bson:"wallet_id,omitempty"`
	Balance  string `bson:"balance,omitempty"`
	Version  uint32 `bson:"version,omitempty"`
}

func (a accountDTO) toEntity() (entity.Account, error) {
	id, err := uuid.Parse(a.ID)
	if err != nil {
		return entity.Account{}, stError.Errorf("%w; %v", entity.ErrorAccountCorrupted, err)
	}

	balance, err := decimal.NewFromString(a.Balance)
	if err != nil {
		return entity.Account{}, stError.Errorf("%w; %v", entity.ErrorAccountCorrupted, err)
	}

	return entity.Account{
		ID:       id,
		WalletID: uint64(a.WalletID),
		Balance:  balance,
		Version:  a.Version,
	}, nil
}

func accountToDTO(account entity.Account) accountDTO {
	return accountDTO{
		ID:       account.ID.String(),
		WalletID: int64(account.WalletID),
		Balance:  account.Balance.String(),
		Version:  account.Version,
	}
}

type decoder interface {
	Decode(interface{}) error
}

func decodeAccount(d decoder) (entity.Account, error) {
	var dto accountDTO

	if err := d.Decode(&dto); err != nil {
		return entity.Account{}, stError.Errorf("account decoding failed; %v", err)
	}

	if dto.ID == "" {
		return entity.Account{}, stError.New(entity.ErrorAccountNotFound)
	}

	account, err := dto.toEntity()
	if err != nil {
		return entity.Account{}, err
	}

	return account, err
}
