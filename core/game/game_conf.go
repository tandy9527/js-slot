package game

import (
	"github.com/tandy9527/js-util/tools"
)

type GameConf struct {
	GameID       int    `yaml:"game_id"`
	GameCode     string `yaml:"game_code"`
	GameName     string `yaml:"game_name"`
	Port         int    `yaml:"port"`
	RoomMaxUsers int    `yaml:"room_max_users"`
	LogPath      string `yaml:"log_path"`
	RouterName   string `yaml:"router_name"`
	Mode         string `yaml:"mode"`
}
type GameConfig struct {
	Game GameConf `yaml:"game"`
}

var GConf *GameConf

func LoadGameConf(path string) {
	// file, err := os.ReadFile(path)
	// if err != nil {
	// 	panic(fmt.Sprintf("read config file error: %v", err))
	// }
	// var cfg GameConfig
	// if err := yaml.Unmarshal(file, &cfg); err != nil {
	// 	panic(fmt.Sprintf("read config file error: %v", err))
	// }
	// GConf = &cfg.Game
	GConf = &tools.Loadyaml[GameConfig](path).Game
}
