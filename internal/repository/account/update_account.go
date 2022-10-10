package account

import (
	"context"

	stError "github.com/go-errors/errors"
	"go.mongodb.org/mongo-driver/bson"

	"github.com/ALarutin/pointpay_test/internal/entity"
)

func (r *Repository) UpdateAccount(ctx context.Context, account entity.Account) error {
	collection := r.db.Collection(r.cfg.CollectionName)

	dto := accountToDTO(account)

	filter := bson.M{
		"_id":     dto.ID,
		"version": bson.M{"$lt": dto.Version},
	}

	update := bson.M{"$set": dto}

	tCtx, cancel := context.WithTimeout(ctx, r.cfg.Timeout)
	defer cancel()

	result, err := collection.UpdateOne(tCtx, filter, update)
	if err != nil {
		return stError.Errorf("account updating failed; %v", err)
	}

	if result.ModifiedCount == 0 {
		return stError.New(entity.ErrorOutdatedAccountToUpdate)
	}

	return nil
}
