package impl

import (
	"github.com/gin-gonic/gin"
	"github.com/janobono/captcha-service/generated/openapi"
	"net/http"
)

type healthController struct {
}

func NewHealthController() openapi.HealthControllerAPI {
	return &healthController{}
}

func (h healthController) Livez(ctx *gin.Context) {
	ctx.JSON(http.StatusOK, openapi.HealthStatus{Status: "UP"})
}

func (h healthController) Readyz(ctx *gin.Context) {
	ctx.JSON(http.StatusOK, openapi.HealthStatus{Status: "READY"})
}
