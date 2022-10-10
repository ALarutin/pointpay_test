package bank

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"

	"github.com/ALarutin/pointpay_test/api/grpc/account/pb"
	"github.com/ALarutin/pointpay_test/internal/entity"
	grpcServer "github.com/ALarutin/pointpay_test/internal/server/grpc/account"
	"github.com/ALarutin/pointpay_test/pkg/logger"
	"github.com/ALarutin/pointpay_test/pkg/stacktrace"
)

func (h *handler) Withdrawal(c *gin.Context) {
	ctx := c.Request.Context()
	accountID := c.Param("id")
	if accountID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "account id is not exist"})
		return
	}

	var dto balanceDTO

	if err := c.Bind(&dto); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "json is not valid"})
		return
	}

	account, err := h.client.Withdrawal(ctx, &pb.ChangeBalanceRequest{
		ID:      accountID,
		Changes: dto.Balance,
	})
	if err == nil {
		c.JSON(http.StatusOK, accountToDTO(account))
		return
	}

	logger.L().
		With(append(logger.GetFields(ctx), stacktrace.Key, stacktrace.Get(err))).
		Error(err.Error())

	if strings.Contains(err.Error(), entity.ErrorAccountNotFound.Error()) {
		c.JSON(http.StatusNotFound, gin.H{"error": entity.ErrorAccountNotFound.Error()})
		return
	}
	if strings.Contains(err.Error(), entity.ErrorNegativeBalance.Error()) {
		c.JSON(http.StatusBadRequest, gin.H{"error": entity.ErrorNegativeBalance.Error()})
		return
	}
	if strings.Contains(err.Error(), grpcServer.ErrorNotValidID.Error()) {
		c.JSON(http.StatusBadRequest, gin.H{"error": grpcServer.ErrorNotValidID.Error()})
		return
	}
	if strings.Contains(err.Error(), grpcServer.ErrorNotValidBalance.Error()) {
		c.JSON(http.StatusBadRequest, gin.H{"error": grpcServer.ErrorNotValidBalance.Error()})
		return
	}

	c.JSON(http.StatusInternalServerError, gin.H{"error": ErrorFromServer.Error()})
}
