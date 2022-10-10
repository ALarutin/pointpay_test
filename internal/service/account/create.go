package account

import (
	"context"

	"github.com/ALarutin/pointpay_test/internal/entity"
)

func (s *Service) Create(ctx context.Context) (entity.Account, error) {
	account := entity.NewAccount()

	if err := s.repo.InsertAccount(ctx, account); err != nil {
		return entity.Account{}, err
	}

	return account, nil
}
