package main

import (
	"context"
	"errors"
	"log"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/expression"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	"github.com/aws/smithy-go"
	"github.com/google/uuid"
)

type Item struct {
	Id          string `json:"id" dynamodbav:"id"`
	Description string `json:"description" dynamodbav:"description"`
	Location    string `json:"location" dynamodbav:"location"`
	Quantity    int    `json:"quantity" dynamodbav:"quantity"`
	// add more fields
}

type NewOrUpdatedItem struct {
	Description string `json:"description" validate:"required"`
	Location    string `json:"location" validate:"required"`
	Quantity    int    `json:"quantity" validate:"required"`
	// add more fields
}

func getClient() (dynamodb.Client, error) {
	sdkConfig, err := config.LoadDefaultConfig(context.TODO())

	dbClient := *dynamodb.NewFromConfig(sdkConfig)

	return dbClient, err
}

func getItem(ctx context.Context, id string) (*Item, error) {
	key, err := attributevalue.Marshal(id)
	if err != nil {
		return nil, err
	}

	input := &dynamodb.GetItemInput{
		TableName: aws.String(TableName),
		Key: map[string]types.AttributeValue{
			"id": key,
		},
	}

	log.Printf("Calling DynamoDB with input: %v", input)
	result, err := db.GetItem(ctx, input)
	if err != nil {
		return nil, err
	}
	log.Printf("Executed GetItem DynamoDb successfully. Result: %#v", result)

	if result.Item == nil {
		return nil, nil
	}

	item := new(Item)
	err = attributevalue.UnmarshalMap(result.Item, item)
	if err != nil {
		return nil, err
	}

	return item, nil
}

func listItems(ctx context.Context) ([]Item, error) {
	items := make([]Item, 0)

	var token map[string]types.AttributeValue

	for {
		input := &dynamodb.ScanInput{
			TableName:         aws.String(TableName),
			ExclusiveStartKey: token,
		}

		result, err := db.Scan(ctx, input)
		if err != nil {
			return nil, err
		}

		var fetchedItem []Item
		err = attributevalue.UnmarshalListOfMaps(result.Items, &fetchedItem)
		if err != nil {
			return nil, err
		}

		items = append(items, fetchedItem...)
		token = result.LastEvaluatedKey
		if token == nil {
			break
		}

	}

	return items, nil
}

func insertItem(ctx context.Context, newItem NewOrUpdatedItem) (*Item, error) {
	item := Item{
		Id:          uuid.NewString(),
		Description: newItem.Description,
		Location:    newItem.Location,
		Quantity:    newItem.Quantity,
		// add more fields
	}

	itemMap, err := attributevalue.MarshalMap(item)
	if err != nil {
		return nil, err
	}

	input := &dynamodb.PutItemInput{
		TableName: aws.String(TableName),
		Item:      itemMap,
	}

	res, err := db.PutItem(ctx, input)
	if err != nil {
		return nil, err
	}

	err = attributevalue.UnmarshalMap(res.Attributes, &item)
	if err != nil {
		return nil, err
	}

	return &item, nil
}

func deleteItem(ctx context.Context, id string) (*Item, error) {
	key, err := attributevalue.Marshal(id)
	if err != nil {
		return nil, err
	}

	input := &dynamodb.DeleteItemInput{
		TableName: aws.String(TableName),
		Key: map[string]types.AttributeValue{
			"id": key,
		},
		ReturnValues: types.ReturnValue(*aws.String("ALL_OLD")),
	}

	res, err := db.DeleteItem(ctx, input)
	if err != nil {
		return nil, err
	}

	if res.Attributes == nil {
		return nil, nil
	}

	item := new(Item)
	err = attributevalue.UnmarshalMap(res.Attributes, item)
	if err != nil {
		return nil, err
	}

	return item, nil
}

func updateItem(ctx context.Context, id string, updatedItem NewOrUpdatedItem) (*Item, error) {
	key, err := attributevalue.Marshal(id)
	if err != nil {
		return nil, err
	}

	expr, err := expression.NewBuilder().WithUpdate(
		expression.Set(
			expression.Name("description"),
			expression.Value(updatedItem.Description),
		).Set(
			expression.Name("location"),
			expression.Value(updatedItem.Location),
		).Set(
			expression.Name("quantity"),
			expression.Value(updatedItem.Quantity),
		),
		// add more fields
	).WithCondition(
		expression.Equal(
			expression.Name("id"),
			expression.Value(id),
		),
	).Build()
	if err != nil {
		return nil, err
	}

	input := &dynamodb.UpdateItemInput{
		Key: map[string]types.AttributeValue{
			"id": key,
		},
		TableName:                 aws.String(TableName),
		UpdateExpression:          expr.Update(),
		ExpressionAttributeNames:  expr.Names(),
		ExpressionAttributeValues: expr.Values(),

		ConditionExpression: expr.Condition(),
		ReturnValues:        types.ReturnValue(*aws.String("ALL_NEW")),
	}

	res, err := db.UpdateItem(ctx, input)
	if err != nil {
		var smErr *smithy.OperationError
		if errors.As(err, &smErr) {
			var condCheckFailed *types.ConditionalCheckFailedException
			if errors.As(err, &condCheckFailed) {
				return nil, nil
			}
		}

		return nil, err
	}

	if res.Attributes == nil {
		return nil, nil
	}

	item := new(Item)
	err = attributevalue.UnmarshalMap(res.Attributes, item)
	if err != nil {
		return nil, err
	}

	return item, nil
}
