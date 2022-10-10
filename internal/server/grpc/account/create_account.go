package account

import (
	"context"

	"google.golang.org/protobuf/types/known/emptypb"

	"github.com/ALarutin/pointpay_test/api/grpc/account/pb"
	"github.com/ALarutin/pointpay_test/pkg/logger"
	"github.com/ALarutin/pointpay_test/pkg/stacktrace"
)

func (s *Server) CreateAccount(ctx context.Context, _ *emptypb.Empty) (*pb.Account, error) {
	account, err := s.svc.Create(ctx)
	if err == nil {
		return accountToPB(account), nil
	}

	logger.L().
		With(append(logger.GetFields(ctx), stacktrace.Key, stacktrace.Get(err))).
		Error(err.Error())

	return nil, ErrorFromServer
}
