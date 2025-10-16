package utils

import "time"

// RunTime 计算运行时间
// return 毫秒
func RunTime(startTime time.Time) int64 {
	return time.Since(startTime).Milliseconds()
}

// StartTime 获取当前时间
func StartTime() time.Time {
	return time.Now()
}
