package metrics

import "time"

type MetricsManagerContract interface {
	SendMeasuredTime(serviceName string, callName string, time time.Duration)
	SendLog(tag string, message string)
}
