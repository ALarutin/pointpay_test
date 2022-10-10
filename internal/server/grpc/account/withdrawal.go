package account

import (
	"context"
	"errors"

	"github.com/ALarutin/pointpay_test/api/grpc/account/pb"
	"github.com/ALarutin/pointpay_test/internal/entity"
	"github.com/ALarutin/pointpay_test/pkg/logger"
	"github.com/ALarutin/pointpay_test/pkg/stacktrace"
)

func (s *Server) Withdrawal(ctx context.Context, req *pb.ChangeBalanceRequest) (*pb.Account, error) {
	id, err := idFromStrToUUID(req.GetID())
	if err != nil {
		return nil, err
	}

	balance, err := balanceFromStrToDecimal(req.GetChanges())
	if err != nil {
		return nil, err
	}

	account, err := s.svc.AddBalance(ctx, id, balance.Neg())
	if err == nil {
		return accountToPB(account), nil
	}

	logger.L().
		With(append(logger.GetFields(ctx), stacktrace.Key, stacktrace.Get(err))).
		Error(err.Error())

	if errors.Is(err, entity.ErrorAccountNotFound) {
		return nil, entity.ErrorAccountNotFound
	}
	if errors.Is(err, entity.ErrorNegativeBalance) {
		return nil, entity.ErrorNegativeBalance
	}

	return nil, ErrorFromServer
}
