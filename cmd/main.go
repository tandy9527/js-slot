package main

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
	// game.LoadGameSetting("config/game_setting.yaml")
	// fmt.Println(game.GSetting.GetString("test", "test2"))
	// game.LoadGameConfig("config/game_info.yaml")
	// fmt.Println(game.GetGameInfo().GetInt("MaxOdds"))
	// fmt.Println(game.GIS.GameInfos[0].GameID)
	// fmt.Println(game.GSetting.GsStringSlice("test2", "test2"))
}

// func init() {
// 	core.LoadGameConf("config/game.yaml")go

// 	logger.LoggerInit(core.GConf.LogPath, 50, 30, 100, true)

// 	core.LoadRedis("config/redis.yaml")

// }
