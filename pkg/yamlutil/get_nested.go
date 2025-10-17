package yamlutil

import "fmt"

func GetNested[T any](c ConfigGetter, path ...string) T {
	m := c.ExtraMap()
	var curr any = m

	for _, p := range path {
		currMap, ok := curr.(map[string]any)
		if !ok {
			panic(fmt.Sprintf("config path error: expected map before %v, got %T", path, curr))
		}
		val, exists := currMap[p]
		if !exists {
			panic(fmt.Sprintf("config not found at path: %v", path))
		}
		curr = val
	}

	if curr == nil {
		panic(fmt.Sprintf("config not found at path: %v", path))
	}

	// 类型断言 + 数字转换
	if val, ok := curr.(T); ok {
		return val
	}

	switch v := curr.(type) {
	case float64:
		var zero T
		switch any(zero).(type) {
		case int:
			return any(int(v)).(T)
		case int64:
			return any(int64(v)).(T)
		case float32:
			return any(float32(v)).(T)
		}
	case int:
		var zero T
		switch any(zero).(type) {
		case float64:
			return any(float64(v)).(T)
		case int64:
			return any(int64(v)).(T)
		}
	}

	panic(fmt.Sprintf("config type error: expected %T at path %v, got %T (value=%v)",
		*new(T), path, curr, curr))
}
