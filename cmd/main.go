package main

import (
	"fmt"

	"github.com/tandy9527/js-slot/core/game"
)

func main() {
	//game.GlobalManager()

	// // 注册 FortuneCat 游戏
	// fc := fortuneCat.NewFortuneCat()
	// manager.RegisterGame(fc)

	// t := test.NewTest()
	// manager.RegisterGame(t)
	// port := fmt.Sprintf(":%d", core.GConf.Port)
	// logger.Infof("start server Game: %v , port: %d", core.GConf.GameCode, core.GConf.Port)
	// http.HandleFunc("/game", ws.WsHandler)
	// if err := http.ListenAndServe(port, nil); err != nil {
	// 	logger.Errorf("server error: %v", err)
	// }
	// server.Start()
	game.LoadGameConfig("config/slot_params.yaml")
	fmt.Println(game.GetGameInfo().GetIntSliceMap("reel_info", "MG", "reel_symbol")["1"][0])
	fmt.Println(game.GetGameInfo().GetBool("test"))
}

// func init() {
// 	core.LoadGameConf("config/game.yaml")

// 	logger.LoggerInit(core.GConf.LogPath, 50, 30, 100, true)

// 	core.LoadRedis("config/redis.yaml")

// }
