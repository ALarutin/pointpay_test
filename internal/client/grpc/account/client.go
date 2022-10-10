package account

import (
	"context"
	"fmt"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	"github.com/ALarutin/pointpay_test/api/grpc/account/pb"
	"github.com/ALarutin/pointpay_test/pkg/logger"
)

type Client struct {
	conn *grpc.ClientConn
	pb.AccountsClient
}

func NewClient(ctx context.Context, port int) (*Client, error) {
	tCtx, cancel := context.WithTimeout(ctx, time.Second*5)
	defer cancel()

	conn, err := grpc.DialContext(tCtx, fmt.Sprintf("localhost:%d", port), grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, err
	}

	return &Client{
		conn:           conn,
		AccountsClient: pb.NewAccountsClient(conn),
	}, nil
}

func (c *Client) Close() {
	if err := c.conn.Close(); err != nil {
		logger.L().Error(err.Error())
	}
}
