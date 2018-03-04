package utils

import (
	"time"
)

const (
	FORMAT_DATE     = "2006-01-02"
	FORMAT_DATETIME = "2006-01-02 15:04:05"
)

func Today() string {
	return time.Now().Format(FORMAT_DATE)
}

func Now() string {
	return time.Now().Format(FORMAT_DATETIME)
}

func SinceMinutes(fromtime string) int {
	fromTime, err := time.ParseInLocation(FORMAT_DATETIME, fromtime, time.Local)
	if err != nil {
		return -1
	}
	return int(time.Since(fromTime).Minutes())
}
