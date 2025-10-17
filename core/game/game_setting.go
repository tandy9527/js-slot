package game

import (
	"fmt"
	"os"

	"github.com/tandy9527/js-slot/pkg/yamlutil"
	"github.com/tandy9527/js-util/logger"
	"gopkg.in/yaml.v3"
)

var GSetting yamlutil.Config[*GameSetting]

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

func (g *GameSetting) ExtraMap() map[string]any {
	return g.Extra
}

func LoadGameSetting(filePath string) {
	data, err := os.ReadFile(filePath)
	if err != nil {
		panic(fmt.Errorf("failed to read config file: %w", err))
	}
	var gsetting GameSetting
	if err := yaml.Unmarshal(data, &gsetting); err != nil {
		panic(fmt.Errorf("failed to unmarshal config: %w", err))
	}
	// GSetting = &gsetting
	GSetting = yamlutil.Config[*GameSetting]{
		Cfg: &gsetting,
	}

	logger.Infof("[gameSetting] load successful")
}

// // GsString 获取嵌套字段字符串
func (g *GameSetting) GsString(path ...string) string {
	return getNested[string](g.Extra, path...)
	// val := getNested(g.Extra, path...)
	// if str, ok := val.(string); ok {
	// 	return str
	// }
	// panic(fmt.Sprintf("config type error: expected string at path %v, got %T (value=%v)", path, val, val))
}

// // GsInt 获取嵌套字段 int
// func (g *GameSetting) GsInt(path ...string) int {
// 	val := getNested(g.Extra, path...)
// 	switch num := val.(type) {
// 	case int:
// 		return num
// 	case int64:
// 		return int(num)
// 	case float64:
// 		return int(num)
// 	default:
// 		panic(fmt.Sprintf("config type error: expected int at path %v, got %T (value=%v)", path, val, val))
// 	}
// }

// // GsFloat 获取嵌套字段 float64
// func (g *GameSetting) GsFloat(path ...string) float64 {
// 	val := getNested(g.Extra, path...)
// 	switch num := val.(type) {
// 	case float64:
// 		return num
// 	case int:
// 		return float64(num)
// 	case int64:
// 		return float64(num)
// 	default:
// 		panic(fmt.Sprintf("config type error: expected float64 at path %v, got %T (value=%v)", path, val, val))
// 	}
// }

// // GsStringSlice 获取 []string
// func (g *GameSetting) GsStringSlice(path ...string) []string {
// 	arr := getNested[[]string](g.Extra, path...)
// 	// // arr, ok := val.([]any)
// 	// // if !ok {
// 	// // 	panic(fmt.Sprintf("config type error: expected []string at path %v, got %T (value=%v)", path, val, val))
// 	// // }

// 	// result := make([]string, 0, len(arr))
// 	// for _, v := range arr {
// 	// 	if s, ok := v.(string); ok {
// 	// 		result = append(result, s)
// 	// 	} else {
// 	// 		panic(fmt.Sprintf("config type error: expected string element at path %v, got %T (value=%v)", path, v, v))
// 	// 	}
// 	// }
// 	return arr
// }

// // GsIntSlice 获取 []int
// func (g *GameSetting) GsIntSlice(path ...string) []int {
// 	val := getNested(g.Extra, path...)
// 	arr, ok := val.([]any)
// 	if !ok {
// 		panic(fmt.Sprintf("config type error: expected []int at path %v, got %T (value=%v)", path, val, val))
// 	}

// 	result := make([]int, 0, len(arr))
// 	for _, v := range arr {
// 		switch num := v.(type) {
// 		case int:
// 			result = append(result, num)
// 		case int64:
// 			result = append(result, int(num))
// 		case float64:
// 			result = append(result, int(num))
// 		default:
// 			panic(fmt.Sprintf("config type error: expected numeric element at path %v, got %T (value=%v)", path, v, v))
// 		}
// 	}
// 	return result
// }

// // GsIntMatrix 获取 [][]int
// func (g *GameSetting) GsIntMatrix(path ...string) [][]int {
// 	val := getNested(g.Extra, path...)
// 	arr, ok := val.([]any)
// 	if !ok {
// 		panic(fmt.Sprintf("config type error: expected [][]int at path %v, got %T (value=%v)", path, val, val))
// 	}

