package metrics

import "time"

func MeasureTime[T any](callName string, metricsManager MetricsManagerContract, f func() *T) *T {
	start := time.Now()
	result := f()
	metricsManager.SendMeasuredTime(callName, time.Since(start))
	return result
}

func MeasureTimeWithError[T any](callName string, metricsManager MetricsManagerContract, f func() (*T, *error)) (*T, *error) {
	start := time.Now()
	result, err := f()
	metricsManager.SendMeasuredTime(callName, time.Since(start))
	return result, err
}
