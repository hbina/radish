package redis

import (
	"strconv"
	"strings"
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

func EscapeString(s string) string {
	s = strings.ReplaceAll(s, "\r", "\\r")
	s = strings.ReplaceAll(s, "\n", "\\n")
	s = strings.ReplaceAll(s, "\t", "\\t")
	return s
}

func CollectArgs(args [][]byte) string {
	result := ""
	for i, arg := range args {
		result += string(arg)
		if i != len(args)-1 {
			result += " "
		}
	}
	return result
}

func ParseIntRange(startStr string, stopStr string) (int, int, error) {
	startExclusive := false
	stopExclusive := false

	if len(startStr) > 0 && startStr[0] == '(' {
		startStr = startStr[1:]
		startExclusive = true
	}

	if len(stopStr) > 0 && stopStr[0] == '(' {
		stopStr = stopStr[1:]
		stopExclusive = true
	}

	start, err := strconv.ParseInt(startStr, 10, 32)

	if err != nil {
		return 0, 0, err
	}

	stop, err := strconv.ParseInt(stopStr, 10, 32)

	if err != nil {
		return 0, 0, err
	}

	if startExclusive {
		start += 1
	}

	if stopExclusive {
		stop -= 1
	}

	return int(start), int(stop), nil
}

func ParseFloatRange(startStr string, stopStr string) (float64, bool, float64, bool, error) {
	startExclusive := false
	stopExclusive := false

	if len(startStr) > 0 && startStr[0] == '(' {
		startStr = startStr[1:]
		startExclusive = true
	}

	if len(stopStr) > 0 && stopStr[0] == '(' {
		stopStr = stopStr[1:]
		stopExclusive = true
	}

	start, err := strconv.ParseFloat(startStr, 64)

	if err != nil {
		return 0, startExclusive, 0, stopExclusive, err
	}

	stop, err := strconv.ParseFloat(stopStr, 64)

	if err != nil {
		return 0, startExclusive, 0, stopExclusive, err
	}

	return start, startExclusive, stop, stopExclusive, nil
}
