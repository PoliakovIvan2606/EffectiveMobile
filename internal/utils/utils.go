package utils

import "time"

func ValidDate(date string) (*time.Time, error) {
	layout := "01-2006"

	time, err := time.Parse(layout, date)
	if err != nil {
		return nil, err
	}
	return &time, nil
}