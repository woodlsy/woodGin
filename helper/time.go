package helper

import (
	"fmt"
	"time"
)

const InitDateTime = "0001-01-01 00:00:00"

func Now(params ...interface{}) string {
	var format string
	if len(params) == 0 {
		format = "2006-01-02 15:04:05"
	}
	if len(params) > 0 {
		format = params[0].(string)
	}
	return time.Now().Format(format)
}

// IsBetweenTime
//
//	@Description: 判断时间c是否在a和b之间
//	@param a
//	@param b
//	@param c
//	@return bool
func IsBetweenTime(a, b, c string) bool {
	layout := "2006-01-02 15:04:05"
	timeA, err := time.Parse(layout, a)
	if err != nil {
		fmt.Println(err)
		return false
	}
	timeB, err := time.Parse(layout, b)
	if err != nil {
		return false
	}
	timeC, err := time.Parse(layout, c)
	if err != nil {
		return false
	}

	return timeC.After(timeA) && timeC.Before(timeB)
}
