package bank

import (
	"encoding/json"
	"io"
	"net/http"

	"github.com/gin-gonic/gin"
	"google.golang.org/protobuf/types/known/emptypb"

	"github.com/ALarutin/pointpay_test/pkg/logger"
	"github.com/ALarutin/pointpay_test/pkg/stacktrace"
)

func (h *handler) GetAccounts(c *gin.Context) {
	ctx := c.Request.Context()

	stream, err := h.client.GetAccounts(ctx, &emptypb.Empty{})
	if err != nil {
		logger.L().
			With(append(logger.GetFields(ctx), stacktrace.Key, stacktrace.Get(err))).
			Error(err.Error())

		c.JSON(http.StatusInternalServerError, gin.H{"error": ErrorFromServer.Error()})
		return
	}

	c.Writer.Header().Add("Content-Type", "application/json")

	if _, err = c.Writer.Write([]byte("[")); err != nil {
		logger.L().
			With(append(logger.GetFields(ctx), stacktrace.Key, stacktrace.Get(err))).
			Error(err.Error())

		c.JSON(http.StatusInternalServerError, gin.H{"error": ErrorFromServer.Error()})
		return
	}

	first := true

	for {
		account, _err := stream.Recv()
		if _err == io.EOF {
			break
		}
		if _err != nil {
			logger.L().
				With(append(logger.GetFields(ctx), stacktrace.Key, stacktrace.Get(_err))).
				Error(err.Error())

			c.JSON(http.StatusInternalServerError, gin.H{"error": ErrorFromServer.Error()})
			return
		}

		if !first {
			if _, err = c.Writer.Write([]byte(",")); err != nil {
				logger.L().
					With(append(logger.GetFields(ctx), stacktrace.Key, stacktrace.Get(err))).
					Error(err.Error())

				c.JSON(http.StatusInternalServerError, gin.H{"error": ErrorFromServer.Error()})
				return
			}
		}

		accountJSON, _err := json.Marshal(account)
		if _err != nil {
			logger.L().
				With(append(logger.GetFields(ctx), stacktrace.Key, stacktrace.Get(_err))).
				Error(err.Error())

			c.JSON(http.StatusInternalServerError, gin.H{"error": ErrorFromServer.Error()})
			return
		}

		if _, err = c.Writer.Write(accountJSON); _err != nil {
			logger.L().
				With(append(logger.GetFields(ctx), stacktrace.Key, stacktrace.Get(_err))).
				Error(err.Error())

			c.JSON(http.StatusInternalServerError, gin.H{"error": ErrorFromServer.Error()})
			return
		}

		if first {
			first = false
		}
	}

	if _, err = c.Writer.Write([]byte("]")); err != nil {
		logger.L().
			With(append(logger.GetFields(ctx), stacktrace.Key, stacktrace.Get(err))).
			Error(err.Error())

		c.JSON(http.StatusInternalServerError, gin.H{"error": ErrorFromServer.Error()})
		return
	}
}
