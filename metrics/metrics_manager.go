package metrics

import (
	"context"
	"fmt"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/cloudwatch"
	"github.com/aws/aws-sdk-go-v2/service/cloudwatch/types"
	"github.com/nicholaspark09/awsgorocket/config"
	"log"
	"time"
)

type MetricsManager struct {
	client      *cloudwatch.Client
	serviceName string
}

const nameSpace = "PipelineFacadeService"

func ProvideMetricsManager(provider config.ConfigProvider, serviceName string) MetricsManagerContract {
	cloudWatchClient := cloudwatch.NewFromConfig(provider.SdkConfig)
	return &MetricsManager{
		client:      cloudWatchClient,
		serviceName: serviceName,
	}
}

func (metricsManager *MetricsManager) SendMeasuredTime(callName string, timeDuration time.Duration) {
	title := "ExecutionTime:" + callName
	elapsed := timeDuration.Milliseconds()
	metricDatum := &types.MetricDatum{
		MetricName: aws.String(title),
		Timestamp:  aws.Time(time.Now().UTC()),
		Dimensions: []types.Dimension{
			{Name: aws.String("ApiLatency"),
				Value: aws.String(fmt.Sprintf("%d", elapsed))},
		},
		Value: aws.Float64(float64(elapsed)),
	}
	_, err := metricsManager.client.PutMetricData(context.TODO(), &cloudwatch.PutMetricDataInput{
		MetricData: []types.MetricDatum{*metricDatum},
		Namespace:  aws.String(metricsManager.serviceName),
	})
	log.Printf("%s, Time: %v", title, elapsed)
	if err != nil {
		log.Printf("Error in sending metrics: %s", err.Error())
	} else {
		fmt.Printf("Method %s took %v to complete", callName, elapsed)
	}
}

func (metricsManager *MetricsManager) Send500Error(callName string, statusCode int, message string) {
	metricDatum := &types.MetricDatum{
		MetricName: aws.String("5XXError:" + callName),
		Timestamp:  aws.Time(time.Now().UTC()),
		Dimensions: []types.Dimension{
			{Name: aws.String("ApiError"),
				Value: aws.String(message)},
		},
		Value: aws.Float64(float64(statusCode)),
	}
	_, err := metricsManager.client.PutMetricData(context.TODO(), &cloudwatch.PutMetricDataInput{
		MetricData: []types.MetricDatum{*metricDatum},
		Namespace:  aws.String(metricsManager.serviceName),
	})
	if err != nil {
		log.Printf("Error in sending metrics: %s", err.Error())
	}
}

func (metricsManager *MetricsManager) Send400Error(callName string, statusCode int, message string) {
	metricDatum := &types.MetricDatum{
		MetricName: aws.String("4XXError:" + callName),
		Timestamp:  aws.Time(time.Now().UTC()),
		Dimensions: []types.Dimension{
			{Name: aws.String("ApiError"),
				Value: aws.String(message)},
		},
		Value: aws.Float64(float64(statusCode)),
	}
	_, err := metricsManager.client.PutMetricData(context.TODO(), &cloudwatch.PutMetricDataInput{
		MetricData: []types.MetricDatum{*metricDatum},
		Namespace:  aws.String(metricsManager.serviceName),
	})
	if err != nil {
		log.Printf("Error in sending metrics: %s", err.Error())
	}
}

func (metricsManager *MetricsManager) SendLog(tag string, message string) {
	// TODO
}
