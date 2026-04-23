package middleware

import (
	"github.com/gin-gonic/gin"
	"go.opentelemetry.io/otel/trace"
)

// RequestTraceLog logs request lifecycle fields and trace correlation IDs.
func RequestTraceLog() gin.HandlerFunc {
	return func(c *gin.Context) {
		//start := time.Now()

		spanCtx := trace.SpanContextFromContext(c.Request.Context())
		traceID := ""
		//spanID := ""
		if spanCtx.IsValid() {
			traceID = spanCtx.TraceID().String()
			//spanID = spanCtx.SpanID().String()
		}

		c.Set("traceID", traceID)

		//otelzap.Ctx(c.Request.Context()).Info("request_start",
		//	zap.String("trace_id", traceID),
		//	zap.String("span_id", spanID),
		//	zap.String("method", c.Request.Method),
		//	zap.String("path", c.FullPath()),
		//	zap.String("client_ip", c.ClientIP()),
		//)

		c.Next()

		//latency := time.Since(start)
		//otelzap.Ctx(c.Request.Context()).Info("request_end",
		//	zap.String("trace_id", traceID),
		//	zap.String("span_id", spanID),
		//	zap.String("method", c.Request.Method),
		//	zap.String("path", c.FullPath()),
		//	zap.Int("status", c.Writer.Status()),
		//	zap.Int64("latency_ms", latency.Milliseconds()),
		//	zap.String("error", c.Errors.ByType(gin.ErrorTypePrivate).String()),
		//)
	}
}
