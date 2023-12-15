package metrics

import "time"

type MetricsManagerContract interface {
	SendMeasuredTime(callName string, time time.Duration)
	SendLog(tag string, message string)
	Send500Error(callName string, statusCode int, message string)
	Send400Error(callName string, statusCode int, message string)
}
