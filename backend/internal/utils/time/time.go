package time

import (
	"strings"
	"time"
)

func EuToTime(StringDate string) (time.Time, error) {
	date, err := time.Parse("02.01.2006", StringDate)
	return date, err
}

// 2023-01-30 01:00:00
func StrToTime(t2 string) (time.Time, error) {
	t2 = strings.ReplaceAll(t2, "T", " ")
	t2 = strings.ReplaceAll(t2, "Z", " ")
	t2 = t2 + " +0300 +03"
	t2Time, _ := time.Parse("2006-01-02 15:04:05 -0700 -07", t2)

	return t2Time, nil
}

func ParseTime(t string) (time.Time, error) {
	tm, err := time.Parse("2006-01-02 15:04:05", t)
	return tm, err
}

func ParseFromFront(rawDate string) (time.Time, error) {
	rawDate = strings.ReplaceAll(rawDate, "T", " ")
	rawDate = strings.ReplaceAll(rawDate, "Z", " ")
	rawDate = rawDate + "+0300 +03"
	newDate, err := time.Parse("2006-01-02 15:04:05 -0700 -07", rawDate)
	return newDate, err
}

func ConvertTime(rawDate string) (time.Time, error) {
	t, err := time.Parse(time.RFC3339Nano, rawDate)
	if err != nil {
		return time.Time{}, err
	}

	loc, _ := time.LoadLocation("Etc/GMT")
	t = t.In(loc)
	return t, nil
}

func ParseDateTime(t string) (time.Time, error) {
	tm, err := time.Parse("15:04:05", t)
	return tm, err
}

func FormatTime(t time.Time) string {
	return t.Format("2006-01-02 15:04:05")
}
