package network_v2

import (
	"bytes"
	json2 "encoding/json"
	"fmt"
	"github.com/nicholaspark09/awsgorocket/utils"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
)

type NetworkManagerV2[T any] struct {
	endpoint    string
	params      map[string]string
	apiKey      *string
	contentType *string
	client      http.Client
}

func ProvideNetworkManagerV2[T any](endpoint string, params map[string]string, apiKey *string, contentType *string) NetworkManagerV2[T] {
	return NetworkManagerV2[T]{
		endpoint:    endpoint,
		params:      params,
		apiKey:      apiKey,
		contentType: contentType,
		client:      http.Client{},
	}
}

func (manager *NetworkManagerV2[T]) GetEndpoint() (string, *error) {
	parsedUrl, err := url.Parse(manager.endpoint)
	if err != nil {
		return manager.endpoint, &err
	}
	queryParams := parsedUrl.Query()
	for key, value := range manager.params {
		queryParams.Add(key, value)
	}
	parsedUrl.RawQuery = queryParams.Encode()
	formedEndpoint := parsedUrl.String()
	return formedEndpoint, nil
}

func Post[T any](manager NetworkManagerV2[T], json []byte) (*T, error) {
	formedEndpoint, queryError := manager.GetEndpoint()
	if queryError != nil {
		log.Println("Error in parsing Query")
		return nil, *queryError
	}
	req, err := http.NewRequest("POST", formedEndpoint, bytes.NewReader(json))
	if err != nil {
		return nil, err
	}
	var contentType string = "application/json"
	if manager.contentType != nil {
		contentType = *manager.contentType
	}
	req.Header.Set("Content-Type", contentType)
	if manager.apiKey != nil {
		req.Header.Set("x-api-key", *manager.apiKey)
	}
	response, err := manager.client.Do(req)
	if err != nil {
		log.Printf("Error in making api request: %s", err.Error())
	}
	defer response.Body.Close()
	if response.StatusCode != http.StatusOK {
		errorMessage := fmt.Sprintf("Error with the request call: %d", response.StatusCode)
		log.Println(errorMessage)
		return nil, utils.GenericError{
			Message:    errorMessage,
			StatusCode: response.StatusCode,
		}
	}
	responseBody, err := ioutil.ReadAll(response.Body)
	if err != nil {
		log.Printf("Error in reading api request: %s", err.Error())
	}
	var result T
	jsonError := json2.Unmarshal(responseBody, &result)
	if jsonError != nil {
		log.Printf("Error in unmarshalling api request: %s", jsonError.Error())
		return nil, jsonError
	}
	return &result, nil
}

func Get[T any](manager NetworkManagerV2[T]) (*T, error) {
	formedEndpoint, queryError := manager.GetEndpoint()
	if queryError != nil {
		log.Println("Error in parsing Query")
		return nil, *queryError
	}
	req, err := http.NewRequest("GET", formedEndpoint, nil)
	if err != nil {
		return nil, err
	}
	var contentType string = "application/json"
	if manager.contentType != nil {
		contentType = *manager.contentType
	}
	req.Header.Set("Content-Type", contentType)
	if manager.apiKey != nil {
		req.Header.Set("x-api-key", *manager.apiKey)
	}
	response, err := manager.client.Do(req)
	if err != nil {
		log.Printf("Error in making api request: %s", err.Error())
	}
	defer response.Body.Close()
	if response.StatusCode != http.StatusOK {
		errorMessage := fmt.Sprintf("Error with the request call: %d", response.StatusCode)
		log.Println(errorMessage)
		return nil, utils.GenericError{
			Message:    errorMessage,
			StatusCode: response.StatusCode,
		}
	}
	responseBody, err := ioutil.ReadAll(response.Body)
	if err != nil {
		log.Printf("Error in reading api request: %s", err.Error())
	}
	var result T
	jsonError := json2.Unmarshal(responseBody, &result)
	if jsonError != nil {
		log.Printf("Error in unmarshalling api request: %s", jsonError.Error())
		return nil, jsonError
	}
	return &result, nil
}
