package controllers

import (
	"github.com/HakimHC/altostratus-golang-auth/config"
	"github.com/HakimHC/altostratus-golang-auth/models"
	"github.com/HakimHC/altostratus-golang-auth/responses"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
	"github.com/aws/aws-sdk-go/service/dynamodb/expression"
	"github.com/labstack/echo/v4"
)

func putItemInDynamoDB(user models.User, tableName string) error {
	av, err := dynamodbattribute.MarshalMap(user)
	if err != nil {
		return err
	}

	input := &dynamodb.PutItemInput{
		Item:      av,
		TableName: aws.String(tableName),
	}

	_, err = config.DynamoDB.PutItem(input)
	return err
}

func getUserByUsername(username string, tableName string) (*models.User, error) {
	filt := expression.Name("username").Equal(expression.Value(username))
	expr, err := expression.NewBuilder().WithFilter(filt).Build()

	params := &dynamodb.ScanInput{
		TableName:                 aws.String(tableName),
		ExpressionAttributeNames:  expr.Names(),
		ExpressionAttributeValues: expr.Values(),
		FilterExpression:          expr.Filter(),
	}

	result, err := config.DynamoDB.Scan(params)
	if err != nil {
		return nil, err
	}
	if len(result.Items) == 0 {
		return nil, nil
	}

	var users []models.User

	err = dynamodbattribute.UnmarshalListOfMaps(result.Items, &users)
	if err != nil {
		return nil, err
	}
	return &users[0], nil
}

func ErrorResponse(c echo.Context, statusCode int, err string) error {
	return c.JSON(statusCode, responses.AuthResponse{
		Status:  statusCode,
		Message: "error",
		Data:    &echo.Map{"data": err},
	})
}
