package util

import "time"

// ParseDatetime parses a string in "yyyy-MM-dd HH:mm:ss" format into time.Time
func ParseDatetime(datetimeStr string) (time.Time, error) {
	layout := "2006-01-02 15:04:05"
	return time.Parse(layout, datetimeStr)
}

// ParseDate parses a string in "yyyy-MM-dd" format into time.Time
func ParseDate(dateStr string) (time.Time, error) {
	layout := "2006-01-02"
	return time.Parse(layout, dateStr)
}
