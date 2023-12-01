package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/aws/aws-lambda-go/events"
	"github.com/go-playground/validator"
)

type ResponseStructure struct {
	Data         interface{} `json:"data"`
	ErrorMessage *string     `json:"errorMessage"` // can be string or nil
}

var validate *validator.Validate = validator.New()

var headers = map[string]string{
	"Access-Control-Allow-Origin":  OriginURL,
	"Access-Control-Allow-Headers": "Content-Type",
}

func router(ctx context.Context, req events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	log.Printf("Received req %#v", req)

	switch req.HTTPMethod {
	case "GET":
		return processGet(ctx, req)
	case "POST":
		return processPost(ctx, req)
	case "DELETE":
		return processDelete(ctx, req)
	case "PUT":
		return processPut(ctx, req)
	case "OPTIONS":
		return processOptions()
	default:
		return clientError(http.StatusMethodNotAllowed)
	}
}

func processOptions() (events.APIGatewayProxyResponse, error) {
	additionalHeaders := map[string]string{
		"Access-Control-Allow-Methods": "OPTIONS, POST, GET, PUT, DELETE",
		"Access-Control-Max-Age":       "3600",
	}
	mergedHeaders := mergeHeaders(headers, additionalHeaders)

	return events.APIGatewayProxyResponse{
		StatusCode: http.StatusOK,
		Headers:    mergedHeaders,
	}, nil
}

func processGet(ctx context.Context, req events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	id, ok := req.PathParameters["id"]
	if !ok {
		return processGetAll(ctx)
	} else {
		return processGetItemById(ctx, id)
	}
}

func processGetItemById(ctx context.Context, id string) (events.APIGatewayProxyResponse, error) {
	log.Printf("Received GET item request with id = %s", id)

	item, err := getItem(ctx, id)
	if err != nil {
		return serverError(err)
	}

	if item == nil {
		return clientError(http.StatusNotFound)
	}

	response := ResponseStructure{
		Data:         item,
		ErrorMessage: nil,
	}

	responseJson, err := json.Marshal(response)
	if err != nil {
		return serverError(err)
	}
	log.Printf("Successfully fetched item %s", response.Data)

	return events.APIGatewayProxyResponse{
		StatusCode: http.StatusOK,
		Body:       string(responseJson),
		Headers:    headers,
	}, nil
}

func processGetAll(ctx context.Context) (events.APIGatewayProxyResponse, error) {
	log.Print("Received GET items request")

	items, err := listItems(ctx)
	if err != nil {
		return serverError(err)
	}

	response := ResponseStructure{
		Data:         items,
		ErrorMessage: nil,
	}

	responseJson, err := json.Marshal(response)
	if err != nil {
		return serverError(err)
	}
	log.Printf("Successfully fetched items: %s", response.Data)

	return events.APIGatewayProxyResponse{
		StatusCode: http.StatusOK,
		Body:       string(responseJson),
		Headers:    headers,
	}, nil
}

func processPost(ctx context.Context, req events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	var createdItem NewOrUpdatedItem
	err := json.Unmarshal([]byte(req.Body), &createdItem)
	if err != nil {
		log.Printf("Can't unmarshal body: %v", err)
		return clientError(http.StatusUnprocessableEntity)
	}

	err = validate.Struct(&createdItem)
	if err != nil {
		log.Printf("Invalid body: %v", err)
		return clientError(http.StatusBadRequest)
	}
	log.Printf("Received POST request with item: %+v", createdItem)

	item, err := insertItem(ctx, createdItem)
	if err != nil {
		return serverError(err)
	}
	log.Printf("Inserted new item: %+v", item)

	response := ResponseStructure{
		Data:         item,
		ErrorMessage: nil,
	}

	responseJson, err := json.Marshal(response)
	if err != nil {
		return serverError(err)
	}
	log.Printf("Successfully fetched item item %s", response.Data)

	additionalHeaders := map[string]string{
		"Location": fmt.Sprintf("/%s/%s", ApiPath, item.Id),
	}
	mergedHeaders := mergeHeaders(headers, additionalHeaders)

	return events.APIGatewayProxyResponse{
		StatusCode: http.StatusCreated,
		Body:       string(responseJson),
		Headers:    mergedHeaders,
	}, nil
}

func processDelete(ctx context.Context, req events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	id, ok := req.PathParameters["id"]
	if !ok {
		return clientError(http.StatusBadRequest)
	}
	log.Printf("Received DELETE request with id = %s", id)

	item, err := deleteItem(ctx, id)
	if err != nil {
		return serverError(err)
	}

	if item == nil {
		return clientError(http.StatusNotFound)
	}

	response := ResponseStructure{
		Data:         item,
		ErrorMessage: nil,
	}

	responseJson, err := json.Marshal(response)
	if err != nil {
		return serverError(err)
	}

	log.Printf("Successfully deleted item item %+v", item)

	return events.APIGatewayProxyResponse{
		StatusCode: http.StatusOK,
		Body:       string(responseJson),
		Headers:    headers,
	}, nil
}

func processPut(ctx context.Context, req events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	id, ok := req.PathParameters["id"]
	if !ok {
		return clientError(http.StatusBadRequest)
	}

	var updatedItem NewOrUpdatedItem
	err := json.Unmarshal([]byte(req.Body), &updatedItem)
	if err != nil {
		log.Printf("Can't unmarshal body: %v", err)
		return clientError(http.StatusUnprocessableEntity)
	}

	err = validate.Struct(&updatedItem)
	if err != nil {
		log.Printf("Invalid body: %v", err)
		return clientError(http.StatusBadRequest)
	}
	log.Printf("Received PUT request with item: %+v", updatedItem)

	item, err := updateItem(ctx, id, updatedItem)
	if err != nil {
		return serverError(err)
	}

	if item == nil {
		return clientError(http.StatusNotFound)
	}

	response := ResponseStructure{
		Data:         item,
		ErrorMessage: nil,
	}

	responseJson, err := json.Marshal(response)
	if err != nil {
		return serverError(err)
	}

	log.Printf("Updated item: %+v", item)

	additionalHeaders := map[string]string{
		"Location": fmt.Sprintf("/%s/%s", ApiPath, item.Id),
	}
	mergedHeaders := mergeHeaders(headers, additionalHeaders)

	return events.APIGatewayProxyResponse{
		StatusCode: http.StatusOK,
		Body:       string(responseJson),
		Headers:    mergedHeaders,
	}, nil
}
