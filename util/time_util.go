package util

import (
	"fmt"
	"ohmycontrolcenter.tech/omcc/internal/infrastructure/config"
	"strconv"
	"time"
)

var GlobalTimeConfig *config.TimeFormatConfig

// InitTimeUtil initTimeUtil
func init() {
	loc, _ := time.LoadLocation("Asia/Taipei")
	GlobalTimeConfig = &config.TimeFormatConfig{
		TimeFormat:   "2006-01-02 15:04:05",
		DateFormat:   "2006-01-02",
		TimeLocation: loc,
	}
}

// static functions

func ToIsoTimeStringFormat(epochString string) (string, error) {
	ms, err := strconv.ParseInt(epochString, 10, 64)
	if err != nil {
		return "", fmt.Errorf("invalid second timestamp: %v", err)
	}
	return time.UnixMilli(ms).In(GlobalTimeConfig.TimeLocation).Format(GlobalTimeConfig.TimeFormat), nil
}

func ToIsoTimeFormat(epochString string) (time.Time, error) {
	ms, err := strconv.ParseInt(epochString, 10, 64)
	if err != nil {
		return time.Time{}, fmt.Errorf("invalid second timestamp: %v", err)
	}
	return time.UnixMilli(ms).In(GlobalTimeConfig.TimeLocation), nil
}

func CurrentDate() string {
	return time.Now().In(GlobalTimeConfig.TimeLocation).Format(GlobalTimeConfig.TimeFormat)
}

func CurrentDateInFormat(format string) string {
	return time.Now().In(GlobalTimeConfig.TimeLocation).Format(format)
}

func FirstDateInCurrentMonth() string {
	now := time.Now().In(GlobalTimeConfig.TimeLocation)
	firstDay := time.Date(now.Year(), now.Month(), 1, 0, 0, 0, 0, GlobalTimeConfig.TimeLocation)
	return firstDay.Format(GlobalTimeConfig.TimeFormat)
}

func CurrentDateToEpoch() string {
	return strconv.FormatInt(time.Now().In(GlobalTimeConfig.TimeLocation).UnixMilli(), 10)
}

func FirstDateInCurrentMonthToEpoch() string {
	now := time.Now().In(GlobalTimeConfig.TimeLocation)
	firstDay := time.Date(now.Year(), now.Month(), 1, 0, 0, 0, 0, GlobalTimeConfig.TimeLocation)
	return strconv.FormatInt(firstDay.UnixMilli(), 10)
}

func FirstAndLastDateOfLastMonth() (time.Time, time.Time) {
	now := time.Now().In(GlobalTimeConfig.TimeLocation)
	firstDayOfCurrentMonth := time.Date(now.Year(), now.Month(), 1, 0, 0, 0, 0, GlobalTimeConfig.TimeLocation)
	firstDayOfLastMonth := firstDayOfCurrentMonth.AddDate(0, -1, 0)
	lastDayOfLastMonth := firstDayOfLastMonth.AddDate(0, 0, -1)
	return firstDayOfLastMonth, lastDayOfLastMonth
}

// some support functions

func FormatTime(t time.Time) string {
	return t.In(GlobalTimeConfig.TimeLocation).Format(GlobalTimeConfig.TimeFormat)
}

func FormatDate(t time.Time) string {
	return t.In(GlobalTimeConfig.TimeLocation).Format(GlobalTimeConfig.DateFormat)
}

func ParseTime(timeStr string) (time.Time, error) {
	return time.ParseInLocation(GlobalTimeConfig.TimeFormat, timeStr, GlobalTimeConfig.TimeLocation)
}

func Location() *time.Location {
	return GlobalTimeConfig.TimeLocation
}
