package database

import (
	"context"
	"errors"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	"github.com/nicholaspark09/awsgorocket/converter"
	"log"
)

type DatabaseHelperContract[T any] interface {
	Create(data T) (*T, *error)
	Fetch(partitionKey string, rangeKey string) (*T, *error)
	FetchAll(partitionKey string, limit int32, lastRangeKey *string) ([]*T, *string)
	Update(data T) bool
	Delete(partitionKey string, rangeKey string) bool
}

type DatabaseHelper[T any] struct {
	Client    *dynamodb.Client
	TableName *string
	Converter converter.ModelConverterContract[T]
}

func (helper *DatabaseHelper[T]) Create(data *T) (*T, *error) {
	item, err := helper.Converter.ConvertToItem(data)
	if err != nil {
		log.Printf("Error in converting object: %s", errors.Unwrap(*err).Error())
		return nil, err
	}
	itemInput := &dynamodb.PutItemInput{
		TableName: helper.TableName,
		Item:      item,
	}
	_, clientError := helper.Client.PutItem(context.TODO(), itemInput)
	if clientError != nil {
		log.Printf("Error in creating an item: %s", clientError.Error())
		return nil, &clientError
	}
	return data, nil
}

func (helper *DatabaseHelper[T]) Fetch(partitionKey string, rangeKey string) (*T, *error) {
	selectedKeys := map[string]types.AttributeValue{
		"partition_key": &types.AttributeValueMemberS{Value: partitionKey},
		"range_key":     &types.AttributeValueMemberS{Value: rangeKey},
	}
	itemOutput, err := helper.Client.GetItem(context.TODO(), &dynamodb.GetItemInput{
		TableName: helper.TableName,
		Key:       selectedKeys,
	})
	if err != nil {
		log.Printf("Error in Fetching an item: %s", err.Error())
		return nil, &err
	}
	if itemOutput.Item == nil {
		log.Printf("No item found: %s", rangeKey)
		return nil, nil
	}
	return helper.Converter.ConvertToModel(itemOutput.Item)
}

func (helper *DatabaseHelper[T]) FetchAll(partitionKey string, limit int32, lastRangeKey *string) ([]*T, *string) {
	input := &dynamodb.QueryInput{
		TableName:              helper.TableName,
		KeyConditionExpression: aws.String("partition_key = :partitionKey"),
		ExpressionAttributeValues: map[string]types.AttributeValue{
			":partitionKey": &types.AttributeValueMemberS{Value: partitionKey},
		},
		Limit: aws.Int32(limit),
	}
	if lastRangeKey != nil && len(*lastRangeKey) > 0 {
		log.Printf("Trying to use a last range key of: %s", *lastRangeKey)
		input.ExclusiveStartKey = map[string]types.AttributeValue{
			"partition_key": &types.AttributeValueMemberS{Value: partitionKey},
			"range_key":     &types.AttributeValueMemberS{Value: *lastRangeKey},
		}
	}
	result, err := helper.Client.Query(context.TODO(), input)
	if err != nil {
		log.Printf("Error in Fetching an item: %s", err.Error())
		return nil, nil
	}
	var items []*T
	for _, item := range result.Items {
		pipeline, err := helper.Converter.ConvertToModel(item)
		if err != nil {
			log.Printf("Error in parsing item: %s", errors.Unwrap(*err).Error())
		}
		items = append(items, pipeline)
	}
	if _, ok := result.LastEvaluatedKey["range_key"]; ok {
		lastKey := converter.ToString("range_key", result.LastEvaluatedKey)
		return items, &lastKey
	}
	return items, nil
}

func (helper *DatabaseHelper[T]) Update(data T) bool {
	item, converterError := helper.Converter.ConvertToItem(&data)
	if converterError != nil {
		log.Printf("Error in converting the model")
		return false
	}
	_, err := helper.Client.PutItem(context.TODO(), &dynamodb.PutItemInput{
		TableName: helper.TableName,
		Item:      item,
	})
	if err != nil {
		log.Printf("Error in updating an item: %s", err.Error())
		return false
	}
	return true
}

func (helper *DatabaseHelper[T]) Delete(partitionKey string, rangeKey string) bool {
	selectedKeys := map[string]types.AttributeValue{
		"partition_key": &types.AttributeValueMemberS{Value: partitionKey},
		"range_key":     &types.AttributeValueMemberS{Value: rangeKey},
	}
	_, err := helper.Client.DeleteItem(context.TODO(), &dynamodb.DeleteItemInput{
		TableName: helper.TableName,
		Key:       selectedKeys,
	})
	if err != nil {
		log.Printf("Error in deleting an item: %s", err.Error())
		return false
	}
	return true
}
