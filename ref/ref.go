package ref

import "time"

func Bool(i bool) *bool {
	return &i
}

func Int(i int) *int {
	return &i
}

func Int64(i int64) *int64 {
	return &i
}

func String(i string) *string {
	return &i
}

func Duration(i time.Duration) *time.Duration {
	return &i
}

func Time(i time.Time) *time.Time {
	return &i
}
