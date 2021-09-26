package standardroutes

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// NotFoundRouteHandler supports a reusable not found setup.
func NotFoundRouteHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.JSON(http.StatusNotFound, NotFoundBody())
	}
}

// NotFoundBody returns the body of the response for not found route.
// Could be extended to have multi lang support here.
func NotFoundBody() gin.H {
	return gin.H{"status": "error", "error": "not found"}
}

// EmptyGoodHandler returns an empty good request for stubbed endpoints.
func EmptyGoodHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "ok"})
	}
}
