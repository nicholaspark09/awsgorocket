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
	client *cloudwatch.Client
}

const nameSpace = "PipelineFacadeService"

func ProvideMetricsManager(provider config.ConfigProvider) MetricsManagerContract {
	cloudWatchClient := cloudwatch.NewFromConfig(provider.SdkConfig)
	return &MetricsManager{
		client: cloudWatchClient,
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
		Namespace:  aws.String(nameSpace),
	})
	log.Printf("%s, Time: %v", title, elapsed)
	if err != nil {
		log.Printf("Error in sending metrics: %s", err.Error())
	} else {
		fmt.Printf("Method %s took %v to complete", callName, elapsed)
	}
}

func (metricsManager *MetricsManager) SendLog(tag string, message string) {
	// TODO
}
