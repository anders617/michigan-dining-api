package date

import (
	"time"
)

const MDiningDateLayout = "2006-01-02T15:04:05-07:00"
const MDiningDateNoTimeLayout = "2006-01-02"

var USEasternLocation, _ = time.LoadLocation("America/Detroit")

func Now() time.Time {
	return time.Now().In(USEasternLocation)
}

func Parse(s *string) (time.Time, error) {
	return time.ParseInLocation(MDiningDateLayout, *s, USEasternLocation)
}

func ParseNoTime(s *string) (time.Time, error) {
	return time.ParseInLocation(MDiningDateNoTimeLayout, *s, USEasternLocation)
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

func FormatNoTime(t time.Time) string {
	return t.In(USEasternLocation).Format(MDiningDateNoTimeLayout)
}

func fetchTimeOnDate(t time.Time) time.Time {
	utc := t.In(time.UTC)
	return time.Date(utc.Year(), utc.Month(), utc.Day(), 6, 30, 0, 0, time.UTC).In(USEasternLocation)
}

func NextFetchTime() time.Time {
	now := Now()
	fetchTime := fetchTimeOnDate(now)
	// If fetch time is still to come today, return it
	if fetchTime.After(now) {
		return fetchTime
	}
	// If fetch time has already happened, return tomorrow fetch time
	return fetchTimeOnDate(now.AddDate(0, 0, 1))

}
