package router

import (
	"gin-user/config"
	"gin-user/controller"
	"gin-user/middleware"
	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"go.opentelemetry.io/contrib/instrumentation/github.com/gin-gonic/gin/otelgin"
)

func NewRouter() *gin.Engine {

	r := gin.Default()

	r.Use(otelgin.Middleware(config.Service))
	r.Use(middleware.RequestTraceLog())

	// 应用Prometheus中间件到所有路由
	//r.Use(middleware.PrometheusMiddleware())

	// 暴露指标端点
	r.GET("/metrics", gin.WrapH(promhttp.Handler()))
	//http.Handle("/metrics", promhttp.Handler())

	v1 := r.Group("/api/v1")
	{
		v1.POST("/register", controller.UserRegisterHandler())
		v1.POST("/login", controller.UserLoginUserHandler())

		auth := v1.Group("/")
		auth.Use(middleware.AuthMiddleware())

		{
			auth.POST("/refresh", controller.UserRefreshTokenHandler())
			auth.POST("/logout", controller.UserLogoutHandler())
			auth.POST("/user/info", controller.UserInfoHandler())
		}

	}

	return r

}
