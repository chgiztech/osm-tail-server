package router

import (
	"github.com/gin-gonic/gin"
	"go.opentelemetry.io/contrib/instrumentation/github.com/gin-gonic/gin/otelgin"
	"net/http"
	"osm-tail/http_handler"
	"osm-tail/middleware"
	"osm-tail/utils/envconf"
)

func RegisterRoute(r *gin.Engine, handler http_handler.Handler) {
	r.Use(middleware.HeaderGuard())
	r.Use(otelgin.Middleware(envconf.AppName))
	r.Use(middleware.TracingMiddleware())

	r.GET("/", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "OSM-GO",
		})
	})

	r.GET("/tiles/:z/:x/:y", handler.HandleGenerate)
	r.POST("/coordinates", handler.HandleCoordinates)
}
