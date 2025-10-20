package game

import (
	"fmt"
	"math/rand"

	"github.com/tandy9527/js-util/tools"
)

var GIS *GameInfos

// slot 游戏数值配置
// 可能有多份不同的数值
// 每次room 会随机取一个
type GameInfos struct {
	GameInfos []GameInfo `yaml:"GameInfos"`
}

// GameInfo数值配置-> slot_game_info.yaml 文件名字需一致
// 具体游戏根据数值配置自行配置
type GameInfo struct {
	GameID   string         `yaml:"game_id"`   // 游戏id -必须
	GameName string         `yaml:"game_name"` // 游戏名称 -必须
	ProbID   string         `yaml:"ProbId"`    // 概率id-会根据该id读取不同的概率配置 -必须
	Extra    map[string]any `yaml:"extra"`     // 扩展字段
}

func LoadGameConfig(filePath string) {
	GIS = tools.Loadyaml[GameInfos](filePath)
}

// GetGameInfoByProbID 根据ProbID获取游戏配置
func (gc *GameInfos) GetGameInfoByProbID(probID string) (*GameInfo, error) {
	for i := range gc.GameInfos {
		if gc.GameInfos[i].ProbID == probID {
			return &gc.GameInfos[i], nil
		}
	}
	return nil, fmt.Errorf("game config not found for ProbID: %s", probID)
}

// GetGameInfo  会随机取一个数值配置
func GetGameInfo() *GameInfo {
	if len(GIS.GameInfos) == 0 {
		panic("no GameInfo available")
	}
	idx := rand.Intn(len(GIS.GameInfos))
	return &GIS.GameInfos[idx]
}

// --------------------------为了操作时方便--------------------------

// // GetString 获取嵌套字段字符串
func (g *GameInfo) GetString(path ...string) string {
	return tools.GetNested[string](g.Extra, path...)
}
func (g *GameInfo) GetBool(path ...string) bool {
	return tools.GetNested[bool](g.Extra, path...)
}

// GetInt 获取嵌套字段 int
func (g *GameInfo) GetInt(path ...string) int {
	return tools.GetNested[int](g.Extra, path...)
}

// GetFloat 获取嵌套字段 float64
func (g *GameInfo) GetFloat(path ...string) float64 {
	return tools.GetNested[float64](g.Extra, path...)
}

// GetStrinGetlice 获取 []string
func (g *GameInfo) GetStrinGetlice(path ...string) []string {
	return tools.GetNested[[]string](g.Extra, path...)
}

// GetIntSlice 获取 []int
func (g *GameInfo) GetIntSlice(path ...string) []int {
	return tools.GetNested[[]int](g.Extra, path...)
}

// GetIntMatrix 获取 [][]int
func (g *GameInfo) GetIntMatrix(path ...string) [][]int {
	return tools.GetNested[[][]int](g.Extra, path...)
}

// GetIntSliceMap 获取 map[string][]int
func (g *GameInfo) GetIntSliceMap(path ...string) map[string][]int {
	return tools.GetNested[map[string][]int](g.Extra, path...)
}

// GetFloatSliceMap 获取 map[string][]float64
func (g *GameInfo) GetFloatSliceMap(path ...string) map[string][]float64 {
	return tools.GetNested[map[string][]float64](g.Extra, path...)
}

// GetStrinGetliceMap 获取 map[string][]string
func (g *GameInfo) GetStrinGetliceMap(path ...string) map[string][]string {
	return tools.GetNested[map[string][]string](g.Extra, path...)
}

// GetIntIntSliceMap 获取 map[int][]int
func (g *GameInfo) GetIntIntSliceMap(path ...string) map[int][]int {
	return tools.GetNested[map[int][]int](g.Extra, path...)
}
