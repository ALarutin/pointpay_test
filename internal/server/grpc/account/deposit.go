package account

import (
	"context"

	"github.com/go-errors/errors"

	"github.com/ALarutin/pointpay_test/api/grpc/account/pb"
	"github.com/ALarutin/pointpay_test/internal/entity"
	"github.com/ALarutin/pointpay_test/pkg/logger"
	"github.com/ALarutin/pointpay_test/pkg/stacktrace"
)

func (s *Server) Deposit(ctx context.Context, req *pb.ChangeBalanceRequest) (*pb.Account, error) {
	id, err := idFromStrToUUID(req.GetID())
	if err != nil {
		return nil, err
	}

	balance, err := balanceFromStrToDecimal(req.GetChanges())
	if err != nil {
		return nil, err
	}

	account, err := s.svc.AddBalance(ctx, id, balance)
	if err == nil {
		return accountToPB(account), nil
	}

	logger.L().
		With(append(logger.GetFields(ctx), stacktrace.Key, stacktrace.Get(err))).
		Error(err.Error())

	if errors.Is(err, entity.ErrorAccountNotFound) {
		return nil, entity.ErrorAccountNotFound
	}

	return nil, ErrorFromServer
}
