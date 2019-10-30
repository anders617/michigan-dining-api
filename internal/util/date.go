package date

import (
	"time"
)

const MDiningDateLayout = "2006-01-02T15:04:05-07:00"

func Parse(s *string) (time.Time, error) {
	return time.Parse(MDiningDateLayout, *s)
}

func Format(dateTime time.Time) string {
	return dateTime.Format(MDiningDateLayout)
}
