package metrics

import "time"

type MetricsManagerContract interface {
	SendMeasuredTime(callName string, time time.Duration)
	SendLog(tag string, message string)
}
