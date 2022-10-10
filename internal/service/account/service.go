package account

import (
	"context"
	"errors"
	"fmt"

	stError "github.com/go-errors/errors"
	"github.com/google/uuid"
	"github.com/shopspring/decimal"

	"github.com/ALarutin/pointpay_test/internal/entity"
)

type Repository interface {
	InsertAccount(ctx context.Context, account entity.Account) error
	UpdateAccount(ctx context.Context, account entity.Account) error
	GetAccount(ctx context.Context, id uuid.UUID) (entity.Account, error)
	GetAllAccountsAsync(ctx context.Context) (accountChan <-chan entity.Account, doneChan <-chan struct{})
}

type Service struct {
	repo Repository
}

func New(repo Repository) (*Service, error) {
	if repo == nil {
		return nil, fmt.Errorf("repository is nil")
	}

	return &Service{repo: repo}, nil
}

var ErrorNilID = errors.New("id is nil")

func validateID(id uuid.UUID) error {
	if id == uuid.Nil {
		return stError.New(ErrorNilID)
	}

	return nil
}

var ErrorZeroBalance = errors.New("balance is zero")

func validateBalance(balance decimal.Decimal) error {
	if balance.IsZero() {
		return stError.New(ErrorZeroBalance)
	}

	return nil
}
