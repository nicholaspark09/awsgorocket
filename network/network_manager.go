package network

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	response "github.com/nicholaspark09/awsgorocket/model"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
)

type NetworkManager[T any] struct {
	endpoint    string
	params      map[string]string
	apiKey      *string
	contentType *string
	client      http.Client
}

func ProvideNetworkManager[T any](endpoint string, params map[string]string, apiKey *string, contentType *string) NetworkManager[T] {
	return NetworkManager[T]{
		endpoint:    endpoint,
		params:      params,
		apiKey:      apiKey,
		contentType: contentType,
		client:      http.Client{},
	}
}

func (manager *NetworkManager[T]) GetEndpoint() (string, *error) {
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

func Get[T any](manager NetworkManager[T]) response.Response[T] {
	formedEndpoint, queryError := manager.GetEndpoint()
	if queryError != nil {
		log.Println("Error in parsing Query")
		return response.Response[T]{
			Error:      queryError,
			Message:    fmt.Sprintf("Error in parsing query: %s", errors.Unwrap(*queryError).Error()),
			Data:       nil,
			StatusCode: 400,
		}
	}
	req, err := http.NewRequest("GET", formedEndpoint, nil)
	if err != nil {
		return response.Response[T]{
			Error:      &err,
			Message:    fmt.Sprintf("Error in parsing query: %s", errors.Unwrap(err).Error()),
			Data:       nil,
			StatusCode: 500,
		}
	}
	var contentType string = "application/json"
	if manager.contentType != nil {
		contentType = *manager.contentType
	}
	req.Header.Set("Content-Type", contentType)
	if manager.apiKey != nil {
		req.Header.Set("x-api-key", *manager.apiKey)
	}
	clientResponse, err := manager.client.Do(req)
	if err != nil {
		log.Printf("Error in making api request: %s", err.Error())
		return response.Response[T]{
			StatusCode: clientResponse.StatusCode,
			Data:       nil,
			Error:      &err,
			Message:    errors.Unwrap(err).Error(),
		}
	}
	defer clientResponse.Body.Close()
	if clientResponse.StatusCode != http.StatusOK {
		errorMessage := fmt.Sprintf("Error with the request call: %d", clientResponse.StatusCode)
		log.Println(errorMessage)
		return response.Response[T]{
			StatusCode: clientResponse.StatusCode,
			Data:       nil,
			Error:      nil,
			Message:    errorMessage,
		}
	}
	responseBody, err := ioutil.ReadAll(clientResponse.Body)
	if err != nil {
		errorMessage := fmt.Sprintf("Error in reading api request: %s", err.Error())
		return response.Response[T]{
			StatusCode: clientResponse.StatusCode,
			Data:       nil,
			Error:      &err,
			Message:    errorMessage,
		}
	}
	var result T
	jsonError := json.Unmarshal(responseBody, &result)
	if jsonError != nil {
		log.Printf("Error in unmarshalling api request: %s", jsonError.Error())
		return response.Response[T]{
			StatusCode: 500,
			Data:       nil,
			Error:      &jsonError,
			Message:    jsonError.Error(),
		}
	}
	return response.Response[T]{
		StatusCode: clientResponse.StatusCode,
		Data:       &result,
		Error:      nil,
		Message:    "Success",
	}
}

func Post[T any](manager NetworkManager[T], jsonBody []byte) response.Response[T] {
	formedEndpoint, queryError := manager.GetEndpoint()
	if queryError != nil {
		log.Println("Error in parsing Query")
		return response.Response[T]{
			Error:      queryError,
			Message:    fmt.Sprintf("Error in parsing query: %s", errors.Unwrap(*queryError).Error()),
			Data:       nil,
			StatusCode: 400,
		}
	}
	req, err := http.NewRequest("POST", formedEndpoint, bytes.NewReader(jsonBody))
	if err != nil {
		return response.Response[T]{
			Error:      &err,
			Message:    fmt.Sprintf("Error in parsing query: %s", errors.Unwrap(err).Error()),
			Data:       nil,
			StatusCode: 500,
		}
	}
	var contentType string = "application/jsonBody"
	if manager.contentType != nil {
		contentType = *manager.contentType
	}
	req.Header.Set("Content-Type", contentType)
	if manager.apiKey != nil {
		req.Header.Set("x-api-key", *manager.apiKey)
	}
	clientResponse, err := manager.client.Do(req)
	if err != nil {
		log.Printf("Error in making api request: %s", err.Error())
		return response.Response[T]{
			StatusCode: clientResponse.StatusCode,
			Data:       nil,
			Error:      &err,
			Message:    errors.Unwrap(err).Error(),
		}
	}
	defer clientResponse.Body.Close()
	if clientResponse.StatusCode != http.StatusOK {
		errorMessage := fmt.Sprintf("Error with the request call: %d", clientResponse.StatusCode)
		log.Println(errorMessage)
		return response.Response[T]{
			StatusCode: clientResponse.StatusCode,
			Data:       nil,
			Error:      nil,
			Message:    errorMessage,
		}
	}
	responseBody, err := ioutil.ReadAll(clientResponse.Body)
	if err != nil {
		log.Printf("Error in reading api request: %s", err.Error())
	}
	var result T
	jsonError := json.Unmarshal(responseBody, &result)
	if jsonError != nil {
		log.Printf("Error in unmarshalling api request: %s", jsonError.Error())
		return response.Response[T]{
			StatusCode: 500,
			Data:       nil,
			Error:      &jsonError,
			Message:    jsonError.Error(),
		}
	}
	return response.Response[T]{
		StatusCode: clientResponse.StatusCode,
		Data:       &result,
		Error:      nil,
		Message:    "Success",
	}
}
