// js-slot/server/server.go
package server

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/tandy9527/js-slot/core"
	"github.com/tandy9527/js-slot/pkg/logger"
	"github.com/tandy9527/js-slot/transport/ws"
)

var srv *http.Server

// Start 启动游戏服务
func Start() error {
	port := fmt.Sprintf(":%d", core.GConf.Port)
	logger.Infof("start server Game: %v , port: %d", core.GConf.GameCode, core.GConf.Port)

	// 配置 HTTP Server
	srv = &http.Server{
		Addr:    port,
		Handler: setupRouter(),
	}

	// 监听退出信号
	go handleShutdown()

	// 启动服务
	if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		logger.Info("server error: %w", err)
		return fmt.Errorf("server error: %w", err)
	}
	return nil
}

// setupRouter 配置路由
func setupRouter() http.Handler {
	mux := http.NewServeMux()
	mux.HandleFunc("/game", ws.WsHandler)
	return mux
}

// handleShutdown 监听退出信号并优雅关闭
func handleShutdown() {
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	<-quit
	logger.Infof("shutdown signal received...")

	// 给 5 秒时间处理完请求
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		logger.Errorf("server forced to shutdown: %v", err)
	}

	close()

	logger.Infof("server exited gracefully")
}
func close() {
	core.CloseRedis()
}

// Init 初始化配置
func init() {
	core.LoadGameConf("config/game.yaml")
	logger.LoggerInit(core.GConf.LogPath, 50, 30, 100, true)
	core.LoadRedis("config/redis.yaml")
}
