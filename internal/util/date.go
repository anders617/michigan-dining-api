package date

import (
	"time"
)

const MDiningDateLayout = "2006-01-02T15:04:05-07:00"

var USEasternLocation, _ = time.LoadLocation("America/Detroit")

func Now() time.Time {
	return time.Now().In(USEasternLocation)
}

func Parse(s *string) (time.Time, error) {
	return time.ParseInLocation(MDiningDateLayout, *s, USEasternLocation)
}

func Format(t time.Time) string {
	return t.In(USEasternLocation).Format(MDiningDateLayout)
}

func DayStart(t time.Time) time.Time {
	return time.Date(t.Year(), t.Month(), t.Day(), 0, 0, 0, 0, USEasternLocation)
}

func DayEnd(t time.Time) time.Time {
	return time.Date(t.Year(), t.Month(), t.Day(), 23, 59, 59, 999999999, USEasternLocation)
}
