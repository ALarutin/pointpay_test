package account

import (
	"context"

	"google.golang.org/protobuf/types/known/emptypb"

	"github.com/ALarutin/pointpay_test/api/grpc/account/pb"
)

func (s *Server) GetAccounts(_ *emptypb.Empty, stream pb.Accounts_GetAccountsServer) error {
	ctx := stream.Context()

	accountChan, doneChan := s.svc.GetAllAsync(ctx)

	for {
		select {
		case <-ctx.Done():
			return context.Canceled
		case <-doneChan:
			return nil
		case account := <-accountChan:
			if err := stream.Send(accountToPB(account)); err != nil {
				return err
			}
		}
	}
}
