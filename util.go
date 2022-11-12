package redis

import (
	"strconv"
	"time"
)

func parseExpiryTime(arg string, multiplier uint64) (time.Time, error) {
	unitTime, err := strconv.ParseUint(string(arg), 10, 64)
	if err != nil {
		return time.Time{}, err
	}
	if unitTime == 0 {
		return time.Time{}, nil
	}
	return time.Now().Add(time.Duration(unitTime * multiplier)), nil
}
