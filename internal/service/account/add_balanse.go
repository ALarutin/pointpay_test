package account

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/shopspring/decimal"

	"github.com/ALarutin/pointpay_test/internal/entity"
)

func (s *Service) AddBalance(ctx context.Context, id uuid.UUID, balance decimal.Decimal) (entity.Account, error) {
	if err := validateID(id); err != nil {
		return entity.Account{}, err
	}

	if err := validateBalance(balance); err != nil {
		return entity.Account{}, err
	}

	account, err := s.repo.GetAccount(ctx, id)
	if err != nil {
		return entity.Account{}, err
	}

	if err = account.AddBalance(balance); err != nil {
		return entity.Account{}, err
	}

	if err = s.repo.UpdateAccount(ctx, account); err != nil {
		return entity.Account{}, err
	}

	fmt.Println(account.Balance)

	return account, nil
}
