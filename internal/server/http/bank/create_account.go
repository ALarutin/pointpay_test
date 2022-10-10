package bank

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"google.golang.org/protobuf/types/known/emptypb"

	"github.com/ALarutin/pointpay_test/pkg/logger"
	"github.com/ALarutin/pointpay_test/pkg/stacktrace"
)

func (h *handler) CreateAccount(c *gin.Context) {
	ctx := c.Request.Context()

	account, err := h.client.CreateAccount(ctx, &emptypb.Empty{})
	if err != nil {
		logger.L().
			With(append(logger.GetFields(ctx), stacktrace.Key, stacktrace.Get(err))).
			Error(err.Error())

		c.JSON(http.StatusInternalServerError, gin.H{"error": ErrorFromServer.Error()})
		return
	}

	c.JSON(http.StatusOK, accountToDTO(account))
}
