package game

import (
	"github.com/tandy9527/js-util/tools"
)

// 游戏设置配置
var GS *GameSetting

// GameSetting 游戏设置
type GameSetting struct {
	// 每局游戏的时间限制
	GameID   string `yaml:"game_id"`
	GameName string `yaml:"game_name"`
	SlotType int    `yaml:"slot_type"`
	// Odds     map[int][]int  `yaml:"odds"`      //符号:赔率- 连线后该符号的赔率,如下是3连线的赔率,4连 10:[0,0,88,888]=3连88 4连888
	// LineData [][]int        `yaml:"line_data"` //连线数据
	Extra map[string]any `yaml:"extra"` //扩展设置
}

func LoadGameSetting(filePath string) {
	GS = tools.Loadyaml[GameSetting](filePath)
}

// --------------------------为了操作时方便--------------------------

// // GetString 获取嵌套字段字符串
func (g *GameSetting) GetString(path ...string) string {
	return tools.GetNested[string](g.Extra, path...)
}
func (g *GameSetting) GetBool(path ...string) bool {
	return tools.GetNested[bool](g.Extra, path...)
}

// GetInt 获取嵌套字段 int
func (g *GameSetting) GetInt(path ...string) int {
	return tools.GetNested[int](g.Extra, path...)
}

// GetFloat 获取嵌套字段 float64
func (g *GameSetting) GetFloat(path ...string) float64 {
	return tools.GetNested[float64](g.Extra, path...)
}

// GetStrinGetlice 获取 []string
func (g *GameSetting) GetStrinGetlice(path ...string) []string {
	return tools.GetNested[[]string](g.Extra, path...)
}

// GetIntSlice 获取 []int
func (g *GameSetting) GetIntSlice(path ...string) []int {
	return tools.GetNested[[]int](g.Extra, path...)
}

// GetIntMatrix 获取 [][]int
func (g *GameSetting) GetIntMatrix(path ...string) [][]int {
	return tools.GetNested[[][]int](g.Extra, path...)
}

// GetIntSliceMap 获取 map[string][]int
func (g *GameSetting) GetIntSliceMap(path ...string) map[string][]int {
	return tools.GetNested[map[string][]int](g.Extra, path...)
}

// GetFloatSliceMap 获取 map[string][]float64
func (g *GameSetting) GetFloatSliceMap(path ...string) map[string][]float64 {
	return tools.GetNested[map[string][]float64](g.Extra, path...)
}

// GetStrinGetliceMap 获取 map[string][]string
func (g *GameSetting) GetStrinGetliceMap(path ...string) map[string][]string {
	return tools.GetNested[map[string][]string](g.Extra, path...)
}

// GetIntIntSliceMap 获取 map[int][]int
func (g *GameSetting) GetIntIntSliceMap(path ...string) map[int][]int {
	return tools.GetNested[map[int][]int](g.Extra, path...)
}
