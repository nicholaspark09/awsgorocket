package metrics

import "time"

func MeasureTime[T any](callName string, metricsManager MetricsManagerContract, f func() *T) *T {
	start := time.Now()
	result := f()
	metricsManager.SendMeasuredTime(callName, time.Since(start))
	return result
}
