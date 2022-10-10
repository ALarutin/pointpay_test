package account

import (
	"context"
	"errors"
	"strings"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	service "github.com/ALarutin/pointpay_test/internal/service/account"
)

var _ service.Repository = (*Repository)(nil)

type Repository struct {
	client *mongo.Client
	db     *mongo.Database
	cfg    *Cfg
}

type Cfg struct {
	URI            string
	Timeout        time.Duration
	CollectionName string
}

var errorNilCfg = errors.New("nil Cfg")

func New(ctx context.Context, cfg *Cfg) (*Repository, error) {
	if cfg == nil {
		return nil, errorNilCfg
	}

	client, err := mongo.NewClient(
		options.Client().
			SetMaxConnecting(10).
			SetMaxConnIdleTime(time.Second * 10).
			ApplyURI(cfg.URI),
	)
	if err != nil {
		return nil, err
	}

	dbCtx, cancel := context.WithTimeout(ctx, time.Second*5)
	defer cancel()

	if err = client.Connect(dbCtx); err != nil {
		return nil, err
	}

	split := strings.Split(cfg.URI, "/")
	database := split[len(split)-1]

	return &Repository{
		client: client,
		db:     client.Database(database),
		cfg:    cfg,
	}, nil
}

func (r *Repository) Close() error {
	dbCtx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	return r.client.Disconnect(dbCtx)
}
