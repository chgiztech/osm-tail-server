package middleware

import (
	"github.com/gin-gonic/gin"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
	"osm-tail/utils/envconf"
)

func HeaderGuard() gin.HandlerFunc {
	return func(c *gin.Context) {
		if envconf.App.EnableHeaderValidation {
			xForwardedFor := c.GetHeader("x-forwarded-for")
			if xForwardedFor == "" {
				c.AbortWithStatusJSON(401, gin.H{"error": "Unauthorized. Missing x-forwarded-for header"})
				return
			}

		}
		c.Next()
	}
}

func TracingMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx := c.Request.Context()
		span := trace.SpanFromContext(ctx)
		if span != nil {
			traceID := c.GetHeader("x-trace-id")
			customTraceID := c.GetHeader("x-custom-trace-id")
			if customTraceID != "" {
				span.SetAttributes(attribute.String("http.x-trace-id", customTraceID))
			} else if traceID != "" {
				span.SetAttributes(attribute.String("http.x-trace-id", traceID))
			}
		}
		c.Next()
	}
}
