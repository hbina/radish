package redis

import (
	"strconv"
	"strings"
	"time"

	"github.com/pkg/errors"
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

func ParseIntRange(startStr string, stopStr string) (int, bool, int, bool, error) {
	startExclusive := false
	stopExclusive := false

	if len(startStr) > 0 && startStr[0] == '[' {
		startStr = startStr[1:]
		startExclusive = false
	} else if len(startStr) > 0 && startStr[0] == '(' {
		startStr = startStr[1:]
		startExclusive = true
	}

	if len(stopStr) > 0 && stopStr[0] == '[' {
		stopStr = stopStr[1:]
		stopExclusive = false
	} else if len(stopStr) > 0 && stopStr[0] == '(' {
		stopStr = stopStr[1:]
		stopExclusive = true
	}

	start, err := strconv.ParseInt(startStr, 10, 32)

	if err != nil {
		return 0, false, 0, false, err
	}

	stop, err := strconv.ParseInt(stopStr, 10, 32)

	if err != nil {
		return 0, false, 0, false, err
	}

	return int(start), startExclusive, int(stop), stopExclusive, nil
}

func ParseFloatRange(startStr string, stopStr string) (float64, bool, float64, bool, error) {
	startExclusive := false
	stopExclusive := false

	if len(startStr) > 0 && startStr[0] == '[' {
		startStr = startStr[1:]
		startExclusive = false
	} else if len(startStr) > 0 && startStr[0] == '(' {
		startStr = startStr[1:]
		startExclusive = true
	}

	if len(stopStr) > 0 && stopStr[0] == '[' {
		stopStr = stopStr[1:]
		stopExclusive = false
	} else if len(stopStr) > 0 && stopStr[0] == '(' {
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

func ParseLexRange(start string, stop string) (string, bool, string, bool, error) {
	startExclusive := false
	stopExclusive := false

	if len(start) > 0 && start[0] == '[' {
		start = start[1:]
		startExclusive = false
	} else if len(start) > 0 && start[0] == '(' {
		start = start[1:]
		startExclusive = true
	} else if start != "+" && start != "-" {
		return start, startExclusive, stop, stopExclusive, errors.New("ERR min or max not valid string range item")
	}

	if len(stop) > 0 && stop[0] == '[' {
		stop = stop[1:]
		stopExclusive = false
	} else if len(stop) > 0 && stop[0] == '(' {
		stop = stop[1:]
		stopExclusive = true
	} else if stop != "+" && stop != "-" {
		return start, startExclusive, stop, stopExclusive, errors.New("ERR min or max not valid string range item")
	}

	return start, startExclusive, stop, stopExclusive, nil
}
