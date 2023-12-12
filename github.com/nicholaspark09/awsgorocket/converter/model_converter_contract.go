package converter

import (
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	"log"
	"strconv"
)

type ModelConverterContract[T any] interface {
	ConvertToItem(data *T) (map[string]types.AttributeValue, *error)
	ConvertToModel(map[string]types.AttributeValue) (*T, *error)
}

func ToUnsafeString(key string, item map[string]types.AttributeValue) *string {
	value, ok := item[key]
	if ok {
		switch v := value.(type) {
		case *types.AttributeValueMemberS:
			return &v.Value
		default:
			log.Printf("Could not convert: %s", value)
			return nil
		}
	}
	log.Printf("Key not found: %s", key)
	return nil
}

func ToString(key string, item map[string]types.AttributeValue) string {
	value, ok := item[key]
	if ok {
		switch v := value.(type) {
		case *types.AttributeValueMemberS:
			return v.Value
		default:
			log.Printf("Could not convert: %s", value)
			return ""
		}
	}
	log.Printf("Key not found: %s", key)
	return ""
}

func ToInt(key string, item map[string]types.AttributeValue) int {
	value, ok := item[key]
	if ok {
		switch v := value.(type) {
		case *types.AttributeValueMemberN:
			num, err := strconv.Atoi(v.Value)
			if err != nil {
				log.Printf("Could not convert: %s", v.Value)
				return -1
			}
			return num
		case *types.AttributeValueMemberS:
			num, err := strconv.Atoi(v.Value)
			if err != nil {
				log.Printf("Could not convert: %s", v.Value)
				return -1
			}
			return num
		default:
			log.Printf("Could not convert: %s", value)
			return -1
		}
	}
	log.Printf("Key not found: %s", key)
	return -1
}

func ToBool(key string, item map[string]types.AttributeValue) bool {
	value, ok := item[key]
	if ok {
		switch v := value.(type) {
		case *types.AttributeValueMemberBOOL:
			return v.Value
		default:
			log.Printf("Could not convert: %s", value)
			return false
		}
	}
	log.Printf("Key not found: %s", key)
	return false
}
