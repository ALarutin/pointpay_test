package account

import (
	"context"

	"github.com/ALarutin/pointpay_test/internal/entity"
)

func (s *Service) GetAllAsync(ctx context.Context) (<-chan entity.Account, <-chan struct{}) {
	return s.repo.GetAllAccountsAsync(ctx)
}
