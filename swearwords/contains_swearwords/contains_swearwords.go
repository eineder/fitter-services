package main

import (
	"context"
	"fmt"
	"strings"

	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
)

type ContainsSwearwordsEvent struct {
	Text string `json:"text"`
}

func main() {
	lambda.Start(HandleRequest)
}

func HandleRequest(ctx context.Context, event *ContainsSwearwordsEvent) (*[]string, error) {
	if event == nil {
		return nil, fmt.Errorf("received nil event")
	}

	if event.Text == "" {
		return nil, fmt.Errorf("received empty text")
	}

	distinctWords := distinct(strings.Split(event.Text, " "))
	statement := aws.String("SELECT * FROM swearwords WHERE word IN ?")
	parameters := getParameters(&distinctWords)
	sess := session.Must(session.NewSession())
	dynamoClient := dynamodb.New(sess)

	swearwords, err := querySwearwords(parameters, dynamoClient, statement)
	if err != nil {
		fmt.Println("Error querying swearwords ", err)
		return nil, err
	}

	return &swearwords, nil
}

func querySwearwords(parameters []*dynamodb.AttributeValue, dynamoClient *dynamodb.DynamoDB, statement *string) ([]string, error) {
	const batchSize = 100
	swearwords := []string{}
	batchParameters := make([]*dynamodb.AttributeValue, batchSize)
	for i, parameter := range parameters {
		batchParameters = append(batchParameters, parameter)
		if i%batchSize == 0 || i == len(parameters)-1 {
			sw, err := queryBatchSwearwords(dynamoClient, statement, batchParameters, swearwords)
			if err != nil {
				fmt.Println("Error querying swearwords ", err)
				return nil, err
			}
			swearwords = append(swearwords, sw...)
			batchParameters = make([]*dynamodb.AttributeValue, batchSize)
		}
	}
	return swearwords, nil
}

func queryBatchSwearwords(svc *dynamodb.DynamoDB, statement *string, batchParameters []*dynamodb.AttributeValue, swearwords []string) ([]string, error) {
	out, err := svc.ExecuteStatement(&dynamodb.ExecuteStatementInput{
		Statement:  statement,
		Parameters: batchParameters,
	})
	if err != nil {
		fmt.Println("Error executing statement ", err)
		return nil, err
	}

	sw := getSwearwords(out.Items)
	return sw, nil
}

func getSwearwords(items []map[string]*dynamodb.AttributeValue) []string {
	sw := []string{}
	for _, item := range items {
		sw = append(sw, *item["word"].S)
	}
	return sw
}

func getParameters(words *[]string) []*dynamodb.AttributeValue {
	parameters := make([]*dynamodb.AttributeValue, len(*words))
	for i, word := range *words {
		parameters[i] = &dynamodb.AttributeValue{
			S: aws.String(word),
		}
	}

	return parameters
}

func distinct(words []string) []string {
	m := make(map[string]bool)
	a := []string{}
	for _, word := range words {
		if m[word] {
			continue
		}
		m[word] = true
		a = append(a, word)
	}
	return a
}
