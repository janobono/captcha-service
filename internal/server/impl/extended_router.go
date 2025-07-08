package impl

import (
	"github.com/janobono/captcha-service/generated/openapi"
	"net/http"

	"github.com/gin-gonic/gin"
)

// NewRouter returns a new router.
func NewRouter(handleFunctions openapi.ApiHandleFunctions, contextPath string) *gin.Engine {
	return NewRouterWithGinEngine(gin.Default(), handleFunctions, contextPath)
}

// NewRouterWithGinEngine add routes to existing gin engine.
func NewRouterWithGinEngine(router *gin.Engine, handleFunctions openapi.ApiHandleFunctions, contextPath string) *gin.Engine {
	group := router.Group(contextPath)

	for _, route := range getRoutes(handleFunctions) {
		handler := route.HandlerFunc
		if handler == nil {
			handler = openapi.DefaultHandleFunc
		}
		switch route.Method {
		case http.MethodGet:
			group.GET(route.Pattern, handler)
		case http.MethodPost:
			group.POST(route.Pattern, handler)
		case http.MethodPut:
			group.PUT(route.Pattern, handler)
		case http.MethodPatch:
			group.PATCH(route.Pattern, handler)
		case http.MethodDelete:
			group.DELETE(route.Pattern, handler)
		}
	}

	return router
}

func getRoutes(handleFunctions openapi.ApiHandleFunctions) []openapi.Route {
	return []openapi.Route{
		{
			"GetCaptcha",
			http.MethodGet,
			"/captcha",
			handleFunctions.CaptchaControllerAPI.GetCaptcha,
		},
		{
			"ValidateCaptcha",
			http.MethodPost,
			"/captcha",
			handleFunctions.CaptchaControllerAPI.ValidateCaptcha,
		},
		{
			"Livez",
			http.MethodGet,
			"/livez",
			handleFunctions.HealthControllerAPI.Livez,
		},
		{
			"Readyz",
			http.MethodGet,
			"/readyz",
			handleFunctions.HealthControllerAPI.Readyz,
		},
	}
}
