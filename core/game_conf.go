package core

import (
	"log"
	"os"

	"gopkg.in/yaml.v3"
)

type GameConf struct {
	GameCode     string `yaml:"game_code"`
	GameName     string `yaml:"game_name"`
	Port         int    `yaml:"port"`
	RoomMaxUsers int    `yaml:"room_max_users"`
	LogPath      string `yaml:"log_path"`
}
type GameConfig struct {
	Game GameConf `yaml:"game"`
}

var GConf *GameConf

func LoadGameConf(path string) error {
	file, err := os.ReadFile(path)
	if err != nil {
		log.Fatalf("read config file error: %v", err)
		return err
	}
	var cfg GameConfig
	if err := yaml.Unmarshal(file, &cfg); err != nil {
		log.Fatalf("unmarshal config error: %v", err)
		return err
	}
	GConf = &cfg.Game
	return nil
}
