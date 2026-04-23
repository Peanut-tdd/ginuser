package main

import (
	"context"
	"gin-user/config"
	"gin-user/router"
	"github.com/uptrace/opentelemetry-go-extra/otelzap"
	"go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"
	"go.uber.org/zap"
	"log"
	"net/http"
)

func main() {
	defer func() {
		if err := config.CloseLogger(); err != nil {
			log.Printf("close logger failed: %v", err)
		}
	}()

	logger := otelzap.New(
		config.ZapLogger,                    // zap实例，按需配置
		otelzap.WithMinLevel(zap.InfoLevel), // 指定日志级别
		otelzap.WithStackTrace(true),        // 在日志中记录 traceI
	)
	defer logger.Sync()

	// 替换全局的logger
	undo := otelzap.ReplaceGlobals(logger)
	defer undo()

	// 主要看这里
	tp, err := config.TracerProvider()

	if err != nil {
		log.Fatal(err)
	}

	// Cleanly shutdown and flush telemetry when the application exits.
	defer func() {
		if err := tp.Shutdown(context.Background()); err != nil {
			log.Fatal(err)
		}
	}()

	r := router.NewRouter()
	handler := otelhttp.NewHandler(r, config.Service)
	server := &http.Server{
		Addr:    config.Conf.App.Port,
		Handler: handler,
	}
	if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		log.Fatal(err)
	}

}

func init() {

	config.InitConfig()
	config.InitLogger()
	config.InitDb()
	config.InitRedis()

}
