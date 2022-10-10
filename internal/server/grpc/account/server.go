package account

import (
	"context"
	"errors"
	"fmt"
	"net"

	stError "github.com/go-errors/errors"
	"github.com/google/uuid"
	"github.com/shopspring/decimal"
	"google.golang.org/grpc"

	"github.com/ALarutin/pointpay_test/api/grpc/account/pb"
	"github.com/ALarutin/pointpay_test/internal/entity"
	service "github.com/ALarutin/pointpay_test/internal/service/account"
	"github.com/ALarutin/pointpay_test/pkg/logger"
)

var ErrorFromServer = errors.New("internal server error")

type Server struct {
	svc *service.Service
	pb.UnimplementedAccountsServer
	grpcSrv *grpc.Server
}

func NewServer(svc *service.Service) (*Server, error) {
	if svc == nil {
		return nil, fmt.Errorf("service is nil")
	}

	return &Server{svc: svc}, nil
}

func (s *Server) Start(ctx context.Context, port int) error {
	listener, err := net.Listen("tcp", fmt.Sprintf("localhost:%d", port))
	if err != nil {
		return err
	}

	grpcSrv := grpc.NewServer()

	s.grpcSrv = grpcSrv

	s.registerGRPC()

	go func() {
		logger.L().Info("server starting")

		if err = grpcSrv.Serve(listener); err != nil {
			logger.L().Fatal(err.Error())
		}
	}()

	<-ctx.Done()

	s.grpcSrv.GracefulStop()

	return nil
}

func (s *Server) registerGRPC() {
	pb.RegisterAccountsServer(s.grpcSrv, s)
}

func accountToPB(account entity.Account) *pb.Account {
	return &pb.Account{
		ID:       account.ID.String(),
		WalletId: account.WalletID,
		Balance:  account.Balance.String(),
	}
}

var ErrorNotValidID = errors.New("id is not valid uuid")

func idFromStrToUUID(id string) (uuid.UUID, error) {
	_id, err := uuid.Parse(id)
	if err != nil {
		return uuid.Nil, stError.Errorf("%w; %v", ErrorNotValidID, err)
	}

	return _id, nil
}

var ErrorNotValidBalance = errors.New("balance is not valid decimal")

func balanceFromStrToDecimal(balance string) (decimal.Decimal, error) {
	_balance, err := decimal.NewFromString(balance)
	if err != nil {
		return decimal.Decimal{}, stError.Errorf("%w; %v", ErrorNotValidBalance, err)
	}

	return _balance, nil
}
