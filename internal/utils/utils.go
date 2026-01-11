package utils

import "time"

var layout = "01-2006"

func ValidDate(date string) (*time.Time, error) {

	time, err := time.Parse(layout, date)
	if err != nil {
		return nil, err
	}
	return &time, nil
}

func ParseDate(time time.Time) string {
	return time.Format(layout)
}