package account

import (
	"context"

	"go.mongodb.org/mongo-driver/bson"

	"github.com/ALarutin/pointpay_test/internal/entity"
	"github.com/ALarutin/pointpay_test/pkg/logger"
	"github.com/ALarutin/pointpay_test/pkg/stacktrace"
)

func (r *Repository) GetAllAccountsAsync(ctx context.Context) (<-chan entity.Account, <-chan struct{}) {
	var (
		accountChan = make(chan entity.Account)
		doneChan    = make(chan struct{})
	)

	go func() {
		defer close(accountChan)
		defer close(doneChan)

		collection := r.db.Collection(r.cfg.CollectionName)

		tCtx, cancel := context.WithTimeout(ctx, r.cfg.Timeout)
		defer cancel()

		cursor, err := collection.Find(tCtx, bson.M{})
		if err != nil {
			logger.L().
				With(append(logger.GetFields(ctx), stacktrace.Key, stacktrace.Get(err))).
				Errorf("accounts finding failed: %v", err)
			return
		}
		defer func() {
			_tCtx, _cancel := context.WithTimeout(ctx, r.cfg.Timeout)
			defer _cancel()

			cursor.Close(_tCtx)
		}()

		for {
			select {
			case <-ctx.Done():
				return
			default:
			}

			isEnd := func() bool {
				_tCtx, _cancel := context.WithTimeout(ctx, r.cfg.Timeout)
				defer _cancel()

				if !cursor.Next(_tCtx) {
					return true
				}

				var account entity.Account

				account, err = decodeAccount(cursor)
				if err != nil {
					logger.L().
						With(append(logger.GetFields(ctx), stacktrace.Key, stacktrace.Get(err))).
						Errorf("accounts decoding failed: %v", err)
					return true
				}

				accountChan <- account

				return false
			}()

			if isEnd {
				break
			}
		}

		if err = cursor.Err(); err != nil {
			logger.L().
				With(append(logger.GetFields(ctx), stacktrace.Key, stacktrace.Get(err))).
				Errorf("accounts finding failed; cursor err: %v", err)
		}
	}()

	return accountChan, doneChan
}
