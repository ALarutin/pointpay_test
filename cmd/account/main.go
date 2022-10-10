package main

import (
	"context"
	"os"
	"os/signal"
	"strconv"
	"syscall"
	"time"

	repository "github.com/ALarutin/pointpay_test/internal/repository/account"
	grpcServer "github.com/ALarutin/pointpay_test/internal/server/grpc/account"
	service "github.com/ALarutin/pointpay_test/internal/service/account"
	"github.com/ALarutin/pointpay_test/pkg/logger"
)

var (
	logLevel            string
	mongoURI            string
	mongoCollectionName string
	grpcServerPort      int
)

func init() {
	logLevel = os.Getenv("LOG_LEVEL")
	mongoURI = os.Getenv("MONGO_URI")
	mongoCollectionName = os.Getenv("MONGO_COLLECTION_NAME")
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

	repo, err := repository.New(ctx, &repository.Cfg{
		URI:            mongoURI,
		Timeout:        time.Second,
		CollectionName: mongoCollectionName,
	})
	if err != nil {
		logger.L().Fatal(err)
	}
	defer repo.Close()

	svc, err := service.New(repo)
	if err != nil {
		logger.L().Error(err)
		return
	}

	srv, err := grpcServer.NewServer(svc)
	if err != nil {
		logger.L().Error(err)
		return
	}

	if err = srv.Start(ctx, grpcServerPort); err != nil {
		logger.L().Error(err)
	}
}
