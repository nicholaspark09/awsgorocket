package config

import (
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/credentials"
)

var configRepository ConfigRepositoryContract

type ConfigProvider struct {
	SdkConfig        aws.Config
	ConfigRepository ConfigRepositoryContract
}

func ProvideConfigProvider() ConfigProvider {
	configRepository := ConfigRepository{}
	region := configRepository.GetString("REGION", "us-east-2")
	accessKey := configRepository.GetString("ACCESS_KEY_ID", "")
	secretKey := configRepository.GetString("SECRET_KEY", "")
	sdkConfig := aws.Config{
		Region: region,
		Credentials: credentials.NewStaticCredentialsProvider(
			accessKey,
			secretKey,
			"",
		),
	}
	return ConfigProvider{
		SdkConfig:        sdkConfig,
		ConfigRepository: &configRepository,
	}
}
