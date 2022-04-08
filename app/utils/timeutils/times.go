package timeutils

import "time"

func GetCurrentMilliTime() int64 {
	return time.Now().Unix()
}
