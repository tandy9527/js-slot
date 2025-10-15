package game

import (
	"fmt"
	"math/rand"
	"os"

	"github.com/tandy9527/js-util/logger"
	"gopkg.in/yaml.v3"
)

type IGame interface {
	GetGameInfo() any
}

var slotParams *SlotParams

// slot 游戏数值配置
type SlotParams struct {
	GameInfo []GameInfo `yaml:"slotParams"`
}

// GameInfo 游戏数值相关信息
// 根据数值重写新增相关结构体
type GameInfo struct {
	GameID   string         `yaml:"game_id"`   // 游戏id -必须
	GameName string         `yaml:"game_name"` // 游戏名称 -必须
	ProbID   string         `yaml:"ProbId"`    // 概率id-会根据该id读取不同的概率配置 -必须
	Extra    map[string]any `yaml:"extra"`     // 扩展字段
}

func LoadGameConfig(filePath string) (*SlotParams, error) {
	data, err := os.ReadFile(filePath)
	if err != nil {
		logger.Errorf("[gameinfo] load failed: %v", err)
		return nil, fmt.Errorf("failed to read config file: %w", err)
	}

	var config SlotParams
	if err := yaml.Unmarshal(data, &config); err != nil {
		logger.Errorf("[gameinfo] load failed: %v", err)
		return nil, fmt.Errorf("failed to unmarshal config: %w", err)
	}
	slotParams = &config
	logger.Infof("[gameinfo] load successful")
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

// GetGameInfo  会随机取一个数值配置
func GetGameInfo() *GameInfo {
	if len(slotParams.GameInfo) == 0 {
		return nil
	}
	idx := rand.Intn(len(slotParams.GameInfo))
	return &slotParams.GameInfo[idx]

}

// GetBool 获取嵌套字段 bool
func (g *GameInfo) GetBool(path ...string) bool {
	val := getNested(g.Extra, path...)

	if b, ok := val.(bool); ok {
		return b
	}
	return false
}

// GetString 获取嵌套字段字符串
func (g *GameInfo) GetString(path ...string) string {
	val := getNested(g.Extra, path...)
	if str, ok := val.(string); ok {
		return str
	}
	return ""
}

// GetInt 获取嵌套字段 int
func (g *GameInfo) GetInt(path ...string) int {
	val := getNested(g.Extra, path...)
	if num, ok := val.(int); ok {
		return num
	}
	// YAML 解析后可能是 float64
	if f, ok := val.(float64); ok {
		return int(f)
	}
	return 0
}
func (g *GameInfo) GetFloat(path ...string) float64 {
	val := getNested(g.Extra, path...)
	if num, ok := val.(float64); ok {
		return num
	}
	// YAML 解析后可能是 float64
	if f, ok := val.(float64); ok {
		return f
	}
	return 0
}

// GetStringSlice 获取嵌套字段的字符串数组
func (g *GameInfo) GetStringSlice(path ...string) []string {
	val := getNested(g.Extra, path...)
	arr, ok := val.([]any)
	if !ok {
		return nil
	}

	result := make([]string, 0, len(arr))
	for _, v := range arr {
		if s, ok := v.(string); ok {
			result = append(result, s)
		}
	}
	return result
}

// GetIntSlice 获取嵌套字段的整数数组
// return []int
func (g *GameInfo) GetIntSlice(path ...string) []int {
	val := getNested(g.Extra, path...)
	arr, ok := val.([]any)
	if !ok {
		return nil
	}

	result := make([]int, 0, len(arr))
	for _, v := range arr {
		switch num := v.(type) {
		case int:
			result = append(result, num)
		case int64:
			result = append(result, int(num))
		case float64: // YAML 默认解析数字为 float64
			result = append(result, int(num))
		}
	}
	return result
}

// return [][]int
func (g *GameInfo) GetIntMatrix(path ...string) [][]int {
	val := getNested(g.Extra, path...)
	arr, ok := val.([]any)
	if !ok {
		return nil
	}

	result := make([][]int, 0, len(arr))
	for _, item := range arr {
		subArr, ok := item.([]any)
		if !ok {
			continue
		}

		row := make([]int, 0, len(subArr))
		for _, v := range subArr {
			switch num := v.(type) {
			case int:
				row = append(row, num)
			case int64:
				row = append(row, int(num))
			case float64:
				row = append(row, int(num)) // YAML 默认是 float64
			}
		}

		result = append(result, row)
	}
	return result
}

// GetIntSliceMap 获取嵌套字段（map[string][]int）结构
func (g *GameInfo) GetIntSliceMap(path ...string) map[string][]int {
	val := getNested(g.Extra, path...)

	// 类型必须是 map[string]any
	rawMap, ok := val.(map[string]any)
	if !ok {
		return nil
	}

	result := make(map[string][]int)
	for key, v := range rawMap {
		arr, ok := v.([]any)
		if !ok {
			continue
		}
		ints := make([]int, 0, len(arr))
		for _, num := range arr {
			switch n := num.(type) {
			case int:
				ints = append(ints, n)
			case float64:
				ints = append(ints, int(n))
			}
		}
		result[key] = ints
	}
	return result
}

// GetFloatSliceMap 获取嵌套字段（map[string][]float64）结构
func (g *GameInfo) GetFloatSliceMap(path ...string) map[string][]float64 {
	val := getNested(g.Extra, path...)

	// 类型必须是 map[string]any
	rawMap, ok := val.(map[string]any)
	if !ok {
		return nil
	}

	result := make(map[string][]float64)
	for key, v := range rawMap {
		arr, ok := v.([]any)
		if !ok {
			continue
		}
		floats := make([]float64, 0, len(arr))
		for _, num := range arr {
			switch n := num.(type) {
			case float64:
				floats = append(floats, n)
			case int:
				floats = append(floats, float64(n))
			}
		}
		result[key] = floats
	}
	return result
}

// GetStringSliceMap 获取嵌套字段（map[string][]string）结构
func (g *GameInfo) GetStringSliceMap(path ...string) map[string][]string {
	val := getNested(g.Extra, path...)

	// 类型必须是 map[string]any
	rawMap, ok := val.(map[string]any)
	if !ok {
		return nil
	}

	result := make(map[string][]string)
	for key, v := range rawMap {
		arr, ok := v.([]any)
		if !ok {
			continue
		}
		strs := make([]string, 0, len(arr))
		for _, s := range arr {
			switch v := s.(type) {
			case string:
				strs = append(strs, v)
			default:
				strs = append(strs, fmt.Sprintf("%v", v)) // 自动转字符串
			}
		}
		result[key] = strs
	}
	return result
}

// 递归读取嵌套 map
func getNested(m map[string]any, path ...string) any {
	var curr any = m
	for _, p := range path {
		if currMap, ok := curr.(map[string]any); ok {
			curr = currMap[p]
		} else {
			return nil
		}
	}
	return curr
}
