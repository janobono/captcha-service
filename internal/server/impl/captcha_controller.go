package impl

import (
	"github.com/gin-gonic/gin"
	"github.com/janobono/captcha-service/generated/openapi"
	"github.com/janobono/captcha-service/internal/service"
	"log/slog"
	"net/http"
)

type captchaController struct {
	captchaService *service.CaptchaService
}

func NewCaptchaController(captchaService *service.CaptchaService) openapi.CaptchaControllerAPI {
	return &captchaController{captchaService}
}

func (c *captchaController) GetCaptcha(ctx *gin.Context) {
	result, err := c.captchaService.Create(ctx)
	if err != nil {
		slog.Error("Failed to generate CAPTCHA", "error", err)
		RespondWithError(ctx, http.StatusInternalServerError, openapi.UNKNOWN, "Failed to generate CAPTCHA")
		return
	}

	ctx.JSON(http.StatusOK, openapi.Captcha{
		CaptchaToken: result.CaptchaToken,
		CaptchaImage: result.CaptchaImage,
	})
}
