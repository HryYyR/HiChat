package common

import "time"

func ParseTime(tstr string) (time.Time, error) {
	t, err := time.Parse("2006-01-02 15:04:05", tstr)
	if err != nil {
		return time.Time{}, err
	}
	return t, nil
}
