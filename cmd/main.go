package main

import (
	"fmt"
	"js-slot/internal/core"
	_ "js-slot/internal/games/fortuneCat"
	"js-slot/internal/transport/ws"
	logger "js-slot/pkg/logger"
	"net/http"
)

func main() {
	//game.GlobalManager()

	// // 注册 FortuneCat 游戏
	// fc := fortuneCat.NewFortuneCat()
	// manager.RegisterGame(fc)

	// t := test.NewTest()
	// manager.RegisterGame(t)
	port := fmt.Sprintf(":%d", core.GConf.Port)
	logger.Infof("start server Game: %v , port: %d", core.GConf.GameCode, core.GConf.Port)
	http.HandleFunc("/game", ws.WsHandler)
	if err := http.ListenAndServe(port, nil); err != nil {
		logger.Errorf("server error: %v", err)
	}

}

func init() {
	core.LoadGameConf("config/game.yaml")

	logger.LoggerInit(core.GConf.LogPath, 50, 30, 100, true)

	core.LoadRedis("config/redis.yaml")

}
