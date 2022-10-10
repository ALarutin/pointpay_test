package main

import (
	"context"
	"os"
	"os/signal"
	"strconv"
	"syscall"

	"github.com/ALarutin/pointpay_test/internal/client/grpc/account"
	"github.com/ALarutin/pointpay_test/internal/server/http/bank"
	"github.com/ALarutin/pointpay_test/pkg/logger"
)

var (
	logLevel       string
	httpServerPort int
	grpcServerPort int
)

func init() {
	logLevel = os.Getenv("LOG_LEVEL")
	httpServerPort, _ = strconv.Atoi(os.Getenv("HTTP_SERVER_PORT"))
	grpcServerPort, _ = strconv.Atoi(os.Getenv("GRPC_SERVER_PORT"))
}

func init() {
	logger.Init(logLevel)
}

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGHUP, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)

	go func() {
		<-sigChan
		cancel()
	}()

	client, err := account.NewClient(ctx, grpcServerPort)
	if err != nil {
		logger.L().Fatal(err)
		return
	}
	defer client.Close()

	srv := bank.NewServer(httpServerPort)

	if err = srv.Start(ctx, client); err != nil {
		logger.L().Error(err)
	}
}
