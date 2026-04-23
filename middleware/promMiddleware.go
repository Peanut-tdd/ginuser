package middleware

import (
	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"strconv"
	"time"
)

// 定义HTTP请求总数指标
var (
	httpRequestsTotal = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Name: "http_requests_total",
			Help: "Total number of HTTP requests",
		},
		[]string{"method", "endpoint", "status"},
	)

	httpRequestDuration = promauto.NewHistogramVec(prometheus.HistogramOpts{
		Name:    "http_request_duration_seconds",
		Help:    "Duration of HTTP requests",
		Buckets: []float64{0.1, 0.3, 0.5, 0.7, 1, 1.5, 2, 3},
	}, []string{"method", "path"})
)

func PrometheusMiddleware() gin.HandlerFunc {
	return func(ctx *gin.Context) {

		// 不处理指标端点路径
		if ctx.Request.URL.Path == "/metrics" {
			// 直接跳过
			ctx.Next()
			return
		}

		startTime := time.Now()

		// 处理请求
		ctx.Next()
		duration := time.Since(startTime).Seconds()

		// 请求处理完成后记录指标
		status := strconv.Itoa(ctx.Writer.Status())

		httpRequestsTotal.WithLabelValues(
			ctx.Request.Method, // HTTP方法
			ctx.FullPath(),     // 路由路径
			status,             // 状态码
		).Inc() // 指标值加1

		httpRequestDuration.WithLabelValues(
			ctx.Request.Method,
			ctx.FullPath()).Observe(duration)
	}
}
