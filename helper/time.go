package helper

import "time"

const InitDateTime = "1990-01-01 00:00:00"

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
