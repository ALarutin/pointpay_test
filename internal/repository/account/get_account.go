package account

import (
	"context"
	"errors"

	stError "github.com/go-errors/errors"
	"github.com/google/uuid"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"

	"github.com/ALarutin/pointpay_test/internal/entity"
)

func (r *Repository) GetAccount(ctx context.Context, id uuid.UUID) (entity.Account, error) {
	collection := r.db.Collection(r.cfg.CollectionName)

	filter := bson.M{"_id": id.String()}

	tCtx, cancel := context.WithTimeout(ctx, r.cfg.Timeout)
	defer cancel()

	result := collection.FindOne(tCtx, filter)

	if err := result.Err(); err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return entity.Account{}, stError.Errorf("%w; %v", entity.ErrorAccountNotFound, err)
		}

		return entity.Account{}, stError.Errorf("account finding failed; %v", err)
	}

	account, err := decodeAccount(result)
	if err != nil {
		return entity.Account{}, err
	}

	return account, nil
}
