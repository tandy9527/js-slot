package game

import (
	"fmt"
	"math/rand"
	"os"

	"gopkg.in/yaml.v3"
)

var slotParams *SlotParams

// slot 游戏数值配置
type SlotParams struct {
	GameInfo []GameInfo `yaml:"slotParams"`
}

// GameInfo 游戏数值相关信息
// 根据数值重写新增相关结构体
type GameInfo struct {
	GameID   string `yaml:"game_id"`   // 游戏id -必须
	GameName string `yaml:"game_name"` // 游戏名称 -必须
	ProbID   string `yaml:"ProbId"`    // 概率id-会根据该id读取不同的概率配置 -必须
}

func LoadGameConfig(filePath string) (*SlotParams, error) {
	data, err := os.ReadFile(filePath)
	if err != nil {
		return nil, fmt.Errorf("failed to read config file: %w", err)
	}

	var config SlotParams
	if err := yaml.Unmarshal(data, &config); err != nil {
		return nil, fmt.Errorf("failed to unmarshal config: %w", err)
	}
	slotParams = &config
	return &config, nil
}

// GetGameInfoByProbID 根据ProbID获取游戏配置
func (gc *SlotParams) GetGameInfoByProbID(probID string) (*GameInfo, error) {
	for i := range gc.GameInfo {
		if gc.GameInfo[i].ProbID == probID {
			return &gc.GameInfo[i], nil
		}
	}
	return nil, fmt.Errorf("game config not found for ProbID: %s", probID)
}

func GetGameInfo() *GameInfo {
	if len(slotParams.GameInfo) == 0 {
		return nil
	}
	idx := rand.Intn(len(slotParams.GameInfo))
	return &slotParams.GameInfo[idx]

}
