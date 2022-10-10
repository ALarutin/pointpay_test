package account

import (
	"context"

	stError "github.com/go-errors/errors"

	"github.com/ALarutin/pointpay_test/internal/entity"
)

func (r *Repository) InsertAccount(ctx context.Context, account entity.Account) error {
	collection := r.db.Collection(r.cfg.CollectionName)

	tCtx, cancel := context.WithTimeout(ctx, r.cfg.Timeout)
	defer cancel()

	if _, err := collection.InsertOne(tCtx, accountToDTO(account)); err != nil {
		return stError.Errorf("account inserting failed; %v", err)
	}

	return nil
}
