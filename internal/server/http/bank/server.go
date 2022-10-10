package bank

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"

	"github.com/ALarutin/pointpay_test/api/grpc/account/pb"
	"github.com/ALarutin/pointpay_test/pkg/logger"
)

type Server struct {
	port   int
	router *gin.Engine
}

func NewServer(port int) *Server {
	return &Server{
		port:   port,
		router: gin.New(),
	}
}

func (s *Server) Start(ctx context.Context, client pb.AccountsClient) error {
	h, err := newHandler(client)
	if err != nil {
		return err
	}

	registerHandler(h, s.router)

	srv := &http.Server{
		Addr:    fmt.Sprintf("localhost:%d", s.port),
		Handler: s.router,
	}

	go func() {
		logger.L().Info("server starting")

		if err = srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logger.L().Fatal(err.Error())
		}
	}()

	<-ctx.Done()

	tCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err = srv.Shutdown(tCtx); err != nil {
		logger.L().Error(err.Error())
	}

	return nil
}

func registerHandler(handler *handler, router *gin.Engine) {
	router.POST("/account", handler.CreateAccount)
	router.GET("/accounts", handler.GetAccounts)
	router.POST("/account/:id/address", handler.GenerateAddress)
	router.PATCH("/account/:id/deposit", handler.Deposit)
	router.PATCH("/account/:id/withdrawal", handler.Withdrawal)
}
