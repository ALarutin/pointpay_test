package account

import (
	"context"

	"github.com/google/uuid"

	"github.com/ALarutin/pointpay_test/internal/entity"
)

func (s *Service) GenerateWalletID(ctx context.Context, id uuid.UUID) (entity.Account, error) {
	if err := validateID(id); err != nil {
		return entity.Account{}, err
	}

	account, err := s.repo.GetAccount(ctx, id)
	if err != nil {
		return entity.Account{}, err
	}

	if err = account.GenerateWallet(); err != nil {
		return entity.Account{}, err
	}

	if err = s.repo.UpdateAccount(ctx, account); err != nil {
		return entity.Account{}, err
	}

	return account, nil
}
