package account

import (
	"context"
	"errors"

	"github.com/ALarutin/pointpay_test/api/grpc/account/pb"
	"github.com/ALarutin/pointpay_test/internal/entity"
	"github.com/ALarutin/pointpay_test/pkg/logger"
	"github.com/ALarutin/pointpay_test/pkg/stacktrace"
)

func (s *Server) GenerateAddress(ctx context.Context, req *pb.GenerateAddressRequest) (*pb.Account, error) {
	accountID, err := idFromStrToUUID(req.GetAccountID())
	if err != nil {
		return nil, err
	}

	account, err := s.svc.GenerateWalletID(ctx, accountID)
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
