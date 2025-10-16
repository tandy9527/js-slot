// js-slot/server/server.go
package server

import (
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"syscall"
	"time"

	"github.com/tandy9527/js-slot/core"
	"github.com/tandy9527/js-slot/core/game"

	"github.com/tandy9527/js-slot/pkg/consts"
	"github.com/tandy9527/js-slot/pkg/utils"
	"github.com/tandy9527/js-slot/transport/ws"

	"github.com/tandy9527/js-util/cache"
	"github.com/tandy9527/js-util/logger"
)

var srv *http.Server

// Start 启动游戏服务
func Start() error {

	port := fmt.Sprintf(":%d", core.GConf.Port)
	logger.Infof("start server Game: %v , port: %d，router:%v", core.GConf.GameCode, core.GConf.Port, core.GConf.RouterName)

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
	router := core.GConf.RouterName
	if utils.IsEmpty(router) {
		router = "game"
	}
	mux := http.NewServeMux()
	mux.HandleFunc("/"+router, ws.WsHandler)
	return mux
}

// handleShutdown 监听退出信号并优雅关闭
func handleShutdown() {
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	sig := <-quit
	logger.Infof("shutdown signal received:%v", sig)

	cleanup()

	PrintShutdownLog(core.GConf.GameCode)
	os.Exit(0)
}

// 清理资源
func cleanup() {
	logger.Info("cleaup")
	cache.CloseRedis()
	core.GlobalConnManager().CloseAll()
}

// Init 初始化配置
func init() {
	core.LoadGameConf("config/game.yaml")
	logger.LoggerInit(core.GConf.LogPath, 50, 30, 100, true)
	PrintStartupLog(core.GConf.GameCode)
	cache.LoadRedis("config/redis.yaml")

	game.LoadGameConfig("config/slot_game_info.yaml")

	cleanGame()

}

func cleanGame() {
	logger.Info("clean game")
	// 在线人数
	cache.GetDB("db2").ZRem(consts.REDIS_GAME_ONLINE, strconv.Itoa(core.GConf.GameID))
	// 游戏在线人
	cache.GetDB("db2").Del(consts.REDIS_GAME_CONN)
}

// PrintStartupLog 启动日志
func PrintStartupLog(instanceID string) {
	logger.Infof("=========================================================================")
	logger.Infof(" GAME STARTED | ID: %s | TIME: %s", instanceID, time.Now().Format("2006-01-02 15:04:05"))
	logger.Infof("=========================================================================")
}

// PrintShutdownLog 退出日志
func PrintShutdownLog(instanceID string) {
	logger.Infof("=========================================================================")
	logger.Infof(" GAME EXITED | ID: %s | TIME: %s", instanceID, time.Now().Format("2006-01-02 15:04:05"))
	logger.Infof("=========================================================================")
}