// 	result := make([][]int, 0, len(arr))
// 	for _, item := range arr {
// 		subArr, ok := item.([]any)
// 		if !ok {
// 			panic(fmt.Sprintf("config type error: expected []int inside [][]int at path %v, got %T (value=%v)", path, item, item))
// 		}

// 		row := make([]int, 0, len(subArr))
// 		for _, v := range subArr {
// 			switch num := v.(type) {
// 			case int:
// 				row = append(row, num)
// 			case int64:
// 				row = append(row, int(num))
// 			case float64:
// 				row = append(row, int(num))
// 			default:
// 				panic(fmt.Sprintf("config type error: expected numeric element inside [][]int at path %v, got %T (value=%v)", path, v, v))
// 			}
// 		}
// 		result = append(result, row)
// 	}
// 	return result
// }

// // GsIntSliceMap 获取 map[string][]int
// func (g *GameSetting) GsIntSliceMap(path ...string) map[string][]int {
// 	val := getNested(g.Extra, path...)
// 	rawMap, ok := val.(map[string]any)
// 	if !ok {
// 		panic(fmt.Sprintf("config type error: expected map[string][]int at path %v, got %T (value=%v)", path, val, val))
// 	}

// 	result := make(map[string][]int)
// 	for key, v := range rawMap {
// 		arr, ok := v.([]any)
// 		if !ok {
// 			panic(fmt.Sprintf("config type error: expected []int at %v[%s], got %T (value=%v)", path, key, v, v))
// 		}
// 		ints := make([]int, 0, len(arr))
// 		for _, num := range arr {
// 			switch n := num.(type) {
// 			case int:
// 				ints = append(ints, n)
// 			case int64:
// 				ints = append(ints, int(n))
// 			case float64:
// 				ints = append(ints, int(n))
// 			default:
// 				panic(fmt.Sprintf("config type error: expected numeric element at %v[%s], got %T (value=%v)", path, key, n, n))
// 			}
// 		}
// 		result[key] = ints
// 	}
// 	return result
// }

// // GsFloatSliceMap 获取 map[string][]float64
// func (g *GameSetting) GsFloatSliceMap(path ...string) map[string][]float64 {
// 	val := getNested(g.Extra, path...)
// 	rawMap, ok := val.(map[string]any)
// 	if !ok {
// 		panic(fmt.Sprintf("config type error: expected map[string][]float64 at path %v, got %T (value=%v)", path, val, val))
// 	}

// 	result := make(map[string][]float64)
// 	for key, v := range rawMap {
// 		arr, ok := v.([]any)
// 		if !ok {
// 			panic(fmt.Sprintf("config type error: expected []float64 at %v[%s], got %T (value=%v)", path, key, v, v))
// 		}
// 		floats := make([]float64, 0, len(arr))
// 		for _, num := range arr {
// 			switch n := num.(type) {
// 			case float64:
// 				floats = append(floats, n)
// 			case int:
// 				floats = append(floats, float64(n))
// 			case int64:
// 				floats = append(floats, float64(n))
// 			default:
// 				panic(fmt.Sprintf("config type error: expected numeric element at %v[%s], got %T (value=%v)", path, key, n, n))
// 			}
// 		}
// 		result[key] = floats
// 	}
// 	return result
// }

// // GsStringSliceMap 获取 map[string][]string
// func (g *GameSetting) GsStringSliceMap(path ...string) map[string][]string {
// 	val := getNested(g.Extra, path...)
// 	rawMap, ok := val.(map[string]any)
// 	if !ok {
// 		panic(fmt.Sprintf("config type error: expected map[string][]string at path %v, got %T (value=%v)", path, val, val))
// 	}

// 	result := make(map[string][]string)
// 	for key, v := range rawMap {
// 		arr, ok := v.([]any)
// 		if !ok {
// 			panic(fmt.Sprintf("config type error: expected []string at %v[%s], got %T (value=%v)", path, key, v, v))
// 		}

// 		strs := make([]string, 0, len(arr))
// 		for _, s := range arr {
// 			switch v := s.(type) {
// 			case string:
// 				strs = append(strs, v)
// 			default:
// 				panic(fmt.Sprintf("config type error: expected string element at %v[%s], got %T (value=%v)", path, key, v, v))
// 			}
// 		}
// 		result[key] = strs
// 	}
// 	return result
// }
