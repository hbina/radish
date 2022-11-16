package redis

import (
	"strconv"
	"time"
)

func ParseTtlFromUnitTime(arg string, multiplier int64) (time.Time, error) {
	unitTime, err := strconv.ParseInt(string(arg), 10, 64)
	if err != nil {
		return time.Time{}, err
	}
	if unitTime == 0 {
		return time.Time{}, nil
	}
	return time.Now().Add(time.Duration(unitTime * multiplier)), nil
}

func ParseTtlFromTimestamp(arg string, multiplier time.Duration) (time.Time, error) {
	unitTime, err := strconv.ParseInt(string(arg), 10, 64)
	if err != nil {
		return time.Time{}, err
	}
	if unitTime == 0 {
		return time.Time{}, nil
	}
	unitTime = int64(float64(unitTime) / (float64(time.Millisecond) / float64(multiplier)))
	return time.UnixMilli(unitTime), nil
}
