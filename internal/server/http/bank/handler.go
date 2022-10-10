package bank

import (
	"errors"
	"fmt"

	"github.com/ALarutin/pointpay_test/api/grpc/account/pb"
)

var ErrorFromServer = errors.New("internal server error")

type handler struct {
	client pb.AccountsClient
}

func newHandler(client pb.AccountsClient) (*handler, error) {
	if client == nil {
		return nil, fmt.Errorf("cleint is nil")
	}

	return &handler{client: client}, nil
}
