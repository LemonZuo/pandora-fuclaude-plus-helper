package util

import "time"

const layout = "2006-01-02 15:04:05"

// ParseTime 从字符串解析时间
func ParseTime(timeStr string) (time.Time, error) {
	return time.Parse(layout, timeStr)
}

// FormatTime 将时间格式化为字符串
func FormatTime(t time.Time) string {
	return t.Format(layout)
}
