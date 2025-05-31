package impl

import (
	"time"

	"github.com/gin-gonic/gin"
	"github.com/janobono/captcha-service/generated/openapi"
)

func RespondWithError(ctx *gin.Context, statusCode int, code openapi.ErrorCode, message string) {
	ctx.JSON(statusCode, openapi.ErrorMessage{
		Code:      code,
		Message:   message,
		Timestamp: time.Now().UTC(),
	})
}
